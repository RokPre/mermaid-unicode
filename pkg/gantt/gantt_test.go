package gantt

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestGanttParseAndRender(t *testing.T) {
	g, err := Parse(`gantt
title Project
dateFormat YYYY-MM-DD
section Build
Task one :a1, 2024-01-01, 7d
Task two :after a1, 3d`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if g.Title != "Project" || g.DateFormat != "YYYY-MM-DD" || len(g.Tasks) != 2 {
		t.Fatalf("gantt = %#v", g)
	}
	output, err := Render(g, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"Project", "[Build]", "Task one", "████████", "after a1"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestGanttValidation(t *testing.T) {
	_, err := Parse("gantt\ntitle Empty")
	if err == nil || !strings.Contains(err.Error(), "no gantt tasks found") {
		t.Fatalf("Parse() error = %v, want no tasks", err)
	}
}

func TestGanttASCII(t *testing.T) {
	g, err := Parse("gantt\nTask :2024-01-01, 1d")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(g, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(output, "########") {
		t.Fatalf("unexpected output:\n%s", output)
	}
}
