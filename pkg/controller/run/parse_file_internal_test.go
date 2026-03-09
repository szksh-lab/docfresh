package run

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFile(t *testing.T) { //nolint:funlen,maintidx
	t.Parallel()

	tests := []struct {
		name    string
		content string
		want    []*Block
		wantErr string
	}{
		{
			name:    "text only",
			content: "Hello world\nThis is plain text.\n",
			want: []*Block{
				{Type: "text", Content: "Hello world\nThis is plain text.\n"},
			},
		},
		{
			name:    "empty content",
			content: "",
			want:    nil,
		},
		{
			name:    "single begin/end pair with surrounding text",
			content: "before\n<!-- docfresh begin\ncommand:\n  command: echo hello\n-->\nold output\n<!-- docfresh end -->\nafter\n", //nolint:dupword
			want: []*Block{
				{Type: "text", Content: "before\n"},
				{
					Type:         "block",
					BeginComment: "<!-- docfresh begin\ncommand:\n  command: echo hello\n-->", //nolint:dupword
					EndComment:   "<!-- docfresh end -->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo hello",
						},
					},
					LineNumber: 2,
				},
				{Type: "text", Content: "\nafter\n"},
			},
		},
		{
			name:    "multiple begin/end pairs",
			content: "# Title\n<!-- docfresh begin\ncommand:\n  command: echo one\n-->\noutput1\n<!-- docfresh end -->\nmiddle\n<!-- docfresh begin\ncommand:\n  command: echo two\n  shell:\n    - sh\n    - -c\n-->\noutput2\n<!-- docfresh end -->\n", //nolint:dupword
			want: []*Block{
				{Type: "text", Content: "# Title\n"},
				{
					Type:         "block",
					BeginComment: "<!-- docfresh begin\ncommand:\n  command: echo one\n-->", //nolint:dupword
					EndComment:   "<!-- docfresh end -->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo one",
						},
					},
					LineNumber: 2,
				},
				{Type: "text", Content: "\nmiddle\n"},
				{
					Type:         "block",
					BeginComment: "<!-- docfresh begin\ncommand:\n  command: echo two\n  shell:\n    - sh\n    - -c\n-->", //nolint:dupword
					EndComment:   "<!-- docfresh end -->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo two",
							Shell:   []string{"sh", "-c"},
						},
					},
					LineNumber: 9,
				},
				{Type: "text", Content: "\n"},
			},
		},
		{
			name:    "end before begin",
			content: "text\n<!-- docfresh end -->\n<!-- docfresh begin\ncommand:\n  command: echo hello\n-->\n", //nolint:dupword
			wantErr: "found <!-- docfresh end --> without a matching <!-- docfresh begin",
		},
		{
			name:    "missing end",
			content: "<!-- docfresh begin\ncommand:\n  command: echo hello\n-->\nsome content\n", //nolint:dupword
			wantErr: "missing <!-- docfresh end --> for begin comment",
		},
		{
			name:    "nested begin",
			content: "<!-- docfresh begin\ncommand:\n  command: echo one\n-->\n<!-- docfresh begin\ncommand:\n  command: echo two\n-->\n<!-- docfresh end -->\n", //nolint:dupword
			wantErr: "nested <!-- docfresh begin found before <!-- docfresh end -->",
		},
		{
			name:    "unclosed begin comment",
			content: "<!-- docfresh begin\ncommand:\n  command: echo hello\n", //nolint:dupword
			wantErr: "unclosed <!-- docfresh begin comment: missing -->",
		},
		{
			name:    "markers inside fenced code block are ignored",
			content: "before\n```\n<!-- docfresh begin\ncommand:\n  command: echo hello\n-->\nsome output\n<!-- docfresh end -->\n```\nafter\n", //nolint:dupword
			want: []*Block{
				{Type: "text", Content: "before\n```\n<!-- docfresh begin\ncommand:\n  command: echo hello\n-->\nsome output\n<!-- docfresh end -->\n```\nafter\n"}, //nolint:dupword
			},
		},
		{
			name:    "markers outside code block still work",
			content: "```\ncode block\n```\n<!-- docfresh begin\ncommand:\n  command: echo hi\n-->\nold\n<!-- docfresh end -->\n", //nolint:dupword
			want: []*Block{
				{Type: "text", Content: "```\ncode block\n```\n"},
				{
					Type:         "block",
					BeginComment: "<!-- docfresh begin\ncommand:\n  command: echo hi\n-->", //nolint:dupword
					EndComment:   "<!-- docfresh end -->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo hi",
						},
					},
					LineNumber: 4,
				},
				{Type: "text", Content: "\n"},
			},
		},
		{
			name:    "markers inside inline code are ignored",
			content: "before `<!-- docfresh begin -->` and `<!-- docfresh end -->` after\n",
			want: []*Block{
				{Type: "text", Content: "before `<!-- docfresh begin -->` and `<!-- docfresh end -->` after\n"},
			},
		},
		{
			name:    "markers inside double-backtick inline code are ignored",
			content: "before `` <!-- docfresh end --> `` after\n",
			want: []*Block{
				{Type: "text", Content: "before `` <!-- docfresh end --> `` after\n"},
			},
		},
		{
			name:    "inline code with markers plus real markers outside",
			content: "See `<!-- docfresh begin -->` for syntax.\n<!-- docfresh begin\ncommand:\n  command: echo real\n-->\nold\n<!-- docfresh end -->\n", //nolint:dupword
			want: []*Block{
				{Type: "text", Content: "See `<!-- docfresh begin -->` for syntax.\n"},
				{
					Type:         "block",
					BeginComment: "<!-- docfresh begin\ncommand:\n  command: echo real\n-->", //nolint:dupword
					EndComment:   "<!-- docfresh end -->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo real",
						},
					},
					LineNumber: 2,
				},
				{Type: "text", Content: "\n"},
			},
		},
		{
			name:    "mixed: code block with markers and real markers outside",
			content: "# Doc\n```markdown\n<!-- docfresh begin\ncommand:\n  command: echo fake\n-->\nfake\n<!-- docfresh end -->\n```\n<!-- docfresh begin\ncommand:\n  command: echo real\n-->\nold output\n<!-- docfresh end -->\n", //nolint:dupword
			want: []*Block{
				{Type: "text", Content: "# Doc\n```markdown\n<!-- docfresh begin\ncommand:\n  command: echo fake\n-->\nfake\n<!-- docfresh end -->\n```\n"}, //nolint:dupword
				{
					Type:         "block",
					BeginComment: "<!-- docfresh begin\ncommand:\n  command: echo real\n-->", //nolint:dupword
					EndComment:   "<!-- docfresh end -->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo real",
						},
					},
					LineNumber: 10,
				},
				{Type: "text", Content: "\n"},
			},
		},
		{
			name:    "single post comment",
			content: "before\n<!-- docfresh post\ncommand: echo cleanup\n-->\nafter\n",
			want: []*Block{
				{Type: "text", Content: "before\n"},
				{
					Type:    "post",
					Content: "<!-- docfresh post\ncommand: echo cleanup\n-->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo cleanup",
						},
					},
					LineNumber: 2,
				},
				{Type: "text", Content: "\nafter\n"},
			},
		},
		{
			name:    "multiple post comments",
			content: "<!-- docfresh post\ncommand: echo first\n-->\n<!-- docfresh post\ncommand: echo second\n-->\n",
			want: []*Block{
				{
					Type:    "post",
					Content: "<!-- docfresh post\ncommand: echo first\n-->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo first",
						},
					},
					LineNumber: 1,
				},
				{Type: "text", Content: "\n"},
				{
					Type:    "post",
					Content: "<!-- docfresh post\ncommand: echo second\n-->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo second",
						},
					},
					LineNumber: 4,
				},
				{Type: "text", Content: "\n"},
			},
		},
		{
			name:    "post comment mixed with begin/end blocks",
			content: "<!-- docfresh post\ncommand: echo cleanup\n-->\n<!-- docfresh begin\ncommand:\n  command: echo hello\n-->\nold\n<!-- docfresh end -->\n", //nolint:dupword
			want: []*Block{
				{
					Type:    "post",
					Content: "<!-- docfresh post\ncommand: echo cleanup\n-->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo cleanup",
						},
					},
					LineNumber: 1,
				},
				{Type: "text", Content: "\n"},
				{
					Type:         "block",
					BeginComment: "<!-- docfresh begin\ncommand:\n  command: echo hello\n-->", //nolint:dupword
					EndComment:   "<!-- docfresh end -->",
					Input: &BlockInput{
						Command: &Command{
							Command: "echo hello",
						},
					},
					LineNumber: 4,
				},
				{Type: "text", Content: "\n"},
			},
		},
		{
			name:    "post comment inside code block is ignored",
			content: "```\n<!-- docfresh post\ncommand: echo cleanup\n-->\n```\nafter\n",
			want: []*Block{
				{Type: "text", Content: "```\n<!-- docfresh post\ncommand: echo cleanup\n-->\n```\nafter\n"},
			},
		},
		{
			name:    "unclosed post comment",
			content: "<!-- docfresh post\ncommand: echo cleanup\n",
			wantErr: "unclosed <!-- docfresh post comment: missing -->",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := ParseFile(tt.content, nil)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatal("expected error but got nil")
				}
				if err.Error() != tt.wantErr {
					t.Fatalf("error mismatch:\n got: %s\nwant: %s", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("blocks mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
