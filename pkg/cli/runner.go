// Package cli provides the command-line interface layer for docfresh.
// This package serves as the main entry point for all CLI operations,
// handling command parsing, flag processing, and routing to appropriate subcommands.
// It orchestrates the overall CLI structure using urfave/cli framework and delegates
// actual business logic to controller packages.
package cli

import (
	"context"

	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/suzuki-shunsuke/urfave-cli-v3-util/urfave"
	"github.com/urfave/cli/v3"
)

func Run(ctx context.Context, logger *slogutil.Logger, env *urfave.Env) error {
	gFlags := &Flags{}
	return urfave.Command(env, &cli.Command{ //nolint:wrapcheck
		Name:  "docfresh",
		Usage: "Make document maintainable, reusable, and testable. https://github.com/suzuki-shunsuke/docfresh",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Usage:       "Log level (debug, info, warn, error)",
				Sources:     cli.EnvVars("DOCFRESH_LOG_LEVEL"),
				Destination: &gFlags.LogLevel,
				Local:       true,
			},
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "configuration file path",
				Sources:     cli.EnvVars("DOCFRESH_CONFIG"),
				Destination: &gFlags.Config,
				Local:       true,
			},
		},
		Commands: []*cli.Command{
			// NewInit(logger, gFlags),
			NewRun(logger, gFlags),
		},
	}).Run(ctx, env.Args)
}
