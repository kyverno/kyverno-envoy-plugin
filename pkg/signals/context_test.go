package signals

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	{
		err := Do(context.Background(), func(ctx context.Context) error {
			return nil
		})
		assert.NoError(t, err)
	}
	{
		err := Do(context.Background(), func(ctx context.Context) error {
			return errors.New("dummy")
		})
		assert.Error(t, err)
	}
}
