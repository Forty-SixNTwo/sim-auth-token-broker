package utilities

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

const (
	REQUEST_ID_HEADER_KEY = "X-Request-ID"
)

type CtxRequestID struct{}

type ErrorResponse struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description,omitempty"`
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func WriteJSONError(w http.ResponseWriter, msg, desc string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: msg, ErrorDescription: desc})
}

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(REQUEST_ID_HEADER_KEY)
		if id == "" {
			id = uuid.NewString()
			w.Header().Set(REQUEST_ID_HEADER_KEY, id)
		}
		ctx := context.WithValue(r.Context(), CtxRequestID{}, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDTransport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	return roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		if id, ok := req.Context().Value(CtxRequestID{}).(string); ok {
			req.Header.Set(REQUEST_ID_HEADER_KEY, id)
		}
		return base.RoundTrip(req)
	})
}

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}
