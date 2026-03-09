package gdocs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	googleauth "github.com/sterlingcodes/alpha-cli/internal/common/google"
	"github.com/sterlingcodes/alpha-cli/pkg/output"
)

var baseURL = "https://docs.googleapis.com/v1/documents"

// DocInfo holds metadata about a document
type DocInfo struct {
	DocumentID string `json:"document_id"`
	Title      string `json:"title"`
	RevisionID string `json:"revision_id"`
}

// NewCmd returns the Google Docs command
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gdocs",
		Aliases: []string{"docs"},
		Short:   "Google Docs commands",
	}

	cmd.AddCommand(newGetCmd())
	cmd.AddCommand(newReadCmd())
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newAppendCmd())
	cmd.AddCommand(newInsertCmd())
	cmd.AddCommand(newReplaceCmd())

	return cmd
}

func newGetCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [doc-id]",
		Short: "Get document metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			return runGet(client, args[0])
		},
	}

	return cmd
}

func runGet(client *googleauth.Client, docID string) error {
	reqURL := fmt.Sprintf("%s/%s?fields=documentId,title,revisionId",
		baseURL, url.PathEscape(docID))

	data, err := client.DoGet(reqURL)
	if err != nil {
		return output.PrintError("request_failed", err.Error(), nil)
	}

	var resp struct {
		DocumentID string `json:"documentId"`
		Title      string `json:"title"`
		RevisionID string `json:"revisionId"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse response: %s", err.Error()), nil)
	}

	result := DocInfo{
		DocumentID: resp.DocumentID,
		Title:      resp.Title,
		RevisionID: resp.RevisionID,
	}

	return output.Print(result)
}

func newReadCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "read [doc-id]",
		Short: "Read document plain text",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			return runRead(client, args[0])
		},
	}

	return cmd
}

func runRead(client *googleauth.Client, docID string) error {
	reqURL := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(docID))

	data, err := client.DoGet(reqURL)
	if err != nil {
		return output.PrintError("request_failed", err.Error(), nil)
	}

	text, err := extractPlainText(data)
	if err != nil {
		return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse document: %s", err.Error()), nil)
	}

	return output.Print(map[string]any{
		"document_id": docID,
		"text":        text,
	})
}

// extractPlainText walks body.content[].paragraph.elements[].textRun.content
func extractPlainText(data []byte) (string, error) {
	var doc struct {
		Body struct {
			Content []struct {
				Paragraph *struct {
					Elements []struct {
						TextRun *struct {
							Content string `json:"content"`
						} `json:"textRun"`
					} `json:"elements"`
				} `json:"paragraph"`
			} `json:"content"`
		} `json:"body"`
	}

	if err := json.Unmarshal(data, &doc); err != nil {
		return "", err
	}

	var sb strings.Builder
	for _, content := range doc.Body.Content {
		if content.Paragraph != nil {
			for _, elem := range content.Paragraph.Elements {
				if elem.TextRun != nil {
					sb.WriteString(elem.TextRun.Content)
				}
			}
		}
	}

	return sb.String(), nil
}

// getEndIndex fetches the document and returns the endOfSegment index
func getEndIndex(client *googleauth.Client, docID string) (int, error) {
	reqURL := fmt.Sprintf("%s/%s", baseURL, url.PathEscape(docID))

	data, err := client.DoGet(reqURL)
	if err != nil {
		return 0, err
	}

	var doc struct {
		Body struct {
			Content []struct {
				EndIndex int `json:"endIndex"`
			} `json:"content"`
		} `json:"body"`
	}

	if err := json.Unmarshal(data, &doc); err != nil {
		return 0, err
	}

	if len(doc.Body.Content) == 0 {
		return 1, nil
	}

	// The last structural element's endIndex - 1 is where we insert
	lastIdx := doc.Body.Content[len(doc.Body.Content)-1].EndIndex
	if lastIdx < 1 {
		return 1, nil
	}

	return lastIdx - 1, nil
}

func newCreateCmd() *cobra.Command {
	var title string

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new document",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			return runCreate(client, title)
		},
	}

	cmd.Flags().StringVar(&title, "title", "", "Document title (required)")
	_ = cmd.MarkFlagRequired("title")

	return cmd
}

func runCreate(client *googleauth.Client, title string) error {
	body := map[string]string{
		"title": title,
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
		DocumentID string `json:"documentId"`
		Title      string `json:"title"`
		RevisionID string `json:"revisionId"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return output.PrintError("parse_failed", err.Error(), nil)
	}

	return output.Print(map[string]any{
		"document_id": resp.DocumentID,
		"title":       resp.Title,
		"revision_id": resp.RevisionID,
	})
}

func newAppendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "append [doc-id] [text]",
		Short: "Append text to the end of a document",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			return runAppend(client, args[0], args[1])
		},
	}

	return cmd
}

func runAppend(client *googleauth.Client, docID, text string) error {
	endIdx, err := getEndIndex(client, docID)
	if err != nil {
		return output.PrintError("request_failed", fmt.Sprintf("Failed to get document end index: %s", err.Error()), nil)
	}

	return doBatchUpdate(client, docID, []map[string]any{
		{
			"insertText": map[string]any{
				"location": map[string]any{
					"index": endIdx,
				},
				"text": text,
			},
		},
	})
}

func newInsertCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "insert [doc-id] [index] [text]",
		Short: "Insert text at a specific index",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			index, err := strconv.Atoi(args[1])
			if err != nil {
				return output.PrintError("invalid_index", fmt.Sprintf("Index must be an integer: %s", err.Error()), nil)
			}

			return runInsert(client, args[0], index, args[2])
		},
	}

	return cmd
}

func runInsert(client *googleauth.Client, docID string, index int, text string) error {
	return doBatchUpdate(client, docID, []map[string]any{
		{
			"insertText": map[string]any{
				"location": map[string]any{
					"index": index,
				},
				"text": text,
			},
		},
	})
}

func newReplaceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "replace [doc-id] [find] [replace]",
		Short: "Replace all occurrences of text",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return err
			}

			return runReplace(client, args[0], args[1], args[2])
		},
	}

	return cmd
}

func runReplace(client *googleauth.Client, docID, find, replace string) error {
	return doBatchUpdate(client, docID, []map[string]any{
		{
			"replaceAllText": map[string]any{
				"containsText": map[string]any{
					"text":      find,
					"matchCase": true,
				},
				"replaceText": replace,
			},
		},
	})
}

func doBatchUpdate(client *googleauth.Client, docID string, requests []map[string]any) error {
	body := map[string]any{
		"requests": requests,
	}
	bodyJSON, err := json.Marshal(body)
	if err != nil {
		return output.PrintError("marshal_failed", err.Error(), nil)
	}

	reqURL := fmt.Sprintf("%s/%s:batchUpdate", baseURL, url.PathEscape(docID))

	data, err := client.DoPost(reqURL, bytes.NewReader(bodyJSON), "application/json")
	if err != nil {
		return output.PrintError("request_failed", err.Error(), nil)
	}

	var resp struct {
		DocumentID string `json:"documentId"`
		Replies    []any  `json:"replies"`
	}

	if err := json.Unmarshal(data, &resp); err != nil {
		return output.PrintError("parse_failed", err.Error(), nil)
	}

	return output.Print(map[string]any{
		"document_id": resp.DocumentID,
		"replies":     resp.Replies,
	})
}
