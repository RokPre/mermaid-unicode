package classdiagram

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/mattn/go-runewidth"
)

const (
	minClassWidth = 10
	classConnPad  = 3
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

func Render(cd *Diagram, config *diagram.Config) (string, error) {
	if cd == nil || len(cd.Classes) == 0 {
		return "", fmt.Errorf("no classes")
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
	for _, relationship := range cd.Relationships {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, renderRelationship(relationship, chars)...)
		renderedInRelationship[relationship.From.ID] = true
		renderedInRelationship[relationship.To.ID] = true
	}

	for _, id := range cd.ClassOrder {
		if renderedInRelationship[id] {
			continue
		}
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		box, _ := renderClassBox(cd.Classes[id], chars)
		lines = append(lines, box...)
	}

	return strings.Join(lines, "\n") + "\n", nil
}

func renderRelationship(relationship *Relationship, chars boxChars) []string {
	leftBox, leftWidth := renderClassBox(relationship.From, chars)
	rightBox, rightWidth := renderClassBox(relationship.To, chars)
	height := max(len(leftBox), len(rightBox))
	leftBox = padLines(leftBox, leftWidth, height)
	rightBox = padLines(rightBox, rightWidth, height)

	connector := relationshipConnector(relationship, chars)
	blankConnector := strings.Repeat(" ", runewidth.StringWidth(connector)+classConnPad*2)
	activeConnector := strings.Repeat(" ", classConnPad) + connector + strings.Repeat(" ", classConnPad)
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
	connector := relationshipOperator(relationship.Operator, chars)
	if relationship.LeftCardinality != "" {
		connector = `"` + relationship.LeftCardinality + `" ` + connector
	}
	if relationship.RightCardinality != "" {
		connector += ` "` + relationship.RightCardinality + `"`
	}
	if relationship.Label != "" {
		connector += " " + relationship.Label
	}
	return connector
}

func relationshipOperator(operator string, chars boxChars) string {
	solid := strings.Repeat(chars.solid, 4)
	dashed := strings.Repeat(chars.dashed, 4)
	switch operator {
	case "<|--":
		return "<|" + solid
	case "*--":
		if chars.solid == "-" {
			return "*" + solid
		}
		return "◆" + solid
	case "o--":
		if chars.solid == "-" {
			return "o" + solid
		}
		return "◇" + solid
	case "-->":
		return solid + ">"
	case "--":
		return solid
	case "..>":
		return dashed + ">"
	case "..|>":
		return dashed + "|>"
	case "..":
		return dashed
	default:
		return operator
	}
}

func renderClassBox(class *Class, chars boxChars) ([]string, int) {
	content := []string{class.Label}
	content = append(content, class.Attributes...)
	content = append(content, class.Operations...)

	innerWidth := minClassWidth
	for _, line := range content {
		if width := runewidth.StringWidth(line); width > innerWidth {
			innerWidth = width
		}
	}
	innerWidth += 2

	lines := []string{
		chars.topLeft + strings.Repeat(chars.horizontal, innerWidth) + chars.topRight,
		chars.vertical + center(class.Label, innerWidth) + chars.vertical,
	}
	if len(class.Attributes) > 0 || len(class.Operations) > 0 {
		lines = append(lines, chars.teeRight+strings.Repeat(chars.horizontal, innerWidth)+chars.teeLeft)
		for _, attr := range class.Attributes {
			lines = append(lines, chars.vertical+padRight(" "+attr, innerWidth)+chars.vertical)
		}
	}
	if len(class.Operations) > 0 {
		lines = append(lines, chars.teeRight+strings.Repeat(chars.horizontal, innerWidth)+chars.teeLeft)
		for _, op := range class.Operations {
			lines = append(lines, chars.vertical+padRight(" "+op, innerWidth)+chars.vertical)
		}
	}
	lines = append(lines, chars.bottomLeft+strings.Repeat(chars.horizontal, innerWidth)+chars.bottomRight)
	return lines, innerWidth + 2
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
