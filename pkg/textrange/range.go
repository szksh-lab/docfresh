package textrange

import (
	"fmt"
	"strings"
)

type Range struct {
	Start *int `json:"start,omitempty" jsonschema_description:"Start line number (0-based, inclusive). Default: 0. Negative values count from the end."`
	End   *int `json:"end,omitempty" jsonschema_description:"End line number (0-based, exclusive). Default: total line count. Negative values count from the end."`
}

func ExtractRange(content string, r *Range) (string, error) {
	if r == nil {
		return content, nil
	}
	// Trim trailing newline before splitting so it doesn't create an extra empty element.
	trailingNewline := strings.HasSuffix(content, "\n")
	trimmed := strings.TrimSuffix(content, "\n")
	lines := strings.Split(trimmed, "\n")
	n := len(lines)
	start, end := 0, n
	if r.Start != nil {
		start = *r.Start
	}
	if r.End != nil {
		end = *r.End
	}
	if start < 0 {
		start = n + start
	}
	if end < 0 {
		end = n + end
	}
	if start < 0 || end > n || start >= end {
		return "", fmt.Errorf("invalid range: start=%d, end=%d, lines=%d", start, end, n)
	}
	result := strings.Join(lines[start:end], "\n")
	if trailingNewline {
		result += "\n"
	}
	return result, nil
}
