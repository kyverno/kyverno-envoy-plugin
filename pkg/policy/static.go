package policy

import (
	"context"
	"fmt"
	"github.com/kyverno/pkg/ext/file"
	"io/fs"
	"sync"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"github.com/kyverno/kyverno-envoy-plugin/pkg/data"
	"github.com/kyverno/pkg/ext/resource/convert"
	"github.com/kyverno/pkg/ext/resource/loader"
	"github.com/kyverno/pkg/ext/yaml"
	"sigs.k8s.io/kubectl-validate/pkg/openapiclient"
)

type staticProvider struct {
	compiled []CompiledPolicy
	err      error
}

func NewStaticProvider(compiler Compiler, policies ...v1alpha1.AuthorizationPolicy) Provider {
	var compiled []CompiledPolicy
	for _, policy := range policies {
		policy, err := compiler.Compile(&policy)
		if err != nil {
			return &staticProvider{err: err.ToAggregate()}
		}
		compiled = append(compiled, policy)
	}
	return &staticProvider{compiled: compiled}
}

func (p *staticProvider) CompiledPolicies(ctx context.Context) ([]CompiledPolicy, error) {
	// TODO: sort based on policy names
	return p.compiled, p.err
}

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

func NewFsProvider(compiler Compiler, f fs.FS) Provider {
	var policies []v1alpha1.AuthorizationPolicy

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

			typed, err := convert.To[v1alpha1.AuthorizationPolicy](untyped)
			if err != nil {
				return &staticProvider{err: fmt.Errorf("failed to convert to AuthorizationPolicy: %w", err)}
			}

			policies = append(policies, *typed)
		}
	}
	return NewStaticProvider(compiler, policies...)
}
