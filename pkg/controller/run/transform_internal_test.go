package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestApplyTransform(t *testing.T) { //nolint:funlen
	t.Parallel()
	tests := []struct {
		name      string
		transform *Transform
		input     *TemplateInput
		expected  *TemplateInput
		wantErr   bool
	}{
		{
			name:      "nil transform",
			transform: nil,
			input: &TemplateInput{
				Stdout: "hello",
			},
			expected: &TemplateInput{
				Stdout: "hello",
			},
		},
		{
			name: "transform stdout only",
			transform: &Transform{
				Stdout: `{{ regexReplaceAll "secret-[a-z]+" .Stdout "***" }}`,
			},
			input: &TemplateInput{
				Stdout:  "token is secret-abc",
				Content: "unchanged",
			},
			expected: &TemplateInput{
				Stdout:  "token is ***",
				Content: "unchanged",
			},
		},
		{
			name: "transform content only",
			transform: &Transform{
				Content: `{{ upper .Content }}`,
			},
			input: &TemplateInput{
				Content: "hello world",
				Stdout:  "unchanged",
			},
			expected: &TemplateInput{
				Content: "HELLO WORLD",
				Stdout:  "unchanged",
			},
		},
		{
			name: "transform multiple fields",
			transform: &Transform{
				Stdout:         `{{ upper .Stdout }}`,
				Stderr:         `{{ lower .Stderr }}`,
				CombinedOutput: `{{ trimSuffix "\n" .CombinedOutput }}`,
			},
			input: &TemplateInput{
				Stdout:         "hello",
				Stderr:         "WARNING",
				CombinedOutput: "output\n",
			},
			expected: &TemplateInput{
				Stdout:         "HELLO",
				Stderr:         "warning",
				CombinedOutput: "output",
			},
		},
		{
			name: "transform with regexReplaceAll",
			transform: &Transform{
				Stdout: `{{ regexReplaceAll "\\d{4}-\\d{2}-\\d{2}" .Stdout "YYYY-MM-DD" }}`,
			},
			input: &TemplateInput{
				Stdout: "created at 2024-01-15",
			},
			expected: &TemplateInput{
				Stdout: "created at YYYY-MM-DD",
			},
		},
		{
			name: "invalid template",
			transform: &Transform{
				Stdout: `{{ invalid`,
			},
			input: &TemplateInput{
				Stdout: "hello",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := applyTransform(tt.transform, tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.expected, tt.input); diff != "" {
				t.Errorf("mismatch (-expected +got):\n%s", diff)
			}
		})
	}
}
