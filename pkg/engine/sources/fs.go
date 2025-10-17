package sources

import (
	"context"
	"fmt"
	"io/fs"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/data"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/compiler"
	"github.com/kyverno/kyverno-envoy-plugin/sdk/core/sources"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/pkg/ext/file"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var (
	vpolGVK = vpol.SchemeGroupVersion.WithKind("ValidatingPolicy")
)

type document = []byte

func defaultLoader(_fs func() (fs.FS, error)) (loader.Loader, error) {
	if _fs == nil {
		_fs = data.Crds
	}
	crdsFs, err := _fs()
	if err != nil {
		return nil, err
	}
	return loader.New(openapiclient.NewLocalCRDFiles(crdsFs))
}

var DefaultLoader = sync.OnceValues(func() (loader.Loader, error) { return defaultLoader(nil) })

func NewFs(compiler vpolcompiler.Compiler, f fs.FS) engine.Source {
	input := sources.NewFs(f, func(_ string, entry fs.DirEntry) bool {
		if entry == nil {
			return false
		}
		// process only files
		if entry.IsDir() {
			return false
		}
		return file.IsYaml(entry.Name()) || file.IsJson(entry.Name())
	})
	transform := sources.NewTransformErr(
		input,
		func(entry sources.FsEntry) ([]document, error) {
			return getDocuments(context.Background(), f, entry.DirEntry)
		},
	)
	flatten := sources.NewFlatten(transform)
	load := sources.NewTransformErr(
		flatten,
		func(document document) (*vpol.ValidatingPolicy, error) {
			ldr, err := DefaultLoader()
			if err != nil {
				return nil, fmt.Errorf("failed to load CRDs: %w", err)
			}
			gvk, untyped, err := ldr.Load(document)
			if err != nil {
				return nil, err
			}
			switch gvk {
			case vpolGVK:
				typed, err := convert.To[vpol.ValidatingPolicy](untyped)
				if err != nil {
					return nil, fmt.Errorf("failed to convert to ValidatingPolicy: %w", err)
				}
				return typed, nil
			}
			return nil, nil
		})
	filter := sources.NewFilter(
		load,
		func(p *vpol.ValidatingPolicy) bool {
			return p != nil
		},
	)
	// TODO: sort by policy name
	compile := sources.NewTransformErr(
		filter,
		func(p *vpol.ValidatingPolicy) (engine.Policy, error) {
			c, errs := compiler.Compile(p)
			if len(errs) > 0 {
				return nil, fmt.Errorf("failed to compile ValidatingPolicy: %w", errs.ToAggregate())
			}
			return c, nil
		},
	)
	return compile
}

func getDocuments(_ context.Context, f fs.FS, entry fs.DirEntry) ([]document, error) {
	if entry == nil {
		return nil, nil
	}
	// process only files
	if entry.IsDir() {
		return nil, nil
	}
	// if it's a yaml file, it can contain multiple documents
	if file.IsYaml(entry.Name()) {
		bytes, err := fs.ReadFile(f, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}
		documents, err := yaml.SplitDocuments(bytes)
		if err != nil {
			return nil, fmt.Errorf("failed to split documents: %w", err)
		}
		return documents, nil
	}
	// if it's a json file, it contains a single document
	if file.IsJson(entry.Name()) {
		doc, err := fs.ReadFile(f, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}
		return []document{doc}, nil
	}
	return nil, nil
}
