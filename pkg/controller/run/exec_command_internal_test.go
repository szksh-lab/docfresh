package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGetCommandDir(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		file    string
		command *Command
		want    string
	}{
		{
			name:    "empty dir uses file directory",
			file:    "/home/user/docs/README.md",
			command: &Command{},
			want:    "/home/user/docs",
		},
		{
			name:    "absolute dir is returned as-is",
			file:    "/home/user/docs/README.md",
			command: &Command{Dir: "/tmp/work"},
			want:    "/tmp/work",
		},
		{
			name:    "relative dir is joined with file directory",
			file:    "/home/user/docs/README.md",
			command: &Command{Dir: "sub/dir"},
			want:    "/home/user/docs/sub/dir",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getCommandDir(tt.file, tt.command)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetShell(t *testing.T) { //nolint:funlen
	t.Parallel()
	langs := map[string]*Language{
		".py": {ScriptShell: []string{"python3"}, Language: "py"},
		".go": {ScriptShell: []string{"go", "run"}, Language: "go"},
	}
	langsByName := map[string]*Language{
		"js": {ScriptShell: []string{"node"}, CommandShell: []string{"node", "-e"}, Language: "js"},
		"py": {ScriptShell: []string{"python3"}, CommandShell: []string{"python3", "-c"}, Language: "py"},
		"go": {ScriptShell: []string{"go", "run"}, Language: "go"},
	}
	tests := []struct {
		name    string
		command *Command
		want    []string
		wantErr bool
	}{
		{
			name:    "explicit shell is returned",
			command: &Command{Shell: []string{"zsh", "-c"}},
			want:    []string{"zsh", "-c"},
		},
		{
			name:    "no script uses bash -c",
			command: &Command{},
			want:    []string{"bash", "-c"},
		},
		{
			name:    "script without embed uses bash",
			command: &Command{Script: "run.sh"},
			want:    []string{"bash"},
		},
		{
			name:    "embed script with known extension",
			command: &Command{Script: "run.py", EmbedScript: true},
			want:    []string{"python3"},
		},
		{
			name:    "embed script with unknown extension",
			command: &Command{Script: "run.rb", EmbedScript: true},
			wantErr: true,
		},
		{
			name:    "embed script with extension that has nil shell",
			command: &Command{Script: "run.unknown", EmbedScript: true},
			wantErr: true,
		},
		{
			name:    "command_language auto-detects command shell",
			command: &Command{Command: "console.log('hello')", CommandLanguage: "js"},
			want:    []string{"node", "-e"},
		},
		{
			name:    "command_language auto-detects script shell",
			command: &Command{Script: "hello.js", CommandLanguage: "js"},
			want:    []string{"node"},
		},
		{
			name:    "command_language with no command shell falls back to bash -c",
			command: &Command{Command: "main()", CommandLanguage: "go"},
			want:    []string{"bash", "-c"},
		},
		{
			name:    "command_language with no script shell falls back to bash",
			command: &Command{Script: "run.txt", CommandLanguage: "json"},
			want:    []string{"bash"},
		},
		{
			name:    "explicit shell overrides command_language",
			command: &Command{Shell: []string{"zsh", "-c"}, CommandLanguage: "js"},
			want:    []string{"zsh", "-c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := getShell(tt.command, langs, langsByName)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		env    map[string]string
		osEnvs []string
		want   []string
	}{
		{
			name:   "nil env returns nil",
			env:    nil,
			osEnvs: []string{"PATH=/usr/bin"},
			want:   nil,
		},
		{
			name:   "empty env returns nil",
			env:    map[string]string{},
			osEnvs: []string{"PATH=/usr/bin"},
			want:   nil,
		},
		{
			name:   "env vars are appended to osEnvs",
			env:    map[string]string{"FOO": "bar"},
			osEnvs: []string{"PATH=/usr/bin"},
			want:   []string{"PATH=/usr/bin", "FOO=bar"},
		},
		{
			name:   "nil osEnvs with env",
			env:    map[string]string{"KEY": "val"},
			osEnvs: nil,
			want:   []string{"KEY=val"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := getEnv(tt.env, tt.osEnvs)
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
