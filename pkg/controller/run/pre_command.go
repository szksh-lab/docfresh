package run

import (
	"context"
	"fmt"
)

func (c *Controller) runPreCommand(ctx context.Context, file string, block *Block) error {
	if block.Input.PreCommand == nil {
		return nil
	}
	result, err := c.execCommand(ctx, file, block.Input.PreCommand)
	if err != nil {
		return err
	}
	if block.Input.PreCommand.Test != "" {
		if err := testResult(c.stderr, block.Input.PreCommand.Test, result); err != nil {
			return fmt.Errorf("test the result of pre_command: %w", err)
		}
	}
	return nil
}
