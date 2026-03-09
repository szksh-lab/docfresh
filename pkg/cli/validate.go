package cli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/docfresh/pkg/controller/validate"
	"github.com/suzuki-shunsuke/slog-util/slogutil"
	"github.com/urfave/cli/v3"
)

type ValidateArgs struct {
	*Flags

	Files             []string
	AllowUnknownField bool
}

func NewValidate(logger *slogutil.Logger, gFlags *Flags) *cli.Command {
	args := &ValidateArgs{
		Flags: gFlags,
	}
	return &cli.Command{
		Name:  "validate",
		Usage: "Validate documents",
		Action: func(ctx context.Context, _ *cli.Command) error {
			return validateAction(ctx, logger, args)
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
				Usage:       "Downgrade unknown field errors to warnings",
				Destination: &args.AllowUnknownField,
			},
		},
	}
}

func validateAction(ctx context.Context, logger *slogutil.Logger, args *ValidateArgs) error {
	if err := logger.SetLevel(args.LogLevel); err != nil {
		return fmt.Errorf("set log level: %w", err)
	}

	fs := afero.NewOsFs()
	ctrl := validate.New(fs, os.Stderr)
	files := make(map[string]struct{}, len(args.Files))
	for _, file := range args.Files {
		files[file] = struct{}{}
	}
	return ctrl.Validate(ctx, logger.Logger, &validate.Input{ //nolint:wrapcheck
		Files:             files,
		AllowUnknownField: args.AllowUnknownField,
	})
}
