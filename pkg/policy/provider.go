package policy

import (
	"context"

	"github.com/kyverno/kyverno-envoy-plugin/apis/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Provider interface {
	CompiledPolicies(context.Context) ([]PolicyFunc, error)
}

func NewKubeProvider(ctx context.Context, config *rest.Config, compiler Compiler) (Provider, error) {
	scheme := runtime.NewScheme()
	if err := v1alpha1.Install(scheme); err != nil {
		return nil, err
	}
	// create kubernetes client
	// TODO: do we want to use a cache ?
	cache, err := cache.New(config, cache.Options{
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	client, err := client.New(config, client.Options{
		Cache: &client.CacheOptions{
			Reader: cache,
		},
		Scheme: scheme,
	})
	if err != nil {
		return nil, err
	}
	go func() {
		if err := cache.Start(ctx); err != nil {
			// TODO: better error handling
			panic(err)
		}
	}()
	// TODO: use the result of the wait
	cache.WaitForCacheSync(ctx)
	return &kubeProvider{
		client:   client,
		compiler: compiler,
	}, nil
}

type kubeProvider struct {
	client   client.Client
	compiler Compiler
}

func (p *kubeProvider) CompiledPolicies(ctx context.Context) ([]PolicyFunc, error) {
	// fetch policies
	var policies v1alpha1.AuthorizationPolicyList
	if err := p.client.List(ctx, &policies, &client.ListOptions{}); err != nil {
		return nil, err
	}
	var out []PolicyFunc
	for _, policy := range policies.Items {
		compiled, err := p.compiler.Compile(policy)
		if err != nil {
			return nil, err
		}
		out = append(out, compiled)
	}
	return out, nil
}
