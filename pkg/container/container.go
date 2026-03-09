package container

import (
	"context"
	"errors"
	"fmt"
)

type Input struct {
	ID        string            `json:"id" yaml:"id" jsonschema_description:"Unique identifier for the container within the file"`
	Engine    string            `json:"engine" yaml:"engine" jsonschema_description:"Container engine to use. Currently only 'docker-cli' is supported"`
	Image     string            `json:"image" yaml:"image" jsonschema_description:"Container image to use"`
	Keep      bool              `json:"keep,omitempty" yaml:"keep" jsonschema_description:"If true, the container is not removed after processing"`
	Workspace string            `json:"workspace,omitempty" yaml:"workspace" jsonschema_description:"Working directory inside the container"`
	Volumes   []string          `json:"volumes,omitempty" yaml:"volumes" jsonschema_description:"Volume mounts in docker format (host:container)"`
	Env       map[string]string `json:"env,omitempty" yaml:"env" jsonschema_description:"Environment variables for the container"`
	CopyFiles map[string]string `json:"copy_files,omitempty" yaml:"copy_files" jsonschema_description:"Files to copy into the container (container_path: host_path)"`
	Command   *Command          `json:"command,omitempty" yaml:"command" jsonschema_description:"Command to run after container creation for setup"`
}

type Command struct {
	Command string `json:"command,omitempty" yaml:"command" jsonschema_description:"Shell command to execute"`
	Script  string `json:"script,omitempty" yaml:"script" jsonschema_description:"Script file path to execute"`
	Dir     string `json:"dir,omitempty" yaml:"dir" jsonschema_description:"Working directory inside the container"`
}

type Ref struct {
	ID string `json:"id" yaml:"id" jsonschema_description:"ID of the container to execute the command in"`
}

type State struct {
	Input       *Input
	ContainerID string // real docker container ID
	Failed      bool
}

type ExecResult struct {
	Command        string
	Stdout         string
	Stderr         string
	CombinedOutput string
	ExitCode       int
}

type Engine interface {
	Create(ctx context.Context, input *Input, file string) (containerID string, err error)
	CopyFiles(ctx context.Context, containerID string, files map[string]string) error
	Exec(ctx context.Context, containerID string, command string, dir string, env map[string]string) (*ExecResult, error)
	ReadFile(ctx context.Context, containerID string, path string, dir string) ([]byte, error)
	Remove(ctx context.Context, containerID string) error
	Name(ctx context.Context, containerID string) (string, error)
}

type FileRunContext struct {
	Containers map[string]*State
	Engine     Engine
}

func NewFileRunContext() *FileRunContext {
	return &FileRunContext{
		Containers: make(map[string]*State),
	}
}

func ValidateInput(input *Input, existing map[string]*State) error {
	if input.ID == "" {
		return errors.New("container id is required")
	}
	if input.Engine == "" {
		return errors.New("container engine is required")
	}
	if input.Engine != "docker-cli" {
		return fmt.Errorf("unsupported container engine: %s (only 'docker-cli' is supported)", input.Engine)
	}
	if input.Image == "" {
		return errors.New("container image is required")
	}
	if _, ok := existing[input.ID]; ok {
		return fmt.Errorf("duplicate container id: %s", input.ID)
	}
	return nil
}

func NewEngine(engine string) (Engine, error) {
	switch engine {
	case "docker-cli":
		return &DockerCLIEngine{}, nil
	default:
		return nil, fmt.Errorf("unsupported container engine: %s", engine)
	}
}
