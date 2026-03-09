package cli

import (
	"context"
	"fmt"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/docfresh/pkg/controller/run"
	"github.com/suzuki-shunsuke/docfresh/pkg/github"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

// RunArgs holds the flag and argument values for the init command.
type RunArgs struct {
	*Flags

	Files             []string // positional argument
	AllowUnknownField bool
}

// NewRun creates a new init command instance with the provided logger.
// It returns a CLI command that can be registered with the main CLI application.
func NewRun(logger *slogutil.Logger, gFlags *Flags) *cli.Command {
	args := &RunArgs{
		Flags: gFlags,
	}
	return &cli.Command{
		Name:  "run",
		Usage: "Update documents",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return runAction(ctx, logger, args)
		},
		Arguments: []cli.Argument{
			&cli.StringArgs{
				Name:        "document file",
				Destination: &args.Files,
				Max:         -1,
			},
		},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:        "allow-unknown-field",
				Usage:       "Allow unknown fields in directive YAML",
				Destination: &args.AllowUnknownField,
			},
		},
	}
}

func runAction(ctx context.Context, logger *slogutil.Logger, args *RunArgs) error {
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	fs := afero.NewOsFs()
	ghtknEnabled, err := github.GetGHTKNEnabledFromEnv()
	if err != nil {
		return fmt.Errorf("check if ghtkn integration is enabled: %w", err)
	}
	gh := github.New(ctx, logger.Logger, github.GetGitHubTokenFromEnv(), ghtknEnabled)
	ctrl, err := run.New(fs, gh)
	if err != nil {
		return fmt.Errorf("initialize a controller: %w", err)
	}
	files := make(map[string]struct{}, len(args.Files))
	for _, file := range args.Files {
		files[file] = struct{}{}
	}
	return ctrl.Run(ctx, logger.Logger, &run.Input{ //nolint:wrapcheck
		Files:             files,
		AllowUnknownField: args.AllowUnknownField,
	})
}
