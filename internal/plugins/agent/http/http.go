package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net"
	"net/http"
	"time"
)

type HTTPClient struct {
	base string
	hc   *http.Client
	auth string
}

func NewUserClient(apiBase, userSessionToken string) *HTTPClient {
	return &HTTPClient{
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

func NewAgentClient(apiBase, agentToken string) *HTTPClient {
	c := NewUserClient(apiBase, agentToken)
	return c
}

func (c *HTTPClient) PostJSON(ctx context.Context, path string, in, out any) (int, error) {
	b, _ := json.Marshal(in)
	req, _ := http.NewRequestWithContext(ctx, "POST", c.base+path, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	// if c.auth != "" {
	// 	req.Header.Set("Authorization", c.auth)
	// }
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

func (c *HTTPClient) GetJSON(ctx context.Context, path string, out any) (int, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", c.base+path, nil)
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
