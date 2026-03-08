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

func (c *Controller) run(ctx context.Context, logger *slog.Logger, tpls *Templates, file string) (gErr error) { //nolint:cyclop
	b, err := afero.ReadFile(c.fs, file)
	if err != nil {
		return fmt.Errorf("read a file: %w", err)
	}
	bs := string(b)
	blocks, err := parseFile(string(b))
	if err != nil {
		return fmt.Errorf("parse a file: %w", err)
	}

	var postBlocks []*Block
	defer func() {
		for i := len(postBlocks) - 1; i >= 0; i-- {
			if err := c.runPostBlock(ctx, logger, tpls, file, postBlocks[i]); err != nil {
				if gErr == nil {
					gErr = err
				} else {
					slogerr.WithError(logger, err).Error("execute post block")
				}
			}
		}
	}()

	var contentBuilder strings.Builder
	for _, block := range blocks {
		if block.Type == "post" {
			contentBuilder.WriteString(block.Content)
			postBlocks = append(postBlocks, block)
			continue
		}
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
	PreCommand    *Command       `json:"pre_command,omitempty" yaml:"pre_command" jsonschema_description:"Execute external commands before the action like command, file, http, and github_content. If it fails, docfresh fails and the action isn't run. The command and output are outputted to the console but the result isn't affected to the document. This is used for setup and checking the requirement"`
	PostCommand   *Command       `json:"post_command,omitempty" yaml:"post_command" jsonschema_description:"Execute external commands after the action like command, file, http, and github_content. If it fails, docfresh fails. The command and output are outputted to the console but the result isn't affected to the document. This is used for testing the action result and cleaning up. post_command is run even if pre_command and action fail."`
	Command       *Command       `json:"command,omitempty" jsonschema_description:"Execute the external command and embed the result to documents"`
	File          *File          `json:"file,omitempty" jsonschema_description:"Read a local file and embed the content to documents"`
	HTTP          *HTTP          `json:"http,omitempty" jsonschema_description:"Call a HTTP request and embed the response to documents"`
	GitHubContent *GitHubContent `json:"github_content,omitempty" yaml:"github_content" jsonschema_description:"Fetch a file by GitHub Contents API and embed it into documents"`
	Template      *Template      `json:"template,omitempty" jsonschema_description:"Customize the template"`
	CodeBlock     *bool          `json:"code_block,omitempty" yaml:"code_block" jsonschema_description:"If this is true, the content is wrapped using markdown's fenced code block"`
	DetailsTag    *DetailsTag    `json:"details_tag,omitempty" yaml:"details_tag" jsonschema_description:"Wrap the output in an HTML details tag"`
}

type DetailsTag struct {
	Summary string `json:"summary,omitempty" jsonschema_description:"The summary text. Defaults: 'Output' for commands, file path for file, URL for http, '<owner>/<repo>/<path>[@<ref>]' for github_content."`
}

func (b *BlockInput) GetCodeBlock() bool {
	if b.CodeBlock != nil {
		return *b.CodeBlock
	}
	if b.Command != nil {
		return true
	}
	return false
}

type TemplateData struct {
	Vars   map[string]any `json:"vars,omitempty" jsonschema_description:"Variables which are passed to template. They can be referred in templates as .Vars.<variable name>"`
	Delims *Delims        `json:"delims,omitempty" jsonschema_description:"The delimiters. The default delimiters are '{{' and '}}'"`
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
	Owner    string        `json:"owner" jsonschema_description:"GitHub repository owner"`
	Repo     string        `json:"repo" jsonschema_description:"GitHub repository name"`
	Ref      string        `json:"ref,omitempty" jsonschema_description:"The ref (branch, tag, SHA). The default branch is used by default."`
	Path     string        `json:"path" jsonschema_description:"The path of GitHub Contents API"`
	Range    *Range        `json:"range,omitempty" jsonschema_description:"Extract a specific range of lines from the content. Uses 0-based indexing with half-open interval [start, end). Negative values count from the end."`
	Template *TemplateData `json:"template,omitempty" jsonschema_description:"If this is set, the file content is rendered as template rather than plain text."`
	Test     string        `json:"test,omitempty" jsonschema_description:"Expr script to test the file content. The evaluation result must be a boolean. If the evaluation result is false, docfresh fails"`
}

type Delims struct {
	Left  string `json:"left" jsonschema_description:"The left delimiter of templates"`
	Right string `json:"right" jsonschema_description:"The right delimiter of templates"`
}

type Template struct {
	Content  string             `json:"content,omitempty" jsonschema_description:"The content of template"`
	Path     string             `json:"path,omitempty" jsonschema_description:"The file path. It's an absolute path or relative path from the current file."`
	Template *template.Template `json:"-" yaml:"-"`
	Vars     map[string]any     `json:"vars,omitempty" jsonschema_description:"Variables which are passed to template. They can be referred in templates as .Vars.<variable name>"`
	Delims   *Delims            `json:"delims,omitempty" jsonschema_description:"The delimiters. The default delimiters are '{{' and '}}'"`
}

func (t *Template) GetVars() map[string]any {
	if t == nil {
		return nil
	}
	return t.Vars
}

type HTTP struct {
	URL      string        `json:"url" jsonschema_description:"URL for HTTP request"`
	Range    *Range        `json:"range,omitempty" jsonschema_description:"Extract a specific range of lines from the response. Uses 0-based indexing with half-open interval [start, end). Negative values count from the end."`
	Template *TemplateData `json:"template,omitempty" jsonschema_description:"If this is set, the response body is rendered as template rather than plain text."`
	Test     string        `json:"test,omitempty" jsonschema_description:"Expr script to test the response. The evaluation result must be a boolean. If the evaluation result is false, docfresh fails"`
	Timeout  int           `json:"timeout,omitempty" jsonschema_description:"HTTP request timeout (seconds). The default value is 5 seconds. If the value is negative, timeout isn't set"`
	Header   http.Header   `json:"header,omitempty" jsonschema_description:"HTTP request header."`
}

type File struct {
	Path     string        `json:"path" jsonschema_description:"The file path. It's an absolute path or relative path from the current file."`
	Range    *Range        `json:"range,omitempty" jsonschema_description:"Extract a specific range of lines from the file. Uses 0-based indexing with half-open interval [start, end). Negative values count from the end."`
	Template *TemplateData `json:"template,omitempty" jsonschema_description:"If this is set, the file is rendered as template rather than plain text."`
	Test     string        `json:"test,omitempty" jsonschema_description:"Expr script to test the file content. The evaluation result must be a boolean. If the evaluation result is false, docfresh fails"`
}

type Command struct {
	Command        string            `json:"command,omitempty" jsonschema_description:"The content of executed script. Either command or script is required"`
	Script         string            `json:"script,omitempty" jsonschema_description:"The file path to executed script. It's an absolute path or relative path from the current file. Either command or script is required"`
	Dir            string            `json:"dir,omitempty" jsonschema_description:"The directory path where commands are executed. It's an absolute path or relative path from the current file. The default value is the directory where the current file is located"`
	Test           string            `json:"test,omitempty" jsonschema_description:"Expr script to test the result of command. The evaluation result must be a boolean. If the evaluation result is false, docfresh fails"`
	ScriptLanguage string            `json:"script_language,omitempty" yaml:"script_language" jsonschema_description:"Language of script. This is used for markdown's fenced code block. This is automatically detected in some languages such as Go and Python"`
	Timeout        int               `json:"timeout,omitempty" jsonschema_description:"The timeout of command. By default, there is no timeout. If timeout is exceeded, the signal SIGINT is sent to the process."`
	TimeoutSigkill int               `json:"timeout_sigkill,omitempty" jsonschema_description:"If this timeout is exceeded, the signal SIGKILL is sent to the process. The default value is 1000 hours, meaning SIGKILL isn't sent usually, so the process should be terminated gracefully by SIGINT."`
	Shell          []string          `json:"shell,omitempty" jsonschema_description:"The command executing command or script. If command is set, the default value is 'bash -c'. If script is set, the default value is decided by script's file extension"`
	Envs           map[string]string `json:"envs,omitempty" jsonschema_description:"Pairs of environment variable names and values"`
	IgnoreFail     bool              `json:"ignore_fail,omitempty" yaml:"ignore_fail" jsonschema_description:"If this is true, docfresh does't fail even if command fails"`
	EmbedScript    bool              `json:"embed_script,omitempty" yaml:"embed_script" jsonschema_description:"If this is true, the content of script is embedded into documents."`
	Quiet          bool              `json:"quiet,omitempty" jsonschema_description:"If this is true, the command output isn't outputted to documents."`
}

type PostCommand struct {
	Command        string            `json:"command,omitempty" jsonschema_description:"The content of executed script. Either command or script is required"`
	Script         string            `json:"script,omitempty" jsonschema_description:"The file path to executed script. It's an absolute path or relative path from the current file. Either command or script is required"`
	Dir            string            `json:"dir,omitempty" jsonschema_description:"The directory path where commands are executed. It's an absolute path or relative path from the current file. The default value is the directory where the current file is located"`
	Test           string            `json:"test,omitempty" jsonschema_description:"Expr script to test the result of command. The evaluation result must be a boolean. If the evaluation result is false, docfresh fails"`
	Timeout        int               `json:"timeout,omitempty" jsonschema_description:"The timeout of command. By default, there is no timeout. If timeout is exceeded, the signal SIGINT is sent to the process."`
	TimeoutSigkill int               `json:"timeout_sigkill,omitempty" jsonschema_description:"If this timeout is exceeded, the signal SIGKILL is sent to the process. The default value is 1000 hours, meaning SIGKILL isn't sent usually, so the process should be terminated gracefully by SIGINT."`
	Shell          []string          `json:"shell,omitempty" jsonschema_description:"The command executing command or script. If command is set, the default value is 'bash -c'. If script is set, the default value is decided by script's file extension"`
	Envs           map[string]string `json:"envs,omitempty" jsonschema_description:"Pairs of environment variable names and values"`
	IgnoreFail     bool              `json:"ignore_fail,omitempty" yaml:"ignore_fail" jsonschema_description:"If this is true, docfresh does't fail even if command fails"`
}

func (p *PostCommand) ToCommand() *Command {
	return &Command{
		Command:        p.Command,
		Script:         p.Script,
		Dir:            p.Dir,
		Test:           p.Test,
		Timeout:        p.Timeout,
		TimeoutSigkill: p.TimeoutSigkill,
		Shell:          p.Shell,
		Envs:           p.Envs,
		IgnoreFail:     p.IgnoreFail,
	}
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
	EmbedScript       bool
	CodeBlock         bool
	Quiet             bool
	DetailsTagSummary string
}
