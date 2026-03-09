package langdetect

import (
	"testing"
)

func TestDefaultScriptLanguages(t *testing.T) { //nolint:gocognit,cyclop,funlen
	t.Parallel()
	langs, langsByName, err := DefaultScriptLanguages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(langs) == 0 {
		t.Fatal("expected non-empty languages map")
	}
	if len(langsByName) == 0 {
		t.Fatal("expected non-empty langsByName map")
	}

	t.Run("by_extension", func(t *testing.T) {
		t.Parallel()
		checks := []struct {
			ext            string
			language       string
			hasScriptShell bool
		}{
			{ext: ".go", language: "go", hasScriptShell: true},
			{ext: ".py", language: "py", hasScriptShell: true},
			{ext: ".sh", language: "sh", hasScriptShell: true},
			{ext: ".bash", language: "sh", hasScriptShell: true},
			{ext: ".js", language: "js", hasScriptShell: true},
			{ext: ".json", language: "json", hasScriptShell: false},
			{ext: ".yaml", language: "yaml", hasScriptShell: false},
			{ext: ".yml", language: "yaml", hasScriptShell: false},
			{ext: ".md", language: "md", hasScriptShell: false},
		}
		for _, c := range checks {
			lang, ok := langs[c.ext]
			if !ok {
				t.Errorf("expected extension %q to be present", c.ext)
				continue
			}
			if lang.Language != c.language {
				t.Errorf("extension %q: language = %q, want %q", c.ext, lang.Language, c.language)
			}
			if c.hasScriptShell && len(lang.ScriptShell) == 0 {
				t.Errorf("extension %q: expected non-empty ScriptShell", c.ext)
			}
			if !c.hasScriptShell && len(lang.ScriptShell) != 0 {
				t.Errorf("extension %q: expected empty ScriptShell, got %v", c.ext, lang.ScriptShell)
			}
		}
	})

	t.Run("by_name", func(t *testing.T) {
		t.Parallel()
		nameChecks := []struct {
			name            string
			hasScriptShell  bool
			hasCommandShell bool
		}{
			{name: "js", hasScriptShell: true, hasCommandShell: true},
			{name: "py", hasScriptShell: true, hasCommandShell: true},
			{name: "sh", hasScriptShell: true, hasCommandShell: true},
			{name: "go", hasScriptShell: true, hasCommandShell: false},
			{name: "json", hasScriptShell: false, hasCommandShell: false},
		}
		for _, c := range nameChecks {
			lang, ok := langsByName[c.name]
			if !ok {
				t.Errorf("expected language %q to be present in langsByName", c.name)
				continue
			}
			if c.hasScriptShell && len(lang.ScriptShell) == 0 {
				t.Errorf("language %q: expected non-empty ScriptShell", c.name)
			}
			if c.hasCommandShell && len(lang.CommandShell) == 0 {
				t.Errorf("language %q: expected non-empty CommandShell", c.name)
			}
		}
	})
}
