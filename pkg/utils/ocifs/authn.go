package ocifs

import (
	"context"
	"io"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/github"
	kauth "github.com/google/go-containerregistry/pkg/authn/kubernetes"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/google"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8scorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func RegistryOpts(lister k8scorev1.SecretInterface,
	insecure bool,
	secrets ...string) (
	[]remote.Option,
	[]name.Option, error) {
	rOpts := make([]remote.Option, 0)
	nOpts := make([]name.Option, 0)
	keychains := make([]authn.Keychain, 0)

	keychains = append(keychains,
		authn.NewKeychainFromHelper(credhelper.NewACRCredentialsHelper()),
		google.Keychain,
		authn.DefaultKeychain,
		github.Keychain,
		authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithLogger(io.Discard))),
	)

	if insecure {
		nOpts = append(nOpts, name.Insecure)
	}

	if len(secrets) != 0 {
		secretsKc, err := NewAutoRefreshSecretsKeychain(lister, secrets...)
		if err != nil {
			return nil, nil, err
		}
		keychains = append(keychains, secretsKc)
	}

	rOpts = append(rOpts, remote.WithAuthFromKeychain(authn.NewMultiKeychain(keychains...)))
	return rOpts, nOpts, nil
}

type autoRefreshSecrets struct {
	lister           k8scorev1.SecretInterface
	imagePullSecrets []string
}

func NewAutoRefreshSecretsKeychain(
	lister k8scorev1.SecretInterface,
	imagePullSecrets ...string,
) (authn.Keychain, error) {
	return &autoRefreshSecrets{
		lister:           lister,
		imagePullSecrets: imagePullSecrets,
	}, nil
}

func (kc *autoRefreshSecrets) Resolve(resource authn.Resource) (authn.Authenticator, error) {
	inner, err := generateKeyChainForPullSecrets(context.TODO(), kc.lister, kc.imagePullSecrets...)
	if err != nil {
		return nil, err
	}
	return inner.Resolve(resource)
}

func generateKeyChainForPullSecrets(ctx context.Context, lister k8scorev1.SecretInterface, imagePullSecrets ...string) (authn.Keychain, error) {
	var secrets []corev1.Secret

	for _, imagePullSecret := range imagePullSecrets {
		secret, err := lister.Get(ctx, imagePullSecret, metav1.GetOptions{})
		if err == nil {
			secrets = append(secrets, *secret)
		} else if !k8serrors.IsNotFound(err) {
			return nil, err
		}
	}
	return kauth.NewFromPullSecrets(context.TODO(), secrets)
}
