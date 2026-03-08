package run

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestCallHTTP(t *testing.T) { //nolint:funlen
	t.Parallel()
	tests := []struct {
		name       string
		handler    http.HandlerFunc
		h          *HTTP
		langs      map[string]*Language
		want       *TemplateInput
		wantErr    bool
		wantHeader http.Header // headers the server should receive
	}{
		{
			name: "basic success",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "hello")
			},
			h:     &HTTP{},
			langs: nil,
			want: &TemplateInput{
				Type:    "http",
				Content: "hello",
				Timeout: 5,
			},
		},
		{
			name: "language from URL extension",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "{}")
			},
			h: &HTTP{},
			langs: map[string]*Language{
				".json": {Language: "json"},
			},
			want: &TemplateInput{
				Type:     "http",
				Content:  "{}",
				Language: "json",
				Timeout:  5,
			},
		},
		{
			name: "explicit language overrides URL detection",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "key: val")
			},
			h: &HTTP{
				Language: "yaml",
			},
			langs: map[string]*Language{
				".json": {Language: "json"},
			},
			want: &TemplateInput{
				Type:     "http",
				Content:  "key: val",
				Language: "yaml",
				Timeout:  5,
			},
		},
		{
			name: "JSON auto-detection",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, `{"key":"value"}`)
			},
			h:     &HTTP{},
			langs: nil,
			want: &TemplateInput{
				Type:     "http",
				Content:  `{"key":"value"}`,
				Language: "json",
				Timeout:  5,
			},
		},
		{
			name: "non-200 status",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			},
			h:       &HTTP{},
			langs:   nil,
			wantErr: true,
		},
		{
			name: "custom headers",
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.Header.Get("X-Custom") != "test-value" {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				fmt.Fprint(w, "ok")
			},
			h: &HTTP{
				Header: http.Header{
					"X-Custom": []string{"test-value"},
				},
			},
			langs: nil,
			want: &TemplateInput{
				Type:    "http",
				Content: "ok",
				Timeout: 5,
			},
		},
		{
			name: "with range",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "line1\nline2\nline3\n")
			},
			h: &HTTP{
				Range: &Range{
					Start: new(1),
					End:   new(2),
				},
			},
			langs: nil,
			want: &TemplateInput{
				Type:    "http",
				Content: "line1\nline2\nline3\n",
				Timeout: 5,
			},
		},
		{
			name: "default timeout",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "data")
			},
			h:     &HTTP{},
			langs: nil,
			want: &TemplateInput{
				Type:    "http",
				Content: "data",
				Timeout: 5,
			},
		},
		{
			name: "custom timeout",
			handler: func(w http.ResponseWriter, _ *http.Request) {
				fmt.Fprint(w, "data")
			},
			h: &HTTP{
				Timeout: 10,
			},
			langs: nil,
			want: &TemplateInput{
				Type:    "http",
				Content: "data",
				Timeout: 10,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ts := httptest.NewServer(tt.handler)
			defer ts.Close()

			// Set URL with a path matching the test case
			urlPath := "/file"
			if tt.langs != nil {
				for ext := range tt.langs {
					urlPath = "/file" + ext
					break
				}
			}
			tt.h.URL = ts.URL + urlPath

			// Set expected URL in want
			if tt.want != nil {
				tt.want.URL = tt.h.URL
			}

			got, err := callHTTP(context.Background(), tt.h, ts.Client(), tt.langs)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestLanguageFromURL(t *testing.T) { //nolint:funlen
	t.Parallel()
	tests := []struct {
		name   string
		rawURL string
		langs  map[string]*Language
		want   string
	}{
		{
			name:   "known extension",
			rawURL: "https://example.com/file.go",
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			want: "go",
		},
		{
			name:   "unknown extension",
			rawURL: "https://example.com/file.xyz",
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			want: "",
		},
		{
			name:   "no extension",
			rawURL: "https://example.com/path/file",
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			want: "",
		},
		{
			name:   "nil langs map",
			rawURL: "https://example.com/file.go",
			langs:  nil,
			want:   "",
		},
		{
			name:   "URL with query params",
			rawURL: "https://example.com/file.json?v=1",
			langs: map[string]*Language{
				".json": {Language: "json"},
			},
			want: "json",
		},
		{
			name:   "empty URL",
			rawURL: "",
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := languageFromURL(tt.rawURL, tt.langs)
			if got != tt.want {
				t.Errorf("languageFromURL() = %q, want %q", got, tt.want)
			}
		})
	}
}
