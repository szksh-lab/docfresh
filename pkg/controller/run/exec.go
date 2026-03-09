package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/suzuki-shunsuke/docfresh/pkg/container"
)

func (c *Controller) exec(ctx context.Context, logger *slog.Logger, file string, input *BlockInput, frc *container.FileRunContext) (*TemplateInput, error) {
	if input.Command != nil {
		if input.Command.Container != nil {
			return execContainerCommand(ctx, frc, input.Command)
		}
		return c.execCommand(ctx, logger, file, input.Command)
	}
	if input.File != nil {
		if input.File.Container != nil {
			return execContainerFile(ctx, frc, input.File, c.langs)
		}
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

func execContainerCommand(ctx context.Context, frc *container.FileRunContext, command *Command) (*TemplateInput, error) {
	ref := command.Container
	state, ok := frc.Containers[ref.ID]
	if !ok {
		return nil, fmt.Errorf("container %q not found", ref.ID)
	}
	result, err := frc.Engine.Exec(ctx, state.ContainerID, command.Command, command.Dir, command.Env)
	if err != nil {
		if !command.IgnoreFail {
			state.Failed = true
			return nil, fmt.Errorf("execute command in container %s: %w", ref.ID, err)
		}
	} else if result.ExitCode != 0 && !command.IgnoreFail {
		state.Failed = true
		return nil, fmt.Errorf("execute command in container %s: exit code %d", ref.ID, result.ExitCode)
	}
	ti := execResultToTemplateInput(result)
	ti.CommandLanguage = command.CommandLanguage
	ti.OutputLanguage = command.OutputLanguage
	ti.EmbedScript = command.EmbedScript
	ti.HideOutput = command.HideOutput
	ti.HideCommand = command.HideCommand
	return ti, nil
}

func execResultToTemplateInput(r *container.ExecResult) *TemplateInput {
	return &TemplateInput{
		Type:           "command",
		Command:        r.Command,
		Stdout:         r.Stdout,
		Stderr:         r.Stderr,
		CombinedOutput: r.CombinedOutput,
		ExitCode:       r.ExitCode,
	}
}

func execContainerFile(ctx context.Context, frc *container.FileRunContext, file *File, langs map[string]*Language) (*TemplateInput, error) {
	ref := file.Container
	state, ok := frc.Containers[ref.ID]
	if !ok {
		return nil, fmt.Errorf("container %q not found", ref.ID)
	}
	workspace := state.Input.Workspace
	b, err := frc.Engine.ReadFile(ctx, state.ContainerID, file.Path, workspace)
	if err != nil {
		return nil, fmt.Errorf("read file from container %s: %w", ref.ID, err)
	}
	sl := file.Language
	if sl == "" {
		ext := filepath.Ext(file.Path)
		if lang, ok := langs[ext]; ok {
			sl = lang.Language
		}
	}
	content := string(b)
	content, err = extractRange(content, file.Range)
	if err != nil {
		return nil, fmt.Errorf("extract range from file: %w", err)
	}
	result := &TemplateInput{
		Type:     "container-file",
		Path:     file.Path,
		Language: sl,
		Content:  content,
		Vars:     file.Template.GetVars(),
	}
	if file.Template != nil {
		if err := renderTemplate(content, result, file.Template.Delims); err != nil {
			return nil, err
		}
	}
	return result, nil
}

func resolvePath(baseFile, file string) string {
	p := filepath.FromSlash(file)
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(filepath.Dir(baseFile), p)
}
