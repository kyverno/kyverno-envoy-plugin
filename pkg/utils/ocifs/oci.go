package ocifs

import (
	"context"
	"fmt"
	"io/fs"
	"net/url"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/hairyhenderson/go-fsimpl"
	"github.com/nlepage/go-tarfs"
)

type ociFS struct {
	ctx context.Context

	repo *url.URL
	img  *v1.Image

	nOpts []name.Option
	rOpts []remote.Option

	tarFS fs.FS
}

var (
	_ fs.FS = (*ociFS)(nil)
)

func (f *ociFS) URL() string {
	return f.repo.String()
}

// ConfigureOCIFS is a function that configures the OCI filesystem with registry options
// and returns a fsimpl.FSProvider
func ConfigureOCIFS(nOpts []name.Option, rOpts []remote.Option) fsimpl.FSProvider {
	return fsimpl.FSProviderFunc(func(u *url.URL) (fs.FS, error) {
		fsys, err := New(u, nOpts, rOpts)
		if err != nil {
			return nil, err
		}
		// Configure the filesystem with registry options
		if ociFS, ok := fsys.(interface {
			WithRegistryOpts([]name.Option, []remote.Option) fs.FS
		}); ok {
			return ociFS.WithRegistryOpts(nOpts, rOpts), nil
		}
		return fsys, nil
	}, "oci")
}

func New(u *url.URL, nOpts []name.Option, rOpts []remote.Option) (fs.FS, error) {
	repoUrl := *u

	fsys := &ociFS{
		ctx:   context.Background(),
		repo:  &repoUrl,
		nOpts: nOpts,
		rOpts: rOpts,
	}

	img, err := fsys.fetchManifest()
	if err != nil {
		return nil, err
	}

	fsys.img = &img

	// Read content of flattened artifact filesystem
	rc := mutate.Extract(img)
	fsys.tarFS, err = tarfs.New(rc)
	if err != nil {
		return nil, err
	}

	return fsys, nil
}

func (f *ociFS) WithContext(ctx context.Context) fs.FS {
	if ctx == nil {
		return f
	}

	fsys := *f
	fsys.ctx = ctx

	return &fsys
}

func (f *ociFS) WithRegistryOpts(nOpts []name.Option, rOpts []remote.Option) fs.FS {
	fsys := *f
	fsys.nOpts = nOpts
	fsys.rOpts = rOpts

	return &fsys
}

func (f *ociFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, &fs.PathError{Op: "open", Path: name, Err: fs.ErrInvalid}
	}

	ociFile, err := f.tarFS.Open(name)
	if err != nil {
		return nil, &fs.PathError{Op: "open", Path: name, Err: err}
	}

	return ociFile, nil
}

func (f *ociFS) fetchManifest() (v1.Image, error) {
	ref := fmt.Sprintf("%s%s", f.repo.Host, f.repo.Path)

	repoRef, err := name.ParseReference(ref, f.nOpts...)
	if err != nil {
		return nil, err
	}

	img, err := remote.Image(repoRef, f.rOpts...)
	if err != nil {
		return nil, err
	}
	f.img = &img

	return img, err
}
