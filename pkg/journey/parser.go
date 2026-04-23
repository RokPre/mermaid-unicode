package journey

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
)

const JourneyKeyword = "journey"

type Diagram struct {
	Title string
	Tasks []Task
}

type Task struct {
	Section string
	Name    string
	Score   int
	Actors  []string
}

func IsJourney(input string) bool {
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") {
			continue
		}
		return strings.HasPrefix(trimmed, JourneyKeyword)
	}
	return false
}

func Parse(input string) (*Diagram, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, fmt.Errorf("empty input")
	}
	lines := diagram.RemoveComments(diagram.SplitLines(input))
	if len(lines) == 0 {
		return nil, fmt.Errorf("no content found")
	}
	if !strings.HasPrefix(strings.TrimSpace(lines[0]), JourneyKeyword) {
		return nil, fmt.Errorf("expected %q keyword", JourneyKeyword)
	}
	j := &Diagram{}
	section := ""
	for i, line := range lines[1:] {
		lineNo := i + 2
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "title ") {
			j.Title = strings.TrimSpace(trimmed[len("title "):])
			continue
		}
		if strings.HasPrefix(lower, "section ") {
			section = strings.TrimSpace(trimmed[len("section "):])
			continue
		}
		parts := strings.Split(trimmed, ":")
		if len(parts) < 3 {
			return nil, fmt.Errorf("line %d: invalid journey task syntax: %q", lineNo, trimmed)
		}
		score, err := strconv.Atoi(strings.TrimSpace(parts[1]))
		if err != nil || score < 1 || score > 5 {
			return nil, fmt.Errorf("line %d: journey score must be 1..5", lineNo)
		}
		actors := []string{}
		for _, actor := range strings.Split(strings.Join(parts[2:], ":"), ",") {
			actor = strings.TrimSpace(actor)
			if actor != "" {
				actors = append(actors, actor)
			}
		}
		j.Tasks = append(j.Tasks, Task{Section: section, Name: strings.TrimSpace(parts[0]), Score: score, Actors: actors})
	}
	if len(j.Tasks) == 0 {
		return nil, fmt.Errorf("no journey tasks found")
	}
	return j, nil
}
