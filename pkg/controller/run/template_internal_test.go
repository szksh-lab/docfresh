package run

import (
	"testing"
	"text/template"

	"github.com/google/go-cmp/cmp"
)

func TestTxtFuncMap(t *testing.T) {
	t.Parallel()
	fns := txtFuncMap()
	for _, name := range []string{"env", "expandenv", "getHostByName"} {
		if _, ok := fns[name]; ok {
			t.Errorf("dangerous function %q should be removed", name)
		}
	}
	// Verify some expected sprig functions still exist
	for _, name := range []string{"upper", "lower", "trim"} {
		if _, ok := fns[name]; !ok {
			t.Errorf("expected function %q to be present", name)
		}
	}
}

func TestExecTpl(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		tplStr  string
		data    any
		want    string
		wantErr bool
	}{
		{
			name:   "simple template",
			tplStr: "hello {{.Name}}",
			data:   map[string]string{"Name": "world"},
			want:   "hello world",
		},
		{
			name:    "invalid template execution",
			tplStr:  "{{.Missing}}",
			data:    "not a map",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tpl := template.Must(template.New("test").Parse(tt.tplStr))
			got, err := execTpl(tpl, tt.data)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRenderTemplate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		content string
		result  *TemplateInput
		delims  *Delims
		want    string
		wantErr bool
	}{
		{
			name:    "basic template rendering",
			content: "path: {{.Path}}",
			result:  &TemplateInput{Path: "README.md"},
			want:    "path: README.md",
		},
		{
			name:    "custom delimiters",
			content: "path: <<.Path>>",
			result:  &TemplateInput{Path: "README.md"},
			delims:  &Delims{Left: "<<", Right: ">>"},
			want:    "path: README.md",
		},
		{
			name:    "invalid template syntax",
			content: "{{.Invalid",
			result:  &TemplateInput{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := *tt.result // copy to avoid mutation across subtests
			err := renderTemplate(tt.content, &result, tt.delims)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, result.Content); diff != "" {
				t.Errorf("mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
