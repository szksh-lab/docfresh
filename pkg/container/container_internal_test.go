package container

import (
	"testing"
)

func TestValidateInput(t *testing.T) { //nolint:funlen
	t.Parallel()
	tests := []struct {
		name     string
		input    *Input
		existing map[string]*State
		wantErr  bool
	}{
		{
			name: "valid input",
			input: &Input{
				ID:     "test",
				Engine: "docker-cli",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*State{},
			wantErr:  false,
		},
		{
			name: "missing id",
			input: &Input{
				Engine: "docker-cli",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*State{},
			wantErr:  true,
		},
		{
			name: "missing engine",
			input: &Input{
				ID:    "test",
				Image: "ubuntu:latest",
			},
			existing: map[string]*State{},
			wantErr:  true,
		},
		{
			name: "unsupported engine",
			input: &Input{
				ID:     "test",
				Engine: "podman",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*State{},
			wantErr:  true,
		},
		{
			name: "missing image",
			input: &Input{
				ID:     "test",
				Engine: "docker-cli",
			},
			existing: map[string]*State{},
			wantErr:  true,
		},
		{
			name: "duplicate id",
			input: &Input{
				ID:     "test",
				Engine: "docker-cli",
				Image:  "ubuntu:latest",
			},
			existing: map[string]*State{
				"test": {},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateInput(tt.input, tt.existing)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewEngine(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		engine  string
		wantErr bool
	}{
		{
			name:    "docker-cli",
			engine:  "docker-cli",
			wantErr: false,
		},
		{
			name:    "unsupported",
			engine:  "podman",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			e, err := NewEngine(tt.engine)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEngine() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && e == nil {
				t.Error("NewEngine() returned nil engine")
			}
		})
	}
}
