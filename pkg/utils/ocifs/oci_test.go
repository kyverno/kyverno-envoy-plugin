package ocifs

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"testing"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/types"
	"github.com/nlepage/go-tarfs"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Valid URL",
			url:     "oci://ghcr.io/lucchmielowski/ivpol:multilayered",
			wantErr: false,
		},
		{
			name:    "Invalid URL",
			url:     "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := url.Parse(tt.url)
			if err != nil != tt.wantErr {
				t.Errorf("url.Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}

			got, err := New(u, []name.Option{}, []remote.Option{})
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got == nil {
				t.Error("New() = nil, want non-nil")
			}
		})
	}
}

func TestOciFS_WithContext(t *testing.T) {
	u, _ := url.Parse("oci://ghcr.io/lucchmielowski/ivpol:multilayered")
	fsys := &ociFS{
		ctx:  context.Background(),
		repo: u,
	}

	tests := []struct {
		name string
		ctx  context.Context
		want bool
	}{
		{
			name: "With valid context",
			ctx:  context.Background(),
			want: true,
		},
		{
			name: "With nil context",
			ctx:  nil,
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fsys.WithContext(tt.ctx)
			if (got != nil) != tt.want {
				t.Errorf("WithContext() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOciFS_Open(t *testing.T) {
	ru, err := url.Parse("oci://ghcr.io/lucchmielowski/ivpol:multilayered")
	if err != nil {
		t.Fatal(err)
	}

	ref := fmt.Sprintf("%s%s", ru.Host, ru.Path)
	if err != nil {
		t.Fatal(err)
	}

	repoRef, err := name.ParseReference(ref)
	if err != nil {
		t.Fatal(err)
	}

	img, err := remote.Image(repoRef, remote.WithAuthFromKeychain(authn.DefaultKeychain))
	if err != nil {
		t.Fatal(err)
	}

	// Read content of flattened artifact filesystem
	setup := func() *ociFS {
		rc := mutate.Extract(img)
		tfs, err := tarfs.New(rc)
		if err != nil {
			log.Fatal(err)
		}

		return &ociFS{
			ctx:   context.Background(),
			tarFS: tfs,
			repo:  ru,
		}
	}

	tests := []struct {
		name    string
		setup   func() *ociFS
		path    string
		wantErr bool
	}{
		{
			name:    "Invalid path with ..",
			setup:   setup,
			path:    "../test",
			wantErr: true,
		},
		{
			name:    "Empty path",
			setup:   setup,
			path:    "",
			wantErr: true,
		},
		{
			name:    "Current directory",
			setup:   setup,
			path:    ".",
			wantErr: false,
		},
		{
			name:    "Layered directory",
			setup:   setup,
			path:    "dir/sample.yaml",
			wantErr: false,
		},
		{
			name:    "Layered directory with dot notation",
			setup:   setup,
			path:    "./dir/sample.yaml",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := tt.setup()
			_, err := fsys.Open(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Open() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestOciFS_URL(t *testing.T) {
	testURL := "oci://ghcr.io/test/image:latest"
	u, _ := url.Parse(testURL)
	fsys := &ociFS{
		ctx:  context.Background(),
		repo: u,
	}

	if got := fsys.URL(); got != testURL {
		t.Errorf("URL() = %v, want %v", got, testURL)
	}
}

func TestFetchManifest(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		wantErr bool
	}{
		{
			name:    "Invalid reference",
			url:     "oci://invalid@registry/test",
			wantErr: true,
		},
		{
			name:    "Valid reference format",
			url:     "oci://ghcr.io/lucchmielowski/ivpol:multilayered",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, _ := url.Parse(tt.url)
			fsys := &ociFS{
				ctx:  context.Background(),
				repo: u,
			}

			_, err := fsys.fetchManifest()
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFS_Implementation(t *testing.T) {
	var _ fs.FS = (*ociFS)(nil)
}

type mockImage struct{}

func (m *mockImage) Digest() (v1.Hash, error) {
	return v1.Hash{}, nil
}

func (m *mockImage) ConfigFile() (*v1.ConfigFile, error) {
	return nil, nil
}

func (m *mockImage) ConfigName() (v1.Hash, error) {
	return v1.Hash{}, nil
}

func (m *mockImage) Layers() ([]v1.Layer, error) {
	return nil, nil
}

func (m *mockImage) LayerByDigest(v1.Hash) (v1.Layer, error) {
	return nil, nil
}

func (m *mockImage) LayerByDiffID(v1.Hash) (v1.Layer, error) {
	return nil, nil
}

func (m *mockImage) Manifest() (*v1.Manifest, error) {
	return nil, nil
}

func (m *mockImage) MediaType() (types.MediaType, error) {
	return "", nil
}

func (m *mockImage) RawManifest() ([]byte, error) {
	return nil, nil
}

func (m *mockImage) Size() (int64, error) {
	return 0, nil
}
