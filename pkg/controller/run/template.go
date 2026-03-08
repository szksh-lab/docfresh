package run

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/afero"
)

func txtFuncMap() template.FuncMap {
	fncs := sprig.TxtFuncMap()
	delete(fncs, "env")
	delete(fncs, "expandenv")
	delete(fncs, "getHostByName")
	return fncs
}

func execTpl(tpl *template.Template, data any) (string, error) {
	buf := &bytes.Buffer{}
	if err := tpl.Execute(buf, data); err != nil {
		return "", fmt.Errorf("execute a template: %w", err)
	}
	return buf.String(), nil
}

func getTemplate(fs afero.Fs, tpls *Templates, block *Block, file string) (*Template, error) {
	if block.Input.Template == nil {
		if block.Input.Command != nil {
			return &Template{
				Template: tpls.Command,
			}, nil
		}
		return nil, nil //nolint:nilnil
	}
	content := block.Input.Template.Content
	if block.Input.Template.Path != "" {
		p := resolvePath(file, block.Input.Template.Path)
		b, err := afero.ReadFile(fs, p)
		if err != nil {
			return nil, fmt.Errorf("read a template file: %w", err)
		}
		content = string(b)
	}
	t := template.New("_").Funcs(tpls.Funcs)
	if d := block.Input.Template.Delims; d != nil {
		t = t.Delims(d.Left, d.Right)
	}
	tpl, err := t.Parse(content)
	if err != nil {
		return nil, fmt.Errorf("parse block template: %w", err)
	}
	return &Template{
		Template: tpl,
		Vars:     block.Input.Template.Vars,
	}, nil
}

func renderTemplate(content string, result *TemplateInput, delims *Delims) error {
	fns := txtFuncMap()
	t := template.New("_").Funcs(fns)
	if delims != nil {
		t = t.Delims(delims.Left, delims.Right)
	}
	tpl, err := t.Parse(content)
	if err != nil {
		return fmt.Errorf("parse file template: %w", err)
	}
	s, err := execTpl(tpl, result)
	if err != nil {
		return fmt.Errorf("render file as template: %w", err)
	}
	result.Content = s
	return nil
}
