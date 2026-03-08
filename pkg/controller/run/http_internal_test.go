package run

import "testing"

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
