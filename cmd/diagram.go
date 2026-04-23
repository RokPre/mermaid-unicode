package cmd

import (
	"fmt"
	"strings"

	"github.com/AlexanderGrooff/mermaid-ascii/pkg/charts"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/classdiagram"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/diagram"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/er"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/gantt"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/journey"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/mindmap"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/requirementdiagram"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/sequence"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/statediagram"
	"github.com/AlexanderGrooff/mermaid-ascii/pkg/timeline"
)

type diagramRegistration struct {
	typeName string
	detect   func(input, firstLine string) bool
	create   func() diagram.Diagram
}

type unsupportedDiagramRegistration struct {
	typeName string
	detect   func(firstLine string) bool
}

var supportedDiagramRegistry = []diagramRegistration{
	{
		typeName: "sequence",
		detect: func(input, _ string) bool {
			return sequence.IsSequenceDiagram(input)
		},
		create: func() diagram.Diagram {
			return &SequenceDiagram{}
		},
	},
	{
		typeName: "graph",
		detect: func(_ string, firstLine string) bool {
			return hasMermaidKeyword(firstLine, "graph") || hasMermaidKeyword(firstLine, "flowchart")
		},
		create: func() diagram.Diagram {
			return &GraphDiagram{}
		},
	},
	{
		typeName: "class",
		detect: func(input, _ string) bool {
			return classdiagram.IsClassDiagram(input)
		},
		create: func() diagram.Diagram {
			return &ClassDiagram{}
		},
	},
	{
		typeName: "state",
		detect: func(input, _ string) bool {
			return statediagram.IsStateDiagram(input)
		},
		create: func() diagram.Diagram {
			return &StateDiagram{}
		},
	},
	{
		typeName: "requirement",
		detect: func(input, _ string) bool {
			return requirementdiagram.IsRequirementDiagram(input)
		},
		create: func() diagram.Diagram {
			return &RequirementDiagram{}
		},
	},
	{
		typeName: "mindmap",
		detect: func(input, _ string) bool {
			return mindmap.IsMindmap(input)
		},
		create: func() diagram.Diagram {
			return &MindmapDiagram{}
		},
	},
	{
		typeName: "timeline",
		detect: func(input, _ string) bool {
			return timeline.IsTimeline(input)
		},
		create: func() diagram.Diagram {
			return &TimelineDiagram{}
		},
	},
	{
		typeName: "journey",
		detect: func(input, _ string) bool {
			return journey.IsJourney(input)
		},
		create: func() diagram.Diagram {
			return &JourneyDiagram{}
		},
	},
	{
		typeName: "pie",
		detect: func(input, _ string) bool {
			return charts.IsPie(input)
		},
		create: func() diagram.Diagram {
			return &PieDiagram{}
		},
	},
	{
		typeName: "gantt",
		detect: func(input, _ string) bool {
			return gantt.IsGantt(input)
		},
		create: func() diagram.Diagram {
			return &GanttDiagram{}
		},
	},
	{
		typeName: "quadrant",
		detect: func(input, _ string) bool {
			return charts.IsQuadrant(input)
		},
		create: func() diagram.Diagram {
			return &QuadrantDiagram{}
		},
	},
	{
		typeName: "er",
		detect: func(input, _ string) bool {
			return er.IsERDiagram(input)
		},
		create: func() diagram.Diagram {
			return &ERDiagram{}
		},
	},
}

var unsupportedDiagramRegistry = []unsupportedDiagramRegistration{
	{typeName: "gitGraph", detect: keywordDetector("gitGraph")},
	{typeName: "zenuml", detect: keywordDetector("zenuml")},
}

func DiagramFactory(input string) (diagram.Diagram, error) {
	input = strings.TrimSpace(input)

	firstLine, ok := firstDiagramLine(input)
	if !ok {
		return nil, fmt.Errorf("missing diagram definition. Supported diagram types: %s", supportedDiagramTypes())
	}

	for _, registered := range supportedDiagramRegistry {
		if registered.detect(input, firstLine) {
			return registered.create(), nil
		}
	}

	for _, registered := range unsupportedDiagramRegistry {
		if registered.detect(firstLine) {
			return nil, fmt.Errorf("unsupported diagram type %q. Supported diagram types: %s", registered.typeName, supportedDiagramTypes())
		}
	}

	return nil, fmt.Errorf("unknown diagram type in first content line %q. Supported diagram types: %s", firstLine, supportedDiagramTypes())
}

func firstDiagramLine(input string) (string, bool) {
	for _, line := range strings.Split(input, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "%%") || isGraphPaddingDirective(trimmed) {
			continue
		}
		return trimmed, true
	}
	return "", false
}

