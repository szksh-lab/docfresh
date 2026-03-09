package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBuildDockerCreateArgs(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		input *ContainerInput
		want  []string
	}{
		{
			name: "basic",
			input: &ContainerInput{
				Image: "ubuntu:latest",
			},
			want: []string{"run", "-d", "--entrypoint=", "ubuntu:latest", "tail", "-f", "/dev/null"},
		},
		{
			name: "with workspace",
			input: &ContainerInput{
				Image:     "ubuntu:latest",
				Workspace: "/app",
			},
			want: []string{"run", "-d", "--entrypoint=", "--workdir=/app", "ubuntu:latest", "tail", "-f", "/dev/null"},
		},
		{
			name: "with volumes",
			input: &ContainerInput{
				Image:   "ubuntu:latest",
				Volumes: []string{"/host:/container", "/data:/data:ro"},
			},
			want: []string{"run", "-d", "--entrypoint=", "-v", "/host:/container", "-v", "/data:/data:ro", "ubuntu:latest", "tail", "-f", "/dev/null"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := buildDockerCreateArgs(tt.input)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("buildDockerCreateArgs() mismatch (-want +got):\n%s", diff)
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
			got := buildDockerExecArgs(tt.containerID, tt.command, tt.dir, tt.env)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("buildDockerExecArgs() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
