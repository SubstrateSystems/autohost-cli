// Package enrollapi provides an HTTP client for the AutoHost enrollment and agent API.
package enrollapi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"
	"time"
)

// Client is a lightweight HTTP client that attaches a Bearer token to every request.
type Client struct {
	base string
	auth string
	hc   *http.Client
}

// NewUserClient creates a client authenticated with a user session token.
func NewUserClient(apiBase, userSessionToken string) *Client {
	return &Client{
		base: apiBase,
		auth: "Bearer " + userSessionToken,
		hc: &http.Client{
			Timeout: 12 * time.Second,
			Transport: &http.Transport{
				DialContext:         (&net.Dialer{Timeout: 5 * time.Second}).DialContext,
				TLSHandshakeTimeout: 5 * time.Second,
				IdleConnTimeout:     30 * time.Second,
				MaxIdleConns:        10,
			},
		},
	}
}

// NewAgentClient creates a client authenticated with an agent token.
func NewAgentClient(apiBase, agentToken string) *Client {
	return NewUserClient(apiBase, agentToken)
}

// PostJSON marshals in, POSTs to base+path, and decodes the response body into out (may be nil).
// On HTTP 4xx/5xx, returns the response body as the error message.
func (c *Client) PostJSON(ctx context.Context, path string, in, out any) (int, error) {
	b, _ := json.Marshal(in)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, c.base+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	if c.auth != "" {
		req.Header.Set("Authorization", c.auth)
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return resp.StatusCode, fmt.Errorf("%s", strings.TrimSpace(string(body)))
	}
	if out != nil {
		_ = json.Unmarshal(body, out)
	}
	return resp.StatusCode, nil
}

// GetJSON GETs base+path and decodes the response body into out (may be nil).
func (c *Client) GetJSON(ctx context.Context, path string, out any) (int, error) {
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, c.base+path, nil)
	if c.auth != "" {
		req.Header.Set("Authorization", c.auth)
	}
	resp, err := c.hc.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	if out != nil {
		_ = json.NewDecoder(resp.Body).Decode(out)
	}
	return resp.StatusCode, nil
}
