package classdiagram

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParseClassDiagram(t *testing.T) {
	input := `classDiagram
direction LR
class BankAccount["Bank Account"] {
  +String owner
  +deposit(amount) bool
}
Customer "1" --> "*" Ticket : owns
Ticket : +String seat`

	cd, err := Parse(input)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if cd.Direction != "LR" {
		t.Fatalf("Direction = %q, want LR", cd.Direction)
	}
	if cd.Classes["BankAccount"].Label != "Bank Account" {
		t.Fatalf("BankAccount label = %q, want Bank Account", cd.Classes["BankAccount"].Label)
	}
	if got := cd.Classes["BankAccount"].Attributes[0]; got != "+String owner" {
		t.Fatalf("BankAccount first attribute = %q, want +String owner", got)
	}
	if got := cd.Classes["BankAccount"].Operations[0]; got != "+deposit(amount) bool" {
		t.Fatalf("BankAccount first operation = %q, want +deposit(amount) bool", got)
	}
	if got := cd.Classes["Ticket"].Attributes[0]; got != "+String seat" {
		t.Fatalf("Ticket attribute = %q, want +String seat", got)
	}
	if len(cd.Relationships) != 1 {
		t.Fatalf("relationships = %d, want 1", len(cd.Relationships))
	}
	rel := cd.Relationships[0]
	if rel.LeftCardinality != "1" || rel.RightCardinality != "*" || rel.Operator != "-->" || rel.Label != "owns" {
		t.Fatalf("relationship = %#v, want Customer \"1\" --> \"*\" Ticket : owns", rel)
	}
}

func TestParseClassValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{
			name: "invalid relationship",
			input: `classDiagram
A --| B`,
			wantErr: "invalid class relationship syntax",
		},
		{
			name: "unterminated class",
			input: `classDiagram
class A {
  +String id`,
			wantErr: "unterminated class block",
		},
		{
			name:    "no classes",
			input:   `classDiagram`,
			wantErr: "no classes found",
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

func TestRenderClassDiagramUnicode(t *testing.T) {
	cd, err := Parse(`classDiagram
class Animal {
  +String name
  +speak() string
}
Animal <|-- Duck
Duck ..> Pond : swims_in
Owner o-- Duck : keeps`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	output, err := Render(cd, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"┌", "Animal", "+String name", "+speak() string", "<|────", "┄┄┄┄> swims_in", "◇──── keeps"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestRenderClassDiagramASCII(t *testing.T) {
	cd, err := Parse(`classDiagram
Customer "1" --> "*" Ticket : owns`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(cd, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"+", "Customer", `"1" ----> "*" owns`} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}
