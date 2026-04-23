package diagram

import (
	"strings"
	"testing"
)

func TestParseStyleMap(t *testing.T) {
	styles := ParseStyleMap("stroke:#ff0000, color: #00ff00, fill:#111111;")

	if styles["stroke"] != "#ff0000" {
		t.Fatalf("stroke = %q, want #ff0000", styles["stroke"])
	}
	if styles["color"] != "#00ff00" {
		t.Fatalf("color = %q, want #00ff00", styles["color"])
	}
	if styles["fill"] != "#111111" {
		t.Fatalf("fill = %q, want #111111", styles["fill"])
	}
}

func TestResolveStylePrecedence(t *testing.T) {
	resolved := ResolveStyle(
		StyleMap{"stroke": "#111111", "fill": "#eeeeee", "color": "#222222"},
		StyleMap{"stroke": "#333333", "fill": "transparent"},
		StyleMap{"color": "#444444"},
	)

	if resolved["stroke"] != "#333333" {
		t.Fatalf("stroke = %q, want class stroke #333333", resolved["stroke"])
	}
	if resolved["fill"] != "#eeeeee" {
		t.Fatalf("fill = %q, want default fill #eeeeee", resolved["fill"])
	}
	if resolved["color"] != "#444444" {
		t.Fatalf("color = %q, want direct color #444444", resolved["color"])
	}
}

func TestResolveStyleUnstyledOutput(t *testing.T) {
	resolved := ResolveStyle(nil, StyleMap{"stroke": "none"}, StyleMap{"fill": "transparent"})
	if len(resolved) != 0 {
		t.Fatalf("resolved style = %#v, want empty", resolved)
	}

	if got := WrapTextInStyle("text", "", "", StyleTypeHTML); got != "text" {
		t.Fatalf("unstyled text = %q, want text", got)
	}
}

func TestWrapTextInStyleHTML(t *testing.T) {
	got := WrapTextInStyle("A", "#ff0000", "#000000", StyleTypeHTML)
	for _, want := range []string{"<span", "color: #ff0000", "background-color: #000000", ">A</span>"} {
		if !strings.Contains(got, want) {
			t.Fatalf("styled HTML = %q, want containing %q", got, want)
		}
	}
}

func TestNormalizeStyleColor(t *testing.T) {
	tests := map[string]string{
		" #abc; ":     "#abc",
		"none":        "",
		"transparent": "",
		"":            "",
	}
	for input, want := range tests {
		if got := NormalizeStyleColor(input); got != want {
			t.Fatalf("NormalizeStyleColor(%q) = %q, want %q", input, got, want)
		}
	}
}
