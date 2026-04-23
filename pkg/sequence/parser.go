package sequence

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const (
	SequenceDiagramKeyword = "sequenceDiagram"
	SolidArrowSyntax       = "->>"
	DottedArrowSyntax      = "-->>"
)

var (
	// participantRegex matches participant declarations: participant [ID] [as Label]
	participantRegex = regexp.MustCompile(`^\s*participant\s+(?:"([^"]+)"|(\S+))(?:\s+as\s+(.+))?$`)

	// actorRegex matches actor declarations: actor [ID] [as Label]
	actorRegex = regexp.MustCompile(`^\s*actor\s+(?:"([^"]+)"|(\S+))(?:\s+as\s+(.+))?$`)

	// messageRegex matches messages: [From]->>[To]: [Label]
	messageRegex = regexp.MustCompile(`^\s*(?:"([^"]+)"|([^\s\->]+))\s*(-->>|->>)\s*(?:"([^"]+)"|([^\s\->]+))\s*:\s*(.*)$`)

	// noteRegex matches notes: Note left of A: text, Note right of A: text, Note over A,B: text
	noteRegex = regexp.MustCompile(`^\s*Note\s+(left of|right of|over)\s+(.+?)\s*:\s*(.*)$`)

	// autonumberRegex matches the autonumber directive
	autonumberRegex = regexp.MustCompile(`^\s*autonumber\s*$`)
)

// SequenceDiagram represents a parsed sequence diagram.
type SequenceDiagram struct {
	Participants []*Participant
	Messages     []*Message
	Notes        []*Note
	Items        []*SequenceItem
	Autonumber   bool
}

type Participant struct {
	ID    string
	Label string
	Index int
}

type Message struct {
	From      *Participant
	To        *Participant
	Label     string
	ArrowType ArrowType
	Number    int // Message number when autonumber is enabled (0 means no number)
}

type SequenceItem struct {
	Message *Message
	Note    *Note
}

type NotePosition string

const (
	NoteLeftOf  NotePosition = "left of"
	NoteRightOf NotePosition = "right of"
	NoteOver    NotePosition = "over"
)

type Note struct {
	Position     NotePosition
	Participants []*Participant
	Text         string
}

type ArrowType int

const (
	SolidArrow ArrowType = iota
	DottedArrow
)

func (a ArrowType) String() string {
	switch a {
	case SolidArrow:
		return "solid"
	case DottedArrow:
		return "dotted"
	default:
		return fmt.Sprintf("ArrowType(%d)", a)
	}
}

func IsSequenceDiagram(input string) bool {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, SequenceDiagramKeyword)
	}
	return false
}

func Parse(input string) (*SequenceDiagram, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	rawLines := diagram.SplitLines(input)
	lines := diagram.RemoveComments(rawLines)
	if len(lines) == 0 {
		return nil, fmt.Errorf("no content found")
	}

	if !strings.HasPrefix(strings.TrimSpace(lines[0]), SequenceDiagramKeyword) {
		return nil, fmt.Errorf("expected %q keyword", SequenceDiagramKeyword)
	}
	lines = lines[1:]

	sd := &SequenceDiagram{
		Participants: []*Participant{},
		Messages:     []*Message{},
		Notes:        []*Note{},
		Items:        []*SequenceItem{},
		Autonumber:   false,
	}
	participantMap := make(map[string]*Participant)

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Check for autonumber directive
		if autonumberRegex.MatchString(trimmed) {
			sd.Autonumber = true
			continue
		}

		if matched, err := sd.parseParticipant(trimmed, participantMap); err != nil {
			return nil, fmt.Errorf("line %d: %w", i+2, err)
		} else if matched {
			continue
		}

		if matched, err := sd.parseNote(trimmed, participantMap); err != nil {
			return nil, fmt.Errorf("line %d: %w", i+2, err)
		} else if matched {
			continue
		}

		if matched, err := sd.parseMessage(trimmed, participantMap); err != nil {
			return nil, fmt.Errorf("line %d: %w", i+2, err)
		} else if matched {
			continue
		}

		return nil, fmt.Errorf("line %d: invalid syntax: %q", i+2, trimmed)
	}

	if len(sd.Participants) == 0 {
		return nil, fmt.Errorf("no participants found")
	}

	return sd, nil
}

