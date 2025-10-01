package providers

import (
	"context"
	"fmt"
	"io/fs"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/data"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	apolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/apol/compiler"
	vpolcompiler "github.com/kyverno/kyverno-envoy-plugin/pkg/engine/vpol/compiler"
	vpol "github.com/kyverno/kyverno/api/policies.kyverno.io/v1alpha1"
	"github.com/kyverno/pkg/ext/file"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

var (
	apolGVK = v1alpha1.SchemeGroupVersion.WithKind("AuthorizationPolicy")
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

type fsProvider struct {
	apolCompiler apolcompiler.Compiler
	vpolCompiler vpolcompiler.Compiler
	fs           fs.FS
}

func NewFsProvider(apolCompiler apolcompiler.Compiler, vpolCompiler vpolcompiler.Compiler, fs fs.FS) engine.Provider {
	return &fsProvider{
		apolCompiler: apolCompiler,
		vpolCompiler: vpolCompiler,
		fs:           fs,
	}
}

func (p *fsProvider) CompiledPolicies(ctx context.Context) ([]engine.CompiledPolicy, error) {
	var policies []engine.CompiledPolicy
	entries, err := fs.ReadDir(p.fs, ".")
	if err != nil {
		return nil, err
	}
	ldr, err := DefaultLoader()
	if err != nil {
		return nil, fmt.Errorf("failed to load CRDs: %w", err)
	}
	for _, entry := range entries {
		// TODO: recursive loading
		documents, err := p.getDocuments(ctx, entry)
		if err != nil {
			return nil, fmt.Errorf("failed to extract documents: %w", err)
		}
		for _, document := range documents {
			gvk, untyped, err := ldr.Load(document)
			if err != nil {
				continue
			}
			switch gvk {
			case apolGVK:
				typed, err := convert.To[v1alpha1.AuthorizationPolicy](untyped)
				if err != nil {
					return nil, fmt.Errorf("failed to convert to AuthorizationPolicy: %w", err)
				}
				compiled, errs := p.apolCompiler.Compile(typed)
				if len(errs) > 0 {
					return nil, fmt.Errorf("failed to compile AuthorizationPolicy: %w", err)
				}
				policies = append(policies, compiled)
			case vpolGVK:
				typed, err := convert.To[vpol.ValidatingPolicy](untyped)
				if err != nil {
					return nil, fmt.Errorf("failed to convert to ValidatingPolicy: %w", err)
				}
				compiled, errs := p.vpolCompiler.Compile(typed)
				if len(errs) > 0 {
					return nil, fmt.Errorf("failed to compile ValidatingPolicy: %w", err)
				}
				policies = append(policies, compiled)
			}
		}
	}
	return policies, nil
}

func (p *fsProvider) getDocuments(ctx context.Context, entry fs.DirEntry) ([]document, error) {
	// process only files
	if entry.IsDir() {
		return nil, nil
	}
	// if it's a yaml file, it can contain multiple documents
	if file.IsYaml(entry.Name()) {
		bytes, err := fs.ReadFile(p.fs, entry.Name())
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
		doc, err := fs.ReadFile(p.fs, entry.Name())
		if err != nil {
			return nil, fmt.Errorf("failed to read file %s: %w", entry.Name(), err)
		}
		return []document{doc}, nil
	}
	return nil, nil
}
