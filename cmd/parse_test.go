package cmd

import "testing"

func TestSplitGraphLines(t *testing.T) {
	input := "graph LR\\nA[\"line1\\nline2\"] --> B\\nC --> D"

	got := splitGraphLines(input)
	want := []string{"graph LR", `A["line1\nline2"] --> B`, "C --> D"}

	if len(got) != len(want) {
		t.Fatalf("line count = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestSplitGraphLinesKeepsLiteralNewlineInsideNodeLabel(t *testing.T) {
	input := "graph LR\nA[\"line1\nline2\"] --> B\nC --> D"

	got := splitGraphLines(input)
	want := []string{"graph LR", "A[\"line1\nline2\"] --> B", "C --> D"}

	if len(got) != len(want) {
		t.Fatalf("line count = %d, want %d", len(got), len(want))
	}

	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("line %d = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestParseNodeWithExplicitLabel(t *testing.T) {
	node := parseNode(`A["line1<br/>line2"]:::primary`)

	if node.name != "A" {
		t.Fatalf("name = %q, want %q", node.name, "A")
	}
	if node.styleClass != "primary" {
		t.Fatalf("styleClass = %q, want %q", node.styleClass, "primary")
	}
	if len(node.label.lines) != 2 {
		t.Fatalf("label lines = %d, want 2", len(node.label.lines))
	}
	if node.label.lines[0] != "line1" || node.label.lines[1] != "line2" {
		t.Fatalf("label lines = %#v, want [line1 line2]", node.label.lines)
	}
	if node.shape != graphNodeShapeSquare {
		t.Fatalf("shape = %q, want %q", node.shape, graphNodeShapeSquare)
	}
	if !node.hasShape {
		t.Fatal("expected square bracket node to have explicit shape")
	}
}

func TestParseNodeShapes(t *testing.T) {
	tests := []struct {
		input     string
		wantName  string
		wantLabel string
		wantShape graphNodeShape
	}{
		{input: `A[Text]`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeSquare},
		{input: `A(Text)`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeRounded},
		{input: `A([Text])`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeStadium},
		{input: `A[[Text]]`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeDouble},
		{input: `A[(Text)]`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeDatabase},
		{input: `A((Text))`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeCircle},
		{input: `A{Text}`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeDecision},
		{input: `A{{Text}}`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeHexagon},
		{input: `A[/Text/]`, wantName: "A", wantLabel: "Text", wantShape: graphNodeShapeParallelogram},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node := parseNode(tt.input)

			if node.name != tt.wantName {
				t.Fatalf("name = %q, want %q", node.name, tt.wantName)
			}
			if len(node.label.lines) != 1 || node.label.lines[0] != tt.wantLabel {
				t.Fatalf("label lines = %#v, want [%s]", node.label.lines, tt.wantLabel)
			}
			if node.shape != tt.wantShape {
				t.Fatalf("shape = %q, want %q", node.shape, tt.wantShape)
			}
			if !node.hasShape {
				t.Fatal("expected node to have explicit shape")
			}
		})
	}
}

func TestParseNodeExpandedShapeSyntax(t *testing.T) {
	tests := []struct {
		input        string
		wantName     string
		wantLabel    string
		wantHasLabel bool
		wantShape    graphNodeShape
	}{
		{input: `A@{ shape: rect }`, wantName: "A", wantLabel: "A", wantShape: graphNodeShapeSquare},
		{input: `B@{ shape: rounded }`, wantName: "B", wantLabel: "B", wantShape: graphNodeShapeRounded},
		{input: `C@{ shape: stadium }`, wantName: "C", wantLabel: "C", wantShape: graphNodeShapeStadium},
		{input: `D@{ shape: subroutine }`, wantName: "D", wantLabel: "D", wantShape: graphNodeShapeDouble},
		{input: `E@{ shape: db }`, wantName: "E", wantLabel: "E", wantShape: graphNodeShapeDatabase},
		{input: `F@{ shape: circle }`, wantName: "F", wantLabel: "F", wantShape: graphNodeShapeCircle},
		{input: `G@{ shape: decision }`, wantName: "G", wantLabel: "G", wantShape: graphNodeShapeDecision},
		{input: `H@{ shape: hexagon }`, wantName: "H", wantLabel: "H", wantShape: graphNodeShapeHexagon},
		{input: `I@{ shape: lean-r }`, wantName: "I", wantLabel: "I", wantShape: graphNodeShapeParallelogram},
		{input: `J@{ shape: rounded, label: "Research" }`, wantName: "J", wantLabel: "Research", wantHasLabel: true, wantShape: graphNodeShapeRounded},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			node := parseNode(tt.input)

			if node.name != tt.wantName {
				t.Fatalf("name = %q, want %q", node.name, tt.wantName)
			}
			if len(node.label.lines) != 1 || node.label.lines[0] != tt.wantLabel {
				t.Fatalf("label lines = %#v, want [%s]", node.label.lines, tt.wantLabel)
			}
			if node.hasLabel != tt.wantHasLabel {
				t.Fatalf("hasLabel = %v, want %v", node.hasLabel, tt.wantHasLabel)
			}
			if node.shape != tt.wantShape {
				t.Fatalf("shape = %q, want %q", node.shape, tt.wantShape)
			}
			if !node.hasShape {
				t.Fatal("expected expanded shape syntax to set explicit shape")
			}
		})
	}
}

