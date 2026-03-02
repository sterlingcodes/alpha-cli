package google

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExchangeRefreshToken(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		if ct := r.Header.Get("Content-Type"); ct != "application/x-www-form-urlencoded" {
			t.Errorf("expected Content-Type application/x-www-form-urlencoded, got %q", ct)
		}

		if err := r.ParseForm(); err != nil {
			t.Fatalf("failed to parse form: %v", err)
		}

		if r.FormValue("client_id") != "test-client-id" {
			t.Errorf("expected client_id 'test-client-id', got %q", r.FormValue("client_id"))
		}
		if r.FormValue("client_secret") != "test-client-secret" {
			t.Errorf("expected client_secret 'test-client-secret', got %q", r.FormValue("client_secret"))
		}
		if r.FormValue("refresh_token") != "test-refresh-token" {
			t.Errorf("expected refresh_token 'test-refresh-token', got %q", r.FormValue("refresh_token"))
		}
		if r.FormValue("grant_type") != "refresh_token" {
			t.Errorf("expected grant_type 'refresh_token', got %q", r.FormValue("grant_type"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"access_token": "mock-access-token",
			"expires_in":   3600,
		})
	}))
	defer srv.Close()

	oldURL := tokenURL
	defer func() { tokenURL = oldURL }()
	tokenURL = srv.URL

	token, err := ExchangeRefreshToken("test-client-id", "test-client-secret", "test-refresh-token")
	if err != nil {
		t.Fatalf("ExchangeRefreshToken failed: %v", err)
	}

	if token != "mock-access-token" {
		t.Errorf("expected 'mock-access-token', got %q", token)
	}
}

func TestExchangeRefreshTokenInvalidGrant(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"error":             "invalid_grant",
			"error_description": "Token has been expired or revoked.",
		})
	}))
	defer srv.Close()

	oldURL := tokenURL
	defer func() { tokenURL = oldURL }()
	tokenURL = srv.URL

	_, err := ExchangeRefreshToken("id", "secret", "bad-token")
	if err == nil {
		t.Fatal("expected error for invalid grant, got nil")
	}
	if !strings.Contains(err.Error(), "Token has been expired or revoked") {
		t.Errorf("expected error to mention expired/revoked token, got: %v", err)
	}
}

func TestExchangeRefreshTokenHTTPError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("server error"))
	}))
	defer srv.Close()

	oldURL := tokenURL
	defer func() { tokenURL = oldURL }()
	tokenURL = srv.URL

	_, err := ExchangeRefreshToken("id", "secret", "token")
	if err == nil {
		t.Fatal("expected error for 500 response, got nil")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("expected error to mention HTTP 500, got: %v", err)
	}
}

func TestDoRequest(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("expected Authorization 'Bearer test-token', got %q", auth)
		}
		if ua := r.Header.Get("User-Agent"); ua != "Alpha-CLI/1.0" {
			t.Errorf("expected User-Agent 'Alpha-CLI/1.0', got %q", ua)
		}
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"result": "ok"})
	}))
	defer srv.Close()

	client := &Client{
		AccessToken: "test-token",
		HTTPClient:  srv.Client(),
	}

	data, err := client.DoGet(srv.URL + "/test")
	if err != nil {
		t.Fatalf("DoGet failed: %v", err)
	}

	var resp map[string]string
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp["result"] != "ok" {
		t.Errorf("expected result 'ok', got %q", resp["result"])
	}
}

func TestDoRequestAPIError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Insufficient Permission",
				"code":    403,
			},
		})
	}))
	defer srv.Close()

	client := &Client{
		AccessToken: "test-token",
		HTTPClient:  srv.Client(),
	}

	_, err := client.DoGet(srv.URL + "/test")
	if err == nil {
		t.Fatal("expected error for 403 response, got nil")
	}
	if !strings.Contains(err.Error(), "Insufficient Permission") {
		t.Errorf("expected error to mention Insufficient Permission, got: %v", err)
	}
}

func TestDoRequestPost(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type 'application/json', got %q", ct)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"created": true})
	}))
	defer srv.Close()

	client := &Client{
		AccessToken: "test-token",
		HTTPClient:  srv.Client(),
	}

	body := strings.NewReader(`{"name":"test"}`)
	data, err := client.DoPost(srv.URL+"/create", body, "application/json")
	if err != nil {
		t.Fatalf("DoPost failed: %v", err)
	}

	var resp map[string]bool
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if !resp["created"] {
		t.Error("expected created=true")
	}
}

func TestDoRequestPatch(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{"updated": true})
	}))
	defer srv.Close()

	client := &Client{
		AccessToken: "test-token",
		HTTPClient:  srv.Client(),
	}

	body := strings.NewReader(`{"name":"updated"}`)
	data, err := client.DoPatch(srv.URL+"/update", body, "application/json")
	if err != nil {
		t.Fatalf("DoPatch failed: %v", err)
	}

	var resp map[string]bool
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if !resp["updated"] {
		t.Error("expected updated=true")
	}
}
