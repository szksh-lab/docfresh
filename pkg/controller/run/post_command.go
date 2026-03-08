package run

import (
	"context"
	"fmt"
	"log/slog"
)

func (c *Controller) runPostCommand(ctx context.Context, logger *slog.Logger, file string, block *Block) error {
	fmt.Fprintf(c.stderr, "> post_command %s:%d\n", file, block.LineNumber)
	result, err := c.execCommand(ctx, logger, file, block.Input.PostCommand)
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
