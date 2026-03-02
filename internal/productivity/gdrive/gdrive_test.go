package gdrive

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	googleauth "github.com/unstablemind/pocket/internal/common/google"
)

func TestNewCmd(t *testing.T) {
	cmd := NewCmd()
	if cmd.Use != "gdrive" {
		t.Errorf("expected Use 'gdrive', got %q", cmd.Use)
	}
	if len(cmd.Commands()) < 9 {
		t.Errorf("expected at least 9 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestFormatTime(t *testing.T) {
	now := time.Now()

	tests := []struct {
		input    string
		expected string
	}{
		{now.Add(-30 * time.Second).Format(time.RFC3339), "now"},
		{now.Add(-5 * time.Minute).Format(time.RFC3339), "5m ago"},
		{now.Add(-2 * time.Hour).Format(time.RFC3339), "2h ago"},
		{now.Add(-3 * 24 * time.Hour).Format(time.RFC3339), "3d ago"},
		{now.Add(-10 * 24 * time.Hour).Format(time.RFC3339), "1w ago"},
		{now.Add(-40 * 24 * time.Hour).Format(time.RFC3339), "1mo ago"},
		{now.Add(-400 * 24 * time.Hour).Format(time.RFC3339), "1y ago"},
		{"", ""},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		result := formatTime(tt.input)
		if result != tt.expected {
			t.Errorf("formatTime(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestFormatSize(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"", ""},
		{"500", "500 B"},
		{"1024", "1.0 KB"},
		{"1048576", "1.0 MB"},
		{"1073741824", "1.0 GB"},
		{"2147483648", "2.0 GB"},
		{"invalid", "invalid"},
	}

	for _, tt := range tests {
		result := formatSize(tt.input)
		if result != tt.expected {
			t.Errorf("formatSize(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestExportMimeForDocs(t *testing.T) {
	tests := []struct {
		format   string
		expected string
	}{
		{"text", "text/plain"},
		{"md", "text/markdown"},
		{"markdown", "text/markdown"},
		{"html", "text/html"},
		{"", "text/plain"},
	}

	for _, tt := range tests {
		result := exportMimeForDocs(tt.format)
		if result != tt.expected {
			t.Errorf("exportMimeForDocs(%q) = %q, expected %q", tt.format, result, tt.expected)
		}
	}
}

func TestExportForDownload(t *testing.T) {
	tests := []struct {
		mimeType    string
		expectedExt string
	}{
		{"application/vnd.google-apps.document", ".docx"},
		{"application/vnd.google-apps.spreadsheet", ".xlsx"},
		{"application/vnd.google-apps.presentation", ".pptx"},
		{"application/vnd.google-apps.drawing", ".png"},
		{"application/vnd.google-apps.form", ".pdf"},
	}

	for _, tt := range tests {
		_, ext := exportForDownload(tt.mimeType)
		if ext != tt.expectedExt {
			t.Errorf("exportForDownload(%q) ext = %q, expected %q", tt.mimeType, ext, tt.expectedExt)
		}
	}
}

func TestIsTextMime(t *testing.T) {
	tests := []struct {
		mime     string
		expected bool
	}{
		{"text/plain", true},
		{"text/csv", true},
		{"text/html", true},
		{"application/json", true},
		{"application/xml", true},
		{"application/pdf", false},
		{"image/png", false},
		{"application/octet-stream", false},
	}

	for _, tt := range tests {
		result := isTextMime(tt.mime)
		if result != tt.expected {
			t.Errorf("isTextMime(%q) = %v, expected %v", tt.mime, result, tt.expected)
		}
	}
}

func newTestClient(srv *httptest.Server) *googleauth.Client {
	return &googleauth.Client{
		AccessToken: "test-token",
		HTTPClient:  srv.Client(),
	}
}

func TestSearchWithOAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if auth := r.Header.Get("Authorization"); auth != "Bearer test-token" {
			t.Errorf("expected Authorization 'Bearer test-token', got %q", auth)
		}
		if ua := r.Header.Get("User-Agent"); ua != "Alpha-CLI/1.0" {
			t.Errorf("expected User-Agent 'Alpha-CLI/1.0', got %q", ua)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"files": []map[string]any{
				{
					"id":           "file1",
					"name":         "test.txt",
					"mimeType":     "text/plain",
					"size":         "1024",
					"createdTime":  "2024-01-01T12:00:00Z",
					"modifiedTime": "2024-01-01T13:00:00Z",
					"webViewLink":  "https://drive.google.com/file/d/file1",
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoGet(srv.URL + "/drive/v3/files?q=test")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var resp struct {
		Files []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"files"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(resp.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(resp.Files))
	}
	if resp.Files[0].ID != "file1" {
		t.Errorf("expected file ID 'file1', got %q", resp.Files[0].ID)
	}
}

func TestListFiles(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		q := params.Get("q")
		if q == "" {
			t.Error("expected query parameter 'q' with parent filter")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"files": []map[string]any{
				{"id": "folder1", "name": "Documents", "mimeType": "application/vnd.google-apps.folder"},
				{"id": "file1", "name": "notes.txt", "mimeType": "text/plain", "size": "256"},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)
	data, err := client.DoGet(srv.URL + "/drive/v3/files?q='root'+in+parents")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var resp struct {
		Files []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"files"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if len(resp.Files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(resp.Files))
	}
}

func TestReadGoogleDoc(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/export") {
			w.Write([]byte("Hello, this is the document content."))
			return
		}
		// Metadata
		json.NewEncoder(w).Encode(map[string]any{
			"id":       "doc1",
			"name":     "My Document",
			"mimeType": "application/vnd.google-apps.document",
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	// Get metadata
	metaData, err := client.DoGet(srv.URL + "/drive/v3/files/doc1")
	if err != nil {
		t.Fatalf("metadata request failed: %v", err)
	}

	var meta struct {
		MimeType string `json:"mimeType"`
	}
	if err := json.Unmarshal(metaData, &meta); err != nil {
		t.Fatalf("failed to unmarshal meta: %v", err)
	}

	if meta.MimeType != "application/vnd.google-apps.document" {
		t.Fatalf("expected google docs mime type, got %q", meta.MimeType)
	}

	// Export
	content, err := client.DoGet(srv.URL + "/drive/v3/files/doc1/export?mimeType=text/plain")
	if err != nil {
		t.Fatalf("export request failed: %v", err)
	}

	if string(content) != "Hello, this is the document content." {
		t.Errorf("unexpected content: %q", string(content))
	}
}

func TestReadRegularFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if params.Get("alt") == "media" {
			w.Write([]byte("file content here"))
			return
		}
		// Metadata
		json.NewEncoder(w).Encode(map[string]any{
			"id":       "file1",
			"name":     "readme.md",
			"mimeType": "text/markdown",
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	content, err := client.DoGet(srv.URL + "/drive/v3/files/file1?alt=media")
	if err != nil {
		t.Fatalf("download request failed: %v", err)
	}

	if string(content) != "file content here" {
		t.Errorf("unexpected content: %q", string(content))
	}
}

func TestMkdir(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var meta map[string]any
		if err := json.Unmarshal(body, &meta); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if meta["name"] != "New Folder" {
			t.Errorf("expected name 'New Folder', got %v", meta["name"])
		}
		if meta["mimeType"] != "application/vnd.google-apps.folder" {
			t.Errorf("expected folder mimeType, got %v", meta["mimeType"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"id":          "folder123",
			"name":        "New Folder",
			"webViewLink": "https://drive.google.com/drive/folders/folder123",
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	metaJSON := `{"name":"New Folder","mimeType":"application/vnd.google-apps.folder"}`
	data, err := client.DoPost(srv.URL+"/drive/v3/files",
		strings.NewReader(metaJSON), "application/json")
	if err != nil {
		t.Fatalf("mkdir request failed: %v", err)
	}

	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.ID != "folder123" {
		t.Errorf("expected ID 'folder123', got %q", resp.ID)
	}
	if resp.Name != "New Folder" {
		t.Errorf("expected name 'New Folder', got %q", resp.Name)
	}
}

func TestDeleteTrash(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var patch map[string]any
		if err := json.Unmarshal(body, &patch); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if patch["trashed"] != true {
			t.Errorf("expected trashed=true, got %v", patch["trashed"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id":      "file1",
			"name":    "deleteme.txt",
			"trashed": true,
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoPatch(srv.URL+"/drive/v3/files/file1",
		strings.NewReader(`{"trashed":true}`), "application/json")
	if err != nil {
		t.Fatalf("delete request failed: %v", err)
	}

	var resp struct {
		ID      string `json:"id"`
		Trashed bool   `json:"trashed"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if !resp.Trashed {
		t.Error("expected trashed=true")
	}
}

func TestForbiddenError(t *testing.T) {
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

	client := newTestClient(srv)

	_, err := client.DoGet(srv.URL + "/drive/v3/files")
	if err == nil {
		t.Error("expected error for 403 response, got nil")
	}
}

func TestNotFoundError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "File not found",
				"code":    404,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	_, err := client.DoGet(srv.URL + "/drive/v3/files/nonexistent")
	if err == nil {
		t.Error("expected error for 404 response, got nil")
	}
}

func TestFileInfoResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"id":           "file123",
			"name":         "Document.docx",
			"mimeType":     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			"size":         "2097152",
			"createdTime":  "2024-01-01T10:00:00Z",
			"modifiedTime": "2024-01-05T15:30:00Z",
			"webViewLink":  "https://drive.google.com/file/d/file123",
			"description":  "Important document",
			"owners": []map[string]any{
				{"displayName": "John Doe"},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoGet(srv.URL + "/drive/v3/files/file123?fields=id,name,mimeType,size,createdTime,modifiedTime,webViewLink,description,owners")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	var resp struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Size        string `json:"size"`
		Description string `json:"description"`
		Owners      []struct {
			DisplayName string `json:"displayName"`
		} `json:"owners"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.ID != "file123" {
		t.Errorf("expected file ID 'file123', got %q", resp.ID)
	}
	if resp.Name != "Document.docx" {
		t.Errorf("expected file name 'Document.docx', got %q", resp.Name)
	}
	if resp.Description != "Important document" {
		t.Errorf("expected description 'Important document', got %q", resp.Description)
	}
	if len(resp.Owners) != 1 || resp.Owners[0].DisplayName != "John Doe" {
		t.Error("expected owner 'John Doe'")
	}
}

func TestUpload(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		ct := r.Header.Get("Content-Type")
		if !strings.Contains(ct, "multipart/") {
			t.Errorf("expected multipart content type, got %q", ct)
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id":          "uploaded1",
			"name":        "test.txt",
			"mimeType":    "text/plain",
			"size":        "100",
			"webViewLink": "https://drive.google.com/file/d/uploaded1",
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	// Simulate multipart upload
	body := strings.NewReader("--boundary\r\nContent-Type: application/json\r\n\r\n{\"name\":\"test.txt\"}\r\n--boundary\r\nContent-Type: text/plain\r\n\r\nfile content\r\n--boundary--")
	data, err := client.DoPost(srv.URL+"/upload/drive/v3/files?uploadType=multipart",
		body, "multipart/related; boundary=boundary")
	if err != nil {
		t.Fatalf("upload request failed: %v", err)
	}

	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.ID != "uploaded1" {
		t.Errorf("expected ID 'uploaded1', got %q", resp.ID)
	}
}

func TestUpdateMetadata(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("expected PATCH, got %s", r.Method)
		}

		body, _ := io.ReadAll(r.Body)
		var meta map[string]any
		if err := json.Unmarshal(body, &meta); err != nil {
			t.Fatalf("failed to unmarshal body: %v", err)
		}

		if meta["name"] != "renamed.txt" {
			t.Errorf("expected name 'renamed.txt', got %v", meta["name"])
		}

		json.NewEncoder(w).Encode(map[string]any{
			"id":           "file1",
			"name":         "renamed.txt",
			"modifiedTime": "2024-06-01T12:00:00Z",
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoPatch(srv.URL+"/drive/v3/files/file1",
		strings.NewReader(`{"name":"renamed.txt"}`), "application/json")
	if err != nil {
		t.Fatalf("update request failed: %v", err)
	}

	var resp struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Name != "renamed.txt" {
		t.Errorf("expected name 'renamed.txt', got %q", resp.Name)
	}
}
