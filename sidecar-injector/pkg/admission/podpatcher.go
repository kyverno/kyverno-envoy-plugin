package admission

import (
	"context"

	corev1 "k8s.io/api/core/v1"
)

type PodPatcher interface {
	PatchPodCreate(ctx context.Context, namespace string, pod corev1.Pod) ([]PatchOperation, error)
}
