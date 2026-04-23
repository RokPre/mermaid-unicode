package er

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const ERDiagramKeyword = "erDiagram"

var (
	entityNameRegex   = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_-]*$`)
	entityAliasRegex  = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_-]*)\s*\[\s*"?([^"\]]+)"?\s*\]\s*$`)
	entityBlockRegex  = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_-]*)(?:\s*\[\s*"?([^"\]]+)"?\s*\])?\s*\{\s*$`)
	relationshipRegex = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_-]*)\s+([|o}{]{2})(--|\.\.)([|o}{]{2})\s+([A-Za-z_][A-Za-z0-9_-]*)(?:\s*:\s*(.*))?$`)
	directionRegex    = regexp.MustCompile(`^\s*direction\s+(LR|RL|TB|TD|BT)\s*$`)
)

type Diagram struct {
	Entities      map[string]*Entity
	EntityOrder   []string
	Relationships []*Relationship
	Direction     string
}

type Entity struct {
	ID         string
	Label      string
	Attributes []Attribute
}

type Attribute struct {
	Type string
	Name string
	Keys []string
}

type Relationship struct {
	From             *Entity
	To               *Entity
	LeftCardinality  string
	RightCardinality string
	Identifying      bool
	Label            string
}

func IsERDiagram(input string) bool {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, ERDiagramKeyword)
	}
	return false
}

func Parse(input string) (*Diagram, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	rawLines := diagram.SplitLines(input)
	lines := diagram.RemoveComments(rawLines)
	if len(lines) == 0 {
		return nil, fmt.Errorf("no content found")
	}

	if !strings.HasPrefix(strings.TrimSpace(lines[0]), ERDiagramKeyword) {
		return nil, fmt.Errorf("expected %q keyword", ERDiagramKeyword)
	}

	er := &Diagram{
		Entities:    map[string]*Entity{},
		EntityOrder: []string{},
		Direction:   "LR",
	}

	var currentEntity *Entity
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if currentEntity != nil {
			if trimmed == "}" {
				currentEntity = nil
				continue
			}
			attr, err := parseAttribute(trimmed)
			if err != nil {
				return nil, fmt.Errorf("line %d: %w", lineNo, err)
			}
			currentEntity.Attributes = append(currentEntity.Attributes, attr)
			continue
		}

		if match := directionRegex.FindStringSubmatch(trimmed); match != nil {
			er.Direction = match[1]
			continue
		}

		if match := entityBlockRegex.FindStringSubmatch(trimmed); match != nil {
			currentEntity = er.ensureEntity(match[1], match[2])
			continue
		}

		if match := relationshipRegex.FindStringSubmatch(trimmed); match != nil {
			er.addRelationship(match)
			continue
		}

		if strings.Contains(trimmed, "--") || strings.Contains(trimmed, "..") {
			return nil, fmt.Errorf("line %d: invalid ER relationship syntax: %q", lineNo, trimmed)
		}

		if match := entityAliasRegex.FindStringSubmatch(trimmed); match != nil {
			er.ensureEntity(match[1], match[2])
			continue
		}

		if entityNameRegex.MatchString(trimmed) {
			er.ensureEntity(trimmed, "")
			continue
		}

		return nil, fmt.Errorf("line %d: invalid ER syntax: %q", lineNo, trimmed)
	}

	if currentEntity != nil {
		return nil, fmt.Errorf("unterminated entity block for %q", currentEntity.ID)
	}
	if len(er.Entities) == 0 {
		return nil, fmt.Errorf("no entities found")
	}

	return er, nil
}

func (er *Diagram) addRelationship(match []string) {
	left := er.ensureEntity(match[1], "")
	right := er.ensureEntity(match[5], "")
	relationship := &Relationship{
		From:             left,
		To:               right,
		LeftCardinality:  match[2],
		RightCardinality: match[4],
		Identifying:      match[3] == "--",
		Label:            strings.TrimSpace(match[6]),
	}
	er.Relationships = append(er.Relationships, relationship)
}

func (er *Diagram) ensureEntity(id, label string) *Entity {
	if entity, ok := er.Entities[id]; ok {
		if label != "" {
			entity.Label = label
		}
		return entity
	}
	if label == "" {
		label = id
	}
	entity := &Entity{ID: id, Label: label}
	er.Entities[id] = entity
	er.EntityOrder = append(er.EntityOrder, id)
	return entity
}

func parseAttribute(line string) (Attribute, error) {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return Attribute{}, fmt.Errorf("invalid ER attribute syntax: %q", line)
	}

	attr := Attribute{
		Type: parts[0],
		Name: strings.Trim(parts[1], `"`),
	}
	for _, part := range parts[2:] {
		part = strings.Trim(part, ",")
		if part == "" || strings.HasPrefix(part, `"`) {
			break
		}
		attr.Keys = append(attr.Keys, part)
	}
	return attr, nil
}
