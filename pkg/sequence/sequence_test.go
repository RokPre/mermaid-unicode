package sequence

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		wantParticipants int
		wantMessages     int
		wantErr          string
	}{
		{"empty input", "", 0, 0, "empty input"},
		{"missing sequenceDiagram keyword", "A->>B: Hello", 0, 0, "expected \"sequenceDiagram\" keyword"},
		{"only comments", "sequenceDiagram\n%% This is a comment\n%% Another comment", 0, 0, "no participants found"},
		{"no participants", "sequenceDiagram", 0, 0, "no participants found"},
		{"duplicate participant ID", "sequenceDiagram\nparticipant Alice\nparticipant Alice\nAlice->>Bob: Hi", 0, 0, "duplicate participant"},
		{"minimal diagram", "sequenceDiagram\nA->>B: Hello", 2, 1, ""},
		{"explicit participants", "sequenceDiagram\nparticipant Alice\nparticipant Bob\nAlice->>Bob: Hi", 2, 1, ""},
		{"dotted arrow", "sequenceDiagram\nA-->>B: Response", 2, 1, ""},
		{"self message", "sequenceDiagram\nA->>A: Self", 1, 1, ""},
		{"multiple messages", "sequenceDiagram\nA->>B: 1\nB->>C: 2\nC-->>A: 3", 3, 3, ""},
		{"with comments", "sequenceDiagram\n%% Comment\nA->>B: Hi %% inline comment", 2, 1, ""},
		{"with note", "sequenceDiagram\nA->>B: Hi\nNote over A,B: Shared context", 2, 1, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd, err := Parse(tt.input)
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Expected error containing %q, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			if len(sd.Participants) != tt.wantParticipants {
				t.Errorf("Expected %d participants, got %d", tt.wantParticipants, len(sd.Participants))
			}
			if len(sd.Messages) != tt.wantMessages {
				t.Errorf("Expected %d messages, got %d", tt.wantMessages, len(sd.Messages))
			}
		})
	}
}

func TestIsSequenceDiagram(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"sequenceDiagram\nA->>B: Hello", true},
		{"graph LR\nA-->B", false},
		{"graph TD\nA-->B", false},
		{"", false},
		{"%% Just a comment", false},
	}

	for _, tt := range tests {
		if got := IsSequenceDiagram(tt.input); got != tt.want {
			t.Errorf("IsSequenceDiagram(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParticipantAlias(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantID    string
		wantLabel string
	}{
		{"simple alias", "sequenceDiagram\nparticipant A as Alice\nA->>A: Hello", "A", "Alice"},
		{"no alias defaults to id", "sequenceDiagram\nparticipant Alice\nAlice->>Alice: Hi", "Alice", "Alice"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(d.Participants) == 0 {
				t.Fatal("expected at least one participant")
			}
			p := d.Participants[0]
			if p.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", p.ID, tt.wantID)
			}
			if p.Label != tt.wantLabel {
				t.Errorf("Label = %q, want %q", p.Label, tt.wantLabel)
			}
			config := diagram.DefaultConfig()
			output, err := Render(d, config)
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			if !strings.Contains(output, tt.wantLabel) {
				t.Errorf("output should contain label %q", tt.wantLabel)
			}
		})
	}
}

func TestActorDeclaration(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantID    string
		wantLabel string
	}{
		{"actor defaults label to id", "sequenceDiagram\nactor Alice\nAlice->>Alice: Think", "Alice", "Alice"},
		{"actor alias", "sequenceDiagram\nactor A as Alice\nA->>A: Think", "A", "Alice"},
		{"quoted actor", "sequenceDiagram\nactor \"External User\" as User\n\"External User\"->>\"External User\": Think", "External User", "User"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := Parse(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(d.Participants) == 0 {
				t.Fatal("expected at least one participant")
			}
			p := d.Participants[0]
			if p.ID != tt.wantID {
				t.Errorf("ID = %q, want %q", p.ID, tt.wantID)
			}
			if p.Label != tt.wantLabel {
				t.Errorf("Label = %q, want %q", p.Label, tt.wantLabel)
			}

			output, err := Render(d, diagram.DefaultConfig())
			if err != nil {
				t.Fatalf("render error: %v", err)
			}
			if !strings.Contains(output, tt.wantLabel) {
				t.Errorf("output should contain actor label %q", tt.wantLabel)
			}
		})
	}
}

