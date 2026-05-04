package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// AR13: handlers stay thin and call service packages, which call generated db queries.
func NewRouter() chi.Router {
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})
	r.Route("/api/v1", func(r chi.Router) {})
	return r
}
