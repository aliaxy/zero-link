//go:build integration

package integration_test

import (
	"net/http"
	"testing"
)

func TestHealth_Liveness(t *testing.T) {
	resp := doRequest(t, http.MethodGet, "/healthz", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	decodeJSON(t, resp, &body)
	if body.Status != "ok" {
		t.Fatalf("status = %q, want %q", body.Status, "ok")
	}
	if body.Message != "api alive" {
		t.Fatalf("message = %q, want %q", body.Message, "api alive")
	}
}

func TestHealth_Readiness(t *testing.T) {
	// /readyz calls link-rpc gRPC health check — proves the full chain is up.
	resp := doRequest(t, http.MethodGet, "/readyz", nil, "")
	assertStatus(t, resp, http.StatusOK)

	var body struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	decodeJSON(t, resp, &body)
	if body.Status != "ok" {
		t.Fatalf("status = %q, want %q", body.Status, "ok")
	}
	if body.Message != "api and rpc ready" {
		t.Fatalf("message = %q, want %q", body.Message, "api and rpc ready")
	}
}
