package gitgraph

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const Keyword = "gitGraph"

type Diagram struct {
	Branches []string
	Events   []Event
}

type Event struct {
	Branch string
	Kind   string
	Label  string
}

func IsGitGraph(input string) bool {
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
	gg := &Diagram{Branches: []string{"main"}}
	current := "main"
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(strings.TrimSuffix(line, ":"))
		if trimmed == "" {
			continue
		}
		fields := strings.Fields(trimmed)
		switch fields[0] {
		case "commit":
			label := parseOptionLabel(strings.TrimSpace(strings.TrimPrefix(trimmed, "commit")))
			if label == "" {
				label = "commit"
			}
			gg.Events = append(gg.Events, Event{Branch: current, Kind: "commit", Label: label})
		case "branch":
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: branch requires a name", lineNo)
			}
			branch := fields[1]
			gg.addBranch(branch)
			gg.Events = append(gg.Events, Event{Branch: branch, Kind: "branch", Label: "branch " + branch})
		case "checkout":
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: checkout requires a branch", lineNo)
			}
			current = fields[1]
			gg.addBranch(current)
			gg.Events = append(gg.Events, Event{Branch: current, Kind: "checkout", Label: "checkout " + current})
		case "merge":
			if len(fields) < 2 {
				return nil, fmt.Errorf("line %d: merge requires a branch", lineNo)
			}
			gg.addBranch(fields[1])
			gg.Events = append(gg.Events, Event{Branch: current, Kind: "merge", Label: "merge " + fields[1]})
		default:
			return nil, fmt.Errorf("line %d: unsupported gitGraph command %q", lineNo, fields[0])
		}
	}
	if len(gg.Events) == 0 {
		return nil, fmt.Errorf("no gitgraph events found")
	}
	return gg, nil
}

func (gg *Diagram) addBranch(branch string) {
	for _, existing := range gg.Branches {
		if existing == branch {
			return
		}
	}
	gg.Branches = append(gg.Branches, branch)
}

func parseOptionLabel(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	for _, prefix := range []string{"id:", "tag:"} {
		if idx := strings.Index(raw, prefix); idx != -1 {
			value := strings.TrimSpace(raw[idx+len(prefix):])
			return strings.Trim(value, `"`)
		}
	}
	return strings.Trim(raw, `"`)
}

func Render(gg *Diagram, config *diagram.Config) (string, error) {
	if gg == nil || len(gg.Events) == 0 {
		return "", fmt.Errorf("no gitgraph events")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	commit, line := "●", "│"
	if config.UseAscii {
		commit, line = "*", "|"
	}
	width := 0
	for _, branch := range gg.Branches {
		if len(branch) > width {
			width = len(branch)
		}
	}
	lines := []string{}
	for _, event := range gg.Events {
		marker := line
		if event.Kind == "commit" {
			marker = commit
		}
		lines = append(lines, fmt.Sprintf("%-*s %s %s", width, event.Branch, marker, event.Label))
	}
	return strings.Join(lines, "\n") + "\n", nil
}
