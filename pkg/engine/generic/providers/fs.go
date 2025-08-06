package providers

import (
	"fmt"
	"io/fs"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/pkg/data"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/engine"
	"github.com/kyverno/pkg/ext/file"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

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

func NewFsProvider[T any](compiler engine.Compiler[T], f fs.FS) engine.Provider {
	var policies []T

	entries, err := fs.ReadDir(f, ".")
	if err != nil {
		return &staticProvider{err: err}
	}

	for _, entry := range entries {
		if entry.IsDir() || !file.IsYaml(entry.Name()) {
			continue
		}

		bytes, err := fs.ReadFile(f, entry.Name())
		if err != nil {
			return &staticProvider{err: fmt.Errorf("failed to read file %s: %w", entry.Name(), err)}
		}

		documents, err := yaml.SplitDocuments(bytes)
		if err != nil {
			return &staticProvider{err: fmt.Errorf("failed to split documents: %w", err)}
		}

		ldr, err := DefaultLoader()
		if err != nil {
			return &staticProvider{err: fmt.Errorf("failed to load CRDs: %w", err)}
		}

		for _, document := range documents {
			_, untyped, err := ldr.Load(document)
			if err != nil {
				// Not an AuthorizationPolicy, skip
				continue
			}

			typed, err := convert.To[T](untyped)
			if err != nil {
				return &staticProvider{err: fmt.Errorf("failed to convert to AuthorizationPolicy: %w", err)}
			}

			policies = append(policies, *typed)
		}
	}
	return NewStaticProvider(compiler, policies...)
}
