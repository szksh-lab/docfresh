package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

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
	result.UseFencedCodeBlockForOutput = block.Input.GetUseFencedCodeBlockForOutput()
	if t := block.Input.Test(); t != "" {
		if err := testResult(c.stderr, t, result); err != nil {
			return "", err
		}
	}
	s, err := render(tpl, result, block.Input.TemplateData())
	if err != nil {
		return "", fmt.Errorf("render template: %w", err)
	}
	return appendEndComment(content, s, block.EndComment), nil
}

func appendEndComment(content, s, endComment string) string {
	if strings.HasSuffix(s, "\n") {
		return content + "\n" + s + endComment
	}
	return content + "\n" + s + "\n" + endComment
}

func render(tpl *template.Template, result *TemplateInput, templateData *TemplateData) (string, error) {
	switch result.Type {
	case "local-file", "http", "github-content":
		return renderFile(tpl, result, templateData)
	case "command":
		return execTpl(tpl, result)
	default:
		return "", fmt.Errorf("unknown type: %s", result.Type)
	}
}

func renderFile(tpl *template.Template, result *TemplateInput, templateData *TemplateData) (string, error) {
	if templateData == nil {
		if tpl != nil {
			return execTpl(tpl, result)
		}
		if !result.UseFencedCodeBlockForOutput {
			return result.Content, nil
		}
		if !strings.HasSuffix(result.Content, "\n") {
			result.Content += "\n"
		}
		return "```" + result.ScriptLanguage + "\n" + result.Content + "```", nil
	}
	fns := txtFuncMap()
	contentTpl, err := template.New("_").Funcs(fns).Parse(result.Content)
	if err != nil {
		return "", fmt.Errorf("parse command template: %w", err)
	}
	result.Vars = templateData.Vars
	content, err := execTpl(contentTpl, result)
	if err != nil {
		return "", err
	}
	if contentTpl == nil {
		return content, nil
	}
	result.Content = content
	return execTpl(contentTpl, result)
}
