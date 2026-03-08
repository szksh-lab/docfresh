package run

import (
	"context"
	"errors"
	"log/slog"
	"path/filepath"
)

func (c *Controller) exec(ctx context.Context, logger *slog.Logger, file string, input *BlockInput) (*TemplateInput, error) {
	if input.Command != nil {
		return c.execCommand(ctx, logger, file, input.Command)
	}
	if input.File != nil {
		return c.readFile(file, input.File)
	}
	if input.HTTP != nil {
		return c.request(ctx, input.HTTP)
	}
	if input.GitHubContent != nil {
		return getGitHubContent(ctx, c.gh, c.langs, input.GitHubContent)
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
