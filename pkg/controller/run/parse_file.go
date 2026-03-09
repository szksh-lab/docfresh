package run

import (
	"errors"
	"fmt"
	"strings"

	"github.com/goccy/go-yaml"
)

// ParseOption controls the behavior of ParseFile.
type ParseOption struct {
	DisallowUnknownField bool
}

func unmarshalYAML(yamlStr string, v any, opt *ParseOption) error {
	if opt != nil && opt.DisallowUnknownField {
		dec := yaml.NewDecoder(strings.NewReader(yamlStr), yaml.DisallowUnknownField())
		return dec.Decode(v) //nolint:wrapcheck
	}
	return yaml.Unmarshal([]byte(yamlStr), v) //nolint:wrapcheck
}

const (
	beginMarker     = "<!-- docfresh begin"
	endMarker       = "<!-- docfresh end -->"
	postMarker      = "<!-- docfresh post"
	containerMarker = "<!-- docfresh container"
)

// lineNumber returns the 1-based line number for the given byte position in content.
func lineNumber(content string, pos int) int {
	return strings.Count(content[:pos], "\n") + 1
}

// ParseFile parses a file and returns a list of blocks.
func ParseFile(content string, opt *ParseOption) ([]*Block, error) { //nolint:cyclop,funlen,gocognit,gocyclo
	codeBlocks := findCodeBlockRanges(content)
	codeBlocks = append(codeBlocks, findInlineCodeRanges(content, codeBlocks)...)
	var blocks []*Block
	pos := 0
	for pos < len(content) {
		beginIdx := indexOutsideCodeBlocks(content, beginMarker, pos, codeBlocks)
		endIdx := indexOutsideCodeBlocks(content, endMarker, pos, codeBlocks)
		postIdx := indexOutsideCodeBlocks(content, postMarker, pos, codeBlocks)
		containerIdx := indexOutsideCodeBlocks(content, containerMarker, pos, codeBlocks)

		// No more markers — emit remaining text and break.
		if beginIdx == -1 && endIdx == -1 && postIdx == -1 && containerIdx == -1 {
			blocks = appendText(blocks, content[pos:])
			break
		}

		// Check if container marker comes first.
		if containerIdx != -1 && (beginIdx == -1 || containerIdx < beginIdx) && (endIdx == -1 || containerIdx < endIdx) && (postIdx == -1 || containerIdx < postIdx) {
			block, newPos, err := parseContainerBlock(content, pos, containerIdx, opt)
			if err != nil {
				return nil, err
			}
			block.LineNumber = lineNumber(content, pos+containerIdx)
			if containerIdx > 0 {
				blocks = appendText(blocks, content[pos:pos+containerIdx])
			}
			blocks = append(blocks, block)
			pos = newPos
			continue
		}

		// Check if post marker comes first.
		if postIdx != -1 && (beginIdx == -1 || postIdx < beginIdx) && (endIdx == -1 || postIdx < endIdx) {
			block, newPos, err := parsePostBlock(content, pos, postIdx, opt)
			if err != nil {
				return nil, err
			}
			block.LineNumber = lineNumber(content, pos+postIdx)
			if postIdx > 0 {
				blocks = appendText(blocks, content[pos:pos+postIdx])
			}
			blocks = append(blocks, block)
			pos = newPos
			continue
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
		if err := unmarshalYAML(yamlStr, &input, opt); err != nil {
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
			Type:         blockTypeBlock,
			Input:        &input,
			BeginComment: beginComment,
			EndComment:   endComment,
			LineNumber:   lineNumber(content, beginStart),
		})

		pos = endCommentEnd
	}
	return blocks, nil
}

func parsePostBlock(content string, pos, postIdx int, opt *ParseOption) (*Block, int, error) {
	postStart := pos + postIdx
	closeIdx := strings.Index(content[postStart+len(postMarker):], "-->")
	if closeIdx == -1 {
		return nil, 0, errors.New("unclosed <!-- docfresh post comment: missing -->")
	}
	postCommentEnd := postStart + len(postMarker) + closeIdx + len("-->")
	postComment := content[postStart:postCommentEnd]

	yamlStr := content[postStart+len(postMarker) : postStart+len(postMarker)+closeIdx]
	yamlStr = strings.TrimSpace(yamlStr)
	cmd := &PostCommand{}
	if err := unmarshalYAML(yamlStr, cmd, opt); err != nil {
		return nil, 0, fmt.Errorf("failed to parse YAML in post comment: %w", err)
	}

	return &Block{
		Type:    blockTypePost,
		Content: postComment,
		Input: &BlockInput{
			Command: cmd.ToCommand(),
		},
	}, postCommentEnd, nil
}

func parseContainerBlock(content string, pos, containerIdx int, opt *ParseOption) (*Block, int, error) {
	containerStart := pos + containerIdx
	closeIdx := strings.Index(content[containerStart+len(containerMarker):], "-->")
	if closeIdx == -1 {
		return nil, 0, errors.New("unclosed <!-- docfresh container comment: missing -->")
	}
	containerCommentEnd := containerStart + len(containerMarker) + closeIdx + len("-->")
	containerComment := content[containerStart:containerCommentEnd]

	yamlStr := content[containerStart+len(containerMarker) : containerStart+len(containerMarker)+closeIdx]
	yamlStr = strings.TrimSpace(yamlStr)
	input := &ContainerInput{}
	if err := unmarshalYAML(yamlStr, input, opt); err != nil {
		return nil, 0, fmt.Errorf("failed to parse YAML in container comment: %w", err)
	}

	return &Block{
		Type:           blockTypeContainer,
		Content:        containerComment,
		ContainerInput: input,
	}, containerCommentEnd, nil
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

// findClosingBackticks scans content starting at pos for a backtick run of exactly openLen.
// Returns the end position (after the closing backticks) or -1 if not found.
func findClosingBackticks(content string, pos, openLen int) int {
	j := pos
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
			return j + closeLen
		}
		j += closeLen
	}
	return -1
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
		end := findClosingBackticks(content, i+openLen, openLen)
		if end >= 0 {
			ranges = append(ranges, [2]int{i, end})
			i = end
		} else {
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
		Type:    blockTypeText,
		Content: text,
	})
}
