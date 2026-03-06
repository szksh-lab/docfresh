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

func getTemplate(fs afero.Fs, tpls *Templates, block *Block, file string) (*template.Template, error) {
	if block.Input.Template == nil {
		if block.Input.Command != nil {
			return tpls.Command, nil
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
	tpl, err := template.New("_").Funcs(tpls.Funcs).Parse(content)
	if err != nil {
		return nil, fmt.Errorf("parse block template: %w", err)
	}
	return tpl, nil
}
