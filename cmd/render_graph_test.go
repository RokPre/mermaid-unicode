package cmd

import (
	"strings"
	"testing"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/mattn/go-runewidth"
)

func TestRenderGraphKeepsDisplayWidthForWideNodeLabels(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph LR\nA[\"中A\"] --> B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	assertUniformDisplayWidth(t, output)
}

func TestRenderGraphKeepsDisplayWidthForWideSubgraphTitles(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph LR\nsubgraph sg [数据库]\nA --> B\nend", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	assertUniformDisplayWidth(t, output)
}

func TestRenderGraphKeepsExplicitTargetLabelAfterBareReference(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph TD\nA[\"Foo\"] --> B[\"Bar\"]\nB --> C[\"Baz\"]", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if !strings.Contains(output, "Bar") {
		t.Fatalf("expected output to contain Bar\noutput:\n%s", output)
	}
	if strings.Contains(output, "\n|  B  |") || strings.Contains(output, "\n| B |\n") {
		t.Fatalf("expected B node to keep explicit label\noutput:\n%s", output)
	}
}

func TestRenderGraphKeepsStandaloneSubgraphLabelWhenReferencedLater(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph TD\nsubgraph one\n    A[\"VcpuManager\"]\nend\nA --> B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if !strings.Contains(output, "VcpuManager") {
		t.Fatalf("expected output to contain VcpuManager\noutput:\n%s", output)
	}
	if strings.Contains(output, "\n| A |\n") || strings.Contains(output, "\n|  A  |\n") {
		t.Fatalf("expected A node to keep standalone explicit label\noutput:\n%s", output)
	}
}

func TestRenderGraphSupportsLiteralNewlineInNodeLabel(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph LR\nA[\"line1\nline2\"] --> B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if !strings.Contains(output, "line1") || !strings.Contains(output, "line2") {
		t.Fatalf("expected output to contain both label lines\noutput:\n%s", output)
	}
	if strings.Contains(output, "A[\"line1") || strings.Contains(output, "line2\"]") {
		t.Fatalf("expected parser to keep literal newline inside the label\noutput:\n%s", output)
	}
}

func TestRenderGraphSeparatesDuplicateEdgeLabels(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph LR\nA -->|miss| B\nA -->|hit| B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if strings.Contains(output, "mhit") {
		t.Fatalf("expected duplicate edge labels not to merge\noutput:\n%s", output)
	}
	if !strings.Contains(output, "miss") || !strings.Contains(output, "hit") {
		t.Fatalf("expected output to contain both duplicate edge labels\noutput:\n%s", output)
	}

	missLine := -1
	hitLine := -1
	for i, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "miss") {
			missLine = i
		}
		if strings.Contains(line, "hit") {
			hitLine = i
		}
	}
	if missLine == -1 || hitLine == -1 || missLine == hitLine {
		t.Fatalf("expected duplicate edge labels on separate lines\noutput:\n%s", output)
	}
}

func TestRenderGraphSeparatesBidirectionalEdgeLabelsLR(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph LR\nA -->|workload exits| B\nB -->|run| A", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if strings.Contains(output, "worklorunexits") {
		t.Fatalf("expected bidirectional edge labels not to merge\noutput:\n%s", output)
	}
	if !strings.Contains(output, "workload") || !strings.Contains(output, "exits") || !strings.Contains(output, "run") {
		t.Fatalf("expected output to contain both bidirectional edge labels\noutput:\n%s", output)
	}

	workloadLine := -1
	runLine := -1
	for i, line := range strings.Split(output, "\n") {
		if strings.Contains(line, "workload") {
			workloadLine = i
		}
		if strings.Contains(line, "run") {
			runLine = i
		}
	}
	if workloadLine == -1 || runLine == -1 || workloadLine == runLine {
		t.Fatalf("expected bidirectional edge labels on separate lines\noutput:\n%s", output)
	}
}

func TestRenderGraphSeparatesBidirectionalEdgeLabelsTD(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph TD\nA -->|forward| B\nB -->|back| A", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if strings.Contains(output, "fbackrd") {
		t.Fatalf("expected bidirectional edge labels not to merge\noutput:\n%s", output)
	}
	if !strings.Contains(output, "forward") || !strings.Contains(output, "back") {
		t.Fatalf("expected output to contain both bidirectional edge labels\noutput:\n%s", output)
	}
}

func TestRenderGraphUsesUnicodeNodeShapeGlyphs(t *testing.T) {
	config := diagram.NewTestConfig(false, "cli")
	output, err := RenderDiagram("graph LR\nA(Rounded)\nB[[Double]]\nC{Decision}\nD([Stadium])\nE{{Hex}}\nF[/Para/]", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	for _, want := range []string{"╭", "╮", "╔", "╗", "◇", "╱", "╲"} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected output to contain %q\noutput:\n%s", want, output)
		}
	}

	assertUniformDisplayWidth(t, output)
}

