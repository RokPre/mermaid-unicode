package journey

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/mattn/go-runewidth"
)

func Render(j *Diagram, config *diagram.Config) (string, error) {
	if j == nil || len(j.Tasks) == 0 {
		return "", fmt.Errorf("no journey tasks")
	}
	if config == nil {
		config = diagram.DefaultConfig()
	}
	barFull, barEmpty := "█", "░"
	if config.UseAscii {
		barFull, barEmpty = "#", "."
	}
	lines := []string{}
	if j.Title != "" {
		lines = append(lines, j.Title)
	}
	current := ""
	for _, task := range j.Tasks {
		if task.Section != current {
			current = task.Section
			if current != "" {
				lines = append(lines, "["+current+"]")
			}
		}
		bar := strings.Repeat(barFull, task.Score) + strings.Repeat(barEmpty, 5-task.Score)
		actors := strings.Join(task.Actors, ", ")
		lines = append(lines, fmt.Sprintf("%s  %d/5  %-5s  %s", padRight(task.Name, 24), task.Score, bar, actors))
	}
	return strings.Join(lines, "\n") + "\n", nil
}

func padRight(text string, width int) string {
	textWidth := runewidth.StringWidth(text)
	if textWidth >= width {
		return text
	}
	return text + strings.Repeat(" ", width-textWidth)
}
