package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
)

func (c *Controller) exec(ctx context.Context, logger *slog.Logger, file string, input *BlockInput, frc *fileRunContext) (*TemplateInput, error) {
	if input.Command != nil {
		if input.Command.Container != nil {
			return execContainerCommand(ctx, frc, input.Command)
		}
		return c.execCommand(ctx, logger, file, input.Command)
	}
	if input.File != nil {
		return readFile(file, input.File, c.fs, c.langs)
	}
	if input.HTTP != nil {
		return callHTTP(ctx, input.HTTP, c.httpClient, c.langs)
	}
	if input.GitHubContent != nil {
		return getGitHubContent(ctx, c.gh, c.langs, input.GitHubContent)
	}
	return nil, errors.New("no command or file specified")
}

func execContainerCommand(ctx context.Context, frc *fileRunContext, command *Command) (*TemplateInput, error) {
	ref := command.Container
	state, ok := frc.containers[ref.ID]
	if !ok {
		return nil, fmt.Errorf("container %q not found", ref.ID)
	}
	result, err := frc.engine.Exec(ctx, state.ContainerID, command.Command, command.Dir, command.Env)
	if err != nil {
		if !command.IgnoreFail {
			state.Failed = true
			return nil, fmt.Errorf("execute command in container %s: %w", ref.ID, err)
		}
	} else if result.ExitCode != 0 && !command.IgnoreFail {
		state.Failed = true
		return nil, fmt.Errorf("execute command in container %s: exit code %d", ref.ID, result.ExitCode)
	}
	result.CommandLanguage = command.CommandLanguage
	result.OutputLanguage = command.OutputLanguage
	result.EmbedScript = command.EmbedScript
	result.Quiet = command.Quiet
	return result, nil
}

func resolvePath(baseFile, file string) string {
	p := filepath.FromSlash(file)
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(filepath.Dir(baseFile), p)
}
