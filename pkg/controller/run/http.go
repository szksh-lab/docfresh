package run

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"
)

func (c *Controller) request(ctx context.Context, h *HTTP) (*TemplateInput, error) { //nolint:cyclop
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

	sl := h.Language
	if sl == "" {
		sl = c.languageFromURL(h.URL)
	}

	content := string(b)
	content, err = extractRange(content, h.Range)
	if err != nil {
		return nil, fmt.Errorf("extract range from http response: %w", err)
	}
	result := &TemplateInput{
		Type:     "http",
		URL:      h.URL,
		Content:  string(b),
		Language: sl,
		Timeout:  h.Timeout,
		Vars:     h.Template.GetVars(),
	}
	if h.Template != nil {
		if err := renderTemplate(content, result, h.Template.Delims); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (c *Controller) languageFromURL(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	ext := path.Ext(u.Path)
	lang, ok := c.langs[ext]
	if !ok {
		return ""
	}
	return lang.Language
}
