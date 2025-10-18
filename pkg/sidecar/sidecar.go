package sidecar

import (
	corev1 "k8s.io/api/core/v1"
)

func Sidecar(image string, controlPlaneAddr string,
	controlPlaneReconnectWait, controlPlaneMaxDialInterval, healthCheckInterval string) corev1.Container {
	container := corev1.Container{
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
			"--control-plane-address=" + controlPlaneAddr,
			"--control-plane-reconnect-wait=" + controlPlaneReconnectWait,
			"--control-plane-max-dial-interval=" + controlPlaneMaxDialInterval,
			"--health-check-interval=" + healthCheckInterval,
		},
		Env: []corev1.EnvVar{
			{
				Name: "POD_IP",
				ValueFrom: &corev1.EnvVarSource{
					FieldRef: &corev1.ObjectFieldSelector{
						FieldPath: "status.podIP",
					},
				},
			},
		},
	}
	return container
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
