package validate

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/docfresh/pkg/controller/run"
)

type Controller struct {
	fs     afero.Fs
	stderr io.Writer
}

type Input struct {
	Files             map[string]struct{}
	AllowUnknownField bool
}

func New(fs afero.Fs, stderr io.Writer) *Controller {
	return &Controller{
		fs:     fs,
		stderr: stderr,
	}
}

func (c *Controller) Validate(_ context.Context, logger *slog.Logger, input *Input) error {
	for file := range input.Files {
		if err := c.validateFile(logger, file, input.AllowUnknownField); err != nil {
			return fmt.Errorf("validate file %s: %w", file, err)
		}
	}
	return nil
}

func (c *Controller) validateFile(logger *slog.Logger, file string, allowUnknownField bool) error {
	b, err := afero.ReadFile(c.fs, file)
	if err != nil {
		return fmt.Errorf("read a file: %w", err)
	}
	content := string(b)
	strictOpt := &run.ParseOption{DisallowUnknownField: true}
	_, err = run.ParseFile(content, strictOpt)
	if err == nil {
		fmt.Fprintf(c.stderr, "%s: valid\n", file)
		return nil
	}
	if !allowUnknownField {
		var yamlErr *run.YAMLError
		if errors.As(err, &yamlErr) {
			fmt.Fprintln(c.stderr, yamlErr)
			return errors.New("parse file failed")
		}
		return fmt.Errorf("parse file: %w", err)
	}
	// Try normal parse to see if the error was about unknown fields.
	if _, normalErr := run.ParseFile(content, nil); normalErr != nil {
		// Normal parse also fails — real YAML error.
		var yamlErr *run.YAMLError
		if errors.As(err, &yamlErr) {
			fmt.Fprintln(c.stderr, yamlErr)
			return errors.New("parse file failed")
		}
		return fmt.Errorf("parse file: %w", err)
	}
	// Normal parse succeeded — the error was about unknown fields.
	var yamlErr *run.YAMLError
	if errors.As(err, &yamlErr) {
		fmt.Fprintln(c.stderr, yamlErr)
	} else {
		logger.Warn("unknown fields detected", "file", file, "error", err)
	}
	fmt.Fprintf(c.stderr, "%s: valid (with warnings)\n", file)
	return nil
}
