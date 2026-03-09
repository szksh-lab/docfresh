package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestValidateContainerInput(t *testing.T) { //nolint:funlen
	t.Parallel()
	tests := []struct {
		name     string
		input    *ContainerInput
		existing map[string]*ContainerState
		wantErr  bool
	}{
		{
			name: "valid input",
			input: &ContainerInput{
				ID:     "test",
				Engine: "docker-cli",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*ContainerState{},
			wantErr:  false,
		},
		{
			name: "missing id",
			input: &ContainerInput{
				Engine: "docker-cli",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*ContainerState{},
			wantErr:  true,
		},
		{
			name: "missing engine",
			input: &ContainerInput{
				ID:    "test",
				Image: "ubuntu:latest",
			},
			existing: map[string]*ContainerState{},
			wantErr:  true,
		},
		{
			name: "unsupported engine",
			input: &ContainerInput{
				ID:     "test",
				Engine: "podman",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*ContainerState{},
			wantErr:  true,
		},
		{
			name: "missing image",
			input: &ContainerInput{
				ID:     "test",
				Engine: "docker-cli",
			},
			existing: map[string]*ContainerState{},
			wantErr:  true,
		},
		{
			name: "duplicate id",
			input: &ContainerInput{
				ID:     "test",
				Engine: "docker-cli",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*ContainerState{
				"test": {},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateContainerInput(tt.input, tt.existing)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateContainerInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewContainerEngine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		engine  string
		wantErr bool
	}{
		{
			name:    "docker-cli",
			engine:  "docker-cli",
			wantErr: false,
		},
		{
			name:    "unsupported",
			engine:  "podman",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e, err := newContainerEngine(tt.engine)
			if (err != nil) != tt.wantErr {
				t.Errorf("newContainerEngine() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && e == nil {
				t.Error("newContainerEngine() returned nil engine")
			}
		})
	}
}

func TestParseContainerBlock(t *testing.T) {
	t.Parallel()
	content := `<!-- docfresh container
id: mycontainer
engine: docker-cli
image: ubuntu:latest
-->`
	block, pos, err := parseContainerBlock(content, 0, 0)
	if err != nil {
		t.Fatalf("parseContainerBlock() error = %v", err)
	}
	if pos != len(content) {
		t.Errorf("parseContainerBlock() pos = %d, want %d", pos, len(content))
	}
	if block.Type != "container" {
		t.Errorf("block.Type = %q, want %q", block.Type, "container")
	}
	if block.ContainerInput == nil {
		t.Fatal("block.ContainerInput is nil")
	}
	if diff := cmp.Diff("mycontainer", block.ContainerInput.ID); diff != "" {
		t.Errorf("ContainerInput.ID mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff("docker-cli", block.ContainerInput.Engine); diff != "" {
		t.Errorf("ContainerInput.Engine mismatch (-want +got):\n%s", diff)
	}
	if diff := cmp.Diff("ubuntu:latest", block.ContainerInput.Image); diff != "" {
		t.Errorf("ContainerInput.Image mismatch (-want +got):\n%s", diff)
	}
}

func TestParseFileWithContainerBlock(t *testing.T) {
	t.Parallel()
	content := `Some text
<!-- docfresh container
id: mycontainer
engine: docker-cli
image: ubuntu:latest
-->
More text
`
	blocks, err := parseFile(content)
	if err != nil {
		t.Fatalf("parseFile() error = %v", err)
	}
	if len(blocks) != 3 {
		t.Fatalf("parseFile() returned %d blocks, want 3", len(blocks))
	}
	if blocks[0].Type != "text" {
		t.Errorf("blocks[0].Type = %q, want %q", blocks[0].Type, "text")
	}
	if blocks[1].Type != "container" {
		t.Errorf("blocks[1].Type = %q, want %q", blocks[1].Type, "container")
	}
	if blocks[2].Type != "text" {
		t.Errorf("blocks[2].Type = %q, want %q", blocks[2].Type, "text")
	}
}
