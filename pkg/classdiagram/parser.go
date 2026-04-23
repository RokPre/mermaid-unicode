package classdiagram

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const ClassDiagramKeyword = "classDiagram"

var (
	classNameRegex     = regexp.MustCompile(`^[A-Za-z_][A-Za-z0-9_.$-]*$`)
	classDeclRegex     = regexp.MustCompile(`^\s*class\s+([A-Za-z_][A-Za-z0-9_.$-]*)(?:\s*\[\s*"?([^"\]]+)"?\s*\])?\s*(\{)?\s*$`)
	colonMemberRegex   = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_.$-]*)\s*:\s*(.+)$`)
	relationshipRegex  = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_.$-]*)(?:\s+"([^"]+)")?\s*(<\|--|\*--|o--|-->|--|\.\.>|\.\.\|>|\.\.)\s*(?:"([^"]+)"\s*)?([A-Za-z_][A-Za-z0-9_.$-]*)(?:\s*:\s*(.*))?$`)
	directionLineRegex = regexp.MustCompile(`^\s*direction\s+(LR|RL|TB|TD|BT)\s*$`)
)

type Diagram struct {
	Classes       map[string]*Class
	ClassOrder    []string
	Relationships []*Relationship
	Direction     string
}

type Class struct {
	ID         string
	Label      string
	Attributes []string
	Operations []string
}

type Relationship struct {
	From             *Class
	To               *Class
	Operator         string
	LeftCardinality  string
	RightCardinality string
	Label            string
}

func IsClassDiagram(input string) bool {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, ClassDiagramKeyword)
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

	if !strings.HasPrefix(strings.TrimSpace(lines[0]), ClassDiagramKeyword) {
		return nil, fmt.Errorf("expected %q keyword", ClassDiagramKeyword)
	}

	cd := &Diagram{
		Classes:    map[string]*Class{},
		ClassOrder: []string{},
		Direction:  "LR",
	}

	var currentClass *Class
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if currentClass != nil {
			if trimmed == "}" {
				currentClass = nil
				continue
			}
			currentClass.addMember(trimmed)
			continue
		}

		if match := directionLineRegex.FindStringSubmatch(trimmed); match != nil {
			cd.Direction = match[1]
			continue
		}

		if match := classDeclRegex.FindStringSubmatch(trimmed); match != nil {
			class := cd.ensureClass(match[1], match[2])
			if match[3] == "{" {
				currentClass = class
			}
			continue
		}

		if match := colonMemberRegex.FindStringSubmatch(trimmed); match != nil {
			cd.ensureClass(match[1], "").addMember(match[2])
			continue
		}

		if match := relationshipRegex.FindStringSubmatch(trimmed); match != nil {
			cd.addRelationship(match)
			continue
		}

		if looksLikeRelationship(trimmed) {
			return nil, fmt.Errorf("line %d: invalid class relationship syntax: %q", lineNo, trimmed)
		}

		if classNameRegex.MatchString(trimmed) {
			cd.ensureClass(trimmed, "")
			continue
		}

		return nil, fmt.Errorf("line %d: invalid class syntax: %q", lineNo, trimmed)
	}

	if currentClass != nil {
		return nil, fmt.Errorf("unterminated class block for %q", currentClass.ID)
	}
	if len(cd.Classes) == 0 {
		return nil, fmt.Errorf("no classes found")
	}

	return cd, nil
}

func (cd *Diagram) addRelationship(match []string) {
	relationship := &Relationship{
		From:             cd.ensureClass(match[1], ""),
		LeftCardinality:  strings.TrimSpace(match[2]),
		Operator:         match[3],
		RightCardinality: strings.TrimSpace(match[4]),
		To:               cd.ensureClass(match[5], ""),
		Label:            strings.TrimSpace(match[6]),
	}
	cd.Relationships = append(cd.Relationships, relationship)
}

func (cd *Diagram) ensureClass(id, label string) *Class {
	if class, ok := cd.Classes[id]; ok {
		if label != "" {
			class.Label = label
		}
		return class
	}
	if label == "" {
		label = id
	}
	class := &Class{ID: id, Label: label}
	cd.Classes[id] = class
	cd.ClassOrder = append(cd.ClassOrder, id)
	return class
}

func (c *Class) addMember(member string) {
	member = strings.TrimSpace(member)
	if member == "" {
		return
	}
	if strings.Contains(member, "(") && strings.Contains(member, ")") {
		c.Operations = append(c.Operations, member)
		return
	}
	c.Attributes = append(c.Attributes, member)
}

func looksLikeRelationship(line string) bool {
	for _, operator := range []string{"<|", "*--", "o--", "-->", "--", ".."} {
		if strings.Contains(line, operator) {
			return true
		}
	}
	return false
}
