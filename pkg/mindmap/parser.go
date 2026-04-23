package mindmap

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const MindmapKeyword = "mindmap"

type Diagram struct {
	Root *Node
}

type Node struct {
	Text     string
	Level    int
	Children []*Node
}

func IsMindmap(input string) bool {
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, MindmapKeyword)
	}
	return false
}

func Parse(input string) (*Diagram, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	lines := diagram.RemoveComments(diagram.SplitLines(input))
	if len(lines) == 0 {
		return nil, fmt.Errorf("no content found")
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), MindmapKeyword) {
		return nil, fmt.Errorf("expected %q keyword", MindmapKeyword)
	}

	var root *Node
	stack := []*Node{}
	indentLevels := []int{}
	for i, line := range lines[1:] {
		lineNo := i + 2
		if strings.TrimSpace(line) == "" {
			continue
		}
		indent := leadingSpaces(line)
		text := normalizeNodeText(strings.TrimSpace(line))
		if text == "" {
			return nil, fmt.Errorf("line %d: empty mindmap node", lineNo)
		}

		if isNewOutdent(indent, indentLevels) {
			return nil, fmt.Errorf("line %d: bad mindmap indentation", lineNo)
		}
		level := levelForIndent(indent, &indentLevels)
		if level > len(stack) {
			return nil, fmt.Errorf("line %d: bad mindmap indentation jump", lineNo)
		}
		node := &Node{Text: text, Level: level}
		if level == 0 {
			if root != nil {
				return nil, fmt.Errorf("line %d: multiple mindmap roots", lineNo)
			}
			root = node
		} else {
			parent := stack[level-1]
			parent.Children = append(parent.Children, node)
		}
		if len(stack) > level {
			stack = stack[:level]
		}
		stack = append(stack, node)
	}
	if root == nil {
		return nil, fmt.Errorf("no mindmap nodes found")
	}
	return &Diagram{Root: root}, nil
}

func leadingSpaces(line string) int {
	count := 0
	for _, r := range line {
		if r != ' ' {
			break
		}
		count++
	}
	return count
}

func isNewOutdent(indent int, levels []int) bool {
	if len(levels) == 0 {
		return false
	}
	for _, known := range levels {
		if known == indent {
			return false
		}
	}
	return indent < levels[len(levels)-1]
}

func levelForIndent(indent int, levels *[]int) int {
	for i, known := range *levels {
		if known == indent {
			return i
		}
	}
	*levels = append(*levels, indent)
	return len(*levels) - 1
}

func normalizeNodeText(text string) string {
	pairs := [][2]string{
		{"((", "))"},
		{"{{", "}}"},
		{"[", "]"},
		{"(", ")"},
		{"{", "}"},
	}
	for _, pair := range pairs {
		if strings.HasPrefix(text, pair[0]) && strings.HasSuffix(text, pair[1]) {
			text = strings.TrimSpace(text[len(pair[0]) : len(text)-len(pair[1])])
			break
		}
	}
	return strings.Trim(text, `"`)
}
