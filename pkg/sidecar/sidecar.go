package sidecar

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"
)

func Sidecar(image string, externalPolicySources ...string) corev1.Container {
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
			"--grpc-address=:9081",
			"--metrics-address=:9082",
			"--kube-policy-source=false",
			"--external-policy-source=file:///data/kyverno-authz-server",
		},
		VolumeMounts: []corev1.VolumeMount{{
			Name:             "kyverno-authz-server",
			ReadOnly:         true,
			MountPropagation: ptr.To(corev1.MountPropagationHostToContainer),
			MountPath:        "/data/kyverno-authz-server",
		},
		},
	}
	for _, source := range externalPolicySources {
		container.Args = append(container.Args, "--external-policy-source="+source)
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
	pod.Spec.Volumes = append(pod.Spec.Volumes, corev1.Volume{
		Name: "kyverno-authz-server",
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: "kyverno-authz-server",
				},
				Optional: ptr.To(true),
			},
		},
	})
	return pod
}
