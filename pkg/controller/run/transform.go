package run

import (
	"fmt"
	"text/template"
)

func applyTransform(transform *Transform, result *TemplateInput) error {
	if transform == nil {
		return nil
	}
	fields := []struct {
		name    string
		tmplStr string
		target  *string
	}{
		{"content", transform.Content, &result.Content},
		{"stdout", transform.Stdout, &result.Stdout},
		{"stderr", transform.Stderr, &result.Stderr},
		{"combined_output", transform.CombinedOutput, &result.CombinedOutput},
	}
	for _, f := range fields {
		if f.tmplStr == "" {
			continue
		}
		tpl, err := template.New("_").Funcs(txtFuncMap()).Parse(f.tmplStr)
		if err != nil {
			return fmt.Errorf("parse transform template for %s: %w", f.name, err)
		}
		s, err := execTpl(tpl, result)
		if err != nil {
			return fmt.Errorf("execute transform template for %s: %w", f.name, err)
		}
		*f.target = s
	}
	return nil
}
