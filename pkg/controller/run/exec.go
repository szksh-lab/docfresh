package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

const waitDelay = 1000 * time.Hour

func setCancel(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		return cmd.Process.Signal(os.Interrupt)
	}
	cmd.WaitDelay = waitDelay
}

func getCommandDir(file string, command *Command) string {
	if command.Dir == "" {
		return filepath.Dir(file)
	}
	if filepath.IsAbs(command.Dir) {
		return command.Dir
	}
	return filepath.Join(filepath.Dir(file), command.Dir)
}

func (c *Controller) execCommand(ctx context.Context, file string, command *Command) (*TemplateInput, error) {
	shell := command.Shell
	if shell == nil {
		shell = []string{"bash", "-c"}
	}
	cmd := exec.CommandContext(ctx, shell[0], append(shell[1:], command.Command)...) //nolint:gosec
	cmd.Dir = getCommandDir(file, command)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combinedOutput)
	setCancel(cmd)
	fmt.Fprintln(os.Stderr, "+", command.Command)
	if err := cmd.Run(); err != nil && !command.IgnoreFail {
		return nil, fmt.Errorf("execute a command: %w", err)
	}
	return &TemplateInput{
		Type:           "command",
		Command:        command.Command,
		Dir:            command.Dir,
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		ExitCode:       cmd.ProcessState.ExitCode(),
	}, nil
}

func (c *Controller) exec(ctx context.Context, file string, input *BlockInput) (*TemplateInput, error) {
	if input.Command != nil {
		return c.execCommand(ctx, file, input.Command)
	}
	if input.File != nil {
		return c.readFile(file, input.File)
	}
	if input.HTTP != nil {
		return c.request(ctx, input.HTTP)
	}
	if input.GitHubContent != nil {
		return c.getGitHubContent(ctx, input.GitHubContent)
	}
	return nil, errors.New("no command or file specified")
}

func (c *Controller) readFile(file string, f *File) (*TemplateInput, error) {
	p := filepath.FromSlash(f.Path)
	if !filepath.IsAbs(p) {
		p = filepath.Join(filepath.Dir(file), p)
	}
	b, err := afero.ReadFile(c.fs, p)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return &TemplateInput{
		Type:    "local-file",
		Path:    p,
		Content: string(b),
	}, nil
}

func (c *Controller) request(ctx context.Context, h *HTTP) (*TemplateInput, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.URL, nil)
	if err != nil {
		return nil, fmt.Errorf("create http request: %w", err)
	}
	resp, err := c.httpClient.Do(req) //nolint:gosec
	if err != nil {
		return nil, fmt.Errorf("send http request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}
	return &TemplateInput{
		Type:    "http",
		URL:     h.URL,
		Content: string(b),
	}, nil
}

func (c *Controller) getGitHubContent(ctx context.Context, content *GitHubContent) (*TemplateInput, error) {
	s, err := c.gh.GetContent(ctx, content.Owner, content.Repo, content.Path, content.Ref)
	if err != nil {
		return nil, fmt.Errorf("get a file by GitHub Content API: %w", err)
	}
	return &TemplateInput{
		Type:    "github-content",
		Content: s,
		Owner:   content.Owner,
		Repo:    content.Repo,
		Path:    content.Path,
		Ref:     content.Ref,
	}, nil
}