func TestParseNodeExpandedShapeKeepsStyleClass(t *testing.T) {
	node := parseNode(`A@{ shape: rounded, label: "Research" }:::important`)

	if node.name != "A" {
		t.Fatalf("name = %q, want A", node.name)
	}
	if node.styleClass != "important" {
		t.Fatalf("styleClass = %q, want important", node.styleClass)
	}
	if node.shape != graphNodeShapeRounded {
		t.Fatalf("shape = %q, want %q", node.shape, graphNodeShapeRounded)
	}
	if len(node.label.lines) != 1 || node.label.lines[0] != "Research" {
		t.Fatalf("label lines = %#v, want [Research]", node.label.lines)
	}
}

func TestParseNodeExpandedShapeIgnoresUnsupportedShape(t *testing.T) {
	node := parseNode(`A@{ shape: cloud }`)

	if node.name != "A" {
		t.Fatalf("name = %q, want A", node.name)
	}
	if node.hasShape {
		t.Fatal("expected unsupported expanded shape to leave shape unset")
	}
	if len(node.label.lines) != 1 || node.label.lines[0] != "A" {
		t.Fatalf("label lines = %#v, want [A]", node.label.lines)
	}
}

func TestMermaidFileToMapPreservesEscapedLabelNewlines(t *testing.T) {
	properties, err := mermaidFileToMap("graph LR\\nA[\"line1\\nline2\"] --> B", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	spec := properties.nodeSpecs["A"]
	if len(spec.label.lines) != 2 {
		t.Fatalf("label lines = %d, want 2", len(spec.label.lines))
	}
	if spec.label.lines[0] != "line1" || spec.label.lines[1] != "line2" {
		t.Fatalf("label lines = %#v, want [line1 line2]", spec.label.lines)
	}
}

func TestMermaidFileToMapPreservesLiteralLabelNewlines(t *testing.T) {
	properties, err := mermaidFileToMap("graph LR\nA[\"line1\nline2\"] --> B", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	spec := properties.nodeSpecs["A"]
	if len(spec.label.lines) != 2 {
		t.Fatalf("label lines = %d, want 2", len(spec.label.lines))
	}
	if spec.label.lines[0] != "line1" || spec.label.lines[1] != "line2" {
		t.Fatalf("label lines = %#v, want [line1 line2]", spec.label.lines)
	}
}

func TestParseSubgraphHeader(t *testing.T) {
	tests := []struct {
		name      string
		header    string
		wantID    string
		wantLabel string
	}{
		{
			name:      "plain title",
			header:    "Frontend",
			wantID:    "",
			wantLabel: "Frontend",
		},
		{
			name:      "explicit id and title",
			header:    "frontend [Frontend Services]",
			wantID:    "frontend",
			wantLabel: "Frontend Services",
		},
		{
			name:      "quoted title",
			header:    `frontend["Frontend Services"]`,
			wantID:    "frontend",
			wantLabel: "Frontend Services",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sg := parseSubgraphHeader(tt.header)
			if sg.id != tt.wantID {
				t.Fatalf("id = %q, want %q", sg.id, tt.wantID)
			}
			if sg.name != tt.wantLabel {
				t.Fatalf("name = %q, want %q", sg.name, tt.wantLabel)
			}
			if len(sg.label.lines) != 1 || sg.label.lines[0] != tt.wantLabel {
				t.Fatalf("label lines = %#v, want [%s]", sg.label.lines, tt.wantLabel)
			}
		})
	}
}

func TestMermaidFileToMapParsesSubgraphIDAndTitle(t *testing.T) {
	properties, err := mermaidFileToMap("graph LR\nsubgraph frontend [Frontend Services]\nA --> B\nend", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	if len(properties.subgraphs) != 1 {
		t.Fatalf("subgraphs = %d, want 1", len(properties.subgraphs))
	}

	sg := properties.subgraphs[0]
	if sg.id != "frontend" {
		t.Fatalf("id = %q, want %q", sg.id, "frontend")
	}
	if sg.name != "Frontend Services" {
		t.Fatalf("name = %q, want %q", sg.name, "Frontend Services")
	}
}

func TestMermaidFileToMapKeepsExplicitNodeLabelAcrossBareReferences(t *testing.T) {
	properties, err := mermaidFileToMap("graph TD\nA[\"Foo\"] --> B[\"Bar\"]\nB --> C[\"Baz\"]", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	spec := properties.nodeSpecs["B"]
	if len(spec.label.lines) != 1 || spec.label.lines[0] != "Bar" {
		t.Fatalf("label lines = %#v, want [Bar]", spec.label.lines)
	}
	if !spec.labelIsExplicit {
		t.Fatal("expected B label to remain explicit")
	}
}

func TestMermaidFileToMapUsesLatestExplicitLabel(t *testing.T) {
	properties, err := mermaidFileToMap("graph TD\nA[\"Old\"] --> B\nA[\"New\"] --> C", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	spec := properties.nodeSpecs["A"]
	if len(spec.label.lines) != 1 || spec.label.lines[0] != "New" {
		t.Fatalf("label lines = %#v, want [New]", spec.label.lines)
	}
	if !spec.labelIsExplicit {
		t.Fatal("expected A label to remain explicit")
	}
}

func TestMermaidFileToMapKeepsExplicitNodeShapeAcrossBareReferences(t *testing.T) {
	properties, err := mermaidFileToMap("graph TD\nA(Rounded) --> B\nA --> C", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	spec := properties.nodeSpecs["A"]
	if spec.shape != graphNodeShapeRounded {
		t.Fatalf("shape = %q, want %q", spec.shape, graphNodeShapeRounded)
	}
	if !spec.shapeIsExplicit {
		t.Fatal("expected A shape to remain explicit")
	}
}

func TestMermaidFileToMapKeepsExpandedNodeShapeAcrossBareReferences(t *testing.T) {
	properties, err := mermaidFileToMap("graph TD\nA@{ shape: rounded, label: \"Research\" } --> B\nA --> C", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	spec := properties.nodeSpecs["A"]
	if spec.shape != graphNodeShapeRounded {
		t.Fatalf("shape = %q, want %q", spec.shape, graphNodeShapeRounded)
	}
	if !spec.shapeIsExplicit {
		t.Fatal("expected A shape to remain explicit")
	}
	if len(spec.label.lines) != 1 || spec.label.lines[0] != "Research" {
		t.Fatalf("label lines = %#v, want [Research]", spec.label.lines)
	}
	if !spec.labelIsExplicit {
		t.Fatal("expected A label to remain explicit")
	}
}

func TestMermaidFileToMapParsesReverseDirections(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{input: "graph RL\nA --> B", want: "RL"},
		{input: "flowchart RL\nA --> B", want: "RL"},
		{input: "graph BT\nA --> B", want: "BT"},
		{input: "flowchart BT\nA --> B", want: "BT"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			properties, err := mermaidFileToMap(tt.input, "cli")
			if err != nil {
				t.Fatalf("mermaidFileToMap() error = %v", err)
			}
			if properties.graphDirection != tt.want {
				t.Fatalf("graphDirection = %q, want %q", properties.graphDirection, tt.want)
			}
		})
	}
}

func TestMermaidFileToMapParsesEdgeLineStyles(t *testing.T) {
	properties, err := mermaidFileToMap("graph LR\nA ==> B\nA ==>|heavy| C\nA -.-> D\nA -.->|dash| E\nA -.- F\nA --- G\nA ---|open| H", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	edges, ok := properties.data.Get("A")
	if !ok {
		t.Fatal("expected node A to have outgoing edges")
	}

	byChild := map[string]textEdge{}
	for _, edge := range edges {
		byChild[edge.child.name] = edge
	}

	tests := []struct {
		child            string
		wantStyle        graphEdgeLineStyle
		wantLabel        string
		wantHasArrowHead bool
	}{
		{child: "B", wantStyle: graphEdgeLineStyleHeavy, wantHasArrowHead: true},
		{child: "C", wantStyle: graphEdgeLineStyleHeavy, wantLabel: "heavy", wantHasArrowHead: true},
		{child: "D", wantStyle: graphEdgeLineStyleDashed, wantHasArrowHead: true},
		{child: "E", wantStyle: graphEdgeLineStyleDashed, wantLabel: "dash", wantHasArrowHead: true},
		{child: "F", wantStyle: graphEdgeLineStyleDashed, wantHasArrowHead: false},
		{child: "G", wantStyle: graphEdgeLineStyleLight, wantHasArrowHead: false},
		{child: "H", wantStyle: graphEdgeLineStyleLight, wantLabel: "open", wantHasArrowHead: false},
	}

	for _, tt := range tests {
		t.Run(tt.child, func(t *testing.T) {
			edge, ok := byChild[tt.child]
			if !ok {
				t.Fatalf("missing edge to %s", tt.child)
			}
			if edge.lineStyle != tt.wantStyle {
				t.Fatalf("lineStyle = %q, want %q", edge.lineStyle, tt.wantStyle)
			}
			if edge.label != tt.wantLabel {
				t.Fatalf("label = %q, want %q", edge.label, tt.wantLabel)
			}
			if edge.hasArrowHead != tt.wantHasArrowHead {
				t.Fatalf("hasArrowHead = %v, want %v", edge.hasArrowHead, tt.wantHasArrowHead)
			}
		})
	}
}

func TestMermaidFileToMapParsesClassAndLinkStyles(t *testing.T) {
	properties, err := mermaidFileToMap("graph LR\nA -->|go| B\nclass A warning\nclassDef warning stroke:#ff0000,color:#00ff00,fill:#111111;\nlinkStyle 0 stroke:#0000ff,color:#ff00ff;", "cli")
	if err != nil {
		t.Fatalf("mermaidFileToMap() error = %v", err)
	}

	spec := properties.nodeSpecs["A"]
	if spec.styleClass != "warning" {
		t.Fatalf("styleClass = %q, want warning", spec.styleClass)
	}

	nodeStyle := (*properties.styleClasses)["warning"].Styles
	if nodeStyle["stroke"] != "#ff0000" {
		t.Fatalf("node stroke = %q, want #ff0000", nodeStyle["stroke"])
	}
	if nodeStyle["color"] != "#00ff00" {
		t.Fatalf("node color = %q, want #00ff00", nodeStyle["color"])
	}
	if nodeStyle["fill"] != "#111111" {
		t.Fatalf("node fill = %q, want #111111", nodeStyle["fill"])
	}

	edgeStyle := properties.edgeStyles[0].Styles
	if edgeStyle["stroke"] != "#0000ff" {
		t.Fatalf("edge stroke = %q, want #0000ff", edgeStyle["stroke"])
	}
	if edgeStyle["color"] != "#ff00ff" {
		t.Fatalf("edge color = %q, want #ff00ff", edgeStyle["color"])
	}
}
