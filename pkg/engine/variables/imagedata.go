package variables

import (
	"context"

	"github.com/kyverno/kyverno/pkg/cel/utils"
	"github.com/kyverno/kyverno/pkg/imageverification/imagedataloader"
	v1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

func ImageData(lister v1.SecretInterface, imageOpts ...imagedataloader.Option) (*imageData, error) {
	// TODO: secrets interface
	idl, err := imagedataloader.New(lister, imageOpts...)
	if err != nil {
		return nil, err
	}
	return &imageData{
		imagedata: idl,
	}, nil
}

type imageData struct {
	imagedata imagedataloader.Fetcher
}

func (cp *imageData) GetImageData(image string) (map[string]any, error) {
	// TODO: get image credentials from image verification policies?
	data, err := cp.imagedata.FetchImageData(context.TODO(), image)
	if err != nil {
		return nil, err
	}
	return utils.GetValue(data.Data())
}
