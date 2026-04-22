package cmd

import "testing"

func TestGraphCharsetMergesMixedLineFamilies(t *testing.T) {
	charset := unicodeLightGraphCharset

	tests := []struct {
		name string
		a    string
		b    string
		want string
	}{
		{name: "light horizontal with heavy vertical", a: "─", b: "┃", want: "╋"},
		{name: "light vertical with heavy horizontal", a: "│", b: "━", want: "╋"},
		{name: "light corner with dashed horizontal", a: "┌", b: "┄", want: "┬"},
		{name: "double vertical with heavy horizontal", a: "║", b: "━", want: "╋"},
		{name: "rounded corner with heavy horizontal", a: "╭", b: "━", want: "┳"},
		{name: "light horizontal with heavy horizontal", a: "─", b: "━", want: "━"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := charset.mergeJunctions(tt.a, tt.b); got != tt.want {
				t.Fatalf("mergeJunctions(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.want)
			}
		})
	}
}

func TestEdgeLineStylePriority(t *testing.T) {
	if edgeLineStylePriority(graphEdgeLineStyleHeavy) <= edgeLineStylePriority(graphEdgeLineStyleLight) {
		t.Fatal("expected heavy edges to have higher priority than light edges")
	}
	if edgeLineStylePriority(graphEdgeLineStyleLight) <= edgeLineStylePriority(graphEdgeLineStyleDashed) {
		t.Fatal("expected light edges to have higher priority than dashed edges")
	}
}

func TestGraphCharsetDoesNotTreatTextAsJunction(t *testing.T) {
	charset := unicodeLightGraphCharset

	if charset.isJunctionChar("A") {
		t.Fatal("expected text to be non-junction")
	}
	if !charset.isJunctionChar("━") {
		t.Fatal("expected heavy line to be junction-capable")
	}
	if got := charset.mergeJunctions("A", "━"); got != "A" {
		t.Fatalf("mergeJunctions with text = %q, want A", got)
	}
}

func TestMergeDrawingsKeepsHigherPriorityLineGlyph(t *testing.T) {
	g := graph{}

	lightBase := mkDrawing(0, 0)
	(*lightBase)[0][0] = "─"
	heavyOverlay := mkDrawing(0, 0)
	(*heavyOverlay)[0][0] = "━"
	merged := g.mergeDrawings(lightBase, drawingCoord{0, 0}, heavyOverlay)
	if got := (*merged)[0][0]; got != "━" {
		t.Fatalf("heavy overlay on light line = %q, want heavy line", got)
	}

	heavyBase := mkDrawing(0, 0)
	(*heavyBase)[0][0] = "━"
	lightOverlay := mkDrawing(0, 0)
	(*lightOverlay)[0][0] = "─"
	merged = g.mergeDrawings(heavyBase, drawingCoord{0, 0}, lightOverlay)
	if got := (*merged)[0][0]; got != "━" {
		t.Fatalf("light overlay on heavy line = %q, want heavy line", got)
	}
}