func isGraphPaddingDirective(line string) bool {
	lower := strings.ToLower(strings.TrimSpace(line))
	return strings.HasPrefix(lower, "paddingx") || strings.HasPrefix(lower, "paddingy")
}

func hasMermaidKeyword(line, keyword string) bool {
	if line == keyword {
		return true
	}
	if !strings.HasPrefix(line, keyword) {
		return false
	}
	if len(line) == len(keyword) {
		return true
	}
	next := line[len(keyword)]
	return next == ' ' || next == '\t' || next == ':'
}

func keywordDetector(keyword string) func(string) bool {
	return func(firstLine string) bool {
		return hasMermaidKeyword(firstLine, keyword)
	}
}

func anyKeywordDetector(keywords ...string) func(string) bool {
	return func(firstLine string) bool {
		for _, keyword := range keywords {
			if hasMermaidKeyword(firstLine, keyword) {
				return true
			}
		}
		return false
	}
}

func supportedDiagramTypes() string {
	types := make([]string, 0, len(supportedDiagramRegistry))
	for _, registered := range supportedDiagramRegistry {
		types = append(types, registered.typeName)
	}
	return strings.Join(types, ", ")
}

type SequenceDiagram struct {
	parsed *sequence.SequenceDiagram
}

func (sd *SequenceDiagram) Parse(input string) error {
	parsed, err := sequence.Parse(input)
	if err != nil {
		return err
	}
	sd.parsed = parsed
	return nil
}

func (sd *SequenceDiagram) Render(config *diagram.Config) (string, error) {
	if sd.parsed == nil {
		return "", fmt.Errorf("sequence diagram not parsed: call Parse() before Render()")
	}
	return sequence.Render(sd.parsed, config)
}

func (sd *SequenceDiagram) Type() string {
	return "sequence"
}

type ERDiagram struct {
	parsed *er.Diagram
}

func (ed *ERDiagram) Parse(input string) error {
	parsed, err := er.Parse(input)
	if err != nil {
		return err
	}
	ed.parsed = parsed
	return nil
}

func (ed *ERDiagram) Render(config *diagram.Config) (string, error) {
	if ed.parsed == nil {
		return "", fmt.Errorf("ER diagram not parsed: call Parse() before Render()")
	}
	return er.Render(ed.parsed, config)
}

func (ed *ERDiagram) Type() string {
	return "er"
}

type ClassDiagram struct {
	parsed *classdiagram.Diagram
}

func (cd *ClassDiagram) Parse(input string) error {
	parsed, err := classdiagram.Parse(input)
	if err != nil {
		return err
	}
	cd.parsed = parsed
	return nil
}

func (cd *ClassDiagram) Render(config *diagram.Config) (string, error) {
	if cd.parsed == nil {
		return "", fmt.Errorf("class diagram not parsed: call Parse() before Render()")
	}
	return classdiagram.Render(cd.parsed, config)
}

func (cd *ClassDiagram) Type() string {
	return "class"
}

type StateDiagram struct {
	parsed *statediagram.Diagram
}

func (sd *StateDiagram) Parse(input string) error {
	parsed, err := statediagram.Parse(input)
	if err != nil {
		return err
	}
	sd.parsed = parsed
	return nil
}

func (sd *StateDiagram) Render(config *diagram.Config) (string, error) {
	if sd.parsed == nil {
		return "", fmt.Errorf("state diagram not parsed: call Parse() before Render()")
	}
	return statediagram.Render(sd.parsed, config)
}

func (sd *StateDiagram) Type() string {
	return "state"
}

type RequirementDiagram struct {
	parsed *requirementdiagram.Diagram
}

func (rd *RequirementDiagram) Parse(input string) error {
	parsed, err := requirementdiagram.Parse(input)
	if err != nil {
		return err
	}
	rd.parsed = parsed
	return nil
}

func (rd *RequirementDiagram) Render(config *diagram.Config) (string, error) {
	if rd.parsed == nil {
		return "", fmt.Errorf("requirement diagram not parsed: call Parse() before Render()")
	}
	return requirementdiagram.Render(rd.parsed, config)
}

func (rd *RequirementDiagram) Type() string {
	return "requirement"
}

type MindmapDiagram struct {
	parsed *mindmap.Diagram
}

func (md *MindmapDiagram) Parse(input string) error {
	parsed, err := mindmap.Parse(input)
	if err != nil {
		return err
	}
	md.parsed = parsed
	return nil
}

func (md *MindmapDiagram) Render(config *diagram.Config) (string, error) {
	if md.parsed == nil {
		return "", fmt.Errorf("mindmap not parsed: call Parse() before Render()")
	}
	return mindmap.Render(md.parsed, config)
}

