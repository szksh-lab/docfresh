package run

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/afero"
)

func readFile(baseFile string, file *File, fs afero.Fs, langs map[string]*Language) (*TemplateInput, error) {
	p := resolvePath(baseFile, file.Path)
	b, err := afero.ReadFile(fs, p)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	sl := file.Language
	if sl == "" {
		ext := filepath.Ext(p)
		if lang, ok := langs[ext]; ok {
			sl = lang.Language
		}
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
