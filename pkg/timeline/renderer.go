package timeline

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func Render(tl *Diagram, config *diagram.Config) (string, error) {
	if tl == nil || len(tl.Items) == 0 {
		return "", fmt.Errorf("no timeline items")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	if strings.EqualFold(tl.Direction, "LR") {
		return renderHorizontal(tl, config.UseAscii), nil
	}
	return renderVertical(tl, config.UseAscii), nil
}

func renderVertical(tl *Diagram, useASCII bool) string {
	lines := []string{}
	if tl.Title != "" {
		lines = append(lines, tl.Title)
	}
	currentSection := ""
	for _, item := range tl.Items {
		if item.Section != currentSection {
			currentSection = item.Section
			if currentSection != "" {
				lines = append(lines, "["+currentSection+"]")
			}
		}
		head := "├─"
		child := "│ "
		if useASCII {
			head = "|-"
			child = "| "
		}
		lines = append(lines, fmt.Sprintf("%s %s: %s", head, item.Period, item.Events[0]))
		for _, event := range item.Events[1:] {
			lines = append(lines, fmt.Sprintf("%s   %s", child, event))
		}
	}
	return strings.Join(lines, "\n") + "\n"
}

func renderHorizontal(tl *Diagram, useASCII bool) string {
	sep := " ── "
	if useASCII {
		sep = " -- "
	}
	parts := []string{}
	for _, item := range tl.Items {
		label := item.Period + ": " + strings.Join(item.Events, " / ")
		if item.Section != "" {
			label = "[" + item.Section + "] " + label
		}
		parts = append(parts, label)
	}
	if tl.Title != "" {
		return tl.Title + "\n" + strings.Join(parts, sep) + "\n"
	}
	return strings.Join(parts, sep) + "\n"
}
