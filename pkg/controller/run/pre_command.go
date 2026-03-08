package run

import (
	"context"
	"fmt"
	"log/slog"
)

func (c *Controller) runPreCommand(ctx context.Context, logger *slog.Logger, file string, block *Block) error {
	if block.Input.PreCommand == nil {
		return nil
	}
	fmt.Fprintf(c.stderr, "> pre_command %s:%d\n", file, block.LineNumber)
	result, err := c.execCommand(ctx, logger, file, block.Input.PreCommand)
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
