package run

import (
	"context"
	"errors"
	"path/filepath"
)

func (c *Controller) exec(ctx context.Context, file string, input *BlockInput) (*TemplateInput, error) {
	if input.Command != nil {
		return c.execCommand(ctx, file, input.Command)
	}
	if input.File != nil {
		return c.readFile(file, input.File.Path)
	}
	if input.HTTP != nil {
		return c.request(ctx, input.HTTP)
	}
	if input.GitHubContent != nil {
		return c.getGitHubContent(ctx, input.GitHubContent)
	}
	return nil, errors.New("no command or file specified")
}

func resolvePath(baseFile, file string) string {
	p := filepath.FromSlash(file)
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(filepath.Dir(baseFile), p)
}
