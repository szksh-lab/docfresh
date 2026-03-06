package run

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

//go:embed command_template.md
var commandTemplate string

type Input struct {
	ConfigFilePath string
	Files          map[string]struct{}
}

type Templates struct {
	Funcs   template.FuncMap
	Command *template.Template
	File    *template.Template
}

func (c *Controller) Run(ctx context.Context, logger *slog.Logger, input *Input) error {
	fns := txtFuncMap()
	cmdTpl, err := template.New("_").Funcs(fns).Parse(commandTemplate)
	if err != nil {
		return fmt.Errorf("parse command template: %w", err)
	}
	tpls := &Templates{
		Command: cmdTpl,
		Funcs:   fns,
	}
	for file := range input.Files {
		logger := logger.With("file", file)
		if err := c.run(ctx, logger, tpls, file); err != nil {
			return slogerr.With(err, "file", file) //nolint:wrapcheck
		}
	}
	return nil
}

func (c *Controller) run(ctx context.Context, logger *slog.Logger, tpls *Templates, file string) error {
	b, err := afero.ReadFile(c.fs, file)
	if err != nil {
		return fmt.Errorf("read a file: %w", err)
	}
	bs := string(b)
	blocks, err := parseFile(string(b))
	if err != nil {
		return fmt.Errorf("parse a file: %w", err)
	}
	var contentBuilder strings.Builder
	for _, block := range blocks {
		s, err := c.renderBlock(ctx, logger, tpls, file, block)
		if err != nil {
			return err
		}
		contentBuilder.WriteString(s)
	}
	content := contentBuilder.String()
	stat, err := c.fs.Stat(file)
	if err != nil {
		return fmt.Errorf("get file info: %w", err)
	}
	if content != bs {
		if err := afero.WriteFile(c.fs, file, []byte(content), stat.Mode()); err != nil {
			return fmt.Errorf("update a file: %w", err)
		}
	}

	return nil
}

type Block struct {
	// text, code block
	Type         string
	Content      string
	Input        *BlockInput
	BeginComment string
	EndComment   string
}

type BlockInput struct {
	PreCommand    *Command       `yaml:"pre_command,omitempty"`
	PostCommand   *Command       `yaml:"post_command,omitempty"`
	Command       *Command       `yaml:",omitempty"`
	File          *File          `yaml:",omitempty"`
	HTTP          *HTTP          `yaml:",omitempty"`
	GitHubContent *GitHubContent `yaml:"github_content,omitempty"`
	Template      *Template      `yaml:",omitempty"`
}

type TemplateData struct {
	Vars map[string]any
}

func (b *BlockInput) TemplateData() *TemplateData {
	if b.File != nil {
		return b.File.TemplateData
	}
	if b.HTTP != nil {
		return b.HTTP.TemplateData
	}
	if b.GitHubContent != nil {
		return b.GitHubContent.TemplateData
	}
	return nil
}

func (b *BlockInput) Test() string {
	if b.File != nil {
		return b.File.Test
	}
	if b.HTTP != nil {
		return b.HTTP.Test
	}
	if b.GitHubContent != nil {
		return b.GitHubContent.Test
	}
	if b.Command != nil {
		return b.Command.Test
	}
	return ""
}

type GitHubContent struct {
	Owner        string
	Repo         string
	Ref          string
	Path         string
	TemplateData *TemplateData `yaml:"template"`
	Test         string
}

type Template struct {
	Content string
}

type HTTP struct {
	URL          string
	TemplateData *TemplateData `yaml:"template"`
	Test         string
}

type File struct {
	Path         string
	TemplateData *TemplateData `yaml:"template"`
	Test         string
}

type Command struct {
	Command    string
	Dir        string   `yaml:",omitempty"`
	Shell      []string `yaml:",omitempty"`
	IgnoreFail bool     `yaml:"ignore_fail,omitempty"`
	Test       string
}

type TemplateInput struct {
	Type string
	// command
	Command        string
	Dir            string
	Stdout         string
	Stderr         string
	CombinedOutput string
	ExitCode       int
	// file
	Path    string
	Content string
	// http
	URL string
	// github content
	Owner string
	Repo  string
	Ref   string
	// template variables
	Vars map[string]any
}
