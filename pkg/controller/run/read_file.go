package run

import (
	"fmt"

	"github.com/spf13/afero"
)

func (c *Controller) readFile(baseFile, file string) (*TemplateInput, error) {
	p := resolvePath(baseFile, file)
	b, err := afero.ReadFile(c.fs, p)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return &TemplateInput{
		Type:    "local-file",
		Path:    file,
		Content: string(b),
	}, nil
}
