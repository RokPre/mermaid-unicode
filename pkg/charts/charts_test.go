package charts

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestPieChart(t *testing.T) {
	p, err := ParsePie("pie showData\ntitle Pets\nDogs : 75\nCats : 25")
	if err != nil {
		t.Fatalf("ParsePie() error = %v", err)
	}
	output, err := RenderPie(p, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("RenderPie() error = %v", err)
	}
	for _, want := range []string{"Pets", "Dogs", "75.00%", "Cats", "25.00%"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestPieValidation(t *testing.T) {
	_, err := ParsePie("pie\nDogs : nope")
	if err == nil || !strings.Contains(err.Error(), "pie value") {
		t.Fatalf("ParsePie() error = %v, want value error", err)
	}
}

func TestQuadrantChart(t *testing.T) {
	q, err := ParseQuadrant(`quadrantChart
title Priorities
x-axis Low --> High
y-axis Low --> High
quadrant-1 Expand
Campaign A: [0.25, 0.75]`)
	if err != nil {
		t.Fatalf("ParseQuadrant() error = %v", err)
	}
	output, err := RenderQuadrant(q, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("RenderQuadrant() error = %v", err)
	}
	for _, want := range []string{"Priorities", "●", "Campaign A: [0.25, 0.75]"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestQuadrantValidation(t *testing.T) {
	_, err := ParseQuadrant("quadrantChart\nBad: [1.5, 0.2]")
	if err == nil || !strings.Contains(err.Error(), "between 0 and 1") {
		t.Fatalf("ParseQuadrant() error = %v, want range error", err)
	}
}
