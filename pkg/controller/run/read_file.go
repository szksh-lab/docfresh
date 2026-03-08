package run

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func (c *Controller) readFile(baseFile string, file *File) (*TemplateInput, error) {
	p := resolvePath(baseFile, file.Path)
	b, err := afero.ReadFile(c.fs, p)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	ext := filepath.Ext(p)
	lang, ok := c.langs[ext]
	sl := ""
	if ok {
		sl = lang.Language
	}
	content := string(b)
	content, err = extractRange(content, file.Range)
	if err != nil {
		return nil, fmt.Errorf("extract range from file: %w", err)
	}

	result := &TemplateInput{
		Type:     "local-file",
		Path:     file.Path,
		Language: sl,
		Content:  content,
		Vars:     file.Template.GetVars(),
	}

	if file.Template != nil {
		if err := renderTemplate(content, result, file.Template.Delims); err != nil {
			return nil, err
		}
	}

	return result, nil
}
