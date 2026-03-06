package handler

import (
	"encoding/json"
	"net/http"

	"github.com/devprep/backend/internal/apperror"
	"github.com/go-chi/chi/v5/middleware"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

type errorResponse struct {
	RequestID string `json:"request_id,omitempty"`
	Code      string `json:"code"`
	Message   string `json:"message"`
	Details   any    `json:"details,omitempty"`
}

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	ae := apperror.As(err)
	if ae == nil {
		ae = apperror.Internal(err)
	}

	writeJSON(w, ae.HTTPStatus, errorResponse{
		RequestID: middleware.GetReqID(r.Context()),
		Code:      ae.Code,
		Message:   ae.Message,
		Details:   ae.Details,
	})
}
