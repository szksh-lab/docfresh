package run

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"github.com/suzuki-shunsuke/slog-error/slogerr"
)

//go:embed command_template.md
var commandTemplate string

type Input struct {
	// ConfigFilePath string
	Files map[string]struct{}
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
	PreCommand                  *Command       `json:"pre_command,omitempty" yaml:"pre_command"`
	PostCommand                 *Command       `json:"post_command,omitempty" yaml:"post_command"`
	Command                     *Command       `json:"command,omitempty"`
	File                        *File          `json:"file,omitempty"`
	HTTP                        *HTTP          `json:"http,omitempty"`
	GitHubContent               *GitHubContent `json:"github_content,omitempty" yaml:"github_content"`
	Template                    *Template      `json:"template,omitempty"`
	UseFencedCodeBlockForOutput *bool          `json:"use_fenced_code_block_for_output,omitempty" yaml:"use_fenced_code_block_for_output"`
}

func (b *BlockInput) GetUseFencedCodeBlockForOutput() bool {
	if b.UseFencedCodeBlockForOutput != nil {
		return *b.UseFencedCodeBlockForOutput
	}
	if b.Command != nil {
		return true
	}
	return false
}

type TemplateData struct {
	Vars map[string]any `json:"vars,omitempty"`
}

func (t *TemplateData) GetVars() map[string]any {
	if t == nil {
		return nil
	}
	return t.Vars
}

func (b *BlockInput) TemplateData() *TemplateData {
	if b.File != nil {
		return b.File.Template
	}
	if b.HTTP != nil {
		return b.HTTP.Template
	}
	if b.GitHubContent != nil {
		return b.GitHubContent.Template
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
	Owner    string        `json:"owner"`
	Repo     string        `json:"repo"`
	Ref      string        `json:"ref,omitempty"`
	Path     string        `json:"path"`
	Template *TemplateData `json:"template,omitempty"`
	Test     string        `json:"test,omitempty"`
}

type Template struct {
	Template *template.Template `json:"-" yaml:"-"`
	Content  string             `json:"content,omitempty"`
	Path     string             `json:"path,omitempty"`
	Vars     map[string]any     `json:"vars,omitempty"`
}

func (t *Template) GetVars() map[string]any {
	if t == nil {
		return nil
	}
	return t.Vars
}

type HTTP struct {
	URL      string        `json:"url"`
	Template *TemplateData `json:"template,omitempty"`
	Test     string        `json:"test,omitempty"`
	Timeout  int           `json:"timeout,omitempty"`
	Header   http.Header   `json:"header,omitempty"`
}

type File struct {
	Path     string        `json:"path"`
	Template *TemplateData `json:"template,omitempty"`
	Test     string        `json:"test,omitempty"`
}

type Command struct {
	Command        string            `json:"command,omitempty"`
	Script         string            `json:"script,omitempty"`
	Dir            string            `json:"dir,omitempty"`
	Test           string            `json:"test,omitempty"`
	ScriptLanguage string            `json:"script_language,omitempty" yaml:"script_language"`
	Timeout        int               `json:"timeout,omitempty"`
	TimeoutSigkill int               `json:"timeout_sigkill,omitempty"`
	Shell          []string          `json:"shell,omitempty"`
	Envs           map[string]string `json:"envs,omitempty"`
	IgnoreFail     bool              `json:"ignore_fail,omitempty" yaml:"ignore_fail"`
	EmbedScript    bool              `json:"embed_script,omitempty" yaml:"embed_script"`
}

type TemplateInput struct {
	Type string
	// command
	Command        string
	Script         string
	Shell          []string
	Dir            string
	Stdout         string
	Stderr         string
	CombinedOutput string
	ScriptLanguage string
	ExitCode       int
	// file
	Path    string
	Content string
	// http
	URL     string
	Timeout int
	// github content
	Owner string
	Repo  string
	Ref   string
	// template variables
	Vars map[string]any
	//
	EmbedScript                 bool
	UseFencedCodeBlockForOutput bool
}
