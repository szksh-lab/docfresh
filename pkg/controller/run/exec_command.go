package run

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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

func getShell(command *Command) []string {
	if len(command.Shell) > 0 {
		return command.Shell
	}
	if command.Script != "" {
		return []string{"bash"}
	}
	return []string{"bash", "-c"}
}

func (c *Controller) execCommand(ctx context.Context, file string, command *Command) (*TemplateInput, error) {
	shell := getShell(command)
	script := command.Script
	if command.Command != "" {
		script = command.Command
	}
	cmd := exec.CommandContext(ctx, shell[0], append(shell[1:], script)...) //nolint:gosec
	cmd.Dir = getCommandDir(file, command)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combinedOutput)
	setCancel(cmd)
	if len(command.Envs) > 0 {
		envs := os.Environ()
		for k, v := range command.Envs {
			envs = append(envs, k+"="+v)
		}
		cmd.Env = envs
	}
	fmt.Fprintln(os.Stderr, "+", command.Command)
	if err := cmd.Run(); err != nil && !command.IgnoreFail {
		return nil, fmt.Errorf("execute a command: %w", err)
	}
	return &TemplateInput{
		Type:           "command",
		Shell:          shell,
		Command:        command.Command,
		Script:         command.Script,
		Dir:            command.Dir,
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		ExitCode:       cmd.ProcessState.ExitCode(),
	}, nil
}
