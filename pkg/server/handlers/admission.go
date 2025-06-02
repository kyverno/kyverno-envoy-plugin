package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gomodules.xyz/jsonpatch/v2"
	admissionv1 "k8s.io/api/admission/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
)

func AdmissionResponse(r *admissionv1.AdmissionRequest, err error, patch ...jsonpatch.Operation) *admissionv1.AdmissionResponse {
	response := admissionv1.AdmissionResponse{
		UID: r.UID,
	}
	var patchBytes []byte
	if err == nil {
		if len(patch) != 0 {
			patchBytes, err = json.Marshal(patch)
		}
	}
	if err != nil {
		response.Allowed = false
		response.Result = &metav1.Status{
			Status:  metav1.StatusFailure,
			Message: err.Error(),
		}
	} else {
		response.Allowed = true
		response.Result = &metav1.Status{
			Status: metav1.StatusSuccess,
		}
		if patchBytes != nil {
			response.PatchType = ptr.To(admissionv1.PatchTypeJSONPatch)
			response.Patch = patchBytes
		}
	}
	return &response
}

func AdmissionReview(inner func(context.Context, *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			HttpError(r.Context(), w, r, errors.New("empty body"), http.StatusBadRequest)
			return
		}
		defer r.Body.Close() //nolint:errcheck
		body, err := io.ReadAll(r.Body)
		if err != nil {
			HttpError(r.Context(), w, r, err, http.StatusBadRequest)
			return
		}
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			HttpError(r.Context(), w, r, errors.New("invalid Content-Type"), http.StatusUnsupportedMediaType)
			return
		}
		var admissionReview admissionv1.AdmissionReview
		if err := json.Unmarshal(body, &admissionReview); err != nil {
			HttpError(r.Context(), w, r, err, http.StatusExpectationFailed)
			return
		}
		admissionResponse := inner(r.Context(), admissionReview.Request)
		admissionReview.Response = admissionResponse
		responseJSON, err := json.Marshal(admissionReview)
		if err != nil {
			HttpError(r.Context(), w, r, err, http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		if _, err := w.Write(responseJSON); err != nil {
			HttpError(r.Context(), w, r, err, http.StatusInternalServerError)
			return
		}
	}
}
