package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/sterlingcodes/alpha-cli/internal/common/config"
)

var tokenURL = "https://oauth2.googleapis.com/token" //nolint:gosec // OAuth endpoint URL, not a credential

const userAgent = "Alpha-CLI/1.0"

// Client is an authenticated Google API client using OAuth 2.0.
type Client struct {
	AccessToken string
	HTTPClient  *http.Client
}

// alpha365CredentialsPath returns the path to the shared credentials file
// written by the Alpha 365 desktop app.
func alpha365CredentialsPath() string {
	if runtime.GOOS == "windows" {
		appData := os.Getenv("APPDATA")
		if appData == "" {
			home, _ := os.UserHomeDir()
			appData = filepath.Join(home, "AppData", "Roaming")
		}
		return filepath.Join(appData, "alpha365", "google-credentials.json")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "alpha365", "google-credentials.json")
}

// loadAlpha365Credentials reads Google OAuth credentials from the shared
// file written by the Alpha 365 desktop app.
func loadAlpha365Credentials() (clientID, clientSecret, refreshToken string, err error) {
	data, err := os.ReadFile(alpha365CredentialsPath())
	if err != nil {
		return "", "", "", err
	}
	var creds struct {
		ClientID     string `json:"client_id"`
		ClientSecret string `json:"client_secret"`
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Unmarshal(data, &creds); err != nil {
		return "", "", "", err
	}
	if creds.ClientID == "" || creds.ClientSecret == "" || creds.RefreshToken == "" {
		return "", "", "", fmt.Errorf("incomplete credentials in shared file")
	}
	return creds.ClientID, creds.ClientSecret, creds.RefreshToken, nil
}

// NewClient reads Google OAuth credentials from config and returns an authenticated client.
// It falls back to shared credentials from the Alpha 365 desktop app if config keys are missing.
func NewClient() (*Client, error) {
	clientID, _ := config.MustGet("google_client_id")
	clientSecret, _ := config.MustGet("google_client_secret")
	refreshToken, _ := config.MustGet("google_refresh_token")

	// Fall back to shared credentials from Alpha 365 desktop app
	if clientID == "" || clientSecret == "" || refreshToken == "" {
		var err error
		clientID, clientSecret, refreshToken, err = loadAlpha365Credentials()
		if err != nil {
			return nil, fmt.Errorf("Google not configured. Connect via Alpha 365 → Services → Google Drive, or run 'alpha setup show google-oauth'")
		}
	}

	accessToken, err := ExchangeRefreshToken(clientID, clientSecret, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get access token: %w", err)
	}

	return &Client{
		AccessToken: accessToken,
		HTTPClient:  &http.Client{Timeout: 30 * time.Second},
	}, nil
}

// ExchangeRefreshToken exchanges a refresh token for an access token.
func ExchangeRefreshToken(clientID, clientSecret, refreshToken string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("refresh_token", refreshToken)
	data.Set("grant_type", "refresh_token")

	req, err := http.NewRequestWithContext(ctx, "POST", tokenURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error            string `json:"error"`
			ErrorDescription string `json:"error_description"`
		}
		if json.Unmarshal(body, &errResp) == nil && errResp.ErrorDescription != "" {
			return "", fmt.Errorf("OAuth error: %s", errResp.ErrorDescription)
		}
		return "", fmt.Errorf("OAuth error (HTTP %d): %s", resp.StatusCode, string(body))
	}

	var tokenResp struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return "", err
	}

	return tokenResp.AccessToken, nil
}

// DoRequest performs an authenticated HTTP request.
func (c *Client) DoRequest(method, reqURL string, body io.Reader, contentType string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, method, reqURL, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.AccessToken)
	req.Header.Set("User-Agent", userAgent)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		var errResp struct {
			Error struct {
				Message string `json:"message"`
				Code    int    `json:"code"`
			} `json:"error"`
		}
		if json.Unmarshal(respBody, &errResp) == nil && errResp.Error.Message != "" {
			return nil, fmt.Errorf("Google API error (%d): %s", resp.StatusCode, errResp.Error.Message)
		}
		return nil, fmt.Errorf("Google API error (HTTP %d): %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

// DoGet performs an authenticated GET request.
func (c *Client) DoGet(reqURL string) ([]byte, error) {
	return c.DoRequest("GET", reqURL, http.NoBody, "")
}

// DoPost performs an authenticated POST request with JSON body.
func (c *Client) DoPost(reqURL string, body io.Reader, contentType string) ([]byte, error) {
	return c.DoRequest("POST", reqURL, body, contentType)
}

// DoPatch performs an authenticated PATCH request with JSON body.
func (c *Client) DoPatch(reqURL string, body io.Reader, contentType string) ([]byte, error) {
	return c.DoRequest("PATCH", reqURL, body, contentType)
}
