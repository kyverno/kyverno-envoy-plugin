package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	admissionv1 "k8s.io/api/admission/v1"
)

func AdmissionReview(inner func(context.Context, *admissionv1.AdmissionRequest) *admissionv1.AdmissionResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Body == nil {
			HttpError(r.Context(), w, r, errors.New("empty body"), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()
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
