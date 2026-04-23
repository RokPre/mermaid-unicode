package cmd

import (
	"sort"

	"github.com/mattn/go-runewidth"
)

func (g *graph) applyColors() {
	if g.drawing == nil {
		return
	}

	for _, n := range g.nodes {
		g.applyNodeColors(n)
	}
	edges := append([]*edge(nil), g.edges...)
	sort.SliceStable(edges, func(i, j int) bool {
		return g.edgeDrawPriority(edges[i]) < g.edgeDrawPriority(edges[j])
	})
	for _, e := range edges {
		g.applyEdgeColors(e)
	}
}

func (g *graph) applyNodeColors(n *node) {
	if n == nil || n.drawing == nil || n.drawingCoord == nil || len(n.styleClass.Styles) == 0 {
		return
	}

	origin := *n.drawingCoord
	width, height := getDrawingSize(n.drawing)
	stroke := n.styleClass.Styles["stroke"]
	textColor := n.styleClass.Styles["color"]
	fill := n.styleClass.Styles["fill"]

	labelCoords := g.nodeLabelCoords(n, origin, width, height)
	labelCoordSet := make(map[drawingCoord]struct{}, len(labelCoords))
	for _, coord := range labelCoords {
		labelCoordSet[coord] = struct{}{}
	}

	if normalizeStyleColor(fill) != "" {
		for x := 1; x < width; x++ {
			for y := 1; y < height; y++ {
				coord := drawingCoord{x: origin.x + x, y: origin.y + y}
				if _, isLabel := labelCoordSet[coord]; isLabel {
					continue
				}
				g.applyCellStyle(coord, "", fill)
			}
		}
	}

	if normalizeStyleColor(stroke) != "" {
		for x := 0; x <= width; x++ {
			g.applyCellStyle(drawingCoord{x: origin.x + x, y: origin.y}, stroke, "")
			g.applyCellStyle(drawingCoord{x: origin.x + x, y: origin.y + height}, stroke, "")
		}
		for y := 1; y < height; y++ {
			g.applyCellStyle(drawingCoord{x: origin.x, y: origin.y + y}, stroke, "")
			g.applyCellStyle(drawingCoord{x: origin.x + width, y: origin.y + y}, stroke, "")
		}
	}

	if normalizeStyleColor(textColor) != "" || normalizeStyleColor(fill) != "" {
		for _, coord := range labelCoords {
			g.applyCellStyle(coord, textColor, fill)
		}
	}
}

func (g *graph) nodeLabelCoords(n *node, origin drawingCoord, width, height int) []drawingCoord {
	innerTop := 1
	innerHeight := height - 1
	contentTop := innerTop + (innerHeight-n.label.contentHeight())/2
	coords := []drawingCoord{}
	for lineIdx, line := range n.label.lines {
		textY := contentTop + lineIdx*(graphLabelLineGap+1)
		textWidth := runewidth.StringWidth(line)
		textX := width/2 - CeilDiv(textWidth, 2) + 1
		for _, r := range line {
			runeWidth := Max(runewidth.RuneWidth(r), 1)
			coords = append(coords, drawingCoord{x: origin.x + textX, y: origin.y + textY})
			textX += runeWidth
		}
	}
	return coords
}

func (g *graph) applyEdgeColors(e *edge) {
	if e == nil || len(e.styleClass.Styles) == 0 {
		return
	}

	stroke := e.styleClass.Styles["stroke"]
	labelColor := e.styleClass.Styles["color"]
	if labelColor == "" {
		labelColor = stroke
	}

	if normalizeStyleColor(stroke) != "" {
		for _, coord := range e.colorCoords {
			g.applyCellStyle(coord, stroke, "")
		}
	}
	if normalizeStyleColor(labelColor) != "" {
		for _, coord := range e.labelCoords {
			g.applyCellStyle(coord, labelColor, "")
		}
	}
}

func (g *graph) applyCellStyle(coord drawingCoord, fg, bg string) {
	if coord.x < 0 || coord.y < 0 || coord.x >= len(*g.drawing) || coord.y >= len((*g.drawing)[0]) {
		return
	}
	cell := (*g.drawing)[coord.x][coord.y]
	if cell == "" {
		return
	}
	(*g.drawing)[coord.x][coord.y] = wrapTextInStyle(cell, fg, bg, g.styleType)
}

func nonSpaceDrawingCoords(drawings ...*drawing) []drawingCoord {
	coords := []drawingCoord{}
	for _, d := range drawings {
		if d == nil {
			continue
		}
		for x := 0; x < len(*d); x++ {
			for y := 0; y < len((*d)[0]); y++ {
				if (*d)[x][y] != " " && (*d)[x][y] != "" {
					coords = append(coords, drawingCoord{x: x, y: y})
				}
			}
		}
	}
	return coords
}
