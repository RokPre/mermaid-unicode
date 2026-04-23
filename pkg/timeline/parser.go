package timeline

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const TimelineKeyword = "timeline"

type Diagram struct {
	Title     string
	Direction string
	Items     []Item
}

type Item struct {
	Section string
	Period  string
	Events  []string
}

func IsTimeline(input string) bool {
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, TimelineKeyword)
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
	first := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(first, TimelineKeyword) {
		return nil, fmt.Errorf("expected %q keyword", TimelineKeyword)
	}
	tl := &Diagram{Direction: "TD"}
	if fields := strings.Fields(first); len(fields) > 1 {
		tl.Direction = strings.ToUpper(fields[1])
	}

	currentSection := ""
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		switch {
		case strings.HasPrefix(lower, "title "):
			tl.Title = strings.TrimSpace(trimmed[len("title "):])
		case strings.HasPrefix(lower, "section "):
			currentSection = strings.TrimSpace(trimmed[len("section "):])
		case strings.HasPrefix(lower, "direction "):
			tl.Direction = strings.ToUpper(strings.TrimSpace(trimmed[len("direction "):]))
		default:
			parts := strings.Split(trimmed, ":")
			if len(parts) < 2 {
				return nil, fmt.Errorf("line %d: invalid timeline item syntax: %q", lineNo, trimmed)
			}
			item := Item{Section: currentSection, Period: strings.TrimSpace(parts[0])}
			for _, event := range parts[1:] {
				event = strings.TrimSpace(event)
				if event != "" {
					item.Events = append(item.Events, event)
				}
			}
			if item.Period == "" || len(item.Events) == 0 {
				return nil, fmt.Errorf("line %d: invalid timeline item syntax: %q", lineNo, trimmed)
			}
			tl.Items = append(tl.Items, item)
		}
	}
	if len(tl.Items) == 0 {
		return nil, fmt.Errorf("no timeline items found")
	}
	return tl, nil
}
