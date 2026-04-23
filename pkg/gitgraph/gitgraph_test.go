package gitgraph

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestGitGraphParseAndRender(t *testing.T) {
	gg, err := Parse(`gitGraph
commit id: "A"
branch develop
checkout develop
commit id: "B"
checkout main
merge develop`)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if len(gg.Branches) != 2 || len(gg.Events) != 6 {
		t.Fatalf("gitgraph = %#v", gg)
	}
	output, err := Render(gg, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	for _, want := range []string{"main", "develop", "● A", "● B", "merge develop"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestGitGraphValidation(t *testing.T) {
	_, err := Parse("gitGraph\nmerge")
	if err == nil || !strings.Contains(err.Error(), "merge requires a branch") {
		t.Fatalf("Parse() error = %v, want merge error", err)
	}
}

func TestGitGraphASCII(t *testing.T) {
	gg, err := Parse("gitGraph\ncommit id: A")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	config := diagram.DefaultConfig()
	config.UseAscii = true
	output, err := Render(gg, config)
	if err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	if !strings.Contains(output, "* A") {
		t.Fatalf("unexpected output:\n%s", output)
	}
}
