package google

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
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

// NewClient reads Google OAuth credentials from config and returns an authenticated client.
func NewClient() (*Client, error) {
	clientID, err := config.MustGet("google_client_id")
	if err != nil {
		return nil, fmt.Errorf("google_client_id not configured. Run 'alpha setup show google-oauth' for setup instructions")
	}

	clientSecret, err := config.MustGet("google_client_secret")
	if err != nil {
		return nil, fmt.Errorf("google_client_secret not configured. Run 'alpha setup show google-oauth' for setup instructions")
	}

	refreshToken, err := config.MustGet("google_refresh_token")
	if err != nil {
		return nil, fmt.Errorf("google_refresh_token not configured. Run 'alpha setup show google-oauth' for setup instructions")
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
