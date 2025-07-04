package policy

import (
	"context"
	"io/fs"
	"path/filepath"
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
	if f, ok := f.(fs.ReadDirFS); ok {
		entries, err := f.ReadDir(".")
		if err != nil {
			// TODO: proper error handling
			return nil
		}
		for _, entry := range entries {
			if !entry.IsDir() {
				ext := filepath.Ext(entry.Name())
				if ext == ".yml" || ext == ".yaml" {
					bytes, err := fs.ReadFile(f, entry.Name())
					if err != nil {
						return nil
					}
					documents, err := yaml.SplitDocuments(bytes)
					if err != nil {
						return nil
					}
					for _, document := range documents {
						loader, err := DefaultLoader()
						if err != nil {
							return nil
						}
						_, untyped, err := loader.Load(document)
						if err != nil {
							return nil
						}
						typed, err := convert.To[v1alpha1.AuthorizationPolicy](untyped)
						if err != nil {
							return nil
						}
						policies = append(policies, *typed)
					}
				}
				// TODO: json
			}
		}
	}
	// TODO: see https://github.com/hairyhenderson/go-fsimpl/issues/1079
	// if f, ok := f.(fs.ReadFileFS); ok {
	// 	data, err := f.ReadFile("..")
	// 	if err != nil {
	// 		// TODO: proper error handling
	// 		return nil
	// 	}
	// 	log.Println(string(data))
	// }
	return NewStaticProvider(compiler, policies...)
}
