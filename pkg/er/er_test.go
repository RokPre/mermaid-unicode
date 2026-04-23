package er

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParseERDiagram(t *testing.T) {
	input := `erDiagram
direction LR
CUSTOMER[Customer]
CUSTOMER {
  string id PK
  string name
}
ORDER {
  int id PK
  string customer_id FK
}
CUSTOMER ||--o{ ORDER : places`

	er, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if er.Direction != "LR" {
		t.Fatalf("Direction = %q, want LR", er.Direction)
	}
	if len(er.Entities) != 2 {
		t.Fatalf("entities = %d, want 2", len(er.Entities))
	}
	if er.Entities["CUSTOMER"].Label != "Customer" {
		t.Fatalf("CUSTOMER label = %q, want Customer", er.Entities["CUSTOMER"].Label)
	}
	if len(er.Entities["ORDER"].Attributes) != 2 {
		t.Fatalf("ORDER attributes = %d, want 2", len(er.Entities["ORDER"].Attributes))
	}
	if got := er.Entities["ORDER"].Attributes[1].Keys[0]; got != "FK" {
		t.Fatalf("ORDER second key = %q, want FK", got)
	}
	if len(er.Relationships) != 1 {
		t.Fatalf("relationships = %d, want 1", len(er.Relationships))
	}
	rel := er.Relationships[0]
	if rel.LeftCardinality != "||" || rel.RightCardinality != "o{" || !rel.Identifying || rel.Label != "places" {
		t.Fatalf("relationship = %#v, want identifying ||--o{ places", rel)
	}
}

func TestParseERValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name: "invalid relationship",
			input: `erDiagram
CUSTOMER -- ORDER : broken`,
			wantErr: "invalid ER relationship syntax",
		},
		{
			name: "unterminated block",
			input: `erDiagram
CUSTOMER {
  string id PK`,
			wantErr: "unterminated entity block",
		},
		{
			name:    "no entities",
			input:   `erDiagram`,
			wantErr: "no entities found",
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

func TestRenderERDiagramUnicode(t *testing.T) {
	er, err := Parse(`erDiagram
CUSTOMER {
  string id PK
}
ORDER {
  int id PK
}
CUSTOMER ||--o{ ORDER : places
ORDER }o..|| CUSTOMER : billed_to`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	output, err := Render(er, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"┌", "CUSTOMER", "string id PK", "||────o{ places", "}o┄┄┄┄|| billed_to"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestRenderERDiagramASCII(t *testing.T) {
	er, err := Parse(`erDiagram
CUSTOMER ||--o{ ORDER : places`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(er, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"+", "| CUSTOMER |", "||----o{ places"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}
