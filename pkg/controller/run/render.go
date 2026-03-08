package run

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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
			if err := c.runPostCommand(ctx, logger, file, block); err != nil {
				if gErr == nil {
					gErr = err
					return
				}
				slogerr.WithError(logger, err).Error("execute post_command")
			}
		}()
	}
	if err := c.runPreCommand(ctx, logger, file, block); err != nil {
		return "", fmt.Errorf("execute pre_command: %w", err)
	}
	result, err := c.exec(ctx, logger, file, block.Input)
	if err != nil {
		return "", fmt.Errorf("execute a command: %w", err)
	}
	result.UseFencedCodeBlockForOutput = block.Input.GetUseFencedCodeBlockForOutput()
	if dt := block.Input.DetailsTag; dt != nil {
		summary := dt.Summary
		if summary == "" {
			summary = defaultDetailsTagSummary(result)
		}
		result.DetailsTagSummary = summary
	}
	if t := block.Input.Test(); t != "" {
		if err := testResult(c.stderr, t, result); err != nil {
			return "", err
		}
	}
	s, err := render(tpl, result)
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

func render(tpl *Template, result *TemplateInput) (string, error) {
	switch result.Type {
	case "local-file", "http", "github-content":
		return renderFile(tpl, result)
	case "command":
		result.Vars = tpl.Vars
		return execTpl(tpl.Template, result)
	default:
		return "", fmt.Errorf("unknown type: %s", result.Type)
	}
}

func renderFile(tpl *Template, result *TemplateInput) (string, error) {
	var s string
	switch {
	case tpl != nil:
		result.Vars = tpl.Vars
		var err error
		s, err = execTpl(tpl.Template, result)
		if err != nil {
			return "", err
		}
	case !result.UseFencedCodeBlockForOutput:
		s = result.Content
	default:
		if !strings.HasSuffix(result.Content, "\n") {
			result.Content += "\n"
		}
		s = "```" + result.ScriptLanguage + "\n" + result.Content + "```"
	}
	if result.DetailsTagSummary != "" {
		s = wrapDetailsTag(s, result.DetailsTagSummary)
	}
	return s, nil
}

func defaultDetailsTagSummary(result *TemplateInput) string {
	switch result.Type {
	case "command":
		return "Output"
	case "local-file":
		return result.Path
	case "http":
		return result.URL
	case "github-content":
		s := result.Owner + "/" + result.Repo + "/" + result.Path
		if result.Ref != "" {
			s += "@" + result.Ref
		}
		return s
	default:
		return "Output"
	}
}

func wrapDetailsTag(content, summary string) string {
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}
	return "<details>\n<summary>" + summary + "</summary>\n\n" + content + "\n</details>"
}
