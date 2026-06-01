//go:build integration

package integration_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// noRedirectHTTPClient never follows 302 redirects.
// Used for waitForHTTP and for redirect-path tests.
var noRedirectHTTPClient = &http.Client{
	Timeout: 10 * time.Second,
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

// doLogin calls POST /admin/login and returns the JWT token.
// It fails the test immediately if the response is not 200 or token is empty.
func doLogin(t *testing.T, username, password string) string {
	t.Helper()
	resp := doRequest(t, http.MethodPost, "/admin/login", map[string]string{
		"username": username,
		"password": password,
	}, "")
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Fatalf("doLogin: status %d, body: %s", resp.StatusCode, body)
	}
	var envelope struct {
		Data struct {
			Token string `json:"token"`
		} `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		t.Fatalf("doLogin: decode: %v", err)
	}
	if envelope.Data.Token == "" {
		t.Fatal("doLogin: empty token in response")
	}
	return envelope.Data.Token
}

// doRequest sends an HTTP request to apiBase+path using the default client
// (follows redirects). body is JSON-encoded when non-nil; token sets Bearer auth.
// The response body is NOT closed — caller must close it.
func doRequest(t *testing.T, method, path string, body any, token string) *http.Response {
	t.Helper()
	var reqBody io.Reader
	if body != nil {
		reqBody = mustJSON(t, body)
	}
	req, err := http.NewRequest(method, apiBase+path, reqBody)
	if err != nil {
		t.Fatalf("doRequest: NewRequest: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("doRequest %s %s: %v", method, path, err)
	}
	return resp
}

// doRequestNoRedirect is identical to doRequest but uses noRedirectHTTPClient.
// Use for GET /:code tests where the 302 must not be auto-followed.
func doRequestNoRedirect(t *testing.T, method, path string, body any, token string) *http.Response {
	t.Helper()
	var reqBody io.Reader
	if body != nil {
		reqBody = mustJSON(t, body)
	}
	req, err := http.NewRequest(method, apiBase+path, reqBody)
	if err != nil {
		t.Fatalf("doRequestNoRedirect: NewRequest: %v", err)
	}
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	resp, err := noRedirectHTTPClient.Do(req)
	if err != nil {
		t.Fatalf("doRequestNoRedirect %s %s: %v", method, path, err)
	}
	return resp
}

// mustJSON encodes v to JSON and returns an io.Reader.
func mustJSON(t *testing.T, v any) io.Reader {
	t.Helper()
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("mustJSON: %v", err)
	}
	return bytes.NewReader(data)
}

// assertStatus fails the test if resp.StatusCode != want.
// It reads and closes the body before failing so the error message is useful.
func assertStatus(t *testing.T, resp *http.Response, want int) {
	t.Helper()
	if resp.StatusCode != want {
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		t.Fatalf("status = %d, want %d; body: %s", resp.StatusCode, want, body)
	}
}

// decodeJSON decodes the response body into v and closes the body.
func decodeJSON(t *testing.T, resp *http.Response, v any) {
	t.Helper()
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
		t.Fatalf("decodeJSON: %v", err)
	}
}

// createLink calls POST /admin/links and returns (id, code).
// It fails the test if the response is not 200.
func createLink(t *testing.T, token, originURL string) (id int64, code string) {
	t.Helper()
	resp := doRequest(t, http.MethodPost, "/admin/links", map[string]string{
		"origin_url": originURL,
	}, token)
	assertStatus(t, resp, http.StatusOK)
	var envelope struct {
		Data struct {
			Id   int64  `json:"id"`
			Code string `json:"code"`
		} `json:"data"`
	}
	decodeJSON(t, resp, &envelope)
	if envelope.Data.Id == 0 {
		t.Fatal("createLink: id is 0")
	}
	if envelope.Data.Code == "" {
		t.Fatal("createLink: code is empty")
	}
	return envelope.Data.Id, envelope.Data.Code
}

// deleteLink calls DELETE /admin/links/:id.
func deleteLink(t *testing.T, token string, id int64) {
	t.Helper()
	resp := doRequest(t, http.MethodDelete, fmt.Sprintf("/admin/links/%d", id), nil, token)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		t.Errorf("deleteLink %d: status %d, body: %s", id, resp.StatusCode, body)
	}
}

// dbExec opens a raw DB connection, runs the query, and closes the connection.
// Used for analytics data cleanup (no HTTP endpoint for visit_event / link_daily_stat).
func dbExec(t *testing.T, query string, args ...any) {
	t.Helper()
	db, err := sql.Open("mysql", testDSN())
	if err != nil {
		t.Errorf("dbExec open: %v", err)
		return
	}
	defer db.Close()
	if _, err := db.Exec(query, args...); err != nil {
		t.Errorf("dbExec %q: %v", query, err)
	}
}