func (md *MindmapDiagram) Type() string {
	return "mindmap"
}

type TimelineDiagram struct {
	parsed *timeline.Diagram
}

func (td *TimelineDiagram) Parse(input string) error {
	parsed, err := timeline.Parse(input)
	if err != nil {
		return err
	}
	td.parsed = parsed
	return nil
}

func (td *TimelineDiagram) Render(config *diagram.Config) (string, error) {
	if td.parsed == nil {
		return "", fmt.Errorf("timeline not parsed: call Parse() before Render()")
	}
	return timeline.Render(td.parsed, config)
}

func (td *TimelineDiagram) Type() string {
	return "timeline"
}

type JourneyDiagram struct {
	parsed *journey.Diagram
}

func (jd *JourneyDiagram) Parse(input string) error {
	parsed, err := journey.Parse(input)
	if err != nil {
		return err
	}
	jd.parsed = parsed
	return nil
}

func (jd *JourneyDiagram) Render(config *diagram.Config) (string, error) {
	if jd.parsed == nil {
		return "", fmt.Errorf("journey not parsed: call Parse() before Render()")
	}
	return journey.Render(jd.parsed, config)
}

func (jd *JourneyDiagram) Type() string {
	return "journey"
}

type PieDiagram struct {
	parsed *charts.PieDiagram
}

func (pd *PieDiagram) Parse(input string) error {
	parsed, err := charts.ParsePie(input)
	if err != nil {
		return err
	}
	pd.parsed = parsed
	return nil
}

func (pd *PieDiagram) Render(config *diagram.Config) (string, error) {
	if pd.parsed == nil {
		return "", fmt.Errorf("pie chart not parsed: call Parse() before Render()")
	}
	return charts.RenderPie(pd.parsed, config)
}

func (pd *PieDiagram) Type() string {
	return "pie"
}

type QuadrantDiagram struct {
	parsed *charts.QuadrantDiagram
}

func (qd *QuadrantDiagram) Parse(input string) error {
	parsed, err := charts.ParseQuadrant(input)
	if err != nil {
		return err
	}
	qd.parsed = parsed
	return nil
}

func (qd *QuadrantDiagram) Render(config *diagram.Config) (string, error) {
	if qd.parsed == nil {
		return "", fmt.Errorf("quadrant chart not parsed: call Parse() before Render()")
	}
	return charts.RenderQuadrant(qd.parsed, config)
}

func (qd *QuadrantDiagram) Type() string {
	return "quadrant"
}

type GanttDiagram struct {
	parsed *gantt.Diagram
}

func (gd *GanttDiagram) Parse(input string) error {
	parsed, err := gantt.Parse(input)
	if err != nil {
		return err
	}
	gd.parsed = parsed
	return nil
}

func (gd *GanttDiagram) Render(config *diagram.Config) (string, error) {
	if gd.parsed == nil {
		return "", fmt.Errorf("gantt not parsed: call Parse() before Render()")
	}
	return gantt.Render(gd.parsed, config)
}

func (gd *GanttDiagram) Type() string {
	return "gantt"
}

type GraphDiagram struct {
	properties *graphProperties
}

func (gd *GraphDiagram) Parse(input string) error {
	properties, err := mermaidFileToMap(input, "cli")
	if err != nil {
		return err
	}
	gd.properties = properties
	return nil
}

func (gd *GraphDiagram) Render(config *diagram.Config) (string, error) {
	if gd.properties == nil {
		return "", fmt.Errorf("graph diagram not parsed: call Parse() before Render()")
	}

	if config == nil {
		config = diagram.DefaultConfig()
	}

	styleType := config.StyleType
	if styleType == "" {
		styleType = "cli"
	}
	gd.properties.boxBorderPadding = config.BoxBorderPadding
	gd.properties.paddingX = config.PaddingBetweenX
	gd.properties.paddingY = config.PaddingBetweenY
	gd.properties.styleType = styleType
	gd.properties.useAscii = config.UseAscii
	gd.properties.graphBoxStyle = config.GraphBoxStyle
	if gd.properties.graphBoxStyle == "" {
		gd.properties.graphBoxStyle = "square"
	}
	gd.properties.graphEdgeStyle = graphEdgeLineStyle(config.GraphEdgeStyle)
	if gd.properties.graphEdgeStyle == "" {
		gd.properties.graphEdgeStyle = graphEdgeLineStyleLight
	}

	return drawMap(gd.properties), nil
}

func (gd *GraphDiagram) Type() string {
	return "graph"
}
