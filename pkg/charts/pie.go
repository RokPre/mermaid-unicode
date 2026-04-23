package charts

import (
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const PieKeyword = "pie"

type PieDiagram struct {
	Title    string
	ShowData bool
	Slices   []PieSlice
}

type PieSlice struct {
	Label string
	Value float64
}

func IsPie(input string) bool {
	return firstKeyword(input, PieKeyword)
}

func ParsePie(input string) (*PieDiagram, error) {
	lines := diagram.RemoveComments(diagram.SplitLines(strings.TrimSpace(input)))
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), PieKeyword) {
		return nil, fmt.Errorf("expected %q keyword", PieKeyword)
	}
	p := &PieDiagram{ShowData: strings.Contains(strings.ToLower(lines[0]), "showdata")}
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(strings.ToLower(trimmed), "title ") {
			p.Title = strings.TrimSpace(trimmed[len("title "):])
			continue
		}
		parts := strings.Split(trimmed, ":")
		if len(parts) != 2 {
			return nil, fmt.Errorf("line %d: invalid pie slice syntax: %q", lineNo, trimmed)
		}
		value, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
		if err != nil || value < 0 {
			return nil, fmt.Errorf("line %d: pie value must be a non-negative number", lineNo)
		}
		p.Slices = append(p.Slices, PieSlice{Label: strings.Trim(strings.TrimSpace(parts[0]), `"`), Value: value})
	}
	if len(p.Slices) == 0 {
		return nil, fmt.Errorf("no pie slices found")
	}
	return p, nil
}

func RenderPie(p *PieDiagram, config *diagram.Config) (string, error) {
	if p == nil || len(p.Slices) == 0 {
		return "", fmt.Errorf("no pie slices")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	full, empty := "█", "░"
	if config.UseAscii {
		full, empty = "#", "."
	}
	total := 0.0
	for _, slice := range p.Slices {
		total += slice.Value
	}
	if total == 0 {
		return "", fmt.Errorf("pie total must be greater than zero")
	}
	lines := []string{}
	if p.Title != "" {
		lines = append(lines, p.Title)
	}
	for _, slice := range p.Slices {
		percent := slice.Value / total * 100
		filled := int(math.Round(percent / 5))
		if filled > 20 {
			filled = 20
		}
		bar := strings.Repeat(full, filled) + strings.Repeat(empty, 20-filled)
		value := ""
		if p.ShowData {
			value = fmt.Sprintf(" %.2f", slice.Value)
		}
		lines = append(lines, fmt.Sprintf("%-16s %6.2f%%%s %s", slice.Label, percent, value, bar))
	}
	return strings.Join(lines, "\n") + "\n", nil
}

func firstKeyword(input, keyword string) bool {
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, keyword)
	}
	return false
}
