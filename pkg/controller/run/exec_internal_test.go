package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestResolvePath(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		baseFile string
		file     string
		want     string
	}{
		{
			name:     "absolute path is returned as-is",
			baseFile: "/home/user/docs/README.md",
			file:     "/tmp/file.txt",
			want:     "/tmp/file.txt",
		},
		{
			name:     "relative path is joined with base file directory",
			baseFile: "/home/user/docs/README.md",
			file:     "sub/file.txt",
			want:     "/home/user/docs/sub/file.txt",
		},
		{
			name:     "file in same directory",
			baseFile: "/home/user/docs/README.md",
			file:     "other.md",
			want:     "/home/user/docs/other.md",
		},
		{
			name:     "forward slash path on any OS",
			baseFile: "/home/user/docs/README.md",
			file:     "a/b/c.txt",
			want:     "/home/user/docs/a/b/c.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := resolvePath(tt.baseFile, tt.file)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
