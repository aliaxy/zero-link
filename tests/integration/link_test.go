//go:build integration

package integration_test

import (
	"fmt"
	"net/http"
	"testing"
)

func TestLogin_Success(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	if len(token) < 10 {
		t.Fatalf("token looks invalid: %q", token)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	resp := doRequest(t, http.MethodPost, "/admin/login", map[string]string{
		"username": testAdminUsername,
		"password": "definitely-wrong-password",
	}, "")
	assertStatus(t, resp, http.StatusUnauthorized)

	var body struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}
	decodeJSON(t, resp, &body)
	if body.Code != "UNAUTHENTICATED" {
		t.Fatalf("code = %q, want %q", body.Code, "UNAUTHENTICATED")
	}
}

func TestProfile_Authenticated(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)

	resp := doRequest(t, http.MethodGet, "/admin/profile", nil, token)
	assertStatus(t, resp, http.StatusOK)

	var envelope struct {
		Data struct {
			Username string `json:"username"`
		} `json:"data"`
	}
	decodeJSON(t, resp, &envelope)
	if envelope.Data.Username != testAdminUsername {
		t.Fatalf("username = %q, want %q", envelope.Data.Username, testAdminUsername)
	}
}

func TestAuth_Required(t *testing.T) {
	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/admin/profile"},
		{http.MethodGet, "/admin/links"},
		{http.MethodPost, "/admin/links"},
		{http.MethodGet, "/admin/links/1"},
		{http.MethodPatch, "/admin/links/1"},
		{http.MethodDelete, "/admin/links/1"},
		{http.MethodGet, "/admin/links/1/stats"},
	}
	for _, ep := range endpoints {
		ep := ep
		t.Run(ep.method+" "+ep.path, func(t *testing.T) {
			resp := doRequest(t, ep.method, ep.path, nil, "")
			defer resp.Body.Close()
			assertStatus(t, resp, http.StatusUnauthorized)
		})
	}
}

func TestLink_FullCRUD(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	originURL := "https://example.com/crud-" + uniqueName(t)

	// CREATE
	id, code := createLink(t, token, originURL)
	t.Cleanup(func() { deleteLink(t, token, id) })

	if len(code) != 6 {
		t.Fatalf("auto-generated code length = %d, want 6", len(code))
	}

	// LIST — created link must appear
	{
		resp := doRequest(t, http.MethodGet, "/admin/links", nil, token)
		assertStatus(t, resp, http.StatusOK)
		var envelope struct {
			Data struct {
				Items []struct {
					Id int64 `json:"id"`
				} `json:"items"`
			} `json:"data"`
		}
		decodeJSON(t, resp, &envelope)
		found := false
		for _, item := range envelope.Data.Items {
			if item.Id == id {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("list: created link id=%d not found in response", id)
		}
	}

	// GET by ID
	{
		resp := doRequest(t, http.MethodGet, fmt.Sprintf("/admin/links/%d", id), nil, token)
		assertStatus(t, resp, http.StatusOK)
		var envelope struct {
			Data struct {
				Id        int64  `json:"id"`
				Code      string `json:"code"`
				OriginUrl string `json:"origin_url"`
				Status    int64  `json:"status"`
			} `json:"data"`
		}
		decodeJSON(t, resp, &envelope)
		if envelope.Data.Id != id {
			t.Fatalf("get: id = %d, want %d", envelope.Data.Id, id)
		}
		if envelope.Data.OriginUrl != originURL {
			t.Fatalf("get: origin_url = %q, want %q", envelope.Data.OriginUrl, originURL)
		}
		if envelope.Data.Status != 1 {
			t.Fatalf("get: status = %d, want 1 (active)", envelope.Data.Status)
		}
	}

	// UPDATE — change origin URL and title
	newURL := "https://example.com/updated-" + uniqueName(t)
	{
		resp := doRequest(t, http.MethodPatch, fmt.Sprintf("/admin/links/%d", id), map[string]string{
			"origin_url": newURL,
			"title":      "Updated Title",
		}, token)
		assertStatus(t, resp, http.StatusOK)
		var envelope struct {
			Data struct {
				OriginUrl string `json:"origin_url"`
				Title     string `json:"title"`
			} `json:"data"`
		}
		decodeJSON(t, resp, &envelope)
		if envelope.Data.OriginUrl != newURL {
			t.Fatalf("update: origin_url = %q, want %q", envelope.Data.OriginUrl, newURL)
		}
		if envelope.Data.Title != "Updated Title" {
			t.Fatalf("update: title = %q, want %q", envelope.Data.Title, "Updated Title")
		}
	}

	// DELETE (t.Cleanup above also deletes; this verifies the response shape)
	{
		resp := doRequest(t, http.MethodDelete, fmt.Sprintf("/admin/links/%d", id), nil, token)
		assertStatus(t, resp, http.StatusOK)
		var envelope struct {
			Data struct {
				Deleted bool `json:"deleted"`
			} `json:"data"`
		}
		decodeJSON(t, resp, &envelope)
		if !envelope.Data.Deleted {
			t.Fatal("delete: deleted = false, want true")
		}
	}
}

func TestLink_NotFound(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	resp := doRequest(t, http.MethodGet, "/admin/links/999999999", nil, token)
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusNotFound)
}

func TestRedirect_Active(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	originURL := "https://go.dev/redirect-active-" + uniqueName(t)
	id, code := createLink(t, token, originURL)
	t.Cleanup(func() { deleteLink(t, token, id) })

	resp := doRequestNoRedirect(t, http.MethodGet, "/"+code, nil, "")
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusFound)

	location := resp.Header.Get("Location")
	if location != originURL {
		t.Fatalf("Location = %q, want %q", location, originURL)
	}
}

func TestRedirect_NotFound(t *testing.T) {
	resp := doRequestNoRedirect(t, http.MethodGet, "/zz_no_such_code_xyz", nil, "")
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusNotFound)
}

func TestRedirect_Disabled(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	originURL := "https://example.com/disabled-" + uniqueName(t)
	id, code := createLink(t, token, originURL)
	t.Cleanup(func() { deleteLink(t, token, id) })

	// Disable the link via PATCH status=2.
	patchResp := doRequest(t, http.MethodPatch, fmt.Sprintf("/admin/links/%d", id), map[string]any{
		"status": 2,
	}, token)
	assertStatus(t, patchResp, http.StatusOK)
	patchResp.Body.Close()

	resp := doRequestNoRedirect(t, http.MethodGet, "/"+code, nil, "")
	defer resp.Body.Close()
	assertStatus(t, resp, http.StatusForbidden)
}

func TestRedirect_Expired(t *testing.T) {
	token := doLogin(t, testAdminUsername, testAdminPassword)
	// The API validates that expire_at must be in the future, so we cannot pass
	// a past timestamp at creation time. Instead, create without expiry and then
	// backdate directly in the DB — the Redis cache has no entry for this code
	// yet (never resolved), so the next FindOneByCode reads from DB.
	id, code := createLink(t, token, "https://example.com/expired-"+uniqueName(t))
	t.Cleanup(func() { deleteLink(t, token, id) })

	dbExec(t, "UPDATE short_link SET expire_at = '2020-01-01 00:00:00' WHERE id = ?", id)

	redirectResp := doRequestNoRedirect(t, http.MethodGet, "/"+code, nil, "")
	defer redirectResp.Body.Close()
	assertStatus(t, redirectResp, http.StatusGone)
}
