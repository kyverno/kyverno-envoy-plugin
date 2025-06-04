package sidecar

import (
	corev1 "k8s.io/api/core/v1"
)

func Sidecar(image string) corev1.Container {
	return corev1.Container{
		Name:            "kyverno-authz-server",
		ImagePullPolicy: corev1.PullIfNotPresent,
		Image:           image,
		Ports: []corev1.ContainerPort{{
			Name:          "http",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 9080,
		}, {
			Name:          "grpc",
			Protocol:      corev1.ProtocolTCP,
			ContainerPort: 9081,
		}},
		Args: []string{
			"serve",
			"authz-server",
			"--probes-address=:9080",
			"--grpc-address=:9081",
			"--metrics-address=:9082",
		},
	}
}

func Inject(pod corev1.Pod, container corev1.Container) corev1.Pod {
	for i, c := range pod.Spec.Containers {
		if c.Name == container.Name {
			pod.Spec.Containers[i] = container
			return pod
		}
	}
	pod.Spec.Containers = append(pod.Spec.Containers, container)
	return pod
}
