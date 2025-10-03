package sidecar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
)

func TestInject(t *testing.T) {
	sidecar := &Sidecar{
		InitContainers: []corev1.Container{{
			Name:  "init",
			Image: "foo:init",
		}},
		Containers: []corev1.Container{{
			Name:  "container",
			Image: "foo:container",
		}},
		ImagePullSecrets: []corev1.LocalObjectReference{{
			Name: "foo",
		}},
		Volumes: []corev1.Volume{{
			Name: "foo",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "foo",
					},
				},
			},
		}},
	}
	tests := []struct {
		name      string
		pod       corev1.Pod
		container *Sidecar
		want      corev1.Pod
	}{{
		name:      "empty",
		pod:       corev1.Pod{},
		container: sidecar,
		want: corev1.Pod{
			Spec: corev1.PodSpec{
				InitContainers:   sidecar.InitContainers,
				Containers:       sidecar.Containers,
				Volumes:          sidecar.Volumes,
				ImagePullSecrets: sidecar.ImagePullSecrets,
			},
		},
	}, {
		name:      "nil",
		pod:       corev1.Pod{},
		container: nil,
		want:      corev1.Pod{},
	}, {
		name: "not empty",
		pod: corev1.Pod{
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{{
					Name: "not-foo",
				}},
				Containers: []corev1.Container{{
					Name: "not-foo",
				}},
				ImagePullSecrets: []corev1.LocalObjectReference{{
					Name: "not-foo",
				}},
				Volumes: []corev1.Volume{{
					Name: "not-foo",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "foo",
							},
						},
					},
				}},
			},
		},
		container: sidecar,
		want: corev1.Pod{
			Spec: corev1.PodSpec{
				InitContainers: []corev1.Container{{
					Name: "not-foo",
				}, {
					Name:  "init",
					Image: "foo:init",
				}},
				Containers: []corev1.Container{{
					Name: "not-foo",
				}, {
					Name:  "container",
					Image: "foo:container",
				}},
				ImagePullSecrets: []corev1.LocalObjectReference{{
					Name: "not-foo",
				}, {
					Name: "foo",
				}},
				Volumes: []corev1.Volume{{
					Name: "not-foo",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "foo",
							},
						},
					},
				}, {
					Name: "foo",
					VolumeSource: corev1.VolumeSource{
						ConfigMap: &corev1.ConfigMapVolumeSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "foo",
							},
						},
					},
				}},
			},
		},
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Inject(tt.pod, tt.container)
			assert.Equal(t, tt.want, got)
		})
	}
}
