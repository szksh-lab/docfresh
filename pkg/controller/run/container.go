package run

import (
	"fmt"

	"github.com/suzuki-shunsuke/docfresh/pkg/container"
)

type (
	ContainerInput   = container.Input
	ContainerCommand = container.Command
	ContainerRef     = container.Ref
	ContainerState   = container.State
	ContainerEngine  = container.Engine
)

func newFileRunContext() *container.FileRunContext {
	return container.NewFileRunContext()
}

func validateContainerInput(input *ContainerInput, existing map[string]*ContainerState) error {
	if err := container.ValidateInput(input, existing); err != nil {
		return fmt.Errorf("validate container input: %w", err)
	}
	return nil
}

func newContainerEngine(engine string) (ContainerEngine, error) {
	e, err := container.NewEngine(engine)
	if err != nil {
		return nil, fmt.Errorf("create container engine: %w", err)
	}
	return e, nil
}
