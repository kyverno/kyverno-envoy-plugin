package impl

import (
	"encoding/json"

	jsoncel "github.com/kyverno/kyverno-envoy-plugin/pkg/cel/libs/json"
)

type JsonImpl struct {
	jsoncel.JsonImpl
}

func (j *JsonImpl) Unmarshal(content []byte) (any, error) {
	var v any
	if err := json.Unmarshal(content, &v); err != nil {
		return nil, err
	}
	return v, nil
}
