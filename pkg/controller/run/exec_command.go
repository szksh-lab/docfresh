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

	"github.com/mattn/go-colorable"
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

func shellFromLanguageName(command *Command, langsByName map[string]*Language) []string {
	if command.CommandLanguage == "" {
		return nil
	}
	lang, ok := langsByName[command.CommandLanguage]
	if !ok {
		return nil
	}
	if command.Script != "" {
		return lang.ScriptShell
	}
	return lang.CommandShell
}

func getShell(command *Command, langs map[string]*Language, langsByName map[string]*Language) ([]string, error) {
	if len(command.Shell) > 0 {
		return command.Shell, nil
	}
	if shell := shellFromLanguageName(command, langsByName); len(shell) > 0 {
		return shell, nil
	}
	if command.Script == "" {
		return []string{"bash", "-c"}, nil
	}
	if !command.EmbedScript {
		return []string{"bash"}, nil
	}
	ext := filepath.Ext(command.Script)
	sl, ok := langs[ext]
	if ok && len(sl.ScriptShell) > 0 {
		return sl.ScriptShell, nil
	}
	return nil, errors.New("shell is required")
}

func (c *Controller) execCommand(ctx context.Context, logger *slog.Logger, file string, command *Command) (*TemplateInput, error) {
	shell, err := getShell(command, c.langs, c.langsByName)
	if err != nil {
		return nil, fmt.Errorf("get command.shell: %w", err)
	}
	script := command.Command
	var content string
	dir := getCommandDir(file, command)
	scriptLanguage := command.CommandLanguage
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
	uncolorizedStdout := colorable.NewNonColorable(stdout)
	uncolorizedStderr := colorable.NewNonColorable(stderr)
	uncolorizedCombinedOutput := colorable.NewNonColorable(combinedOutput)
	cmd.Stdout = io.MultiWriter(os.Stdout, uncolorizedStdout, uncolorizedCombinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, uncolorizedStderr, uncolorizedCombinedOutput)
	setCancel(logger, cmd, time.Duration(command.TimeoutSigkill)*time.Second)
	cmd.Env = getEnv(command.Env, c.environ)
	fmt.Fprintln(os.Stderr, "+", command.Command)
	if err := cmd.Run(); err != nil && !command.IgnoreFail {
		return nil, fmt.Errorf("execute a command: %w", err)
	}
	return &TemplateInput{
		Type:            "command",
		Shell:           shell,
		Command:         command.Command,
		Script:          command.Script,
		CommandLanguage: scriptLanguage,
		OutputLanguage:  command.OutputLanguage,
		EmbedScript:     command.EmbedScript,
		Dir:             command.Dir,
		Stdout:          stdout.String(),
		Stderr:          stderr.String(),
		CombinedOutput:  combinedOutput.String(),
		ExitCode:        cmd.ProcessState.ExitCode(),
		Content:         content,
		Quiet:           command.Quiet,
	}, nil
}

func getEnv(env map[string]string, osEnvs []string) []string {
	arr := osEnvs
	if len(env) == 0 {
		return nil
	}
	for k, v := range env {
		arr = append(arr, k+"="+v)
	}
	return arr
}