func (sd *SequenceDiagram) parseParticipant(line string, participants map[string]*Participant) (bool, error) {
	match := participantRegex.FindStringSubmatch(line)
	if match == nil {
		match = actorRegex.FindStringSubmatch(line)
	}
	if match == nil {
		return false, nil
	}

	id := match[2]
	if match[1] != "" {
		id = match[1]
	}
	label := match[3]
	if label == "" {
		label = id
	}
	label = strings.Trim(label, `"`)

	if _, exists := participants[id]; exists {
		return true, fmt.Errorf("duplicate participant %q", id)
	}

	p := &Participant{
		ID:    id,
		Label: label,
		Index: len(sd.Participants),
	}
	sd.Participants = append(sd.Participants, p)
	participants[id] = p
	return true, nil
}

func (sd *SequenceDiagram) parseNote(line string, participants map[string]*Participant) (bool, error) {
	match := noteRegex.FindStringSubmatch(line)
	if match == nil {
		return false, nil
	}

	position := NotePosition(match[1])
	participantIDs := splitNoteParticipants(match[2])
	if len(participantIDs) == 0 {
		return true, fmt.Errorf("note must reference at least one participant")
	}
	if position != NoteOver && len(participantIDs) > 1 {
		return true, fmt.Errorf("note %s supports exactly one participant", position)
	}

	noteParticipants := make([]*Participant, 0, len(participantIDs))
	for _, id := range participantIDs {
		noteParticipants = append(noteParticipants, sd.getParticipant(id, participants))
	}

	note := &Note{
		Position:     position,
		Participants: noteParticipants,
		Text:         strings.TrimSpace(match[3]),
	}
	sd.Notes = append(sd.Notes, note)
	sd.Items = append(sd.Items, &SequenceItem{Note: note})
	return true, nil
}

func splitNoteParticipants(raw string) []string {
	ids := []string{}
	for _, part := range strings.Split(raw, ",") {
		id := strings.TrimSpace(part)
		id = strings.Trim(id, `"`)
		if id != "" {
			ids = append(ids, id)
		}
	}
	return ids
}

func (sd *SequenceDiagram) parseMessage(line string, participants map[string]*Participant) (bool, error) {
	match := messageRegex.FindStringSubmatch(line)
	if match == nil {
		return false, nil
	}

	fromID := match[2]
	if match[1] != "" {
		fromID = match[1]
	}

	arrow := match[3]

	toID := match[5]
	if match[4] != "" {
		toID = match[4]
	}

	label := strings.TrimSpace(match[6])

	from := sd.getParticipant(fromID, participants)
	to := sd.getParticipant(toID, participants)

	aType := DottedArrow
	if arrow == SolidArrowSyntax {
		aType = SolidArrow
	}

	msgNumber := 0
	if sd.Autonumber {
		msgNumber = len(sd.Messages) + 1
	}

	msg := &Message{
		From:      from,
		To:        to,
		Label:     label,
		ArrowType: aType,
		Number:    msgNumber,
	}
	sd.Messages = append(sd.Messages, msg)
	sd.Items = append(sd.Items, &SequenceItem{Message: msg})
	return true, nil
}

func (sd *SequenceDiagram) getParticipant(id string, participants map[string]*Participant) *Participant {
	if p, exists := participants[id]; exists {
		return p
	}

	p := &Participant{
		ID:    id,
		Label: id,
		Index: len(sd.Participants),
	}
	sd.Participants = append(sd.Participants, p)
	participants[id] = p
	return p
}
