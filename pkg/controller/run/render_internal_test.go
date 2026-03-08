package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWrapDetailsTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
		summary string
		want    string
	}{
		{
			name:    "content with trailing newline",
			content: "hello\n",
			summary: "Output",
			want:    "<details>\n<summary>Output</summary>\n\nhello\n\n</details>",
		},
		{
			name:    "content without trailing newline",
			content: "hello",
			summary: "Output",
			want:    "<details>\n<summary>Output</summary>\n\nhello\n\n</details>",
		},
		{
			name:    "multiline content",
			content: "line1\nline2\n",
			summary: "Details",
			want:    "<details>\n<summary>Details</summary>\n\nline1\nline2\n\n</details>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := wrapDetailsTag(tt.content, tt.summary)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("wrapDetailsTag() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestDefaultDetailsTagSummary(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		result *TemplateInput
		want   string
	}{
		{
			name:   "command",
			result: &TemplateInput{Type: "command"},
			want:   "Output",
		},
		{
			name:   "local-file",
			result: &TemplateInput{Type: "local-file", Path: "README.md"},
			want:   "README.md",
		},
		{
			name:   "http",
			result: &TemplateInput{Type: "http", URL: "https://example.com"},
			want:   "https://example.com",
		},
		{
			name:   "github-content with ref",
			result: &TemplateInput{Type: "github-content", Owner: "org", Repo: "repo", Ref: "main", Path: "file.txt"},
			want:   "org/repo/file.txt@main",
		},
		{
			name:   "github-content without ref",
			result: &TemplateInput{Type: "github-content", Owner: "org", Repo: "repo", Path: "file.txt"},
			want:   "org/repo/file.txt",
		},
		{
			name:   "unknown type",
			result: &TemplateInput{Type: "unknown"},
			want:   "Output",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := defaultDetailsTagSummary(tt.result)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("defaultDetailsTagSummary() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRenderFile_DetailsTag(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		result *TemplateInput
		want   string
	}{
		{
			name: "plain content with details tag",
			result: &TemplateInput{
				Content:           "hello\n",
				DetailsTagSummary: "Summary",
			},
			want: "<details>\n<summary>Summary</summary>\n\nhello\n\n</details>",
		},
		{
			name: "fenced code block with details tag",
			result: &TemplateInput{
				Content:           "hello\n",
				CodeBlock:         true,
				DetailsTagSummary: "Summary",
			},
			want: "<details>\n<summary>Summary</summary>\n\n```\nhello\n```\n\n</details>",
		},
		{
			name: "no details tag",
			result: &TemplateInput{
				Content: "hello",
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := renderFile(nil, tt.result)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("renderFile() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
