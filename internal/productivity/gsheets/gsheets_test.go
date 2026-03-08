package gsheets

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
	if cmd.Use != "gsheets" {
		t.Errorf("expected Use 'gsheets', got %q", cmd.Use)
	}
	if len(cmd.Commands()) < 7 {
		t.Errorf("expected at least 7 subcommands, got %d", len(cmd.Commands()))
	}
}

func TestColToLetter(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "A"},
		{1, "B"},
		{25, "Z"},
		{26, "AA"},
		{27, "AB"},
		{51, "AZ"},
		{52, "BA"},
		{701, "ZZ"},
		{702, "AAA"},
	}

	for _, tt := range tests {
		result := colToLetter(tt.input)
		if result != tt.expected {
			t.Errorf("colToLetter(%d) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestGetMetadataWithOAuth(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify OAuth bearer token
		if r.Header.Get("Authorization") != "Bearer test-oauth-token" {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"properties": map[string]any{
				"title":  "Test Spreadsheet",
				"locale": "en_US",
			},
			"sheets": []map[string]any{
				{
					"properties": map[string]any{
						"sheetId": 0,
						"title":   "Sheet1",
						"index":   0,
						"gridProperties": map[string]any{
							"rowCount":    1000,
							"columnCount": 26,
						},
					},
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoGet(srv.URL + "/spreadsheet123?fields=properties,sheets.properties")
	if err != nil {
		t.Fatalf("DoGet failed: %v", err)
	}

	var resp struct {
		Properties struct {
			Title string `json:"title"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Properties.Title != "Test Spreadsheet" {
		t.Errorf("expected title 'Test Spreadsheet', got %q", resp.Properties.Title)
	}
}

func TestReadRangeResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-oauth-token" {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"range": "Sheet1!A1:C3",
			"values": [][]interface{}{
				{"Name", "Age", "City"},
				{"Alice", 30, "NYC"},
				{"Bob", 25, "LA"},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoGet(srv.URL + "/spreadsheet123/values/Sheet1!A1:C3")
	if err != nil {
		t.Fatalf("DoGet failed: %v", err)
	}

	var resp struct {
		Range  string          `json:"range"`
		Values [][]interface{} `json:"values"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if resp.Range != "Sheet1!A1:C3" {
		t.Errorf("expected range 'Sheet1!A1:C3', got %q", resp.Range)
	}

	if len(resp.Values) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(resp.Values))
	}

	if resp.Values[0][0] != "Name" {
		t.Errorf("expected first cell 'Name', got %v", resp.Values[0][0])
	}
}

func TestBatchGetResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		if len(params["ranges"]) < 1 {
			t.Error("expected at least one 'ranges' query parameter")
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"valueRanges": []map[string]any{
				{
					"range": "Sheet1!A:ZZ",
					"values": [][]interface{}{
						{"Apple", "Banana"},
						{"Cat", "Dog"},
					},
				},
				{
					"range": "Sheet2!A:ZZ",
					"values": [][]interface{}{
						{"Test", "Value"},
					},
				},
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoGet(srv.URL + "/spreadsheet123/values:batchGet?ranges=Sheet1!A:ZZ&ranges=Sheet2!A:ZZ")
	if err != nil {
		t.Fatalf("DoGet failed: %v", err)
	}

	var resp struct {
		ValueRanges []struct {
			Range  string          `json:"range"`
			Values [][]interface{} `json:"values"`
		} `json:"valueRanges"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	if len(resp.ValueRanges) != 2 {
		t.Fatalf("expected 2 value ranges, got %d", len(resp.ValueRanges))
	}

	if resp.ValueRanges[0].Range != "Sheet1!A:ZZ" {
		t.Errorf("expected range 'Sheet1!A:ZZ', got %q", resp.ValueRanges[0].Range)
	}
}

func TestWriteValues(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.Header.Get("Authorization") != "Bearer test-oauth-token" {
			t.Errorf("expected Bearer token, got %q", r.Header.Get("Authorization"))
		}
		if r.URL.Query().Get("valueInputOption") != "USER_ENTERED" {
			t.Errorf("expected valueInputOption=USER_ENTERED, got %q", r.URL.Query().Get("valueInputOption"))
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		if body["range"] != "Sheet1!A1:B2" {
			t.Errorf("expected range 'Sheet1!A1:B2', got %v", body["range"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"updatedRange":   "Sheet1!A1:B2",
			"updatedRows":    2,
			"updatedColumns": 2,
			"updatedCells":   4,
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	bodyJSON := []byte(`{"range":"Sheet1!A1:B2","values":[["X","Y"],["Z","W"]]}`)
	data, err := client.DoRequest("PUT", srv.URL+"/spreadsheet123/values/Sheet1!A1:B2?valueInputOption=USER_ENTERED",
		bytes.NewReader(bodyJSON), "application/json")
	if err != nil {
		t.Fatalf("DoRequest failed: %v", err)
	}

	var resp struct {
		UpdatedCells int `json:"updatedCells"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.UpdatedCells != 4 {
		t.Errorf("expected 4 updated cells, got %d", resp.UpdatedCells)
	}
}

func TestAppendRows(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Query().Get("insertDataOption") != "INSERT_ROWS" {
			t.Errorf("expected insertDataOption=INSERT_ROWS, got %q", r.URL.Query().Get("insertDataOption"))
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"updates": map[string]any{
				"updatedRange":   "Sheet1!A4:B5",
				"updatedRows":    2,
				"updatedColumns": 2,
				"updatedCells":   4,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	bodyJSON := []byte(`{"range":"Sheet1","values":[["new1","new2"]]}`)
	data, err := client.DoPost(srv.URL+"/spreadsheet123/values/Sheet1:append?valueInputOption=USER_ENTERED&insertDataOption=INSERT_ROWS",
		bytes.NewReader(bodyJSON), "application/json")
	if err != nil {
		t.Fatalf("DoPost failed: %v", err)
	}

	var resp struct {
		Updates struct {
			UpdatedRows int `json:"updatedRows"`
		} `json:"updates"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.Updates.UpdatedRows != 2 {
		t.Errorf("expected 2 updated rows, got %d", resp.Updates.UpdatedRows)
	}
}

func TestClearRange(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"clearedRange": "Sheet1!A1:D10",
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	data, err := client.DoPost(srv.URL+"/spreadsheet123/values/Sheet1!A1:D10:clear",
		bytes.NewReader([]byte("{}")), "application/json")
	if err != nil {
		t.Fatalf("DoPost failed: %v", err)
	}

	var resp struct {
		ClearedRange string `json:"clearedRange"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.ClearedRange != "Sheet1!A1:D10" {
		t.Errorf("expected 'Sheet1!A1:D10', got %q", resp.ClearedRange)
	}
}

func TestCreateSpreadsheet(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		var body map[string]any
		json.NewDecoder(r.Body).Decode(&body)
		props := body["properties"].(map[string]any)
		if props["title"] != "New Sheet" {
			t.Errorf("expected title 'New Sheet', got %v", props["title"])
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"spreadsheetId":  "new-id-123",
			"spreadsheetUrl": "https://docs.google.com/spreadsheets/d/new-id-123",
			"properties": map[string]any{
				"title": "New Sheet",
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	bodyJSON := []byte(`{"properties":{"title":"New Sheet"}}`)
	data, err := client.DoPost(srv.URL, bytes.NewReader(bodyJSON), "application/json")
	if err != nil {
		t.Fatalf("DoPost failed: %v", err)
	}

	var resp struct {
		SpreadsheetID string `json:"spreadsheetId"`
		Properties    struct {
			Title string `json:"title"`
		} `json:"properties"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}
	if resp.SpreadsheetID != "new-id-123" {
		t.Errorf("expected id 'new-id-123', got %q", resp.SpreadsheetID)
	}
	if resp.Properties.Title != "New Sheet" {
		t.Errorf("expected title 'New Sheet', got %q", resp.Properties.Title)
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

	_, err := client.DoGet(srv.URL + "/spreadsheet123")
	if err == nil {
		t.Error("expected error for 403 response, got nil")
	}
}

func TestNotFoundError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"error": map[string]any{
				"message": "Spreadsheet not found",
				"code":    404,
			},
		})
	}))
	defer srv.Close()

	client := newTestClient(srv)

	_, err := client.DoGet(srv.URL + "/spreadsheet123")
	if err == nil {
		t.Error("expected error for 404 response, got nil")
	}
}

func TestSearchLogic(t *testing.T) {
	valueRanges := []struct {
		Range  string
		Values [][]interface{}
	}{
		{
			Range: "'Sheet1'!A1:Z100",
			Values: [][]interface{}{
				{"This is a test", "other"},
				{"value", "another test"},
			},
		},
	}

	queryLower := "test"
	var matches []CellMatch

	for _, vr := range valueRanges {
		sheetName := vr.Range
		if idx := len(sheetName); idx > 0 {
			sheetName = "Sheet1"
		}

		for rowIdx, row := range vr.Values {
			for colIdx, cell := range row {
				cellStr := ""
				if cell != nil {
					cellStr = cell.(string)
				}
				if len(cellStr) > 0 && len(queryLower) > 0 {
					contains := false
					for i := 0; i <= len(cellStr)-len(queryLower); i++ {
						if cellStr[i:i+len(queryLower)] == queryLower {
							contains = true
							break
						}
					}
					if contains {
						cellRef := colToLetter(colIdx) + string(rune('1'+rowIdx))
						matches = append(matches, CellMatch{
							Sheet: sheetName,
							Cell:  cellRef,
							Value: cellStr,
							Row:   rowIdx + 1,
							Col:   colIdx + 1,
						})
					}
				}
			}
		}
	}

	if len(matches) != 2 {
		t.Fatalf("expected 2 matches, got %d", len(matches))
	}

	if matches[0].Cell != "A1" {
		t.Errorf("expected cell 'A1', got %q", matches[0].Cell)
	}
	if matches[1].Cell != "B2" {
		t.Errorf("expected cell 'B2', got %q", matches[1].Cell)
	}
}
