package policy

import (
	"context"
	"io/fs"
	"log"
	"path/filepath"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
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

func NewFsProvider(compiler Compiler, f fs.FS) Provider {
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
					log.Println(entry.Name())
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
	return NewStaticProvider(compiler)
}
