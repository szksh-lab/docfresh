package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseContainerBlock(t *testing.T) {
	t.Parallel()
	content := `<!-- docfresh container
id: mycontainer
engine: docker-cli
image: ubuntu:latest
-->`
	block, pos, err := parseContainerBlock(content, 0, 0, nil)
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
	blocks, err := ParseFile(content, nil)
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