func TestMessageRegex(t *testing.T) {
	tests := []struct {
		input     string
		wantFrom  string
		wantArrow string
		wantAct   string
		wantTo    string
		wantLabel string
		wantMatch bool
	}{
		{"A->>B: Hello", "A", "->>", "", "B", "Hello", true},
		{"A-->>B: Response", "A", "-->>", "", "B", "Response", true},
		{"A->>+B: Start", "A", "->>", "+", "B", "Start", true},
		{"B-->>-A: Done", "B", "-->>", "-", "A", "Done", true},
		{`"My Service"->>B: Test`, "My Service", "->>", "", "B", "Test", true},
		{"A->>B: ", "A", "->>", "", "B", "", true},
		{"A->B: Test", "", "", "", "", "", false},
		{"A->>B", "", "", "", "", "", false},
	}

	for _, tt := range tests {
		match := messageRegex.FindStringSubmatch(tt.input)
		if !tt.wantMatch {
			if match != nil {
				t.Errorf("messageRegex should not match %q", tt.input)
			}
			continue
		}
		if match == nil {
			t.Fatalf("messageRegex failed to match: %q", tt.input)
		}
		gotFrom := match[2]
		if match[1] != "" {
			gotFrom = match[1]
		}
		gotArrow := match[3]
		gotAct := match[4]
		gotTo := match[6]
		if match[5] != "" {
			gotTo = match[5]
		}
		gotLabel := match[7]

		if gotFrom != tt.wantFrom || gotArrow != tt.wantArrow || gotAct != tt.wantAct || gotTo != tt.wantTo || gotLabel != tt.wantLabel {
			t.Errorf("messageRegex(%q) = (%q, %q, %q, %q, %q), want (%q, %q, %q, %q, %q)",
				tt.input, gotFrom, gotArrow, gotAct, gotTo, gotLabel, tt.wantFrom, tt.wantArrow, tt.wantAct, tt.wantTo, tt.wantLabel)
		}
	}
}

