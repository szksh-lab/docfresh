package run

import (
	"context"
	"fmt"
)

func (c *Controller) runPostCommand(ctx context.Context, file string, block *Block) error {
	result, err := c.execCommand(ctx, file, block.Input.PostCommand)
	if err != nil {
		return fmt.Errorf("execute post command: %w", err)
	}
	if block.Input.PostCommand.Test == "" {
		return nil
	}
	if err := testResult(c.stderr, block.Input.PostCommand.Test, result); err != nil {
		return fmt.Errorf("test the result of post command: %w", err)
	}
	return nil
}
