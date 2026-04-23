package mindmap

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func Render(mm *Diagram, config *diagram.Config) (string, error) {
	if mm == nil || mm.Root == nil {
		return "", fmt.Errorf("no mindmap nodes")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	lines := []string{mm.Root.Text}
	for i, child := range mm.Root.Children {
		last := i == len(mm.Root.Children)-1
		renderNode(child, "", last, config.UseAscii, &lines)
	}
	return strings.Join(lines, "\n") + "\n", nil
}

func renderNode(node *Node, prefix string, last bool, useASCII bool, lines *[]string) {
	branch, nextPrefix := "└── ", prefix+"    "
	if !last {
		branch, nextPrefix = "├── ", prefix+"│   "
	}
	if useASCII {
		branch, nextPrefix = "`-- ", prefix+"    "
		if !last {
			branch, nextPrefix = "|-- ", prefix+"|   "
		}
	}
	*lines = append(*lines, prefix+branch+node.Text)
	for i, child := range node.Children {
		renderNode(child, nextPrefix, i == len(node.Children)-1, useASCII, lines)
	}
}
