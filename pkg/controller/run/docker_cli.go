package run

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type DockerCLIEngine struct{}

func (d *DockerCLIEngine) Create(ctx context.Context, input *ContainerInput, file string) (string, error) {
	args := buildDockerCreateArgs(input, file)
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stderr = os.Stderr
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("docker run: %w", err)
	}
	containerID := strings.TrimSpace(string(out))
	return containerID, nil
}

func buildDockerCreateArgs(input *ContainerInput, file string) []string {
	args := []string{"run", "-d", "--entrypoint="}
	if input.Workspace != "" {
		args = append(args, "--workdir="+input.Workspace)
	}
	for _, v := range input.Volumes {
		args = append(args, "-v", v)
	}
	for k, v := range input.Env {
		args = append(args, "-e", k+"="+v)
	}
	args = append(args, "--label=docfresh.file_path="+file)
	absPath, err := filepath.Abs(file)
	if err == nil {
		args = append(args, "--label=docfresh.absolute_file_path="+absPath)
	}
	if input.ID != "" {
		args = append(args, "--label=docfresh.id="+input.ID)
	}
	args = append(args, input.Image, "tail", "-f", "/dev/null")
	return args
}

func (d *DockerCLIEngine) CopyFiles(ctx context.Context, containerID string, files map[string]string) error {
	for dest, src := range files {
		cmd := exec.CommandContext(ctx, "docker", "cp", src, containerID+":"+dest)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("docker cp %s -> %s: %w", src, dest, err)
		}
	}
	return nil
}

func (d *DockerCLIEngine) Exec(ctx context.Context, containerID string, command string, dir string, env map[string]string) (*TemplateInput, error) {
	args := buildDockerExecArgs(containerID, command, dir, env)
	cmd := exec.CommandContext(ctx, "docker", args...)
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	combinedOutput := &bytes.Buffer{}
	cmd.Stdout = io.MultiWriter(os.Stdout, stdout, combinedOutput)
	cmd.Stderr = io.MultiWriter(os.Stderr, stderr, combinedOutput)
	fmt.Fprintln(os.Stderr, "+", command)
	exitCode := 0
	if err := cmd.Run(); err != nil {
		var exitErr *exec.ExitError
		if ok := errors.As(err, &exitErr); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("docker exec: %w", err)
		}
	}
	if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}
	return &TemplateInput{
		Type:           "command",
		Command:        command,
		Stdout:         stdout.String(),
		Stderr:         stderr.String(),
		CombinedOutput: combinedOutput.String(),
		ExitCode:       exitCode,
	}, nil
}

func buildDockerExecArgs(containerID, command, dir string, env map[string]string) []string {
	args := []string{"exec"}
	if dir != "" {
		args = append(args, "-w", dir)
	}
	for k, v := range env {
		args = append(args, "-e", k+"="+v)
	}
	args = append(args, containerID, "bash", "-c", command)
	return args
}

func (d *DockerCLIEngine) Name(ctx context.Context, containerID string) (string, error) {
	cmd := exec.CommandContext(ctx, "docker", "inspect", "--format", "{{.Name}}", containerID)
	out, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("docker inspect %s: %w", containerID, err)
	}
	name := strings.TrimSpace(string(out))
	return strings.TrimPrefix(name, "/"), nil
}

func (d *DockerCLIEngine) Remove(ctx context.Context, containerID string) error {
	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", containerID)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker rm -f %s: %w", containerID, err)
	}
	return nil
}
