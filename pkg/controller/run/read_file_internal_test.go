package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/spf13/afero"
)

func TestReadFile(t *testing.T) { //nolint:funlen,cyclop
	t.Parallel()
	tests := []struct {
		name     string
		baseFile string
		file     *File
		langs    map[string]*Language
		setup    func(t *testing.T, fs afero.Fs)
		want     *TemplateInput
		wantErr  bool
	}{
		{
			name:     "basic success",
			baseFile: "/docs/README.md",
			file:     &File{Path: "hello.txt"},
			setup: func(t *testing.T, fs afero.Fs) {
				t.Helper()
				if err := afero.WriteFile(fs, "/docs/hello.txt", []byte("hello world"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: &TemplateInput{
				Type:    "local-file",
				Path:    "hello.txt",
				Content: "hello world",
			},
		},
		{
			name:     "language from extension",
			baseFile: "/docs/README.md",
			file:     &File{Path: "main.go"},
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			setup: func(t *testing.T, fs afero.Fs) {
				t.Helper()
				if err := afero.WriteFile(fs, "/docs/main.go", []byte("package main"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: &TemplateInput{
				Type:     "local-file",
				Path:     "main.go",
				Content:  "package main",
				Language: "go",
			},
		},
		{
			name:     "explicit language overrides extension",
			baseFile: "/docs/README.md",
			file:     &File{Path: "main.go", Language: "golang"},
			langs: map[string]*Language{
				".go": {Language: "go"},
			},
			setup: func(t *testing.T, fs afero.Fs) {
				t.Helper()
				if err := afero.WriteFile(fs, "/docs/main.go", []byte("package main"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: &TemplateInput{
				Type:     "local-file",
				Path:     "main.go",
				Content:  "package main",
				Language: "golang",
			},
		},
		{
			name:     "file not found",
			baseFile: "/docs/README.md",
			file:     &File{Path: "missing.txt"},
			setup:    func(_ *testing.T, _ afero.Fs) {},
			wantErr:  true,
		},
		{
			name:     "with range",
			baseFile: "/docs/README.md",
			file: &File{
				Path: "lines.txt",
				Range: &Range{
					Start: new(1),
					End:   new(2),
				},
			},
			setup: func(t *testing.T, fs afero.Fs) {
				t.Helper()
				if err := afero.WriteFile(fs, "/docs/lines.txt", []byte("line1\nline2\nline3\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: &TemplateInput{
				Type:    "local-file",
				Path:    "lines.txt",
				Content: "line2\n",
			},
		},
		{
			name:     "range error",
			baseFile: "/docs/README.md",
			file: &File{
				Path: "short.txt",
				Range: &Range{
					Start: new(5),
					End:   new(10),
				},
			},
			setup: func(t *testing.T, fs afero.Fs) {
				t.Helper()
				if err := afero.WriteFile(fs, "/docs/short.txt", []byte("one line\n"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			wantErr: true,
		},
		{
			name:     "relative path resolution",
			baseFile: "/dir/doc.md",
			file:     &File{Path: "sub/file.txt"},
			setup: func(t *testing.T, fs afero.Fs) {
				t.Helper()
				if err := afero.WriteFile(fs, "/dir/sub/file.txt", []byte("nested content"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: &TemplateInput{
				Type:    "local-file",
				Path:    "sub/file.txt",
				Content: "nested content",
			},
		},
		{
			name:     "absolute path",
			baseFile: "/dir/doc.md",
			file:     &File{Path: "/abs/file.txt"},
			setup: func(t *testing.T, fs afero.Fs) {
				t.Helper()
				if err := afero.WriteFile(fs, "/abs/file.txt", []byte("absolute content"), 0o644); err != nil {
					t.Fatal(err)
				}
			},
			want: &TemplateInput{
				Type:    "local-file",
				Path:    "/abs/file.txt",
				Content: "absolute content",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fs := afero.NewMemMapFs()
			tt.setup(t, fs)
			got, err := readFile(tt.baseFile, tt.file, fs, tt.langs)
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
