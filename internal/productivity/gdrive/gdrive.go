package gdrive

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mime"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	googleauth "github.com/unstablemind/pocket/internal/common/google"
	"github.com/unstablemind/pocket/pkg/output"
)

var baseURL = "https://www.googleapis.com"

const maxReadSize = 100 * 1024 // 100KB cap for read command

// DriveSearchResult holds search results
type DriveSearchResult struct {
	Query string      `json:"query"`
	Files []DriveFile `json:"files"`
	Count int         `json:"count"`
}

// DriveFile holds file metadata from search
type DriveFile struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	MimeType   string `json:"mime_type"`
	Size       string `json:"size,omitempty"`
	CreatedAt  string `json:"created_at,omitempty"`
	ModifiedAt string `json:"modified_at,omitempty"`
	URL        string `json:"url,omitempty"`
}

// FileInfo holds detailed file metadata
type FileInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	MimeType    string `json:"mime_type"`
	Size        string `json:"size,omitempty"`
	CreatedAt   string `json:"created_at,omitempty"`
	ModifiedAt  string `json:"modified_at,omitempty"`
	URL         string `json:"url,omitempty"`
	Description string `json:"description,omitempty"`
	Owner       string `json:"owner,omitempty"`
}

// NewCmd returns the Google Drive command
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "gdrive",
		Aliases: []string{"drive"},
		Short:   "Google Drive commands (full read/write access)",
	}

	cmd.AddCommand(newSearchCmd())
	cmd.AddCommand(newInfoCmd())
	cmd.AddCommand(newListCmd())
	cmd.AddCommand(newReadCmd())
	cmd.AddCommand(newDownloadCmd())
	cmd.AddCommand(newUploadCmd())
	cmd.AddCommand(newMkdirCmd())
	cmd.AddCommand(newUpdateCmd())
	cmd.AddCommand(newDeleteCmd())

	return cmd
}

func newSearchCmd() *cobra.Command {
	var limit int

	cmd := &cobra.Command{
		Use:   "search [query]",
		Short: "Search files in Google Drive",
		Long:  `Search files using Google Drive query syntax: name contains 'test', mimeType='application/pdf', etc.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			query := args[0]
			params := url.Values{
				"q":        {query},
				"fields":   {"files(id,name,mimeType,size,createdTime,modifiedTime,webViewLink)"},
				"pageSize": {fmt.Sprintf("%d", limit)},
			}

			reqURL := fmt.Sprintf("%s/drive/v3/files?%s", baseURL, params.Encode())
			data, err := client.DoGet(reqURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var resp struct {
				Files []struct {
					ID           string `json:"id"`
					Name         string `json:"name"`
					MimeType     string `json:"mimeType"`
					Size         string `json:"size"`
					CreatedTime  string `json:"createdTime"`
					ModifiedTime string `json:"modifiedTime"`
					WebViewLink  string `json:"webViewLink"`
				} `json:"files"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse response: %s", err.Error()), nil)
			}

			files := make([]DriveFile, 0, len(resp.Files))
			for _, f := range resp.Files {
				files = append(files, DriveFile{
					ID:         f.ID,
					Name:       f.Name,
					MimeType:   f.MimeType,
					Size:       formatSize(f.Size),
					CreatedAt:  formatTime(f.CreatedTime),
					ModifiedAt: formatTime(f.ModifiedTime),
					URL:        f.WebViewLink,
				})
			}

			result := DriveSearchResult{
				Query: query,
				Files: files,
				Count: len(files),
			}

			return output.Print(result)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "l", 10, "Maximum number of results")

	return cmd
}

func newInfoCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info [file-id]",
		Short: "Get file metadata",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			fileID := args[0]
			params := url.Values{
				"fields": {"id,name,mimeType,size,createdTime,modifiedTime,webViewLink,description,owners"},
			}

			reqURL := fmt.Sprintf("%s/drive/v3/files/%s?%s", baseURL, url.PathEscape(fileID), params.Encode())
			data, err := client.DoGet(reqURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var resp struct {
				ID           string `json:"id"`
				Name         string `json:"name"`
				MimeType     string `json:"mimeType"`
				Size         string `json:"size"`
				CreatedTime  string `json:"createdTime"`
				ModifiedTime string `json:"modifiedTime"`
				WebViewLink  string `json:"webViewLink"`
				Description  string `json:"description"`
				Owners       []struct {
					DisplayName string `json:"displayName"`
				} `json:"owners"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", fmt.Sprintf("Failed to parse response: %s", err.Error()), nil)
			}

			ownerName := ""
			if len(resp.Owners) > 0 {
				ownerName = resp.Owners[0].DisplayName
			}

			result := FileInfo{
				ID:          resp.ID,
				Name:        resp.Name,
				MimeType:    resp.MimeType,
				Size:        formatSize(resp.Size),
				CreatedAt:   formatTime(resp.CreatedTime),
				ModifiedAt:  formatTime(resp.ModifiedTime),
				URL:         resp.WebViewLink,
				Description: resp.Description,
				Owner:       ownerName,
			}

			return output.Print(result)
		},
	}

	return cmd
}

func newListCmd() *cobra.Command {
	var folderID string
	var limit int

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List files in a folder (default: root)",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			query := fmt.Sprintf("'%s' in parents and trashed = false", folderID)
			params := url.Values{
				"q":        {query},
				"fields":   {"files(id,name,mimeType,size,createdTime,modifiedTime,webViewLink)"},
				"pageSize": {fmt.Sprintf("%d", limit)},
				"orderBy":  {"folder,name"},
			}

			reqURL := fmt.Sprintf("%s/drive/v3/files?%s", baseURL, params.Encode())
			data, err := client.DoGet(reqURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var resp struct {
				Files []struct {
					ID           string `json:"id"`
					Name         string `json:"name"`
					MimeType     string `json:"mimeType"`
					Size         string `json:"size"`
					CreatedTime  string `json:"createdTime"`
					ModifiedTime string `json:"modifiedTime"`
					WebViewLink  string `json:"webViewLink"`
				} `json:"files"`
			}

			if err := json.Unmarshal(data, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			files := make([]DriveFile, 0, len(resp.Files))
			for _, f := range resp.Files {
				files = append(files, DriveFile{
					ID:         f.ID,
					Name:       f.Name,
					MimeType:   f.MimeType,
					Size:       formatSize(f.Size),
					CreatedAt:  formatTime(f.CreatedTime),
					ModifiedAt: formatTime(f.ModifiedTime),
					URL:        f.WebViewLink,
				})
			}

			return output.Print(map[string]any{
				"folder": folderID,
				"files":  files,
				"count":  len(files),
			})
		},
	}

	cmd.Flags().StringVar(&folderID, "folder", "root", "Folder ID to list (default: root)")
	cmd.Flags().IntVarP(&limit, "limit", "l", 20, "Maximum number of results")

	return cmd
}

func newReadCmd() *cobra.Command {
	var format string

	cmd := &cobra.Command{
		Use:   "read [file-id]",
		Short: "Read file content as text",
		Long:  "Exports Google Workspace files (Docs, Sheets, Slides) or downloads regular text files. Output capped at 100KB.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			fileID := args[0]

			// First get file metadata to determine type
			metaURL := fmt.Sprintf("%s/drive/v3/files/%s?fields=id,name,mimeType,size", baseURL, url.PathEscape(fileID))
			metaData, err := client.DoGet(metaURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var meta struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				MimeType string `json:"mimeType"`
				Size     string `json:"size"`
			}
			if err := json.Unmarshal(metaData, &meta); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			var content []byte

			switch {
			case strings.HasPrefix(meta.MimeType, "application/vnd.google-apps.document"):
				exportMime := exportMimeForDocs(format)
				content, err = exportFile(client, fileID, exportMime)
			case strings.HasPrefix(meta.MimeType, "application/vnd.google-apps.spreadsheet"):
				content, err = exportFile(client, fileID, "text/csv")
			case strings.HasPrefix(meta.MimeType, "application/vnd.google-apps.presentation"):
				content, err = exportFile(client, fileID, "text/plain")
			case isTextMime(meta.MimeType):
				content, err = downloadFile(client, fileID)
			default:
				return output.PrintError("binary_file", fmt.Sprintf("File '%s' is binary (%s). Use 'gdrive download' instead.", meta.Name, meta.MimeType), nil)
			}

			if err != nil {
				return output.PrintError("read_failed", err.Error(), nil)
			}

			// Cap output
			truncated := false
			if len(content) > maxReadSize {
				content = content[:maxReadSize]
				truncated = true
			}

			return output.Print(map[string]any{
				"id":        meta.ID,
				"name":      meta.Name,
				"mime_type": meta.MimeType,
				"content":   string(content),
				"truncated": truncated,
			})
		},
	}

	cmd.Flags().StringVar(&format, "format", "text", "Export format for Google Docs: text, md, html")

	return cmd
}

func newDownloadCmd() *cobra.Command {
	var outputPath string

	cmd := &cobra.Command{
		Use:   "download [file-id]",
		Short: "Download file to local disk",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			fileID := args[0]

			// Get metadata for filename
			metaURL := fmt.Sprintf("%s/drive/v3/files/%s?fields=id,name,mimeType", baseURL, url.PathEscape(fileID))
			metaData, err := client.DoGet(metaURL)
			if err != nil {
				return output.PrintError("request_failed", err.Error(), nil)
			}

			var meta struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				MimeType string `json:"mimeType"`
			}
			if err := json.Unmarshal(metaData, &meta); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			var content []byte

			// Google Workspace files need export
			if strings.HasPrefix(meta.MimeType, "application/vnd.google-apps.") {
				exportMime, ext := exportForDownload(meta.MimeType)
				content, err = exportFile(client, fileID, exportMime)
				if err != nil {
					return output.PrintError("download_failed", err.Error(), nil)
				}
				// Append extension if needed
				if !strings.HasSuffix(meta.Name, ext) {
					meta.Name += ext
				}
			} else {
				content, err = downloadFile(client, fileID)
				if err != nil {
					return output.PrintError("download_failed", err.Error(), nil)
				}
			}

			dest := outputPath
			if dest == "" {
				dest = meta.Name
			}

			if err := os.WriteFile(dest, content, 0o644); err != nil {
				return output.PrintError("write_failed", fmt.Sprintf("Failed to write file: %s", err.Error()), nil)
			}

			return output.Print(map[string]any{
				"downloaded": true,
				"file_id":    meta.ID,
				"name":       meta.Name,
				"path":       dest,
				"size":       formatSize(fmt.Sprintf("%d", len(content))),
			})
		},
	}

	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "Output file path (default: original filename)")

	return cmd
}

func newUploadCmd() *cobra.Command {
	var name string
	var folderID string

	cmd := &cobra.Command{
		Use:   "upload [path]",
		Short: "Upload a local file to Google Drive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			filePath := args[0]

			fileData, err := os.ReadFile(filePath)
			if err != nil {
				return output.PrintError("read_failed", fmt.Sprintf("Failed to read file: %s", err.Error()), nil)
			}

			uploadName := name
			if uploadName == "" {
				uploadName = filepath.Base(filePath)
			}

			// Detect MIME type
			mimeType := mime.TypeByExtension(filepath.Ext(filePath))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			// Build multipart upload
			metadata := map[string]any{
				"name": uploadName,
			}
			if folderID != "" {
				metadata["parents"] = []string{folderID}
			}

			metaJSON, err := json.Marshal(metadata)
			if err != nil {
				return output.PrintError("marshal_failed", err.Error(), nil)
			}

			var body bytes.Buffer
			writer := multipart.NewWriter(&body)

			// Metadata part
			metaHeader := make(textproto.MIMEHeader)
			metaHeader.Set("Content-Type", "application/json; charset=UTF-8")
			metaPart, err := writer.CreatePart(metaHeader)
			if err != nil {
				return output.PrintError("upload_failed", err.Error(), nil)
			}
			metaPart.Write(metaJSON)

			// File part
			fileHeader := make(textproto.MIMEHeader)
			fileHeader.Set("Content-Type", mimeType)
			filePart, err := writer.CreatePart(fileHeader)
			if err != nil {
				return output.PrintError("upload_failed", err.Error(), nil)
			}
			filePart.Write(fileData)

			writer.Close()

			reqURL := fmt.Sprintf("%s/upload/drive/v3/files?uploadType=multipart&fields=id,name,mimeType,size,webViewLink", baseURL)
			respData, err := client.DoPost(reqURL, &body, writer.FormDataContentType())
			if err != nil {
				return output.PrintError("upload_failed", err.Error(), nil)
			}

			var resp struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				MimeType    string `json:"mimeType"`
				Size        string `json:"size"`
				WebViewLink string `json:"webViewLink"`
			}
			if err := json.Unmarshal(respData, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			return output.Print(map[string]any{
				"uploaded": true,
				"id":       resp.ID,
				"name":     resp.Name,
				"mime_type": resp.MimeType,
				"size":     formatSize(resp.Size),
				"url":      resp.WebViewLink,
			})
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "Override file name in Drive")
	cmd.Flags().StringVar(&folderID, "folder", "", "Parent folder ID")

	return cmd
}

func newMkdirCmd() *cobra.Command {
	var parentID string

	cmd := &cobra.Command{
		Use:   "mkdir [name]",
		Short: "Create a folder in Google Drive",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			metadata := map[string]any{
				"name":     args[0],
				"mimeType": "application/vnd.google-apps.folder",
			}
			if parentID != "" {
				metadata["parents"] = []string{parentID}
			}

			metaJSON, err := json.Marshal(metadata)
			if err != nil {
				return output.PrintError("marshal_failed", err.Error(), nil)
			}

			reqURL := fmt.Sprintf("%s/drive/v3/files?fields=id,name,webViewLink", baseURL)
			respData, err := client.DoPost(reqURL, bytes.NewReader(metaJSON), "application/json")
			if err != nil {
				return output.PrintError("create_failed", err.Error(), nil)
			}

			var resp struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				WebViewLink string `json:"webViewLink"`
			}
			if err := json.Unmarshal(respData, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			return output.Print(map[string]any{
				"created": true,
				"id":      resp.ID,
				"name":    resp.Name,
				"url":     resp.WebViewLink,
			})
		},
	}

	cmd.Flags().StringVar(&parentID, "parent", "", "Parent folder ID")

	return cmd
}

func newUpdateCmd() *cobra.Command {
	var name string
	var contentPath string

	cmd := &cobra.Command{
		Use:   "update [file-id]",
		Short: "Update file metadata or content",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			fileID := args[0]

			if contentPath != "" {
				// Update content via upload
				fileData, err := os.ReadFile(contentPath)
				if err != nil {
					return output.PrintError("read_failed", err.Error(), nil)
				}

				mimeType := mime.TypeByExtension(filepath.Ext(contentPath))
				if mimeType == "" {
					mimeType = "application/octet-stream"
				}

				reqURL := fmt.Sprintf("%s/upload/drive/v3/files/%s?uploadType=media&fields=id,name,mimeType,size,modifiedTime", baseURL, url.PathEscape(fileID))
				respData, err := client.DoPatch(reqURL, bytes.NewReader(fileData), mimeType)
				if err != nil {
					return output.PrintError("update_failed", err.Error(), nil)
				}

				var resp struct {
					ID           string `json:"id"`
					Name         string `json:"name"`
					MimeType     string `json:"mimeType"`
					Size         string `json:"size"`
					ModifiedTime string `json:"modifiedTime"`
				}
				if err := json.Unmarshal(respData, &resp); err != nil {
					return output.PrintError("parse_failed", err.Error(), nil)
				}

				return output.Print(map[string]any{
					"updated":     true,
					"id":          resp.ID,
					"name":        resp.Name,
					"size":        formatSize(resp.Size),
					"modified_at": formatTime(resp.ModifiedTime),
				})
			}

			if name != "" {
				// Update metadata only
				metadata := map[string]any{"name": name}
				metaJSON, err := json.Marshal(metadata)
				if err != nil {
					return output.PrintError("marshal_failed", err.Error(), nil)
				}

				reqURL := fmt.Sprintf("%s/drive/v3/files/%s?fields=id,name,modifiedTime", baseURL, url.PathEscape(fileID))
				respData, err := client.DoPatch(reqURL, bytes.NewReader(metaJSON), "application/json")
				if err != nil {
					return output.PrintError("update_failed", err.Error(), nil)
				}

				var resp struct {
					ID           string `json:"id"`
					Name         string `json:"name"`
					ModifiedTime string `json:"modifiedTime"`
				}
				if err := json.Unmarshal(respData, &resp); err != nil {
					return output.PrintError("parse_failed", err.Error(), nil)
				}

				return output.Print(map[string]any{
					"updated":     true,
					"id":          resp.ID,
					"name":        resp.Name,
					"modified_at": formatTime(resp.ModifiedTime),
				})
			}

			return output.PrintError("no_changes", "Specify --name or --content to update", nil)
		},
	}

	cmd.Flags().StringVar(&name, "name", "", "New file name")
	cmd.Flags().StringVar(&contentPath, "content", "", "Path to new file content")

	return cmd
}

func newDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [file-id]",
		Short: "Move file to trash (safe, not permanent)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := googleauth.NewClient()
			if err != nil {
				return output.PrintError("auth_failed", err.Error(), nil)
			}

			fileID := args[0]
			body := bytes.NewReader([]byte(`{"trashed":true}`))

			reqURL := fmt.Sprintf("%s/drive/v3/files/%s?fields=id,name,trashed", baseURL, url.PathEscape(fileID))
			respData, err := client.DoPatch(reqURL, body, "application/json")
			if err != nil {
				return output.PrintError("delete_failed", err.Error(), nil)
			}

			var resp struct {
				ID      string `json:"id"`
				Name    string `json:"name"`
				Trashed bool   `json:"trashed"`
			}
			if err := json.Unmarshal(respData, &resp); err != nil {
				return output.PrintError("parse_failed", err.Error(), nil)
			}

			return output.Print(map[string]any{
				"trashed": true,
				"id":      resp.ID,
				"name":    resp.Name,
			})
		},
	}

	return cmd
}

// exportFile exports a Google Workspace file in the given MIME type.
func exportFile(client *googleauth.Client, fileID, exportMime string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/drive/v3/files/%s/export?mimeType=%s", baseURL, url.PathEscape(fileID), url.QueryEscape(exportMime))
	return client.DoGet(reqURL)
}

// downloadFile downloads a regular (non-Workspace) file.
func downloadFile(client *googleauth.Client, fileID string) ([]byte, error) {
	reqURL := fmt.Sprintf("%s/drive/v3/files/%s?alt=media", baseURL, url.PathEscape(fileID))
	return client.DoGet(reqURL)
}

// exportMimeForDocs returns the export MIME type for Google Docs based on requested format.
func exportMimeForDocs(format string) string {
	switch strings.ToLower(format) {
	case "md", "markdown":
		return "text/markdown"
	case "html":
		return "text/html"
	default:
		return "text/plain"
	}
}

// exportForDownload returns the export MIME type and file extension for downloading Google Workspace files.
func exportForDownload(mimeType string) (string, string) {
	switch {
	case strings.Contains(mimeType, "document"):
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document", ".docx"
	case strings.Contains(mimeType, "spreadsheet"):
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", ".xlsx"
	case strings.Contains(mimeType, "presentation"):
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation", ".pptx"
	case strings.Contains(mimeType, "drawing"):
		return "image/png", ".png"
	default:
		return "application/pdf", ".pdf"
	}
}

// isTextMime returns true if the MIME type represents a text-readable file.
func isTextMime(mimeType string) bool {
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}
	textTypes := []string{
		"application/json",
		"application/xml",
		"application/javascript",
		"application/x-yaml",
		"application/toml",
		"application/x-sh",
	}
	for _, t := range textTypes {
		if mimeType == t {
			return true
		}
	}
	return false
}

func formatTime(isoTime string) string {
	if isoTime == "" {
		return ""
	}
	t, err := time.Parse(time.RFC3339, isoTime)
	if err != nil {
		return isoTime
	}

	diff := time.Since(t)
	switch {
	case diff < time.Minute:
		return "now"
	case diff < time.Hour:
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	case diff < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	case diff < 7*24*time.Hour:
		return fmt.Sprintf("%dd ago", int(diff.Hours()/24))
	case diff < 30*24*time.Hour:
		return fmt.Sprintf("%dw ago", int(diff.Hours()/(24*7)))
	case diff < 365*24*time.Hour:
		return fmt.Sprintf("%dmo ago", int(diff.Hours()/(24*30)))
	default:
		return fmt.Sprintf("%dy ago", int(diff.Hours()/(24*365)))
	}
}

func formatSize(sizeStr string) string {
	if sizeStr == "" {
		return ""
	}

	var size int64
	if _, err := fmt.Sscanf(sizeStr, "%d", &size); err != nil {
		return sizeStr
	}

	switch {
	case size >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(size)/float64(1<<30))
	case size >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(size)/float64(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(size)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", size)
	}
}
