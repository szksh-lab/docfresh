package run

import (
	"testing"
)

func TestDefaultScriptLanguages(t *testing.T) {
	t.Parallel()
	langs, err := defaultScriptLanguages()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(langs) == 0 {
		t.Fatal("expected non-empty languages map")
	}

	// Verify some known extensions from languages.yaml
	checks := []struct {
		ext      string
		language string
		hasShell bool
	}{
		{ext: ".go", language: "go", hasShell: true},
		{ext: ".py", language: "py", hasShell: true},
		{ext: ".sh", language: "sh", hasShell: true},
		{ext: ".bash", language: "sh", hasShell: true},
		{ext: ".js", language: "js", hasShell: true},
		{ext: ".json", language: "json", hasShell: false},
		{ext: ".yaml", language: "yaml", hasShell: false},
		{ext: ".yml", language: "yaml", hasShell: false},
		{ext: ".md", language: "md", hasShell: false},
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
		if c.hasShell && len(lang.Shell) == 0 {
			t.Errorf("extension %q: expected non-empty shell", c.ext)
		}
		if !c.hasShell && len(lang.Shell) != 0 {
			t.Errorf("extension %q: expected empty shell, got %v", c.ext, lang.Shell)
		}
	}
}
