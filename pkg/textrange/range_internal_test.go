package textrange

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestExtractRange(t *testing.T) { //nolint:funlen
	t.Parallel()
	tests := []struct {
		name    string
		content string
		r       *Range
		want    string
		wantErr bool
	}{
		{
			name:    "nil range returns content as-is",
			content: "line0\nline1\nline2",
			r:       nil,
			want:    "line0\nline1\nline2",
		},
		{
			name:    "normal range",
			content: "line0\nline1\nline2\nline3",
			r:       &Range{Start: new(1), End: new(3)},
			want:    "line1\nline2",
		},
		{
			name:    "default start",
			content: "line0\nline1\nline2",
			r:       &Range{End: new(2)},
			want:    "line0\nline1",
		},
		{
			name:    "default end",
			content: "line0\nline1\nline2",
			r:       &Range{Start: new(1)},
			want:    "line1\nline2",
		},
		{
			name:    "negative start",
			content: "line0\nline1\nline2\nline3",
			r:       &Range{Start: new(-2)},
			want:    "line2\nline3",
		},
		{
			name:    "negative end",
			content: "line0\nline1\nline2\nline3",
			r:       &Range{End: new(-1)},
			want:    "line0\nline1\nline2",
		},
		{
			name:    "negative start and end",
			content: "line0\nline1\nline2\nline3",
			r:       &Range{Start: new(-3), End: new(-1)},
			want:    "line1\nline2",
		},
		{
			name:    "start equals end",
			content: "line0\nline1",
			r:       &Range{Start: new(1), End: new(1)},
			wantErr: true,
		},
		{
			name:    "start greater than end",
			content: "line0\nline1",
			r:       &Range{Start: new(2), End: new(1)},
			wantErr: true,
		},
		{
			name:    "end exceeds line count",
			content: "line0\nline1",
			r:       &Range{End: new(5)},
			wantErr: true,
		},
		{
			name:    "negative start out of bounds",
			content: "line0\nline1",
			r:       &Range{Start: new(-5)},
			wantErr: true,
		},
		{
			name:    "trailing newline preserved",
			content: "line0\nline1\nline2\n",
			r:       &Range{Start: new(-2)},
			want:    "line1\nline2\n",
		},
		{
			name:    "trailing newline with normal range",
			content: "line0\nline1\nline2\n",
			r:       &Range{Start: new(0), End: new(2)},
			want:    "line0\nline1\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ExtractRange(tt.content, tt.r)
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
