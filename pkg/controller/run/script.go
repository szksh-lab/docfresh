package run

import (
	_ "embed"
	"fmt"

	"github.com/goccy/go-yaml"
)

//go:embed languages.yaml
var languagesYAML []byte

type ScriptLanguage struct {
	Extensions []string
	Shell      []string
}

type Language struct {
	Shell    []string
	Language string
}

func defaultScriptLanguages() (map[string]*Language, error) {
	langs := map[string]*ScriptLanguage{}
	if err := yaml.Unmarshal(languagesYAML, &langs); err != nil {
		return nil, fmt.Errorf("unmrshal languages.yaml: %w", err)
	}
	ret := make(map[string]*Language, len(langs))
	for langName, lang := range langs {
		for _, ext := range lang.Extensions {
			ret[ext] = &Language{
				Shell:    lang.Shell,
				Language: langName,
			}
		}
	}
	return ret, nil
}
