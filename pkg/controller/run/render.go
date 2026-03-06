package run

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"text/template"

	"github.com/expr-lang/expr"
	"github.com/suzuki-shunsuke/go-error-with-exit-code/ecerror"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

func (c *Controller) renderBlock(ctx context.Context, logger *slog.Logger, tpls *Templates, file string, block *Block) (gS string, gErr error) { //nolint:cyclop
	if block.Type == "text" {
		return block.Content, nil
	}
	if block.Input == nil {
		return "", errors.New("block input is nil")
	}
	tpl, err := getTemplate(c.fs, tpls, block, file)
	if err != nil {
		return "", err
	}
	content := block.BeginComment
	if block.Input.PostCommand != nil {
		defer func() {
			if err := c.runPostCommand(ctx, file, block); err != nil {
				if gErr == nil {
					gErr = err
					return
				}
				slogerr.WithError(logger, err).Error("execute post_command")
			}
		}()
	}
	if err := c.runPreCommand(ctx, file, block); err != nil {
		return "", fmt.Errorf("execute pre_command: %w", err)
	}
	result, err := c.exec(ctx, file, block.Input)
	if err != nil {
		return "", fmt.Errorf("execute a command: %w", err)
	}
	if t := block.Input.Test(); t != "" {
		if err := testResult(c.stderr, t, result); err != nil {
			return "", err
		}
	}
	s, err := c.render(tpl, result, block.Input.TemplateData())
	if err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return appendEndComment(content, s, block.EndComment), nil
}

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

func appendEndComment(content, s, endComment string) string {
	if strings.HasSuffix(s, "\n") {
		return content + "\n" + s + endComment
	}
	return content + "\n" + s + "\n" + endComment
}

func testResult(stderr io.Writer, testCode string, result *TemplateInput) error {
	prog, err := expr.Compile(testCode, expr.Env(result), expr.AsBool())
	if err != nil {
		fmt.Fprintf(stderr, `[ERROR] compile an expression
%v`, err)
		return ecerror.Wrap(nil, 1)
	}
	output, err := expr.Run(prog, result)
	if err != nil {
		fmt.Fprintf(stderr, `[ERROR] evaluate an expression
%v`, err)
		return ecerror.Wrap(nil, 1)
	}
	f, ok := output.(bool)
	if !ok {
		return errors.New("the test result must be boolean")
	}
	if !f {
		return slogerr.With(errors.New("test failed"), "test", testCode) //nolint:wrapcheck
	}
	return nil
}

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

func (c *Controller) render(tpl *template.Template, result *TemplateInput, templateData *TemplateData) (string, error) {
	switch result.Type {
	case "local-file", "http", "github-content":
		if templateData == nil {
			return result.Content, nil
		}
		fns := txtFuncMap()
		tpl, err := template.New("_").Funcs(fns).Parse(result.Content)
		if err != nil {
			return "", fmt.Errorf("parse command template: %w", err)
		}
		result.Vars = templateData.Vars
		content, err := execTpl(tpl, result)
		if err != nil {
			return "", err
		}
		if tpl == nil {
			return content, nil
		}
		result.Content = content
		return execTpl(tpl, result)
	case "command":
		return execTpl(tpl, result)
	default:
		return "", fmt.Errorf("unknown type: %s", result.Type)
	}
}
