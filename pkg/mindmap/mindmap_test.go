package mindmap

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParseMindmap(t *testing.T) {
	mm, err := Parse(`mindmap
  Root
    [Idea]
      Research
    Plan`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if mm.Root.Text != "Root" {
		t.Fatalf("root = %q, want Root", mm.Root.Text)
	}
	if len(mm.Root.Children) != 2 {
		t.Fatalf("root children = %d, want 2", len(mm.Root.Children))
	}
	if mm.Root.Children[0].Text != "Idea" || mm.Root.Children[0].Children[0].Text != "Research" {
		t.Fatalf("parsed tree = %#v", mm.Root)
	}
}

func TestParseMindmapValidation(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr string
	}{
		{"multiple roots", "mindmap\nRoot\nOther", "multiple mindmap roots"},
		{"bad indentation", "mindmap\nRoot\n    Child\n  Other", "bad mindmap indentation"},
		{"no nodes", "mindmap", "no mindmap nodes found"},
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

func TestRenderMindmapUnicode(t *testing.T) {
	mm, err := Parse(`mindmap
Root
  Idea
    Research
  Plan`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	output, err := Render(mm, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"Root", "├── Idea", "│   └── Research", "└── Plan"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestRenderMindmapASCII(t *testing.T) {
	mm, err := Parse("mindmap\nRoot\n  A\n  B")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(mm, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"|-- A", "`-- B"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}
