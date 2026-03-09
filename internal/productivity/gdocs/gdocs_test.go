package gdocs

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	googleauth "github.com/sterlingcodes/alpha-cli/internal/common/google"
)

func newTestClient(srv *httptest.Server) *googleauth.Client {
	return &googleauth.Client{
		AccessToken: "test-oauth-token",
		HTTPClient:  srv.Client(),
	}
}

func TestNewCmd(t *testing.T) {
	cmd := NewCmd()
	if cmd.Use != "gdocs" {
		t.Errorf("expected Use 'gdocs', got %q", cmd.Use)
	}
	if len(cmd.Commands()) != 6 {
		t.Errorf("expected 6 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestGetMetadata(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-oauth-token" {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"documentId": "doc-123",
			"title":      "Test Document",
			"revisionId": "rev-456",
		})
	}))
	defer srv.Close()

	// Override baseURL for testing
	origURL := baseURL
	baseURL = srv.URL
	defer func() { baseURL = origURL }()

	client := newTestClient(srv)

	data, err := client.DoGet(srv.URL + "/doc-123?fields=documentId,title,revisionId")
	if err != nil {
		t.Fatalf("DoGet failed: %v", err)
	}

	var resp struct {
		DocumentID string `json:"documentId"`
		Title      string `json:"title"`
		RevisionID string `json:"revisionId"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.DocumentID != "doc-123" {
		t.Errorf("expected documentId 'doc-123', got %q", resp.DocumentID)
	}
	if resp.Title != "Test Document" {
		t.Errorf("expected title 'Test Document', got %q", resp.Title)
	}
	if resp.RevisionID != "rev-456" {
		t.Errorf("expected revisionId 'rev-456', got %q", resp.RevisionID)
	}
}

func TestReadDocumentText(t *testing.T) {
	docJSON := map[string]any{
		"body": map[string]any{
			"content": []map[string]any{
				{
					"paragraph": map[string]any{
						"elements": []map[string]any{
							{
								"textRun": map[string]any{
									"content": "Hello ",
								},
							},
							{
								"textRun": map[string]any{
									"content": "World\n",
								},
							},
						},
					},
				},
				{
					"paragraph": map[string]any{
						"elements": []map[string]any{
							{
								"textRun": map[string]any{
									"content": "Second paragraph\n",
								},
							},
						},
					},
				},
			},
		},
	}

	data, err := json.Marshal(docJSON)
	if err != nil {
		t.Fatalf("failed to marshal test doc: %v", err)
	}

	text, err := extractPlainText(data)
	if err != nil {
		t.Fatalf("extractPlainText failed: %v", err)
	}

	expected := "Hello World\nSecond paragraph\n"
	if text != expected {
		t.Errorf("expected text %q, got %q", expected, text)
	}
}

func TestCreateDocument(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["title"] != "New Doc" {
			t.Errorf("expected title 'New Doc', got %v", body["title"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"documentId": "new-doc-123",
			"title":      "New Doc",
			"revisionId": "rev-1",
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	bodyJSON := []byte(`{"title":"New Doc"}`)
	data, err := client.DoPost(srv.URL, bytes.NewReader(bodyJSON), "application/json")
	if err != nil {
		t.Fatalf("DoPost failed: %v", err)
	}

	var resp struct {
		DocumentID string `json:"documentId"`
		Title      string `json:"title"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.DocumentID != "new-doc-123" {
		t.Errorf("expected id 'new-doc-123', got %q", resp.DocumentID)
	}
	if resp.Title != "New Doc" {
		t.Errorf("expected title 'New Doc', got %q", resp.Title)
	}
}

func TestBatchUpdateReplaceAllText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-oauth-token" {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)

		requests, ok := body["requests"].([]any)
		if !ok || len(requests) != 1 {
			t.Fatalf("expected 1 request, got %v", body["requests"])
		}

		req := requests[0].(map[string]any)
		replaceAll, ok := req["replaceAllText"].(map[string]any)
		if !ok {
			t.Fatal("expected replaceAllText request")
		}

		containsText := replaceAll["containsText"].(map[string]any)
		if containsText["text"] != "foo" {
			t.Errorf("expected find text 'foo', got %v", containsText["text"])
		}
		if replaceAll["replaceText"] != "bar" {
			t.Errorf("expected replace text 'bar', got %v", replaceAll["replaceText"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"documentId": "doc-123",
			"replies": []map[string]any{
				{
					"replaceAllText": map[string]any{
						"occurrencesChanged": 3,
					},
				},
			},
		})
	}))
	defer srv.Close()

	origURL := baseURL
	baseURL = srv.URL
	defer func() { baseURL = origURL }()

	client := newTestClient(srv)

	err := runReplace(client, "doc-123", "foo", "bar")
	if err != nil {
		t.Fatalf("runReplace failed: %v", err)
	}
}

func TestReadDocumentEmptyBody(t *testing.T) {
	docJSON := map[string]any{
		"body": map[string]any{
			"content": []map[string]any{},
		},
	}

	data, err := json.Marshal(docJSON)
	if err != nil {
		t.Fatalf("failed to marshal test doc: %v", err)
	}

	text, err := extractPlainText(data)
	if err != nil {
		t.Fatalf("extractPlainText failed: %v", err)
	}

	if text != "" {
		t.Errorf("expected empty text, got %q", text)
	}
}

func TestForbiddenError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Access denied",
				"code":    403,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	_, err := client.DoGet(srv.URL + "/doc-123")
	if err == nil {
		t.Error("expected error for 403 response, got nil")
	}
}
