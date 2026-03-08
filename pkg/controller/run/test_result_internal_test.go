package run

import (
	"bytes"
	"testing"
)

func TestTestResult(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		testCode string
		result   *TemplateInput
		wantErr  bool
	}{
		{
			name:     "passing test",
			testCode: "Stdout == \"hello\"",
			result:   &TemplateInput{Stdout: "hello"},
		},
		{
			name:     "failing test",
			testCode: "Stdout == \"world\"",
			result:   &TemplateInput{Stdout: "hello"},
			wantErr:  true,
		},
		{
			name:     "compile error",
			testCode: "invalid ===",
			result:   &TemplateInput{},
			wantErr:  true,
		},
		{
			name:     "test with ExitCode",
			testCode: "ExitCode == 0",
			result:   &TemplateInput{ExitCode: 0},
		},
		{
			name:     "test with ExitCode failure",
			testCode: "ExitCode == 0",
			result:   &TemplateInput{ExitCode: 1},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			stderr := &bytes.Buffer{}
			err := testResult(stderr, tt.testCode, tt.result)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}
