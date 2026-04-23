package requirementdiagram

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const RequirementDiagramKeyword = "requirementDiagram"

var (
	blockStartRegex   = regexp.MustCompile(`^\s*(requirement|functionalRequirement|interfaceRequirement|performanceRequirement|physicalRequirement|designConstraint|element)\s+([A-Za-z_][A-Za-z0-9_-]*)\s*\{\s*$`)
	fieldRegex        = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_-]*)\s*:\s*(.*)$`)
	relationshipRegex = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_-]*)\s*-\s*(contains|copies|derives|satisfies|verifies|refines|traces)\s*->\s*([A-Za-z_][A-Za-z0-9_-]*)\s*$`)
)

type Diagram struct {
	Items         map[string]*Item
	ItemOrder     []string
	Relationships []*Relationship
}

type Item struct {
	ID     string
	Kind   string
	Fields map[string]string
}

type Relationship struct {
	From *Item
	To   *Item
	Kind string
}

func IsRequirementDiagram(input string) bool {
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, RequirementDiagramKeyword)
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
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), RequirementDiagramKeyword) {
		return nil, fmt.Errorf("expected %q keyword", RequirementDiagramKeyword)
	}

	rd := &Diagram{Items: map[string]*Item{}, ItemOrder: []string{}}
	var current *Item
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if current != nil {
			if trimmed == "}" {
				current = nil
				continue
			}
			match := fieldRegex.FindStringSubmatch(trimmed)
			if match == nil {
				return nil, fmt.Errorf("line %d: invalid requirement field syntax: %q", lineNo, trimmed)
			}
			current.Fields[strings.ToLower(match[1])] = strings.TrimSpace(match[2])
			continue
		}
		if match := blockStartRegex.FindStringSubmatch(trimmed); match != nil {
			current = rd.ensureItem(match[2], match[1])
			continue
		}
		if match := relationshipRegex.FindStringSubmatch(trimmed); match != nil {
			from := rd.ensureItem(match[1], "element")
			to := rd.ensureItem(match[3], "requirement")
			rd.Relationships = append(rd.Relationships, &Relationship{From: from, To: to, Kind: match[2]})
			continue
		}
		if strings.Contains(trimmed, "->") {
			return nil, fmt.Errorf("line %d: invalid requirement relationship syntax: %q", lineNo, trimmed)
		}
		return nil, fmt.Errorf("line %d: invalid requirement syntax: %q", lineNo, trimmed)
	}
	if current != nil {
		return nil, fmt.Errorf("unterminated %s block for %q", current.Kind, current.ID)
	}
	if len(rd.Items) == 0 {
		return nil, fmt.Errorf("no requirements or elements found")
	}
	return rd, nil
}

func (rd *Diagram) ensureItem(id, kind string) *Item {
	if item, ok := rd.Items[id]; ok {
		if item.Kind == "" || item.Kind == "element" && kind != "" {
			item.Kind = kind
		}
		return item
	}
	item := &Item{ID: id, Kind: kind, Fields: map[string]string{}}
	rd.Items[id] = item
	rd.ItemOrder = append(rd.ItemOrder, id)
	return item
}
