package er

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/mattn/go-runewidth"
)

const (
	minEntityWidth = 8
	connectorPad   = 3
)

type boxChars struct {
	topLeft     string
	topRight    string
	bottomLeft  string
	bottomRight string
	horizontal  string
	vertical    string
	teeRight    string
	teeLeft     string
	solid       string
	dashed      string
}

var asciiChars = boxChars{
	topLeft:     "+",
	topRight:    "+",
	bottomLeft:  "+",
	bottomRight: "+",
	horizontal:  "-",
	vertical:    "|",
	teeRight:    "+",
	teeLeft:     "+",
	solid:       "-",
	dashed:      ".",
}

var unicodeChars = boxChars{
	topLeft:     "┌",
	topRight:    "┐",
	bottomLeft:  "└",
	bottomRight: "┘",
	horizontal:  "─",
	vertical:    "│",
	teeRight:    "├",
	teeLeft:     "┤",
	solid:       "─",
	dashed:      "┄",
}

func Render(er *Diagram, config *diagram.Config) (string, error) {
	if er == nil || len(er.Entities) == 0 {
		return "", fmt.Errorf("no entities")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}

	chars := unicodeChars
	if config.UseAscii {
		chars = asciiChars
	}

	lines := []string{}
	renderedInRelationship := map[string]bool{}
	for _, relationship := range er.Relationships {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, renderRelationship(relationship, chars)...)
		renderedInRelationship[relationship.From.ID] = true
		renderedInRelationship[relationship.To.ID] = true
	}

	for _, id := range er.EntityOrder {
		if renderedInRelationship[id] {
			continue
		}
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		box, _ := renderEntityBox(er.Entities[id], chars)
		lines = append(lines, box...)
	}

	return strings.Join(lines, "\n") + "\n", nil
}

func renderRelationship(relationship *Relationship, chars boxChars) []string {
	leftBox, leftWidth := renderEntityBox(relationship.From, chars)
	rightBox, rightWidth := renderEntityBox(relationship.To, chars)
	height := max(len(leftBox), len(rightBox))
	leftBox = padLines(leftBox, leftWidth, height)
	rightBox = padLines(rightBox, rightWidth, height)

	connector := relationshipConnector(relationship, chars)
	blankConnector := strings.Repeat(" ", runewidth.StringWidth(connector)+connectorPad*2)
	activeConnector := strings.Repeat(" ", connectorPad) + connector + strings.Repeat(" ", connectorPad)
	mid := height / 2

	lines := make([]string, 0, height)
	for i := 0; i < height; i++ {
		sep := blankConnector
		if i == mid {
			sep = activeConnector
		}
		lines = append(lines, padRight(leftBox[i], leftWidth)+sep+padRight(rightBox[i], rightWidth))
	}
	return lines
}

func relationshipConnector(relationship *Relationship, chars boxChars) string {
	line := chars.solid
	if !relationship.Identifying {
		line = chars.dashed
	}
	connector := relationship.LeftCardinality + strings.Repeat(line, 4) + relationship.RightCardinality
	if relationship.Label != "" {
		connector += " " + relationship.Label
	}
	return connector
}

func renderEntityBox(entity *Entity, chars boxChars) ([]string, int) {
	content := []string{entity.Label}
	for _, attr := range entity.Attributes {
		content = append(content, attr.String())
	}

	innerWidth := minEntityWidth
	for _, line := range content {
		if width := runewidth.StringWidth(line); width > innerWidth {
			innerWidth = width
		}
	}
	innerWidth += 2

	lines := []string{
		chars.topLeft + strings.Repeat(chars.horizontal, innerWidth) + chars.topRight,
		chars.vertical + center(entity.Label, innerWidth) + chars.vertical,
	}
	if len(entity.Attributes) > 0 {
		lines = append(lines, chars.teeRight+strings.Repeat(chars.horizontal, innerWidth)+chars.teeLeft)
		for _, attr := range entity.Attributes {
			lines = append(lines, chars.vertical+padRight(" "+attr.String(), innerWidth)+chars.vertical)
		}
	}
	lines = append(lines, chars.bottomLeft+strings.Repeat(chars.horizontal, innerWidth)+chars.bottomRight)
	return lines, innerWidth + 2
}

func (a Attribute) String() string {
	text := a.Type + " " + a.Name
	if len(a.Keys) > 0 {
		text += " " + strings.Join(a.Keys, ",")
	}
	return text
}

func padLines(lines []string, width, height int) []string {
	padded := append([]string(nil), lines...)
	for len(padded) < height {
		padded = append(padded, strings.Repeat(" ", width))
	}
	return padded
}

func center(text string, width int) string {
	textWidth := runewidth.StringWidth(text)
	if textWidth >= width {
		return text
	}
	left := (width - textWidth) / 2
	right := width - textWidth - left
	return strings.Repeat(" ", left) + text + strings.Repeat(" ", right)
}

func padRight(text string, width int) string {
	textWidth := runewidth.StringWidth(text)
	if textWidth >= width {
		return text
	}
	return text + strings.Repeat(" ", width-textWidth)
}
