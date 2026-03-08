package run

import (
	"context"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type mockGitHub struct {
	content string
	err     error
}

func (m *mockGitHub) GetContent(_ context.Context, _, _, _, _ string) (string, error) {
	return m.content, m.err
}

func TestGetGitHubContent(t *testing.T) { //nolint:funlen
	t.Parallel()
	tests := []struct {
		name    string
		gh      *mockGitHub
		langs   map[string]*Language
		content *GitHubContent
		want    *TemplateInput
		wantErr bool
	}{
		{
			name: "basic success",
			gh: &mockGitHub{
				content: "hello world",
			},
			langs: nil,
			content: &GitHubContent{
				Owner: "owner",
				Repo:  "repo",
				Path:  "README.md",
			},
			want: &TemplateInput{
				Type:    "github-content",
				Content: "hello world",
				Owner:   "owner",
				Repo:    "repo",
				Path:    "README.md",
			},
		},
		{
			name: "with ref",
			gh: &mockGitHub{
				content: "content at ref",
			},
			langs: nil,
			content: &GitHubContent{
				Owner: "owner",
				Repo:  "repo",
				Path:  "file.txt",
				Ref:   "v1.0.0",
			},
			want: &TemplateInput{
				Type:    "github-content",
				Content: "content at ref",
				Owner:   "owner",
				Repo:    "repo",
				Path:    "file.txt",
				Ref:     "v1.0.0",
			},
		},
		{
			name: "language auto-detection from extension",
			gh: &mockGitHub{
				content: "package main",
			},
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			content: &GitHubContent{
				Owner: "owner",
				Repo:  "repo",
				Path:  "main.go",
			},
			want: &TemplateInput{
				Type:     "github-content",
				Content:  "package main",
				Owner:    "owner",
				Repo:     "repo",
				Path:     "main.go",
				Language: "go",
			},
		},
		{
			name: "explicit language overrides auto-detection",
			gh: &mockGitHub{
				content: "package main",
			},
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			content: &GitHubContent{
				Owner:    "owner",
				Repo:     "repo",
				Path:     "main.go",
				Language: "golang",
			},
			want: &TemplateInput{
				Type:     "github-content",
				Content:  "package main",
				Owner:    "owner",
				Repo:     "repo",
				Path:     "main.go",
				Language: "golang",
			},
		},
		{
			name: "with range",
			gh: &mockGitHub{
				content: "line1\nline2\nline3\n",
			},
			langs: nil,
			content: &GitHubContent{
				Owner: "owner",
				Repo:  "repo",
				Path:  "file.txt",
				Range: &Range{
					Start: new(1),
					End:   new(2),
				},
			},
			want: &TemplateInput{
				Type:    "github-content",
				Content: "line2\n",
				Owner:   "owner",
				Repo:    "repo",
				Path:    "file.txt",
			},
		},
		{
			name: "API error",
			gh: &mockGitHub{
				err: errors.New("not found"),
			},
			langs: nil,
			content: &GitHubContent{
				Owner: "owner",
				Repo:  "repo",
				Path:  "missing.txt",
			},
			wantErr: true,
		},
		{
			name: "range error",
			gh: &mockGitHub{
				content: "line1\n",
			},
			langs: nil,
			content: &GitHubContent{
				Owner: "owner",
				Repo:  "repo",
				Path:  "file.txt",
				Range: &Range{
					Start: new(5),
					End:   new(10),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getGitHubContent(context.Background(), tt.gh, tt.langs, tt.content)
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
