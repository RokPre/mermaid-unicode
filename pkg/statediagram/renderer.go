package statediagram

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/mattn/go-runewidth"
)

const stateConnPad = 3

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
	arrowRight  string
	start       string
	end         string
	choice      string
	fork        string
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
	arrowRight:  ">",
	start:       "(*)",
	end:         "((*))",
	choice:      "<>",
	fork:        "===",
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
	arrowRight:  "▶",
	start:       "●",
	end:         "◉",
	choice:      "◇",
	fork:        "━━━",
}

func Render(sd *Diagram, config *diagram.Config) (string, error) {
	if sd == nil || len(sd.States) == 0 {
		return "", fmt.Errorf("no states")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}

	chars := unicodeChars
	if config.UseAscii {
		chars = asciiChars
	}

	lines := []string{}
	renderedInTransition := map[string]bool{}
	for _, transition := range sd.Transitions {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, renderTransition(transition, chars)...)
		renderedInTransition[transition.From.ID] = true
		renderedInTransition[transition.To.ID] = true
	}

	for _, note := range sd.Notes {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, renderNote(note, chars)...)
	}

	for _, composite := range sd.Composites {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, renderComposite(composite, chars)...)
	}

	for _, id := range sd.StateOrder {
		if renderedInTransition[id] || id == StartEndID {
			continue
		}
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		box, _ := renderStateBox(sd.States[id], "", chars)
		lines = append(lines, box...)
	}

	return strings.Join(lines, "\n") + "\n", nil
}

func renderTransition(transition *Transition, chars boxChars) []string {
	fromRole := ""
	toRole := ""
	if transition.From.ID == StartEndID {
		fromRole = "start"
	}
	if transition.To.ID == StartEndID {
		toRole = "end"
	}

	leftBox, leftWidth := renderStateBox(transition.From, fromRole, chars)
	rightBox, rightWidth := renderStateBox(transition.To, toRole, chars)
	height := max(len(leftBox), len(rightBox))
	leftBox = padLines(leftBox, leftWidth, height)
	rightBox = padLines(rightBox, rightWidth, height)

	connector := strings.Repeat(chars.solid, 4) + chars.arrowRight
	if transition.Label != "" {
		connector += " " + transition.Label
	}
	blankConnector := strings.Repeat(" ", runewidth.StringWidth(connector)+stateConnPad*2)
	activeConnector := strings.Repeat(" ", stateConnPad) + connector + strings.Repeat(" ", stateConnPad)
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

func renderStateBox(state *State, role string, chars boxChars) ([]string, int) {
	if role == "start" {
		return []string{chars.start}, runewidth.StringWidth(chars.start)
	}
	if role == "end" {
		return []string{chars.end}, runewidth.StringWidth(chars.end)
	}

	label := state.Label
	if state.Kind == "choice" {
		label = chars.choice + " " + label
	}
	if state.Kind == "fork" || state.Kind == "join" {
		label = chars.fork + " " + label
	}
	content := []string{label}
	if state.Description != "" {
		content = append(content, state.Description)
	}

	innerWidth := 8
	for _, line := range content {
		if width := runewidth.StringWidth(line); width > innerWidth {
			innerWidth = width
		}
	}
	innerWidth += 2

	lines := []string{
		chars.topLeft + strings.Repeat(chars.horizontal, innerWidth) + chars.topRight,
		chars.vertical + center(label, innerWidth) + chars.vertical,
	}
	if state.Description != "" {
		lines = append(lines, chars.teeRight+strings.Repeat(chars.horizontal, innerWidth)+chars.teeLeft)
		lines = append(lines, chars.vertical+padRight(" "+state.Description, innerWidth)+chars.vertical)
	}
	lines = append(lines, chars.bottomLeft+strings.Repeat(chars.horizontal, innerWidth)+chars.bottomRight)
	return lines, innerWidth + 2
}

func renderNote(note *Note, chars boxChars) []string {
	text := fmt.Sprintf("note %s %s: %s", note.Position, note.State.Label, note.Text)
	innerWidth := runewidth.StringWidth(text) + 2
	return []string{
		chars.topLeft + strings.Repeat(chars.horizontal, innerWidth) + chars.topRight,
		chars.vertical + " " + text + " " + chars.vertical,
		chars.bottomLeft + strings.Repeat(chars.horizontal, innerWidth) + chars.bottomRight,
	}
}

func renderComposite(composite *Composite, chars boxChars) []string {
	title := " " + composite.Label + " "
	innerWidth := runewidth.StringWidth(title)
	for _, child := range composite.Children {
		if width := runewidth.StringWidth(child) + 2; width > innerWidth {
			innerWidth = width
		}
	}
	if innerWidth < 12 {
		innerWidth = 12
	}

	topFill := innerWidth - runewidth.StringWidth(title)
	lines := []string{chars.topLeft + title + strings.Repeat(chars.horizontal, topFill) + chars.topRight}
	for _, child := range composite.Children {
		lines = append(lines, chars.vertical+padRight(" "+child, innerWidth)+chars.vertical)
	}
	if len(composite.Children) == 0 {
		lines = append(lines, chars.vertical+strings.Repeat(" ", innerWidth)+chars.vertical)
	}
	lines = append(lines, chars.bottomLeft+strings.Repeat(chars.horizontal, innerWidth)+chars.bottomRight)
	return lines
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
