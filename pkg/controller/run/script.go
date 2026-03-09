package run

import (
	"fmt"

	"github.com/suzuki-shunsuke/docfresh/pkg/langdetect"
)

type Language = langdetect.Language

func defaultScriptLanguages() (map[string]*Language, map[string]*Language, error) {
	byExt, byName, err := langdetect.DefaultScriptLanguages()
	if err != nil {
		return nil, nil, fmt.Errorf("load default script languages: %w", err)
	}
	return byExt, byName, nil
}
