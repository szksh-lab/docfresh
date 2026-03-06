package run

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"
)

func (c *Controller) request(ctx context.Context, h *HTTP) (*TemplateInput, error) {
	if h.Timeout == 0 {
		h.Timeout = 5
	}
	if h.Timeout > 0 {
		requestCtx, cancel := context.WithTimeout(ctx, time.Duration(h.Timeout)*time.Second)
		defer cancel()
		ctx = requestCtx
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}
	if len(h.Header) > 0 {
		req.Header = h.Header
	}

	resp, err := c.httpClient.Do(req) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return &TemplateInput{
		Type:    "http",
		URL:     h.URL,
		Content: string(b),
		Timeout: h.Timeout,
	}, nil
}
