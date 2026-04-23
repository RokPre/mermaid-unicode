package requirementdiagram

import (
	"fmt"
	"sort"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/mattn/go-runewidth"
)

type boxChars struct {
	tl, tr, bl, br, h, v, r, l, arrow string
}

var asciiChars = boxChars{"+", "+", "+", "+", "-", "|", "+", "+", ">"}
var unicodeChars = boxChars{"┌", "┐", "└", "┘", "─", "│", "├", "┤", "▶"}

func Render(rd *Diagram, config *diagram.Config) (string, error) {
	if rd == nil || len(rd.Items) == 0 {
		return "", fmt.Errorf("no requirements or elements")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	chars := unicodeChars
	if config.UseAscii {
		chars = asciiChars
	}

	lines := []string{}
	rendered := map[string]bool{}
	for _, rel := range rd.Relationships {
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		lines = append(lines, renderRelationship(rel, chars)...)
		rendered[rel.From.ID] = true
		rendered[rel.To.ID] = true
	}
	for _, id := range rd.ItemOrder {
		if rendered[id] {
			continue
		}
		if len(lines) > 0 {
			lines = append(lines, "")
		}
		box, _ := renderItemBox(rd.Items[id], chars)
		lines = append(lines, box...)
	}
	return strings.Join(lines, "\n") + "\n", nil
}

func renderRelationship(rel *Relationship, chars boxChars) []string {
	left, leftWidth := renderItemBox(rel.From, chars)
	right, rightWidth := renderItemBox(rel.To, chars)
	height := max(len(left), len(right))
	left = padLines(left, leftWidth, height)
	right = padLines(right, rightWidth, height)
	conn := strings.Repeat(chars.h, 4) + chars.arrow + " " + rel.Kind
	blank := strings.Repeat(" ", runewidth.StringWidth(conn)+6)
	active := "   " + conn + "   "
	mid := height / 2
	lines := []string{}
	for i := 0; i < height; i++ {
		sep := blank
		if i == mid {
			sep = active
		}
		lines = append(lines, padRight(left[i], leftWidth)+sep+padRight(right[i], rightWidth))
	}
	return lines
}

func renderItemBox(item *Item, chars boxChars) ([]string, int) {
	rows := []string{item.Kind + " " + item.ID}
	keys := make([]string, 0, len(item.Fields))
	for key := range item.Fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		rows = append(rows, key+": "+item.Fields[key])
	}
	width := 12
	for _, row := range rows {
		if w := runewidth.StringWidth(row); w > width {
			width = w
		}
	}
	width += 2
	lines := []string{
		chars.tl + strings.Repeat(chars.h, width) + chars.tr,
		chars.v + center(rows[0], width) + chars.v,
	}
	if len(rows) > 1 {
		lines = append(lines, chars.r+strings.Repeat(chars.h, width)+chars.l)
		for _, row := range rows[1:] {
			lines = append(lines, chars.v+padRight(" "+row, width)+chars.v)
		}
	}
	lines = append(lines, chars.bl+strings.Repeat(chars.h, width)+chars.br)
	return lines, width + 2
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
