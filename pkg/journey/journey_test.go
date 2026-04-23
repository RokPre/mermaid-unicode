package journey

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParseJourney(t *testing.T) {
	j, err := Parse(`journey
title My day
section Work
  Write code: 5: Me, Pair
  Review: 3: Me`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if j.Title != "My day" || len(j.Tasks) != 2 {
		t.Fatalf("journey = %#v, want title and two tasks", j)
	}
	if j.Tasks[0].Score != 5 || len(j.Tasks[0].Actors) != 2 {
		t.Fatalf("first task = %#v", j.Tasks[0])
	}
}

func TestParseJourneyValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"bad syntax", "journey\nTask: 3", "invalid journey task syntax"},
		{"bad score", "journey\nTask: 9: Me", "journey score must be 1..5"},
		{"no tasks", "journey\ntitle Empty", "no journey tasks found"},
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

func TestRenderJourneyUnicode(t *testing.T) {
	j, err := Parse("journey\ntitle My day\nsection Work\nWrite code: 5: Me")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	output, err := Render(j, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"My day", "[Work]", "Write code", "5/5", "█████"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestRenderJourneyASCII(t *testing.T) {
	j, err := Parse("journey\nTask: 2: Me")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(j, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(output, "##...") {
		t.Fatalf("unexpected output:\n%s", output)
	}
}
