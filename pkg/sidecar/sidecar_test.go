package sidecar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestSidecar(t *testing.T) {
	tests := []struct {
		name  string
		image string
		want  corev1.Container
	}{{
		image: "foo:bar",
		want: corev1.Container{
			Name:            "kyverno-authz-server",
			ImagePullPolicy: corev1.PullIfNotPresent,
			Image:           "foo:bar",
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
				"--http-address=:9080",
				"--grpc-address=:9081",
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Sidecar(tt.image)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInject(t *testing.T) {
	sidecar := Sidecar("foo:bar")
	tests := []struct {
		name      string
		pod       corev1.Pod
		container corev1.Container
		want      corev1.Pod
	}{{
		name:      "no containers",
		pod:       corev1.Pod{},
		container: sidecar,
		want: corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					sidecar,
				},
			},
		},
	}, {
		name: "found",
		pod: corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: sidecar.Name,
					},
				},
			},
		},
		container: sidecar,
		want: corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					sidecar,
				},
			},
		},
	}, {
		name: "not found",
		pod: corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "not-" + sidecar.Name,
					},
				},
			},
		},
		container: sidecar,
		want: corev1.Pod{
			Spec: corev1.PodSpec{
				Containers: []corev1.Container{
					{
						Name: "not-" + sidecar.Name,
					},
					sidecar,
				},
			},
		},
	},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Inject(tt.pod, tt.container)
			assert.Equal(t, tt.want, got)
		})
	}
}
