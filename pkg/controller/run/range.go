package run

import (
	"fmt"

	"github.com/suzuki-shunsuke/docfresh/pkg/textrange"
)

type Range = textrange.Range

func extractRange(content string, r *Range) (string, error) {
	s, err := textrange.ExtractRange(content, r)
	if err != nil {
		return "", fmt.Errorf("extract range: %w", err)
	}
	return s, nil
}
