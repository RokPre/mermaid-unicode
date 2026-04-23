package statediagram

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const (
	StateDiagramKeyword   = "stateDiagram"
	StateDiagramV2Keyword = "stateDiagram-v2"
	StartEndID            = "[*]"
)

var (
	transitionRegex       = regexp.MustCompile(`^\s*(\[\*\]|[A-Za-z_][A-Za-z0-9_.$-]*)\s+-->\s+(\[\*\]|[A-Za-z_][A-Za-z0-9_.$-]*)(?:\s*:\s*(.*))?$`)
	stateAliasRegex       = regexp.MustCompile(`^\s*state\s+"([^"]+)"\s+as\s+([A-Za-z_][A-Za-z0-9_.$-]*)\s*$`)
	stateBlockRegex       = regexp.MustCompile(`^\s*state\s+([A-Za-z_][A-Za-z0-9_.$-]*)(?:\s+"([^"]+)")?\s*\{\s*$`)
	stateKindRegex        = regexp.MustCompile(`^\s*state\s+([A-Za-z_][A-Za-z0-9_.$-]*)\s+<<(choice|fork|join)>>\s*$`)
	stateDescriptionRegex = regexp.MustCompile(`^\s*state\s+([A-Za-z_][A-Za-z0-9_.$-]*)\s*:\s*(.+)$`)
	stateDeclarationRegex = regexp.MustCompile(`^\s*state\s+([A-Za-z_][A-Za-z0-9_.$-]*)\s*$`)
	descriptionRegex      = regexp.MustCompile(`^\s*([A-Za-z_][A-Za-z0-9_.$-]*)\s*:\s*(.+)$`)
	noteRegex             = regexp.MustCompile(`^\s*note\s+(left of|right of)\s+([A-Za-z_][A-Za-z0-9_.$-]*)\s*:\s*(.*)$`)
	directionRegex        = regexp.MustCompile(`^\s*direction\s+(LR|RL|TB|TD|BT)\s*$`)
)

type Diagram struct {
	States      map[string]*State
	StateOrder  []string
	Transitions []*Transition
	Notes       []*Note
	Composites  []*Composite
	Direction   string
}

type State struct {
	ID          string
	Label       string
	Description string
	Kind        string
}

type Transition struct {
	From  *State
	To    *State
	Label string
}

type Note struct {
	State    *State
	Position string
	Text     string
}

type Composite struct {
	ID       string
	Label    string
	Children []string
}

func IsStateDiagram(input string) bool {
	lines := strings.Split(input, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, StateDiagramKeyword) || strings.HasPrefix(trimmed, StateDiagramV2Keyword)
	}
	return false
}

func Parse(input string) (*Diagram, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}

	rawLines := diagram.SplitLines(input)
	lines := diagram.RemoveComments(rawLines)
	if len(lines) == 0 {
		return nil, fmt.Errorf("no content found")
	}

	first := strings.TrimSpace(lines[0])
	if !strings.HasPrefix(first, StateDiagramKeyword) && !strings.HasPrefix(first, StateDiagramV2Keyword) {
		return nil, fmt.Errorf("expected %q or %q keyword", StateDiagramKeyword, StateDiagramV2Keyword)
	}

	sd := &Diagram{
		States:     map[string]*State{},
		StateOrder: []string{},
		Direction:  "LR",
	}
	compositeStack := []*Composite{}

	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if trimmed == "}" {
			if len(compositeStack) == 0 {
				return nil, fmt.Errorf("line %d: composite state end without matching start", lineNo)
			}
			compositeStack = compositeStack[:len(compositeStack)-1]
			continue
		}

		if match := directionRegex.FindStringSubmatch(trimmed); match != nil {
			sd.Direction = match[1]
			continue
		}
		if match := stateBlockRegex.FindStringSubmatch(trimmed); match != nil {
			state := sd.ensureState(match[1], match[2])
			composite := &Composite{ID: state.ID, Label: state.Label}
			sd.Composites = append(sd.Composites, composite)
			compositeStack = append(compositeStack, composite)
			sd.addToCurrentComposite(compositeStack, state.ID)
			continue
		}
		if match := stateAliasRegex.FindStringSubmatch(trimmed); match != nil {
			state := sd.ensureState(match[2], match[1])
			sd.addToCurrentComposite(compositeStack, state.ID)
			continue
		}
		if match := stateKindRegex.FindStringSubmatch(trimmed); match != nil {
			state := sd.ensureState(match[1], "")
			state.Kind = match[2]
			sd.addToCurrentComposite(compositeStack, state.ID)
			continue
		}
		if match := stateDescriptionRegex.FindStringSubmatch(trimmed); match != nil {
			state := sd.ensureState(match[1], "")
			state.Description = strings.TrimSpace(match[2])
			sd.addToCurrentComposite(compositeStack, state.ID)
			continue
		}
		if match := noteRegex.FindStringSubmatch(trimmed); match != nil {
			state := sd.ensureState(match[2], "")
			sd.Notes = append(sd.Notes, &Note{State: state, Position: match[1], Text: strings.TrimSpace(match[3])})
			sd.addToCurrentComposite(compositeStack, state.ID)
			continue
		}
		if match := transitionRegex.FindStringSubmatch(trimmed); match != nil {
			from := sd.ensureState(match[1], "")
			to := sd.ensureState(match[2], "")
			sd.Transitions = append(sd.Transitions, &Transition{From: from, To: to, Label: strings.TrimSpace(match[3])})
			sd.addToCurrentComposite(compositeStack, from.ID)
			sd.addToCurrentComposite(compositeStack, to.ID)
			continue
		}
		if strings.Contains(trimmed, "-->") {
			return nil, fmt.Errorf("line %d: invalid state transition syntax: %q", lineNo, trimmed)
		}
		if match := stateDeclarationRegex.FindStringSubmatch(trimmed); match != nil {
			state := sd.ensureState(match[1], "")
			sd.addToCurrentComposite(compositeStack, state.ID)
			continue
		}
		if match := descriptionRegex.FindStringSubmatch(trimmed); match != nil {
			state := sd.ensureState(match[1], "")
			state.Description = strings.TrimSpace(match[2])
			sd.addToCurrentComposite(compositeStack, state.ID)
			continue
		}

		return nil, fmt.Errorf("line %d: invalid state syntax: %q", lineNo, trimmed)
	}

	if len(compositeStack) != 0 {
		return nil, fmt.Errorf("unterminated composite state %q", compositeStack[len(compositeStack)-1].ID)
	}
	if len(sd.States) == 0 {
		return nil, fmt.Errorf("no states found")
	}
	return sd, nil
}

func (sd *Diagram) ensureState(id, label string) *State {
	if state, ok := sd.States[id]; ok {
		if label != "" {
			state.Label = label
		}
		return state
	}
	if label == "" {
		label = id
	}
	state := &State{ID: id, Label: label}
	sd.States[id] = state
	sd.StateOrder = append(sd.StateOrder, id)
	return state
}

func (sd *Diagram) addToCurrentComposite(stack []*Composite, stateID string) {
	if len(stack) == 0 || stateID == StartEndID {
		return
	}
	current := stack[len(stack)-1]
	if current.ID == stateID {
		return
	}
	for _, child := range current.Children {
		if child == stateID {
			return
		}
	}
	current.Children = append(current.Children, stateID)
}
