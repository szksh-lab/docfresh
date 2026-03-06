package run

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func (c *Controller) readFile(baseFile, file string) (*TemplateInput, error) {
	p := resolvePath(baseFile, file)
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
	return &TemplateInput{
		Type:           "local-file",
		Path:           file,
		ScriptLanguage: sl,
		Content:        string(b),
	}, nil
}
