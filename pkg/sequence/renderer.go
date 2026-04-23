package sequence

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/mattn/go-runewidth"
)

const (
	defaultSelfMessageWidth   = 4
	defaultMessageSpacing     = 1
	defaultParticipantSpacing = 5
	boxPaddingLeftRight       = 2
	minBoxWidth               = 3
	boxBorderWidth            = 2
	labelLeftMargin           = 2
	labelBufferSpace          = 10
)

type diagramLayout struct {
	participantWidths  []int
	participantCenters []int
	totalWidth         int
	messageSpacing     int
	selfMessageWidth   int
}

func calculateLayout(sd *SequenceDiagram, config *diagram.Config) *diagramLayout {
	participantSpacing := config.SequenceParticipantSpacing
	if participantSpacing <= 0 {
		participantSpacing = defaultParticipantSpacing
	}

	widths := make([]int, len(sd.Participants))
	for i, p := range sd.Participants {
		w := runewidth.StringWidth(p.Label) + boxPaddingLeftRight
		if w < minBoxWidth {
			w = minBoxWidth
		}
		widths[i] = w
	}

	centers := make([]int, len(sd.Participants))
	currentX := 0
	for i := range sd.Participants {
		boxWidth := widths[i] + boxBorderWidth
		if i == 0 {
			centers[i] = boxWidth / 2
			currentX = boxWidth
		} else {
			currentX += participantSpacing
			centers[i] = currentX + boxWidth/2
			currentX += boxWidth
		}
	}

	last := len(sd.Participants) - 1
	totalWidth := centers[last] + (widths[last]+boxBorderWidth)/2

	msgSpacing := config.SequenceMessageSpacing
	if msgSpacing <= 0 {
		msgSpacing = defaultMessageSpacing
	}
	selfWidth := config.SequenceSelfMessageWidth
	if selfWidth <= 0 {
		selfWidth = defaultSelfMessageWidth
	}

	return &diagramLayout{
		participantWidths:  widths,
		participantCenters: centers,
		totalWidth:         totalWidth,
		messageSpacing:     msgSpacing,
		selfMessageWidth:   selfWidth,
	}
}

func Render(sd *SequenceDiagram, config *diagram.Config) (string, error) {
	if sd == nil || len(sd.Participants) == 0 {
		return "", fmt.Errorf("no participants")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}

	chars := Unicode
	if config.UseAscii {
		chars = ASCII
	}

	layout := calculateLayout(sd, config)
	var lines []string

	lines = append(lines, buildLine(sd.Participants, layout, func(i int) string {
		return string(chars.TopLeft) + strings.Repeat(string(chars.Horizontal), layout.participantWidths[i]) + string(chars.TopRight)
	}))

	lines = append(lines, buildLine(sd.Participants, layout, func(i int) string {
		w := layout.participantWidths[i]
		labelLen := runewidth.StringWidth(sd.Participants[i].Label)
		pad := (w - labelLen) / 2
		return string(chars.Vertical) + strings.Repeat(" ", pad) + sd.Participants[i].Label +
			strings.Repeat(" ", w-pad-labelLen) + string(chars.Vertical)
	}))

	lines = append(lines, buildLine(sd.Participants, layout, func(i int) string {
		w := layout.participantWidths[i]
		return string(chars.BottomLeft) + strings.Repeat(string(chars.Horizontal), w/2) +
			string(chars.TeeDown) + strings.Repeat(string(chars.Horizontal), w-w/2-1) +
			string(chars.BottomRight)
	}))

	items := sd.Items
	if len(items) == 0 {
		items = make([]*SequenceItem, 0, len(sd.Messages))
		for _, msg := range sd.Messages {
			items = append(items, &SequenceItem{Message: msg})
		}
	}

	active := map[int]bool{}
	for _, item := range items {
		for i := 0; i < layout.messageSpacing; i++ {
			lines = append(lines, buildLifelineWithActivation(layout, chars, active))
		}

		switch {
		case item.Message != nil && item.Message.From == item.Message.To:
			lines = append(lines, renderSelfMessage(item.Message, layout, chars, active)...)
			applyMessageActivation(item.Message, active)
		case item.Message != nil:
			lines = append(lines, renderMessage(item.Message, layout, chars, active)...)
			applyMessageActivation(item.Message, active)
		case item.Note != nil:
			lines = append(lines, renderNote(item.Note, layout, chars, active)...)
		case item.Activation != nil:
			if item.Activation.Active {
				active[item.Activation.Participant.Index] = true
				lines = append(lines, buildLifelineWithActivation(layout, chars, active))
			} else {
				lines = append(lines, buildLifelineWithActivation(layout, chars, active))
				delete(active, item.Activation.Participant.Index)
			}
		}
	}

	lines = append(lines, buildLifelineWithActivation(layout, chars, active))
	return strings.Join(lines, "\n") + "\n", nil
}

