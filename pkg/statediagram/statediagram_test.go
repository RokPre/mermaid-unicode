package statediagram

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParseStateDiagram(t *testing.T) {
	input := `stateDiagram-v2
direction LR
state "Still State" as Still
Still : waiting
state choice <<choice>>
[*] --> Still
Still --> choice : decide
note right of Still: local note
state Moving {
  Walk --> Run : faster
}`

	sd, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if sd.Direction != "LR" {
		t.Fatalf("Direction = %q, want LR", sd.Direction)
	}
	if sd.States["Still"].Label != "Still State" || sd.States["Still"].Description != "waiting" {
		t.Fatalf("Still = %#v, want alias and description", sd.States["Still"])
	}
	if sd.States["choice"].Kind != "choice" {
		t.Fatalf("choice kind = %q, want choice", sd.States["choice"].Kind)
	}
	if len(sd.Transitions) != 3 {
		t.Fatalf("transitions = %d, want 3", len(sd.Transitions))
	}
	if len(sd.Notes) != 1 {
		t.Fatalf("notes = %d, want 1", len(sd.Notes))
	}
	if len(sd.Composites) != 1 || len(sd.Composites[0].Children) != 2 {
		t.Fatalf("composites = %#v, want Moving with Walk and Run", sd.Composites)
	}
}

func TestParseStateValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name: "invalid transition",
			input: `stateDiagram-v2
A -> B`,
			wantErr: "invalid state syntax",
		},
		{
			name: "unexpected composite end",
			input: `stateDiagram-v2
}`,
			wantErr: "composite state end without matching start",
		},
		{
			name: "unterminated composite",
			input: `stateDiagram-v2
state Parent {
  A --> B`,
			wantErr: "unterminated composite state",
		},
		{
			name:    "no states",
			input:   `stateDiagram-v2`,
			wantErr: "no states found",
		},
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

func TestRenderStateDiagramUnicode(t *testing.T) {
	sd, err := Parse(`stateDiagram-v2
state "Still State" as Still
Still : waiting
state choice <<choice>>
[*] --> Still
Still --> choice : decide
Still --> [*] : done
note right of Still: local note
state Moving {
  Walk --> Run : faster
}`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	output, err := Render(sd, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"●", "◉", "Still State", "waiting", "────▶ decide", "◇ choice", "note right of Still State", "Moving", "Walk"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestRenderStateDiagramASCII(t *testing.T) {
	sd, err := Parse(`stateDiagram
[*] --> Still
Still --> [*] : done`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(sd, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"(*)", "((*))", "----> done"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}
