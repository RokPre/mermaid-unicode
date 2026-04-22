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
		{name: "light horizontal with heavy vertical", a: "─", b: "┃", want: "┼"},
		{name: "light vertical with heavy horizontal", a: "│", b: "━", want: "┼"},
		{name: "light corner with dashed horizontal", a: "┌", b: "┄", want: "┬"},
		{name: "double vertical with heavy horizontal", a: "║", b: "━", want: "┼"},
		{name: "rounded corner with heavy horizontal", a: "╭", b: "━", want: "┬"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := charset.mergeJunctions(tt.a, tt.b); got != tt.want {
				t.Fatalf("mergeJunctions(%q, %q) = %q, want %q", tt.a, tt.b, got, tt.want)
			}
		})
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
