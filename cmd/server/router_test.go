package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHealthzReturnsOKJSON(t *testing.T) {
	r := NewRouter()
	ctx := context.Background()

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("/healthz status = %d, want %d", rr.Code, http.StatusOK)
	}

	contentType := rr.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		t.Fatalf("/healthz Content-Type = %q, want application/json prefix", contentType)
	}

	body, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	var payload map[string]string
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("/healthz body is not valid JSON: %v (body=%q)", err, body)
	}
	if payload["status"] != "ok" {
		t.Fatalf("/healthz body.status = %q, want %q", payload["status"], "ok")
	}
}

func TestUnknownRouteReturns404(t *testing.T) {
	r := NewRouter()
	ctx := context.Background()

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/does-not-exist", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("unknown route status = %d, want %d", rr.Code, http.StatusNotFound)
	}
}

func TestApiV1GroupIsMountedAndEmpty(t *testing.T) {
	r := NewRouter()
	ctx := context.Background()

	req := httptest.NewRequestWithContext(ctx, http.MethodGet, "/api/v1/anything", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Fatalf("/api/v1/anything status = %d, want 404 (group is a placeholder)", rr.Code)
	}
}

func TestHealthzRejectsNonGETMethods(t *testing.T) {
	r := NewRouter()

	for _, method := range []string{http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequestWithContext(context.Background(), method, "/healthz", nil)
			rr := httptest.NewRecorder()

			r.ServeHTTP(rr, req)

			if rr.Code != http.StatusMethodNotAllowed {
				t.Fatalf("/healthz %s status = %d, want %d", method, rr.Code, http.StatusMethodNotAllowed)
			}
		})
	}
}
