package cmd

import (
	"strings"
	"testing"
)

func TestDiagramFactoryDetectsSupportedTypes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name: "sequence diagram",
			input: `sequenceDiagram
    A->>B: Test`,
			want: "sequence",
		},
		{
			name: "sequence diagram after comments",
			input: `%% comment

sequenceDiagram
    A->>B: Test`,
			want: "sequence",
		},
		{
			name: "graph diagram",
			input: `graph LR
    A-->B`,
			want: "graph",
		},
		{
			name: "flowchart diagram",
			input: `flowchart TD
    A-->B`,
			want: "graph",
		},
		{
			name: "graph diagram after padding directives",
			input: `paddingX=8
paddingY=2
graph LR
    A-->B`,
			want: "graph",
		},
		{
			name: "er diagram",
			input: `erDiagram
CUSTOMER ||--o{ ORDER : places`,
			want: "er",
		},
		{
			name: "class diagram",
			input: `classDiagram
Animal <|-- Duck`,
			want: "class",
		},
		{
			name: "state diagram",
			input: `stateDiagram-v2
[*] --> Still`,
			want: "state",
		},
		{
			name: "requirement diagram",
			input: `requirementDiagram
requirement req {
  id: 1
}`,
			want: "requirement",
		},
		{
			name: "mindmap",
			input: `mindmap
Root
  Idea`,
			want: "mindmap",
		},
		{
			name: "timeline",
			input: `timeline
2024 : Idea`,
			want: "timeline",
		},
		{
			name: "journey",
			input: `journey
Task: 3: Me`,
			want: "journey",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag, err := DiagramFactory(tt.input)
			if err != nil {
				t.Fatalf("DiagramFactory() error = %v", err)
			}
			if diag.Type() != tt.want {
				t.Fatalf("DiagramFactory().Type() = %q, want %q", diag.Type(), tt.want)
			}
		})
	}
}

func TestDiagramFactoryRejectsKnownUnsupportedTypes(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		wantType    string
		wantMessage string
	}{
		{name: "gantt", input: "gantt\ntitle Plan", wantType: "gantt"},
		{name: "pie", input: "pie showData\ntitle Shares", wantType: "pie"},
		{name: "quadrant", input: "quadrantChart\ntitle Risk", wantType: "quadrantChart"},
		{name: "gitgraph", input: "gitGraph LR:\ncommit", wantType: "gitGraph"},
		{name: "zenuml", input: "zenuml\nA.method()", wantType: "zenuml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diag, err := DiagramFactory(tt.input)
			if err == nil {
				t.Fatalf("DiagramFactory() returned diagram type %q, want error", diag.Type())
			}
			message := err.Error()
			if !strings.Contains(message, tt.wantType) {
				t.Fatalf("error = %q, want unsupported type %q", message, tt.wantType)
			}
			if !strings.Contains(message, "graph") || !strings.Contains(message, "sequence") {
				t.Fatalf("error = %q, want supported graph and sequence types", message)
			}
		})
	}
}

func TestDiagramFactoryRejectsUnknownInput(t *testing.T) {
	_, err := DiagramFactory("not a diagram")
	if err == nil {
		t.Fatal("DiagramFactory() error = nil, want unknown diagram error")
	}

	message := err.Error()
	if !strings.Contains(message, "unknown diagram type") {
		t.Fatalf("error = %q, want unknown diagram type", message)
	}
	if !strings.Contains(message, "graph") || !strings.Contains(message, "sequence") {
		t.Fatalf("error = %q, want supported graph and sequence types", message)
	}
}

func TestDiagramFactoryRejectsEmptyInput(t *testing.T) {
	_, err := DiagramFactory("%% only comments\n\n")
	if err == nil {
		t.Fatal("DiagramFactory() error = nil, want missing definition error")
	}

	message := err.Error()
	if !strings.Contains(message, "missing diagram definition") {
		t.Fatalf("error = %q, want missing diagram definition", message)
	}
}
