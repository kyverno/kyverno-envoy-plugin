package sidecar

import (
	"os"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/yaml"
)

type Sidecar struct {
	InitContainers   []corev1.Container            `json:"initContainers,omitempty"`
	Containers       []corev1.Container            `json:"containers,omitempty"`
	Volumes          []corev1.Volume               `json:"volumes,omitempty"`
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
}

func Load(file string) (*Sidecar, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var sidecar *Sidecar
	if err := yaml.UnmarshalStrict(data, &sidecar); err != nil {
		return nil, err
	}
	return sidecar, nil
}

func Inject(pod corev1.Pod, sidecar *Sidecar) corev1.Pod {
	if sidecar != nil {
		pod.Spec.ImagePullSecrets = append(pod.Spec.ImagePullSecrets, sidecar.ImagePullSecrets...)
		pod.Spec.InitContainers = append(pod.Spec.InitContainers, sidecar.InitContainers...)
		pod.Spec.Containers = append(pod.Spec.Containers, sidecar.Containers...)
		pod.Spec.Volumes = append(pod.Spec.Volumes, sidecar.Volumes...)
	}
	return pod
}