func applyMessageActivation(msg *Message, active map[int]bool) {
	if msg.PostActivation != nil {
		active[msg.PostActivation.Participant.Index] = true
	}
	if msg.PostDeactivation != nil {
		delete(active, msg.PostDeactivation.Participant.Index)
	}
}

func buildLine(participants []*Participant, layout *diagramLayout, draw func(int) string) string {
	var sb strings.Builder
	for i := range participants {
		boxWidth := layout.participantWidths[i] + boxBorderWidth
		left := layout.participantCenters[i] - boxWidth/2

		needed := left - len([]rune(sb.String()))
		if needed > 0 {
			sb.WriteString(strings.Repeat(" ", needed))
		}
		sb.WriteString(draw(i))
	}
	return sb.String()
}

func buildLifeline(layout *diagramLayout, chars BoxChars) string {
	return buildLifelineWithActivation(layout, chars, nil)
}

func buildLifelineWithActivation(layout *diagramLayout, chars BoxChars, active map[int]bool) string {
	line := make([]rune, layout.totalWidth+1)
	for i := range line {
		line[i] = ' '
	}
	for i, c := range layout.participantCenters {
		if c < len(line) {
			line[c] = chars.Vertical
			if active != nil && active[i] {
				line[c] = chars.Activation
			}
		}
	}
	return strings.TrimRight(string(line), " ")
}

func renderNote(note *Note, layout *diagramLayout, chars BoxChars, active map[int]bool) []string {
	left, width := noteBounds(note, layout)
	minWidth := runewidth.StringWidth(note.Text) + boxBorderWidth + 2
	if width < minWidth {
		width = minWidth
	}

	top := string(chars.TopLeft) + strings.Repeat(string(chars.Horizontal), width-2) + string(chars.TopRight)
	textWidth := runewidth.StringWidth(note.Text)
	leftPad := (width - boxBorderWidth - textWidth) / 2
	rightPad := width - boxBorderWidth - textWidth - leftPad
	middle := string(chars.Vertical) + strings.Repeat(" ", leftPad) + note.Text + strings.Repeat(" ", rightPad) + string(chars.Vertical)
	bottom := string(chars.BottomLeft) + strings.Repeat(string(chars.Horizontal), width-2) + string(chars.BottomRight)

	return []string{
		placeOverlay(buildLifelineWithActivation(layout, chars, active), left, top),
		placeOverlay(buildLifelineWithActivation(layout, chars, active), left, middle),
		placeOverlay(buildLifelineWithActivation(layout, chars, active), left, bottom),
	}
}

func noteBounds(note *Note, layout *diagramLayout) (int, int) {
	firstCenter := layout.participantCenters[note.Participants[0].Index]
	lastCenter := firstCenter
	for _, p := range note.Participants[1:] {
		center := layout.participantCenters[p.Index]
		if center < firstCenter {
			firstCenter = center
		}
		if center > lastCenter {
			lastCenter = center
		}
	}

	textWidth := runewidth.StringWidth(note.Text) + boxBorderWidth + 2
	switch note.Position {
	case NoteLeftOf:
		left := firstCenter - textWidth
		if left < 0 {
			left = 0
		}
		return left, textWidth
	case NoteRightOf:
		return firstCenter + 1, textWidth
	default:
		left := firstCenter - 1
		if left < 0 {
			left = 0
		}
		width := lastCenter - left + 2
		if width < textWidth {
			width = textWidth
		}
		return left, width
	}
}

