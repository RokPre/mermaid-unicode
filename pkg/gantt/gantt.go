package gantt

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const Keyword = "gantt"

type Diagram struct {
	Title      string
	DateFormat string
	Tasks      []Task
}

type Task struct {
	Section string
	Name    string
	Spec    string
}

func IsGantt(input string) bool {
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, Keyword)
	}
	return false
}

func Parse(input string) (*Diagram, error) {
	lines := diagram.RemoveComments(diagram.SplitLines(strings.TrimSpace(input)))
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), Keyword) {
		return nil, fmt.Errorf("expected %q keyword", Keyword)
	}
	g := &Diagram{}
	section := ""
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		switch {
		case strings.HasPrefix(lower, "title "):
			g.Title = strings.TrimSpace(trimmed[len("title "):])
		case strings.HasPrefix(lower, "dateformat "):
			g.DateFormat = strings.TrimSpace(trimmed[len("dateFormat "):])
		case strings.HasPrefix(lower, "section "):
			section = strings.TrimSpace(trimmed[len("section "):])
		case strings.Contains(trimmed, ":"):
			parts := strings.SplitN(trimmed, ":", 2)
			name, spec := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
			if name == "" || spec == "" {
				return nil, fmt.Errorf("line %d: invalid gantt task syntax: %q", lineNo, trimmed)
			}
			g.Tasks = append(g.Tasks, Task{Section: section, Name: name, Spec: spec})
		default:
			return nil, fmt.Errorf("line %d: invalid gantt syntax: %q", lineNo, trimmed)
		}
	}
	if len(g.Tasks) == 0 {
		return nil, fmt.Errorf("no gantt tasks found")
	}
	return g, nil
}

func Render(g *Diagram, config *diagram.Config) (string, error) {
	if g == nil || len(g.Tasks) == 0 {
		return "", fmt.Errorf("no gantt tasks")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	full := "█"
	if config.UseAscii {
		full = "#"
	}
	lines := []string{}
	if g.Title != "" {
		lines = append(lines, g.Title)
	}
	if g.DateFormat != "" {
		lines = append(lines, "dateFormat: "+g.DateFormat)
	}
	section := ""
	for _, task := range g.Tasks {
		if task.Section != section {
			section = task.Section
			if section != "" {
				lines = append(lines, "["+section+"]")
			}
		}
		lines = append(lines, fmt.Sprintf("%-24s %s %s", task.Name, strings.Repeat(full, 8), task.Spec))
	}
	return strings.Join(lines, "\n") + "\n", nil
}
