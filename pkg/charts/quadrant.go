package charts

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const QuadrantKeyword = "quadrantChart"

var pointRegex = regexp.MustCompile(`^\s*([^:]+)\s*:\s*\[\s*([0-9.]+)\s*,\s*([0-9.]+)\s*\]\s*$`)

type QuadrantDiagram struct {
	Title     string
	XAxis     string
	YAxis     string
	Quadrants map[int]string
	Points    []Point
}

type Point struct {
	Label string
	X     float64
	Y     float64
}

func IsQuadrant(input string) bool {
	return firstKeyword(input, QuadrantKeyword)
}

func ParseQuadrant(input string) (*QuadrantDiagram, error) {
	lines := diagram.RemoveComments(diagram.SplitLines(strings.TrimSpace(input)))
	if len(lines) == 0 || !strings.HasPrefix(strings.TrimSpace(lines[0]), QuadrantKeyword) {
		return nil, fmt.Errorf("expected %q keyword", QuadrantKeyword)
	}
	q := &QuadrantDiagram{Quadrants: map[int]string{}}
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		switch {
		case strings.HasPrefix(lower, "title "):
			q.Title = strings.TrimSpace(trimmed[len("title "):])
		case strings.HasPrefix(lower, "x-axis "):
			q.XAxis = strings.TrimSpace(trimmed[len("x-axis "):])
		case strings.HasPrefix(lower, "y-axis "):
			q.YAxis = strings.TrimSpace(trimmed[len("y-axis "):])
		case strings.HasPrefix(lower, "quadrant-"):
			parts := strings.SplitN(trimmed, " ", 2)
			if len(parts) == 2 {
				n, _ := strconv.Atoi(strings.TrimPrefix(strings.ToLower(parts[0]), "quadrant-"))
				q.Quadrants[n] = strings.TrimSpace(parts[1])
			}
		default:
			match := pointRegex.FindStringSubmatch(trimmed)
			if match == nil {
				return nil, fmt.Errorf("line %d: invalid quadrant point syntax: %q", lineNo, trimmed)
			}
			x, _ := strconv.ParseFloat(match[2], 64)
			y, _ := strconv.ParseFloat(match[3], 64)
			if x < 0 || x > 1 || y < 0 || y > 1 {
				return nil, fmt.Errorf("line %d: quadrant point values must be between 0 and 1", lineNo)
			}
			q.Points = append(q.Points, Point{Label: strings.TrimSpace(match[1]), X: x, Y: y})
		}
	}
	if len(q.Points) == 0 {
		return nil, fmt.Errorf("no quadrant points found")
	}
	return q, nil
}

func RenderQuadrant(q *QuadrantDiagram, config *diagram.Config) (string, error) {
	if q == nil || len(q.Points) == 0 {
		return "", fmt.Errorf("no quadrant points")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	point := "●"
	h, v := "─", "│"
	if config.UseAscii {
		point, h, v = "*", "-", "|"
	}
	size := 11
	grid := make([][]string, size)
	for y := range grid {
		grid[y] = make([]string, size)
		for x := range grid[y] {
			grid[y][x] = " "
		}
	}
	mid := size / 2
	for i := 0; i < size; i++ {
		grid[mid][i] = h
		grid[i][mid] = v
	}
	for _, p := range q.Points {
		x := int(math.Round(p.X * float64(size-1)))
		y := size - 1 - int(math.Round(p.Y*float64(size-1)))
		grid[y][x] = point
	}
	lines := []string{}
	if q.Title != "" {
		lines = append(lines, q.Title)
	}
	for _, row := range grid {
		lines = append(lines, strings.Join(row, ""))
	}
	for _, p := range q.Points {
		lines = append(lines, fmt.Sprintf("%s: [%.2f, %.2f]", p.Label, p.X, p.Y))
	}
	return strings.Join(lines, "\n") + "\n", nil
}
