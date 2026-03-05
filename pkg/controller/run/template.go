package run

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/Masterminds/sprig/v3"
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
