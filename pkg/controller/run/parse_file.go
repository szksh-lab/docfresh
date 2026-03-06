package run

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

const (
	beginMarker = "<!-- docfresh begin"
	endMarker   = "<!-- docfresh end -->"
)

// parseFile parses a file and returns a list of blocks.
func parseFile(content string) ([]*Block, error) { //nolint:cyclop
	codeBlocks := findCodeBlockRanges(content)
	var blocks []*Block
	pos := 0
	for pos < len(content) {
		beginIdx := indexOutsideCodeBlocks(content, beginMarker, pos, codeBlocks)
		endIdx := indexOutsideCodeBlocks(content, endMarker, pos, codeBlocks)

		// No more markers — emit remaining text and break.
		if beginIdx == -1 && endIdx == -1 {
			blocks = appendText(blocks, content[pos:])
			break
		}

		// end before begin (or end without begin).
		if beginIdx == -1 || (endIdx != -1 && endIdx < beginIdx) {
			return nil, errors.New("found <!-- docfresh end --> without a matching <!-- docfresh begin")
		}

		// Emit text before the begin marker.
		if beginIdx > 0 {
			blocks = appendText(blocks, content[pos:pos+beginIdx])
		}

		// Find closing --> of the begin comment.
		beginStart := pos + beginIdx
		closeIdx := strings.Index(content[beginStart+len(beginMarker):], "-->")
		if closeIdx == -1 {
			return nil, errors.New("unclosed <!-- docfresh begin comment: missing -->")
		}
		beginCommentEnd := beginStart + len(beginMarker) + closeIdx + len("-->")
		beginComment := content[beginStart:beginCommentEnd]

		// Extract YAML from inside the begin comment.
		yamlStr := content[beginStart+len(beginMarker) : beginStart+len(beginMarker)+closeIdx]
		yamlStr = strings.TrimSpace(yamlStr)
		var input BlockInput
		if err := yaml.Unmarshal([]byte(yamlStr), &input); err != nil {
			return nil, fmt.Errorf("failed to parse YAML in begin comment: %w", err)
		}

		// Find matching end marker after the begin comment.
		endIdx = indexOutsideCodeBlocks(content, endMarker, beginCommentEnd, codeBlocks)
		if endIdx == -1 {
			return nil, fmt.Errorf("missing %s for begin comment", endMarker)
		}

		// Check for nested begin markers between this begin and the end.
		nestedIdx := indexOutsideCodeBlocks(content, beginMarker, beginCommentEnd, codeBlocks)
		if nestedIdx != -1 && nestedIdx < endIdx {
			return nil, errors.New("nested <!-- docfresh begin found before <!-- docfresh end -->")
		}

		endCommentEnd := beginCommentEnd + endIdx + len(endMarker)
		endComment := content[beginCommentEnd+endIdx : endCommentEnd]

		blocks = append(blocks, &Block{
			Type:         "block",
			Input:        &input,
			BeginComment: beginComment,
			EndComment:   endComment,
		})

		pos = endCommentEnd
	}
	return blocks, nil
}

// findCodeBlockRanges returns byte-offset ranges [start, end) for fenced code blocks.
func findCodeBlockRanges(content string) [][2]int {
	var ranges [][2]int
	pos := 0
	inCodeBlock := false
	openBackticks := 0
	blockStart := 0

	for pos < len(content) {
		lineEnd := strings.Index(content[pos:], "\n")
		var line string
		if lineEnd == -1 {
			line = content[pos:]
			lineEnd = len(content)
		} else {
			lineEnd += pos
			line = content[pos:lineEnd]
		}

		trimmed := strings.TrimLeft(line, " \t")
		if strings.HasPrefix(trimmed, "```") {
			backtickCount := 0
			for _, c := range trimmed {
				if c == '`' {
					backtickCount++
				} else {
					break
				}
			}

			if !inCodeBlock {
				inCodeBlock = true
				openBackticks = backtickCount
				blockStart = pos
			} else if backtickCount >= openBackticks {
				// Closing fence must contain only backticks and optional whitespace.
				afterBackticks := strings.TrimRight(trimmed[backtickCount:], " \t")
				if afterBackticks == "" {
					end := lineEnd
					if end < len(content) {
						end++ // include \n
					}
					ranges = append(ranges, [2]int{blockStart, end})
					inCodeBlock = false
				}
			}
		}

		if lineEnd < len(content) {
			pos = lineEnd + 1
		} else {
			break
		}
	}

	if inCodeBlock {
		ranges = append(ranges, [2]int{blockStart, len(content)})
	}

	return ranges
}

// indexOutsideCodeBlocks works like strings.Index(content[start:], substr)
// but skips matches whose absolute position falls within a code block range.
func indexOutsideCodeBlocks(content, substr string, start int, ranges [][2]int) int {
	offset := 0
	s := content[start:]
	for {
		idx := strings.Index(s[offset:], substr)
		if idx == -1 {
			return -1
		}
		absPos := start + offset + idx
		if !insideCodeBlock(absPos, ranges) {
			return offset + idx
		}
		offset += idx + 1
	}
}

func insideCodeBlock(pos int, ranges [][2]int) bool {
	for _, r := range ranges {
		if pos >= r[0] && pos < r[1] {
			return true
		}
	}
	return false
}

func appendText(blocks []*Block, text string) []*Block {
	if text == "" {
		return blocks
	}
	return append(blocks, &Block{
		Type:    "text",
		Content: text,
	})
}