func TestSequenceNotes(t *testing.T) {
	input := `sequenceDiagram
participant A as Alice
participant B as Bob
Note left of A: Local note
A->>B: Hello
Note over A,B: Shared context
Note right of B: Remote note`

	d, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.Notes) != 3 {
		t.Fatalf("notes = %d, want 3", len(d.Notes))
	}
	if len(d.Items) != 4 {
		t.Fatalf("items = %d, want 4", len(d.Items))
	}
	if d.Notes[0].Position != NoteLeftOf || d.Notes[0].Text != "Local note" {
		t.Fatalf("first note = %#v, want left local note", d.Notes[0])
	}
	if len(d.Notes[1].Participants) != 2 || d.Notes[1].Position != NoteOver {
		t.Fatalf("second note = %#v, want note over two participants", d.Notes[1])
	}
	if d.Notes[2].Position != NoteRightOf || d.Notes[2].Participants[0].ID != "B" {
		t.Fatalf("third note = %#v, want right of B", d.Notes[2])
	}

	output, err := Render(d, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	for _, want := range []string{"Local note", "Hello", "Shared context", "Remote note"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestSequenceActivations(t *testing.T) {
	input := `sequenceDiagram
participant A as Alice
participant B as Bob
activate B
A->>B: Work
deactivate B
B-->>A: Done`

	d, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.Activations) != 2 {
		t.Fatalf("activations = %d, want 2", len(d.Activations))
	}
	if len(d.Items) != 4 {
		t.Fatalf("items = %d, want 4", len(d.Items))
	}
	if !d.Activations[0].Active || d.Activations[0].Participant.ID != "B" {
		t.Fatalf("first activation = %#v, want activate B", d.Activations[0])
	}
	if d.Activations[1].Active || d.Activations[1].Participant.ID != "B" {
		t.Fatalf("second activation = %#v, want deactivate B", d.Activations[1])
	}

	output, err := Render(d, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	for _, want := range []string{"Work", "Done", "┃"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestSequenceActivationShorthand(t *testing.T) {
	input := `sequenceDiagram
A->>+B: Start
B-->>-A: Done`

	d, err := Parse(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(d.Activations) != 2 {
		t.Fatalf("activations = %d, want 2", len(d.Activations))
	}
	first := d.Messages[0]
	if first.PostActivation == nil || first.PostActivation.Participant.ID != "B" {
		t.Fatalf("first message post activation = %#v, want activate B", first.PostActivation)
	}
	second := d.Messages[1]
	if second.PostDeactivation == nil || second.PostDeactivation.Participant.ID != "B" {
		t.Fatalf("second message post deactivation = %#v, want deactivate B", second.PostDeactivation)
	}

	output, err := Render(d, diagram.DefaultConfig())
	if err != nil {
		t.Fatalf("render error: %v", err)
	}
	for _, want := range []string{"Start", "Done", "┃"} {
		if !strings.Contains(output, want) {
			t.Fatalf("output missing %q\noutput:\n%s", want, output)
		}
	}
}

func TestActivationRegex(t *testing.T) {
	tests := []struct {
		input      string
		wantAction string
		wantID     string
	}{
		{"activate Alice", "activate", "Alice"},
		{"deactivate Alice", "deactivate", "Alice"},
		{`activate "External User"`, "activate", "External User"},
	}

	for _, tt := range tests {
		match := activationRegex.FindStringSubmatch(tt.input)
		if match == nil {
			t.Fatalf("activationRegex failed to match %q", tt.input)
		}
		gotID := match[3]
		if match[2] != "" {
			gotID = match[2]
		}
		if match[1] != tt.wantAction || gotID != tt.wantID {
			t.Fatalf("activationRegex(%q) = action %q id %q, want action %q id %q",
				tt.input, match[1], gotID, tt.wantAction, tt.wantID)
		}
	}
}

func TestNoteRegex(t *testing.T) {
	tests := []struct {
		input        string
		wantPos      string
		wantTargets  string
		wantNoteText string
	}{
		{"Note left of A: Local note", "left of", "A", "Local note"},
		{"Note right of A: Remote note", "right of", "A", "Remote note"},
		{"Note over A,B: Shared context", "over", "A,B", "Shared context"},
	}

	for _, tt := range tests {
		match := noteRegex.FindStringSubmatch(tt.input)
		if match == nil {
			t.Fatalf("noteRegex failed to match %q", tt.input)
		}
		if match[1] != tt.wantPos || match[2] != tt.wantTargets || match[3] != tt.wantNoteText {
			t.Fatalf("noteRegex(%q) = (%q, %q, %q), want (%q, %q, %q)",
				tt.input, match[1], match[2], match[3], tt.wantPos, tt.wantTargets, tt.wantNoteText)
		}
	}
}

func TestParticipantRegex(t *testing.T) {
	tests := []struct {
		input     string
		wantID    string
		wantAlias string
	}{
		{"participant Alice", "Alice", ""},
		{"participant Alice as A", "Alice", "A"},
		{`participant "My Service"`, "My Service", ""},
		{`participant "My Service" as Service`, "My Service", "Service"},
	}

	for _, tt := range tests {
		match := participantRegex.FindStringSubmatch(tt.input)
		if match == nil {
			t.Fatalf("participantRegex failed to match: %q", tt.input)
		}
		gotID := match[2]
		if match[1] != "" {
			gotID = match[1]
		}
		gotAlias := match[3]

		if gotID != tt.wantID || gotAlias != tt.wantAlias {
			t.Errorf("participantRegex(%q) = (%q, %q), want (%q, %q)",
				tt.input, gotID, gotAlias, tt.wantID, tt.wantAlias)
		}
	}
}

func TestActorRegex(t *testing.T) {
	tests := []struct {
		input     string
		wantID    string
		wantAlias string
	}{
		{"actor Alice", "Alice", ""},
		{"actor Alice as A", "Alice", "A"},
		{`actor "External User"`, "External User", ""},
		{`actor "External User" as User`, "External User", "User"},
	}

	for _, tt := range tests {
		match := actorRegex.FindStringSubmatch(tt.input)
		if match == nil {
			t.Fatalf("actorRegex failed to match: %q", tt.input)
		}
		gotID := match[2]
		if match[1] != "" {
			gotID = match[1]
		}
		gotAlias := match[3]

		if gotID != tt.wantID || gotAlias != tt.wantAlias {
			t.Errorf("actorRegex(%q) = ID %q alias %q, want ID %q alias %q",
				tt.input, gotID, gotAlias, tt.wantID, tt.wantAlias)
		}
	}
}

func TestSplitLines(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"A->>B: Hello", []string{"A->>B: Hello"}},
		{"line1\nline2\nline3", []string{"line1", "line2", "line3"}},
		{"line1\\nline2\\nline3", []string{"line1", "line2", "line3"}},
		{"", []string{""}},
	}

	for _, tt := range tests {
		result := diagram.SplitLines(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("SplitLines(%q) len = %d, want %d", tt.input, len(result), len(tt.expected))
		}
	}
}

func TestRemoveComments(t *testing.T) {
	tests := []struct {
		input    []string
		expected []string
	}{
		{[]string{"A->>B: Hello", "B-->>A: Hi"}, []string{"A->>B: Hello", "B-->>A: Hi"}},
		{[]string{"%% This is a comment", "A->>B: Hello"}, []string{"A->>B: Hello"}},
		{[]string{"A->>B: Hello %% inline comment", "B-->>A: Hi"}, []string{"A->>B: Hello", "B-->>A: Hi"}},
	}

	for _, tt := range tests {
		result := diagram.RemoveComments(tt.input)
		if len(result) != len(tt.expected) {
			t.Errorf("RemoveComments() len = %d, want %d", len(result), len(tt.expected))
		}
	}
}

func TestArrowTypeString(t *testing.T) {
	if SolidArrow.String() != "solid" {
		t.Errorf("SolidArrow.String() = %q, want \"solid\"", SolidArrow.String())
	}
	if DottedArrow.String() != "dotted" {
		t.Errorf("DottedArrow.String() = %q, want \"dotted\"", DottedArrow.String())
	}
}

func FuzzParseSequenceDiagram(f *testing.F) {
	f.Add("sequenceDiagram\nA->>B: Hello")
	f.Add("sequenceDiagram\nparticipant Alice\nAlice->>Bob: Hi")
	f.Add("sequenceDiagram\nA-->>B: Response")
	f.Add("sequenceDiagram\nA->>A: Self")

	f.Fuzz(func(t *testing.T, input string) {
		sd, err := Parse(input)
		if err != nil {
			return
		}

		for i, p := range sd.Participants {
			if p.Index != i {
				t.Errorf("Participant %q has incorrect index: got %d, expected %d", p.ID, p.Index, i)
			}
			if p.ID == "" {
				t.Errorf("Participant at index %d has empty ID", i)
			}
			if p.Label == "" {
				t.Errorf("Participant %q has empty label", p.ID)
			}
		}

		for i, msg := range sd.Messages {
			if msg.From == nil || msg.To == nil {
				t.Errorf("Message %d has nil participant", i)
			}
		}

		seen := make(map[string]bool)
		for _, p := range sd.Participants {
			if seen[p.ID] {
				t.Errorf("Duplicate participant ID: %q", p.ID)
			}
			seen[p.ID] = true
		}

		config := diagram.DefaultConfig()
		_, _ = Render(sd, config)
	})
}

func FuzzRenderSequenceDiagram(f *testing.F) {
	seeds := []string{
		"sequenceDiagram\nA->>B: Test",
		"sequenceDiagram\nA->>A: Self",
		"sequenceDiagram\nA->>B: 1\nB->>C: 2\nC->>A: 3",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, input string) {
		sd, err := Parse(input)
		if err != nil {
			return
		}

		for _, useAscii := range []bool{true, false} {
			config := diagram.DefaultConfig()
			config.UseAscii = useAscii

			output, err := Render(sd, config)
			if err != nil {
				return
			}

			if strings.TrimSpace(output) == "" {
				t.Error("Renderer produced empty output for valid diagram")
			}

			for _, p := range sd.Participants {
				if !strings.Contains(output, p.Label) {
					t.Errorf("Rendered output missing participant label: %q", p.Label)
				}
			}

			if !utf8.ValidString(output) {
				t.Error("Rendered output contains invalid UTF-8")
			}
		}
	})
}

func BenchmarkParse(b *testing.B) {
	tests := []struct {
		name         string
		participants int
		messages     int
	}{
		{"small_2p_5m", 2, 5},
		{"medium_5p_20m", 5, 20},
		{"large_10p_50m", 10, 50},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			input := generateDiagram(tt.participants, tt.messages)
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := Parse(input)
				if err != nil {
					b.Fatalf("parse failed: %v", err)
				}
			}
		})
	}
}

