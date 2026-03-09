package container

import (
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildDockerCreateArgs(t *testing.T) {
	t.Parallel()
	absPath, err := filepath.Abs("README.md")
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name  string
		input *Input
		file  string
		want  []string
	}{
		{
			name: "basic",
			input: &Input{
				ID:    "test",
				Image: "ubuntu:latest",
			},
			file: "README.md",
			want: []string{"run", "-d", "--entrypoint=", "--label=docfresh.file_path=README.md", "--label=docfresh.absolute_file_path=" + absPath, "--label=docfresh.id=test", "ubuntu:latest", "tail", "-f", "/dev/null"},
		},
		{
			name: "with workspace",
			input: &Input{
				ID:        "mycontainer",
				Image:     "ubuntu:latest",
				Workspace: "/app",
			},
			file: "README.md",
			want: []string{"run", "-d", "--entrypoint=", "--workdir=/app", "--label=docfresh.file_path=README.md", "--label=docfresh.absolute_file_path=" + absPath, "--label=docfresh.id=mycontainer", "ubuntu:latest", "tail", "-f", "/dev/null"},
		},
		{
			name: "with volumes",
			input: &Input{
				ID:      "vol",
				Image:   "ubuntu:latest",
				Volumes: []string{"/host:/container", "/data:/data:ro"},
			},
			file: "README.md",
			want: []string{"run", "-d", "--entrypoint=", "-v", "/host:/container", "-v", "/data:/data:ro", "--label=docfresh.file_path=README.md", "--label=docfresh.absolute_file_path=" + absPath, "--label=docfresh.id=vol", "ubuntu:latest", "tail", "-f", "/dev/null"},
		},
		{
			name: "with absolute file path",
			input: &Input{
				ID:    "abs",
				Image: "alpine:latest",
			},
			file: "/tmp/docs/README.md",
			want: []string{"run", "-d", "--entrypoint=", "--label=docfresh.file_path=/tmp/docs/README.md", "--label=docfresh.absolute_file_path=/tmp/docs/README.md", "--label=docfresh.id=abs", "alpine:latest", "tail", "-f", "/dev/null"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := BuildDockerCreateArgs(tt.input, tt.file)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("BuildDockerCreateArgs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestBuildDockerExecArgs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		containerID string
		command     string
		dir         string
		env         map[string]string
		want        []string
	}{
		{
			name:        "basic",
			containerID: "abc123",
			command:     "echo hello",
			want:        []string{"exec", "abc123", "bash", "-c", "echo hello"},
		},
		{
			name:        "with dir",
			containerID: "abc123",
			command:     "ls",
			dir:         "/app",
			want:        []string{"exec", "-w", "/app", "abc123", "bash", "-c", "ls"},
		},
		{
			name:        "with env",
			containerID: "abc123",
			command:     "env",
			env:         map[string]string{"FOO": "bar"},
			want:        []string{"exec", "-e", "FOO=bar", "abc123", "bash", "-c", "env"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := BuildDockerExecArgs(tt.containerID, tt.command, tt.dir, tt.env)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("BuildDockerExecArgs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
