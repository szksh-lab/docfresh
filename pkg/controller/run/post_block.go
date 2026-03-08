package run

import (
	"context"
	"fmt"
	"log/slog"
)

func (c *Controller) runPostBlock(ctx context.Context, logger *slog.Logger, _ *Templates, file string, block *Block) error {
	fmt.Fprintf(c.stderr, "> post_comment %s:%d\n", file, block.LineNumber)
	result, err := c.execCommand(ctx, logger, file, block.Input.Command)
	if err != nil {
		return fmt.Errorf("execute post block command: %w", err)
	}
	if block.Input.Command.Test == "" {
		return nil
	}
	if err := testResult(c.stderr, block.Input.Command.Test, result); err != nil {
		return fmt.Errorf("test the result of post block command: %w", err)
	}
	return nil
}
