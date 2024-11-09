package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
	admissionregistrationv1 "k8s.io/api/admissionregistration/v1"
	"k8s.io/utils/ptr"
)

func TestAuthorizationPolicySpec_GetFailurePolicy(t *testing.T) {
	tests := []struct {
		name          string
		failurePolicy *admissionregistrationv1.FailurePolicyType
		want          admissionregistrationv1.FailurePolicyType
	}{{
		name:          "not set",
		failurePolicy: nil,
		want:          admissionregistrationv1.Fail,
	}, {
		name:          "fail",
		failurePolicy: ptr.To(admissionregistrationv1.Fail),
		want:          admissionregistrationv1.Fail,
	}, {
		name:          "ignore",
		failurePolicy: ptr.To(admissionregistrationv1.Ignore),
		want:          admissionregistrationv1.Ignore,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &AuthorizationPolicySpec{
				FailurePolicy: tt.failurePolicy,
			}
			got := s.GetFailurePolicy()
			assert.Equal(t, tt.want, got)
		})
	}
}
