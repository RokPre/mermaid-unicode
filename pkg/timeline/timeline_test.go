package timeline

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParseTimeline(t *testing.T) {
	tl, err := Parse(`timeline
title Project
section Discover
2024 : Idea : Research
section Build
2025 : Ship`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if tl.Title != "Project" || len(tl.Items) != 2 {
		t.Fatalf("timeline = %#v, want title and two items", tl)
	}
	if tl.Items[0].Section != "Discover" || len(tl.Items[0].Events) != 2 {
		t.Fatalf("first item = %#v, want section and two events", tl.Items[0])
	}
}

func TestParseTimelineValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"invalid item", "timeline\n2024", "invalid timeline item syntax"},
		{"no items", "timeline\ntitle Empty", "no timeline items found"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("Parse() error = %v, want containing %q", err, tt.wantErr)
			}
		})
	}
}

func TestRenderTimelineUnicode(t *testing.T) {
	tl, err := Parse("timeline\ntitle Project\nsection Discover\n2024 : Idea : Research")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	output, err := Render(tl, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"Project", "[Discover]", "├─ 2024: Idea", "│    Research"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestRenderTimelineASCIIHorizontal(t *testing.T) {
	tl, err := Parse("timeline LR\n2024 : Idea\n2025 : Ship")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(tl, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(output, "2024: Idea -- 2025: Ship") {
		t.Fatalf("unexpected output:\n%s", output)
	}
}
