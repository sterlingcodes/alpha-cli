package gsheets

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/spf13/cobra"

	googleauth "github.com/unstablemind/pocket/internal/common/google"
	"github.com/unstablemind/pocket/pkg/output"
)

var baseURL = "https://sheets.googleapis.com/v4/spreadsheets"

// SpreadsheetInfo holds metadata about a spreadsheet
type SpreadsheetInfo struct {
	ID     string      `json:"id"`
	Title  string      `json:"title"`
	Locale string      `json:"locale"`
	Sheets []SheetInfo `json:"sheets"`
}

// SheetInfo holds metadata about a single sheet
type SheetInfo struct {
	ID       int    `json:"id"`
	Title    string `json:"title"`
	Index    int    `json:"index"`
	RowCount int    `json:"row_count"`
	ColCount int    `json:"col_count"`
}

// SheetData holds cell values from a range
type SheetData struct {
	SpreadsheetID string     `json:"spreadsheet_id"`
	Range         string     `json:"range"`
	Rows          [][]string `json:"rows"`
	RowCount      int        `json:"row_count"`
	ColCount      int        `json:"col_count"`
}

// SearchResult holds search matches
type SearchResult struct {
	Query   string      `json:"query"`
	Matches []CellMatch `json:"matches"`
	Count   int         `json:"count"`
}

// CellMatch holds a single cell match
type CellMatch struct {
	Sheet string `json:"sheet"`
	Cell  string `json:"cell"`
	Value string `json:"value"`
	Row   int    `json:"row"`
	Col   int    `json:"col"`
}

// NewCmd returns the Google Sheets command
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gsheets",
		Aliases: []string{"sheets", "spreadsheet"},
		Short:   "Google Sheets commands",
	}

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newReadCmd())
	cmd.AddCommand(newSearchCmd())
	cmd.AddCommand(newWriteCmd())
	cmd.AddCommand(newAppendCmd())
	cmd.AddCommand(newClearCmd())
	cmd.AddCommand(newCreateCmd())

	return cmd
}

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [spreadsheet-id]",
		Short: "Get spreadsheet metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			spreadsheetID := args[0]
			reqURL := fmt.Sprintf("%s/%s?fields=properties,sheets.properties",
				baseURL, url.PathEscape(spreadsheetID))

			data, err := client.DoGet(reqURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var resp struct {
				Properties struct {
					Title  string `json:"title"`
					Locale string `json:"locale"`
				} `json:"properties"`
				Sheets []struct {
					Properties struct {
						SheetID   int    `json:"sheetId"`
						Title     string `json:"title"`
						Index     int    `json:"index"`
						GridProps struct {
							RowCount    int `json:"rowCount"`
							ColumnCount int `json:"columnCount"`
						} `json:"gridProperties"`
					} `json:"properties"`
				} `json:"sheets"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse response: %s", err.Error()), nil)
			}

			sheets := make([]SheetInfo, 0, len(resp.Sheets))
			for _, s := range resp.Sheets {
				sheets = append(sheets, SheetInfo{
					ID:       s.Properties.SheetID,
					Title:    s.Properties.Title,
					Index:    s.Properties.Index,
					RowCount: s.Properties.GridProps.RowCount,
					ColCount: s.Properties.GridProps.ColumnCount,
				})
			}

			result := SpreadsheetInfo{
				ID:     spreadsheetID,
				Title:  resp.Properties.Title,
				Locale: resp.Properties.Locale,
				Sheets: sheets,
			}

			return output.Print(result)
		},
	}

	return cmd
}

func newReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read [spreadsheet-id] [range]",
		Short: "Read cell values from a range",
		Long:  `Read cell values. Range format: "Sheet1!A1:D10" or "A1:D10"`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			spreadsheetID := args[0]
			rangeStr := args[1]

			reqURL := fmt.Sprintf("%s/%s/values/%s",
				baseURL, url.PathEscape(spreadsheetID),
				url.PathEscape(rangeStr))

			data, err := client.DoGet(reqURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var resp struct {
				Range  string          `json:"range"`
				Values [][]interface{} `json:"values"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse response: %s", err.Error()), nil)
			}

			// Convert to string slices
			rows := make([][]string, 0, len(resp.Values))
			maxCols := 0
			for _, row := range resp.Values {
				strRow := make([]string, 0, len(row))
				for _, cell := range row {
					strRow = append(strRow, fmt.Sprintf("%v", cell))
				}
				rows = append(rows, strRow)
				if len(strRow) > maxCols {
					maxCols = len(strRow)
				}
			}

			result := SheetData{
				SpreadsheetID: spreadsheetID,
				Range:         resp.Range,
				Rows:          rows,
				RowCount:      len(rows),
				ColCount:      maxCols,
			}

			return output.Print(result)
		},
	}

	return cmd
}

func newSearchCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search [spreadsheet-id] [query]",
		Short: "Search for a value across all sheets",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			spreadsheetID := args[0]
			query := args[1]

			// First, get sheet list
			metaURL := fmt.Sprintf("%s/%s?fields=sheets.properties.title",
				baseURL, url.PathEscape(spreadsheetID))

			metaData, err := client.DoGet(metaURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var metaResp struct {
				Sheets []struct {
					Properties struct {
						Title string `json:"title"`
					} `json:"properties"`
				} `json:"sheets"`
			}

			if err := json.Unmarshal(metaData, &metaResp); err != nil {
				return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse metadata: %s", err.Error()), nil)
			}

			if len(metaResp.Sheets) == 0 {
				return output.Print(SearchResult{
					Query:   query,
					Matches: []CellMatch{},
					Count:   0,
				})
			}

			// Build batch read URL with ranges for all sheets
			params := url.Values{}
			for _, sheet := range metaResp.Sheets {
				params.Add("ranges", sheet.Properties.Title+"!A:ZZ")
			}

			batchURL := fmt.Sprintf("%s/%s/values:batchGet?%s",
				baseURL, url.PathEscape(spreadsheetID), params.Encode())

			batchData, err := client.DoGet(batchURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var batchResp struct {
				ValueRanges []struct {
					Range  string          `json:"range"`
					Values [][]interface{} `json:"values"`
				} `json:"valueRanges"`
			}

			if err := json.Unmarshal(batchData, &batchResp); err != nil {
				return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse batch data: %s", err.Error()), nil)
			}

			queryLower := strings.ToLower(query)
			var matches []CellMatch

			for _, vr := range batchResp.ValueRanges {
				// Extract sheet name from range like "Sheet1!A1:ZZ1000"
				sheetName := vr.Range
				if idx := strings.Index(sheetName, "!"); idx >= 0 {
					sheetName = sheetName[:idx]
				}
				// Remove surrounding single quotes if present
				sheetName = strings.Trim(sheetName, "'")

				for rowIdx, row := range vr.Values {
					for colIdx, cell := range row {
						cellStr := fmt.Sprintf("%v", cell)
						if strings.Contains(strings.ToLower(cellStr), queryLower) {
							cellRef := fmt.Sprintf("%s%d", colToLetter(colIdx), rowIdx+1)
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

			if matches == nil {
				matches = []CellMatch{}
			}

			result := SearchResult{
				Query:   query,
				Matches: matches,
				Count:   len(matches),
			}

			return output.Print(result)
		},
	}

	return cmd
}

func newWriteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "write [spreadsheet-id] [range] [values-json]",
		Short: "Write values to a range",
		Long:  `Write cell values. Values as JSON array of arrays: '[["A","B"],["C","D"]]'. Range format: "Sheet1!A1:B2"`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			spreadsheetID := args[0]
			rangeStr := args[1]
			valuesJSON := args[2]

			var values [][]interface{}
			if err := json.Unmarshal([]byte(valuesJSON), &values); err != nil {
				return output.PrintError("invalid_values", fmt.Sprintf("Values must be a JSON array of arrays: %s", err.Error()), nil)
			}

			body := map[string]any{
				"range":  rangeStr,
				"values": values,
			}
			bodyJSON, err := json.Marshal(body)
			if err != nil {
				return output.PrintError("marshal_failed", err.Error(), nil)
			}

			reqURL := fmt.Sprintf("%s/%s/values/%s?valueInputOption=USER_ENTERED",
				baseURL, url.PathEscape(spreadsheetID), url.PathEscape(rangeStr))

			data, err := client.DoRequest("PUT", reqURL, bytes.NewReader(bodyJSON), "application/json")
			if err != nil {
				return output.PrintError("write_failed", err.Error(), nil)
			}

			var resp struct {
				UpdatedRange string `json:"updatedRange"`
				UpdatedRows  int    `json:"updatedRows"`
				UpdatedCols  int    `json:"updatedColumns"`
				UpdatedCells int    `json:"updatedCells"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			return output.Print(map[string]any{
				"spreadsheet_id": spreadsheetID,
				"updated_range":  resp.UpdatedRange,
				"updated_rows":   resp.UpdatedRows,
				"updated_cols":   resp.UpdatedCols,
				"updated_cells":  resp.UpdatedCells,
			})
		},
	}

	return cmd
}

func newAppendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append [spreadsheet-id] [range] [values-json]",
		Short: "Append rows to a sheet",
		Long:  `Append rows after the last row with data. Values as JSON array of arrays: '[["A","B"],["C","D"]]'. Range format: "Sheet1"`,
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			spreadsheetID := args[0]
			rangeStr := args[1]
			valuesJSON := args[2]

			var values [][]interface{}
			if err := json.Unmarshal([]byte(valuesJSON), &values); err != nil {
				return output.PrintError("invalid_values", fmt.Sprintf("Values must be a JSON array of arrays: %s", err.Error()), nil)
			}

			body := map[string]any{
				"range":  rangeStr,
				"values": values,
			}
			bodyJSON, err := json.Marshal(body)
			if err != nil {
				return output.PrintError("marshal_failed", err.Error(), nil)
			}

			reqURL := fmt.Sprintf("%s/%s/values/%s:append?valueInputOption=USER_ENTERED&insertDataOption=INSERT_ROWS",
				baseURL, url.PathEscape(spreadsheetID), url.PathEscape(rangeStr))

			data, err := client.DoPost(reqURL, bytes.NewReader(bodyJSON), "application/json")
			if err != nil {
				return output.PrintError("append_failed", err.Error(), nil)
			}

			var resp struct {
				Updates struct {
					UpdatedRange string `json:"updatedRange"`
					UpdatedRows  int    `json:"updatedRows"`
					UpdatedCols  int    `json:"updatedColumns"`
					UpdatedCells int    `json:"updatedCells"`
				} `json:"updates"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			return output.Print(map[string]any{
				"spreadsheet_id": spreadsheetID,
				"updated_range":  resp.Updates.UpdatedRange,
				"updated_rows":   resp.Updates.UpdatedRows,
				"updated_cols":   resp.Updates.UpdatedCols,
				"updated_cells":  resp.Updates.UpdatedCells,
			})
		},
	}

	return cmd
}

func newClearCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "clear [spreadsheet-id] [range]",
		Short: "Clear values from a range",
		Long:  `Clear cell values in a range. Range format: "Sheet1!A1:D10"`,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			spreadsheetID := args[0]
			rangeStr := args[1]

			reqURL := fmt.Sprintf("%s/%s/values/%s:clear",
				baseURL, url.PathEscape(spreadsheetID), url.PathEscape(rangeStr))

			data, err := client.DoPost(reqURL, bytes.NewReader([]byte("{}")), "application/json")
			if err != nil {
				return output.PrintError("clear_failed", err.Error(), nil)
			}

			var resp struct {
				ClearedRange string `json:"clearedRange"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			return output.Print(map[string]any{
				"spreadsheet_id": spreadsheetID,
				"cleared_range":  resp.ClearedRange,
			})
		},
	}

	return cmd
}

func newCreateCmd() *cobra.Command {
	var title string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new spreadsheet",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			body := map[string]any{
				"properties": map[string]string{
					"title": title,
				},
			}
			bodyJSON, err := json.Marshal(body)
			if err != nil {
				return output.PrintError("marshal_failed", err.Error(), nil)
			}

			data, err := client.DoPost(baseURL, bytes.NewReader(bodyJSON), "application/json")
			if err != nil {
				return output.PrintError("create_failed", err.Error(), nil)
			}

			var resp struct {
				SpreadsheetID  string `json:"spreadsheetId"`
				SpreadsheetURL string `json:"spreadsheetUrl"`
				Properties     struct {
					Title string `json:"title"`
				} `json:"properties"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			return output.Print(map[string]any{
				"id":    resp.SpreadsheetID,
				"title": resp.Properties.Title,
				"url":   resp.SpreadsheetURL,
			})
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Spreadsheet title (required)")
	_ = cmd.MarkFlagRequired("title")

	return cmd
}

func colToLetter(colIndex int) string {
	result := ""
	for colIndex >= 0 {
		result = string(rune('A'+colIndex%26)) + result
		colIndex = colIndex/26 - 1
	}
	return result
}