func BenchmarkRender(b *testing.B) {
	tests := []struct {
		name         string
		participants int
		messages     int
	}{
		{"small_2p_5m", 2, 5},
		{"medium_5p_20m", 5, 20},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			input := generateDiagram(tt.participants, tt.messages)
			sd, err := Parse(input)
			if err != nil {
				b.Fatalf("parse failed: %v", err)
			}
			config := diagram.DefaultConfig()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, renderErr := Render(sd, config)
				if renderErr != nil {
					b.Fatalf("render error: %v", renderErr)
				}
			}
		})
	}
}

func generateDiagram(numParticipants, numMessages int) string {
	var sb strings.Builder
	sb.WriteString("sequenceDiagram\n")
	for i := 0; i < numParticipants; i++ {
		sb.WriteString("    participant P")
		sb.WriteString(string(rune('0' + i)))
		sb.WriteString("\n")
	}
	for i := 0; i < numMessages; i++ {
		from := i % numParticipants
		to := (i + 1) % numParticipants
		arrow := "-"
		if i%2 == 0 {
			arrow = "--"
		}
		sb.WriteString("    P")
		sb.WriteString(string(rune('0' + from)))
		sb.WriteString(arrow)
		sb.WriteString(">>P")
		sb.WriteString(string(rune('0' + to)))
		sb.WriteString(": Message\n")
	}
	return sb.String()
}
