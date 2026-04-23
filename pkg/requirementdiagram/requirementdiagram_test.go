package requirementdiagram

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParseRequirementDiagram(t *testing.T) {
	input := `requirementDiagram
requirement test_req {
  id: 1
  text: the system shall work
  risk: high
  verifymethod: test
}
element simulator {
  type: simulation
  docref: sim.md
}
simulator - satisfies -> test_req`

	rd, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(rd.Items) != 2 {
		t.Fatalf("items = %d, want 2", len(rd.Items))
	}
	if rd.Items["test_req"].Fields["risk"] != "high" {
		t.Fatalf("risk = %q, want high", rd.Items["test_req"].Fields["risk"])
	}
	if len(rd.Relationships) != 1 || rd.Relationships[0].Kind != "satisfies" {
		t.Fatalf("relationships = %#v, want satisfies", rd.Relationships)
	}
}

func TestParseRequirementValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"invalid relationship", "requirementDiagram\nA satisfies B", "invalid requirement syntax"},
		{"unterminated block", "requirementDiagram\nrequirement req {\nid: 1", "unterminated requirement block"},
		{"no items", "requirementDiagram", "no requirements or elements found"},
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

func TestRenderRequirementDiagramUnicode(t *testing.T) {
	rd, err := Parse(`requirementDiagram
requirement req {
  id: 1
  text: must work
}
element tester {
  type: person
}
tester - verifies -> req`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	output, err := Render(rd, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"┌", "requirement req", "text: must work", "element tester", "────▶ verifies"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestRenderRequirementDiagramASCII(t *testing.T) {
	rd, err := Parse(`requirementDiagram
element tester {
  type: person
}
requirement req {
  id: 1
}
tester - traces -> req`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(rd, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"+", "element tester", "----> traces"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}
