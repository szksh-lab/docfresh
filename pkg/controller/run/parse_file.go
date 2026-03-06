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
	codeBlocks = append(codeBlocks, findInlineCodeRanges(content, codeBlocks)...)
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

// countLeadingBackticks returns the number of leading backtick characters in s.
func countLeadingBackticks(s string) int {
	n := 0
	for _, c := range s {
		if c != '`' {
			break
		}
		n++
	}
	return n
}

// isClosingFence reports whether trimmed is a valid closing fence
// with at least minBackticks backticks.
func isClosingFence(trimmed string, minBackticks int) bool {
	n := countLeadingBackticks(trimmed)
	if n < minBackticks {
		return false
	}
	// Closing fence must contain only backticks and optional whitespace.
	return strings.TrimRight(trimmed[n:], " \t") == ""
}

// nextLine returns the line text and the absolute end-of-line position.
// lineEnd points to the '\n' or len(content) if there is no trailing newline.
func nextLine(content string, pos int) (string, int) {
	idx := strings.Index(content[pos:], "\n")
	if idx == -1 {
		return content[pos:], len(content)
	}
	return content[pos : pos+idx], pos + idx
}

// findCodeBlockRanges returns byte-offset ranges [start, end) for fenced code blocks.
func findCodeBlockRanges(content string) [][2]int {
	var ranges [][2]int
	pos := 0
	inCodeBlock := false
	openBackticks := 0
	blockStart := 0

	for pos < len(content) {
		line, lineEnd := nextLine(content, pos)
		trimmed := strings.TrimLeft(line, " \t")

		if strings.HasPrefix(trimmed, "```") {
			ranges, inCodeBlock, openBackticks, blockStart = handleFenceLine(
				trimmed, pos, lineEnd, len(content),
				ranges, inCodeBlock, openBackticks, blockStart,
			)
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

func handleFenceLine(
	trimmed string, lineStart, lineEnd, contentLen int,
	ranges [][2]int, inCodeBlock bool, openBackticks, blockStart int,
) ([][2]int, bool, int, int) {
	backtickCount := countLeadingBackticks(trimmed)

	if !inCodeBlock {
		return ranges, true, backtickCount, lineStart
	}

	if !isClosingFence(trimmed, openBackticks) {
		return ranges, inCodeBlock, openBackticks, blockStart
	}

	end := lineEnd
	if end < contentLen {
		end++ // include \n
	}
	ranges = append(ranges, [2]int{blockStart, end})
	return ranges, false, 0, 0
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

// findInlineCodeRanges returns byte-offset ranges [start, end) for inline code spans,
// skipping positions that are already inside fenced code block ranges.
func findInlineCodeRanges(content string, fencedRanges [][2]int) [][2]int {
	var ranges [][2]int
	i := 0
	for i < len(content) {
		if insideCodeBlock(i, fencedRanges) {
			i++
			continue
		}
		if content[i] != '`' {
			i++
			continue
		}
		// Count the opening backtick sequence length.
		openLen := 0
		for i+openLen < len(content) && content[i+openLen] == '`' {
			openLen++
		}
		start := i
		// Search for a closing backtick sequence of exactly openLen.
		j := i + openLen
		found := false
		for j < len(content) {
			if content[j] != '`' {
				j++
				continue
			}
			closeLen := 0
			for j+closeLen < len(content) && content[j+closeLen] == '`' {
				closeLen++
			}
			if closeLen == openLen {
				ranges = append(ranges, [2]int{start, j + closeLen})
				i = j + closeLen
				found = true
				break
			}
			j += closeLen
		}
		if !found {
			// No closing sequence found; not a code span.
			i += openLen
		}
	}
	return ranges
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