func placeOverlay(base string, left int, overlay string) string {
	line := []rune(base)
	width := left + len([]rune(overlay))
	if len(line) < width {
		padding := make([]rune, width-len(line))
		for i := range padding {
			padding[i] = ' '
		}
		line = append(line, padding...)
	}
	for i, r := range []rune(overlay) {
		line[left+i] = r
	}
	return strings.TrimRight(string(line), " ")
}

func renderMessage(msg *Message, layout *diagramLayout, chars BoxChars, active map[int]bool) []string {
	var lines []string
	from, to := layout.participantCenters[msg.From.Index], layout.participantCenters[msg.To.Index]

	label := msg.Label
	if msg.Number > 0 {
		label = fmt.Sprintf("%d. %s", msg.Number, msg.Label)
	}

	if label != "" {
		start := min(from, to) + labelLeftMargin
		labelWidth := runewidth.StringWidth(label)
		w := max(layout.totalWidth, start+labelWidth) + labelBufferSpace
		line := []rune(buildLifelineWithActivation(layout, chars, active))
		if len(line) < w {
			padding := make([]rune, w-len(line))
			for k := range padding {
				padding[k] = ' '
			}
			line = append(line, padding...)
		}

		col := start
		for _, r := range label {
			if col < len(line) {
				line[col] = r
				col++
			}
		}
		lines = append(lines, strings.TrimRight(string(line), " "))
	}

	line := []rune(buildLifelineWithActivation(layout, chars, active))
	style := chars.SolidLine
	if msg.ArrowType == DottedArrow {
		style = chars.DottedLine
	}

	if from < to {
		line[from] = chars.TeeRight
		for i := from + 1; i < to; i++ {
			line[i] = style
		}
		line[to-1] = chars.ArrowRight
		line[to] = chars.Vertical
	} else {
		line[to] = chars.Vertical
		line[to+1] = chars.ArrowLeft
		for i := to + 2; i < from; i++ {
			line[i] = style
		}
		line[from] = chars.TeeLeft
	}
	lines = append(lines, strings.TrimRight(string(line), " "))
	return lines
}

func renderSelfMessage(msg *Message, layout *diagramLayout, chars BoxChars, active map[int]bool) []string {
	var lines []string
	center := layout.participantCenters[msg.From.Index]
	width := layout.selfMessageWidth

	ensureWidth := func(l string) []rune {
		target := layout.totalWidth + width + 1
		r := []rune(l)
		if len(r) < target {
			pad := make([]rune, target-len(r))
			for i := range pad {
				pad[i] = ' '
			}
			r = append(r, pad...)
		}
		return r
	}

	label := msg.Label
	if msg.Number > 0 {
		label = fmt.Sprintf("%d. %s", msg.Number, msg.Label)
	}

	if label != "" {
		line := ensureWidth(buildLifelineWithActivation(layout, chars, active))
		start := center + labelLeftMargin
		labelWidth := runewidth.StringWidth(label)
		needed := start + labelWidth + labelBufferSpace
		if len(line) < needed {
			pad := make([]rune, needed-len(line))
			for i := range pad {
				pad[i] = ' '
			}
			line = append(line, pad...)
		}
		col := start
		for _, c := range label {
			if col < len(line) {
				line[col] = c
				col++
			}
		}
		lines = append(lines, strings.TrimRight(string(line), " "))
	}

	l1 := ensureWidth(buildLifelineWithActivation(layout, chars, active))
	l1[center] = chars.TeeRight
	for i := 1; i < width; i++ {
		l1[center+i] = chars.Horizontal
	}
	l1[center+width-1] = chars.SelfTopRight
	lines = append(lines, strings.TrimRight(string(l1), " "))

	l2 := ensureWidth(buildLifelineWithActivation(layout, chars, active))
	l2[center+width-1] = chars.Vertical
	lines = append(lines, strings.TrimRight(string(l2), " "))

	l3 := ensureWidth(buildLifelineWithActivation(layout, chars, active))
	l3[center] = chars.Vertical
	l3[center+1] = chars.ArrowLeft
	for i := 2; i < width-1; i++ {
		l3[center+i] = chars.Horizontal
	}
	l3[center+width-1] = chars.SelfBottom
	lines = append(lines, strings.TrimRight(string(l3), " "))

	return lines
}
