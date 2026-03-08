package run

import (
	_ "embed"
	"fmt"

	"github.com/goccy/go-yaml"
)

//go:embed languages.yaml
var languagesYAML []byte

type ScriptLanguage struct {
	Extensions   []string `yaml:"extensions"`
	ScriptShell  []string `yaml:"script_shell"`
	CommandShell []string `yaml:"command_shell"`
}

type Language struct {
	ScriptShell  []string
	CommandShell []string
	Language     string
}

func defaultScriptLanguages() (map[string]*Language, map[string]*Language, error) {
	langs := map[string]*ScriptLanguage{}
	if err := yaml.Unmarshal(languagesYAML, &langs); err != nil {
		return nil, nil, fmt.Errorf("unmrshal languages.yaml: %w", err)
	}
	byExt := make(map[string]*Language, len(langs))
	byName := make(map[string]*Language, len(langs))
	for langName, lang := range langs {
		l := &Language{
			ScriptShell:  lang.ScriptShell,
			CommandShell: lang.CommandShell,
			Language:     langName,
		}
		byName[langName] = l
		for _, ext := range lang.Extensions {
			byExt[ext] = l
		}
	}
	return byExt, byName, nil
}