func TestRenderGraphShapeGlyphsRespectAsciiFallback(t *testing.T) {
	config := diagram.NewTestConfig(true, "cli")
	output, err := RenderDiagram("graph LR\nA(Rounded)\nB[[Double]]\nC{Decision}", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	for _, disallowed := range []string{"╭", "╮", "╔", "╗", "◇"} {
		if strings.Contains(output, disallowed) {
			t.Fatalf("expected ASCII fallback not to contain %q\noutput:\n%s", disallowed, output)
		}
	}
	if !strings.Contains(output, "+") {
		t.Fatalf("expected ASCII fallback to use plus corners\noutput:\n%s", output)
	}
}

func TestRenderGraphUsesStyledEdgeGlyphs(t *testing.T) {
	config := diagram.NewTestConfig(false, "cli")

	heavyOutput, err := RenderDiagram("graph LR\nA ==> B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() heavy error = %v", err)
	}
	if !strings.Contains(heavyOutput, "━") || !strings.Contains(heavyOutput, "►") {
		t.Fatalf("expected heavy edge output to contain heavy line and arrowhead\noutput:\n%s", heavyOutput)
	}

	dashedOutput, err := RenderDiagram("graph LR\nA -.-> B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() dashed error = %v", err)
	}
	if !strings.Contains(dashedOutput, "┄") || !strings.Contains(dashedOutput, "►") {
		t.Fatalf("expected dashed arrow output to contain dashed line and arrowhead\noutput:\n%s", dashedOutput)
	}

	openDashedOutput, err := RenderDiagram("graph LR\nA -.- B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() open dashed error = %v", err)
	}
	if !strings.Contains(openDashedOutput, "┄") {
		t.Fatalf("expected open dashed output to contain dashed line\noutput:\n%s", openDashedOutput)
	}
	if strings.Contains(openDashedOutput, "►") {
		t.Fatalf("expected open dashed output not to contain an arrowhead\noutput:\n%s", openDashedOutput)
	}
}

func TestRenderGraphUsesConfiguredDefaultStyles(t *testing.T) {
	config := diagram.NewTestConfig(false, "cli")
	config.GraphBoxStyle = "rounded"
	config.GraphEdgeStyle = "heavy"

	output, err := RenderDiagram("graph LR\nA --> B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if !strings.Contains(output, "╭") || !strings.Contains(output, "╮") {
		t.Fatalf("expected configured rounded box style\noutput:\n%s", output)
	}
	if !strings.Contains(output, "━") || !strings.Contains(output, "►") {
		t.Fatalf("expected configured heavy edge style\noutput:\n%s", output)
	}
}

func TestRenderGraphExplicitEdgeStyleOverridesConfiguredDefault(t *testing.T) {
	config := diagram.NewTestConfig(false, "cli")
	config.GraphEdgeStyle = "dashed"

	output, err := RenderDiagram("graph LR\nA ==> B", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if !strings.Contains(output, "━") {
		t.Fatalf("expected explicit heavy edge style to win\noutput:\n%s", output)
	}
	if strings.Contains(output, "┄") {
		t.Fatalf("expected explicit heavy edge style not to use dashed default\noutput:\n%s", output)
	}
}

func TestRenderGraphUsesConfiguredSubgraphFrameStyle(t *testing.T) {
	config := diagram.NewTestConfig(false, "cli")
	config.GraphBoxStyle = "double"

	output, err := RenderDiagram("graph LR\nsubgraph Parser\nA --> B\nend", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	if !strings.Contains(output, "╔") || !strings.Contains(output, "╗") {
		t.Fatalf("expected configured double-line subgraph frame\noutput:\n%s", output)
	}
	if !strings.Contains(output, "Parser") {
		t.Fatalf("expected subgraph label to remain readable\noutput:\n%s", output)
	}
}

func TestRenderGraphAppliesNodeAndEdgeColors(t *testing.T) {
	config := diagram.NewTestConfig(false, "html")
	output, err := RenderDiagram("graph LR\nA[Start]:::warning -->|go| B[Done]\nclassDef warning stroke:#ff0000,color:#00ff00,fill:#111111;\nlinkStyle 0 stroke:#0000ff,color:#ff00ff;", config)
	if err != nil {
		t.Fatalf("RenderDiagram() error = %v", err)
	}

	for _, want := range []string{
		"color: #ff0000",
		"color: #00ff00",
		"background-color: #111111",
		"color: #0000ff",
		"color: #ff00ff",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("expected colored output to contain %q\noutput:\n%s", want, output)
		}
	}
}

func assertUniformDisplayWidth(t *testing.T, output string) {
	t.Helper()

	lines := strings.Split(output, "\n")
	if len(lines) == 0 {
		t.Fatal("expected rendered output")
	}

	want := runewidth.StringWidth(lines[0])
	for i, line := range lines[1:] {
		if got := runewidth.StringWidth(line); got != want {
			t.Fatalf("line %d display width = %d, want %d\noutput:\n%s", i+2, got, want, output)
		}
	}
}
