package admission

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	admissionv1 "k8s.io/api/admission/v1"
	corev1 "k8s.io/api/core/v1"
)

// PodAdmissionRequestHandler PodAdmissionRequest handler
type PodAdmissionRequestHandler struct {
	PodHandler PodPatcher
}

func (handler *PodAdmissionRequestHandler) handleAdmissionCreate(ctx context.Context, request *admissionv1.AdmissionRequest) ([]PatchOperation, error) {
	pod, err := unmarshalPod(request.Object.Raw)
	if err != nil {
		return nil, err
	}
	return handler.PodHandler.PatchPodCreate(ctx, pod)
}

func unmarshalPod(rawObject []byte) (corev1.Pod, error) {
	var pod corev1.Pod
	err := json.Unmarshal(rawObject, &pod)
	return pod, errors.Wrapf(err, "error unmarshalling object")
}
