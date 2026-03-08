package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
)

const defaultWaitDelay = 1000 * time.Hour

func setCancel(logger *slog.Logger, cmd *exec.Cmd, waitDelay time.Duration) {
	cmd.Cancel = func() error {
		logger.Warn("SIGINT is sent to cancel the command")
		return cmd.Process.Signal(os.Interrupt)
	}
	if waitDelay == 0 {
		waitDelay = defaultWaitDelay
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

func getShell(command *Command, langs map[string]*Language) ([]string, error) {
	if len(command.Shell) > 0 {
		return command.Shell, nil
	}
	if command.Script == "" {
		return []string{"bash", "-c"}, nil
	}
	if !command.EmbedScript {
		return []string{"bash"}, nil
	}
	ext := filepath.Ext(command.Script)
	sl, ok := langs[ext]
	if ok && sl.Shell != nil {
		return sl.Shell, nil
	}
	return nil, errors.New("shell is required")
}

func (c *Controller) execCommand(ctx context.Context, logger *slog.Logger, file string, command *Command) (*TemplateInput, error) {
	shell, err := getShell(command, c.langs)
	if err != nil {
		return nil, fmt.Errorf("get command.shell: %w", err)
	}
	script := command.Command
	var content string
	dir := getCommandDir(file, command)
	scriptLanguage := command.Language
	if command.Script != "" {
		script = command.Script
		if scriptLanguage == "" {
			sl, ok := c.langs[filepath.Ext(command.Script)]
			if ok {
				scriptLanguage = sl.Language
			}
		}
		b, err := afero.ReadFile(c.fs, filepath.Join(dir, command.Script))
		if err != nil {
			return nil, fmt.Errorf("read a command.script: %w", err)
		}
		content = string(b)
	}
	if command.Timeout > 0 {
		requestCtx, cancel := context.WithTimeout(ctx, time.Duration(command.Timeout)*time.Second)
		defer cancel()
		ctx = requestCtx
	}
	cmd := exec.CommandContext(ctx, shell[0], append(shell[1:], script)...) //nolint:gosec
	cmd.Dir = dir
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combinedOutput)
	setCancel(logger, cmd, time.Duration(command.TimeoutSigkill)*time.Second)
	cmd.Env = getEnvs(command.Envs, c.environ)
	fmt.Fprintln(os.Stderr, "+", command.Command)
	if err := cmd.Run(); err != nil && !command.IgnoreFail {
		return nil, fmt.Errorf("execute a command: %w", err)
	}
	return &TemplateInput{
		Type:           "command",
		Shell:          shell,
		Command:        command.Command,
		Script:         command.Script,
		Language:       scriptLanguage,
		EmbedScript:    command.EmbedScript,
		Dir:            command.Dir,
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		ExitCode:       cmd.ProcessState.ExitCode(),
		Content:        content,
		Quiet:          command.Quiet,
	}, nil
}

func getEnvs(envs map[string]string, osEnvs []string) []string {
	arr := osEnvs
	if len(envs) == 0 {
		return nil
	}
	for k, v := range envs {
		arr = append(arr, k+"="+v)
	}
	return arr
}
