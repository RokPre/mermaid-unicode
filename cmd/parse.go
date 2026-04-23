package cmd

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
	log "github.com/sirupsen/logrus"
)

type graphProperties struct {
	data             *orderedmap.OrderedMap[string, []textEdge]
	nodeSpecs        map[string]graphNodeSpec
	edgeStyles       map[int]styleClass
	styleClasses     *map[string]styleClass
	edgeIndex        int
	boxBorderPadding int
	graphDirection   string
	graphBoxStyle    string
	graphEdgeStyle   graphEdgeLineStyle
	styleType        string
	paddingX         int
	paddingY         int
	subgraphs        []*textSubgraph
	useAscii         bool
}

type textNode struct {
	name       string
	label      graphLabel
	hasLabel   bool
	shape      graphNodeShape
	hasShape   bool
	styleClass string
}

type graphNodeShape string

const (
	graphNodeShapeSquare        graphNodeShape = "square"
	graphNodeShapeRounded       graphNodeShape = "rounded"
	graphNodeShapeStadium       graphNodeShape = "stadium"
	graphNodeShapeDouble        graphNodeShape = "double"
	graphNodeShapeDatabase      graphNodeShape = "database"
	graphNodeShapeCircle        graphNodeShape = "circle"
	graphNodeShapeDecision      graphNodeShape = "decision"
	graphNodeShapeHexagon       graphNodeShape = "hexagon"
	graphNodeShapeParallelogram graphNodeShape = "parallelogram"
)

type graphNodeSpec struct {
	label           graphLabel
	labelIsExplicit bool
	shape           graphNodeShape
	shapeIsExplicit bool
	styleClass      string
}

type graphEdgeLineStyle string

const (
	graphEdgeLineStyleLight  graphEdgeLineStyle = "light"
	graphEdgeLineStyleHeavy  graphEdgeLineStyle = "heavy"
	graphEdgeLineStyleDashed graphEdgeLineStyle = "dashed"
)

var graphExpandedNodePropertyRegex = regexp.MustCompile(`(?i)([a-z][a-z0-9_-]*)\s*:\s*(?:"([^"]*)"|'([^']*)'|([^,}]+))`)

type textEdge struct {
	parent       textNode
	child        textNode
	label        string
	index        int
	lineStyle    graphEdgeLineStyle
	lineStyleSet bool
	hasArrowHead bool
}

type textSubgraph struct {
	id       string
	name     string
	label    graphLabel
	nodes    []string
	parent   *textSubgraph
	children []*textSubgraph
}

func parseSubgraphHeader(header string) textSubgraph {
	trimmed := strings.TrimSpace(header)
	labelText := trimmed
	id := ""

	if match := regexp.MustCompile(`^(\S+)\s*\[(.+)\]$`).FindStringSubmatch(trimmed); match != nil {
		id = strings.TrimSpace(match[1])
		labelText = strings.TrimSpace(match[2])
		labelText = strings.Trim(labelText, `"`)
	}

	return textSubgraph{
		id:    id,
		name:  labelText,
		label: newGraphLabel(labelText),
		nodes: []string{},
	}
}

func splitGraphLines(mermaid string) []string {
	lines := []string{}
	var current strings.Builder
	bracketDepth := 0
	inQuotes := false

	for i := 0; i < len(mermaid); i++ {
		switch mermaid[i] {
		case '"':
			inQuotes = !inQuotes
		case '[':
			if !inQuotes {
				bracketDepth++
			}
		case ']':
			if !inQuotes && bracketDepth > 0 {
				bracketDepth--
			}
		case '\n':
			if bracketDepth == 0 {
				lines = append(lines, current.String())
				current.Reset()
				continue
			}
		case '\\':
			if i+1 < len(mermaid) && mermaid[i+1] == 'n' && bracketDepth == 0 {
				lines = append(lines, current.String())
				current.Reset()
				i++
				continue
			}
		}

		current.WriteByte(mermaid[i])
	}

	return append(lines, current.String())
}

func parseNode(line string) textNode {
	// Trim any whitespace from the line that might be left after comment removal
	trimmedLine := strings.TrimSpace(line)
	styleClass := ""
	if idx := strings.LastIndex(trimmedLine, ":::"); idx != -1 {
		styleClass = strings.TrimSpace(trimmedLine[idx+3:])
		trimmedLine = strings.TrimSpace(trimmedLine[:idx])
	}

	name := trimmedLine
	labelText := trimmedLine
	if node, ok := parseExpandedShapeNode(trimmedLine, styleClass); ok {
		return node
	}
	for _, shape := range nodeShapeSyntaxes() {
		if open := strings.Index(trimmedLine, shape.open); open > 0 && strings.HasSuffix(trimmedLine, shape.close) {
			name = strings.TrimSpace(trimmedLine[:open])
			labelText = strings.TrimSpace(trimmedLine[open+len(shape.open) : len(trimmedLine)-len(shape.close)])
			labelText = strings.Trim(labelText, `"`)
			return textNode{
				name:       name,
				label:      newGraphLabel(labelText),
				hasLabel:   true,
				shape:      shape.shape,
				hasShape:   true,
				styleClass: styleClass,
			}
		}
	}

	return textNode{name: name, label: newGraphLabel(labelText), styleClass: styleClass}
}

func parseExpandedShapeNode(line, styleClass string) (textNode, bool) {
	open := strings.Index(line, "@{")
	if open <= 0 || !strings.HasSuffix(line, "}") {
		return textNode{}, false
	}

	name := strings.TrimSpace(line[:open])
	if name == "" {
		return textNode{}, false
	}

	metadata := strings.TrimSpace(line[open+len("@{") : len(line)-len("}")])
	properties := parseExpandedNodeProperties(metadata)
	labelText, hasLabel := properties["label"]
	if !hasLabel {
		labelText = name
	}

	node := textNode{
		name:       name,
		label:      newGraphLabel(labelText),
		hasLabel:   hasLabel,
		styleClass: styleClass,
	}

	if shapeName, ok := properties["shape"]; ok {
		if shape, ok := expandedNodeShape(shapeName); ok {
			node.shape = shape
			node.hasShape = true
		}
	}

	return node, true
}

func parseExpandedNodeProperties(metadata string) map[string]string {
	properties := map[string]string{}
	for _, match := range graphExpandedNodePropertyRegex.FindAllStringSubmatch(metadata, -1) {
		key := strings.ToLower(strings.TrimSpace(match[1]))
		value := ""
		for _, candidate := range match[2:] {
			if candidate != "" {
				value = strings.TrimSpace(candidate)
				break
			}
		}
		if key != "" && value != "" {
			properties[key] = value
		}
	}
	return properties
}

func expandedNodeShape(shapeName string) (graphNodeShape, bool) {
	normalized := strings.ToLower(strings.TrimSpace(shapeName))
	normalized = strings.Trim(normalized, `"'`)

	switch normalized {
	case "rect", "rectangle", "proc", "process":
		return graphNodeShapeSquare, true
	case "rounded", "event", "delay", "half-rounded-rectangle":
		return graphNodeShapeRounded, true
	case "stadium", "terminal", "pill":
		return graphNodeShapeStadium, true
	case "subroutine", "subprocess", "subproc", "fr-rect", "framed-rectangle", "framed-rect":
		return graphNodeShapeDouble, true
	case "database", "db", "cyl", "cylinder":
		return graphNodeShapeDatabase, true
	case "circle", "circ", "small-circle", "sm-circ", "start", "double-circle", "dbl-circ", "framed-circle", "fr-circ", "stop":
		return graphNodeShapeCircle, true
	case "decision", "diamond", "diam", "question":
		return graphNodeShapeDecision, true
	case "hexagon", "hex", "prepare":
		return graphNodeShapeHexagon, true
	case "lean-r", "lean-right", "in-out", "lean-l", "lean-left", "out-in", "parallelogram", "parallelogram-alt":
		return graphNodeShapeParallelogram, true
	default:
		return "", false
	}
}

type nodeShapeSyntax struct {
	open  string
	close string
	shape graphNodeShape
}

func nodeShapeSyntaxes() []nodeShapeSyntax {
	return []nodeShapeSyntax{
		{open: "([", close: "])", shape: graphNodeShapeStadium},
		{open: "[[", close: "]]", shape: graphNodeShapeDouble},
		{open: "[(", close: ")]", shape: graphNodeShapeDatabase},
		{open: "((", close: "))", shape: graphNodeShapeCircle},
		{open: "{{", close: "}}", shape: graphNodeShapeHexagon},
		{open: "[/", close: "/]", shape: graphNodeShapeParallelogram},
		{open: "[", close: "]", shape: graphNodeShapeSquare},
		{open: "(", close: ")", shape: graphNodeShapeRounded},
		{open: "{", close: "}", shape: graphNodeShapeDecision},
	}
}

func parseStyleClass(matchedLine []string) styleClass {
	className := matchedLine[0]
	return styleClass{className, parseStyleMap(matchedLine[1])}
}

func parseStyleMap(styles string) map[string]string {
	styleMap := make(map[string]string)
	for _, style := range strings.Split(styles, ",") {
		style = strings.TrimSpace(strings.TrimSuffix(style, ";"))
		if style == "" {
			continue
		}
		kv := strings.SplitN(style, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(strings.TrimSuffix(kv[1], ";"))
		if key != "" && value != "" {
			styleMap[key] = value
		}
	}
	return styleMap
}

func setEdgeWithLabel(lhs, rhs []textNode, label string, lineStyle graphEdgeLineStyle, hasArrowHead bool, gp *graphProperties) []textNode {
	log.Debug("Setting arrow from ", lhs, " to ", rhs, " with label ", label)
	for _, l := range lhs {
		for _, r := range rhs {
			setData(l, textEdge{
				parent:       l,
				child:        r,
				label:        label,
				index:        gp.edgeIndex,
				lineStyle:    lineStyle,
				lineStyleSet: lineStyle != graphEdgeLineStyleLight,
				hasArrowHead: hasArrowHead,
			}, gp.data, gp.nodeSpecs)
			gp.edgeIndex++
		}
	}
	return rhs
}

func setArrowWithLabel(lhs, rhs []textNode, label string, gp *graphProperties) []textNode {
	return setEdgeWithLabel(lhs, rhs, label, graphEdgeLineStyleLight, true, gp)
}

func setArrow(lhs, rhs []textNode, gp *graphProperties) []textNode {
	return setArrowWithLabel(lhs, rhs, "", gp)
}

func rememberNode(node textNode, nodeSpecs map[string]graphNodeSpec) {
	spec := nodeSpecs[node.name]
	if node.hasLabel || len(spec.label.lines) == 0 {
		spec.label = node.label
		spec.labelIsExplicit = node.hasLabel
	}
	if node.hasShape || spec.shape == "" {
		if node.hasShape {
			spec.shape = node.shape
		} else {
			spec.shape = graphNodeShapeSquare
		}
		spec.shapeIsExplicit = node.hasShape
	}
	if node.styleClass != "" {
		spec.styleClass = node.styleClass
	}
	nodeSpecs[node.name] = spec
}

func addNode(node textNode, data *orderedmap.OrderedMap[string, []textEdge], nodeSpecs map[string]graphNodeSpec) {
	rememberNode(node, nodeSpecs)
	if _, ok := data.Get(node.name); !ok {
		data.Set(node.name, []textEdge{})
	}
}

func setData(parent textNode, edge textEdge, data *orderedmap.OrderedMap[string, []textEdge], nodeSpecs map[string]graphNodeSpec) {
	rememberNode(parent, nodeSpecs)
	rememberNode(edge.child, nodeSpecs)
	// Check if the parent is in the map
	if children, ok := data.Get(parent.name); ok {
		// If it is, append the child to the list of children
		data.Set(parent.name, append(children, edge))
	} else {
		// If it isn't, add it to the map
		data.Set(parent.name, []textEdge{edge})
	}
	// Check if the child is in the map
	if _, ok := data.Get(edge.child.name); ok {
		// If it is, do nothing
	} else {
		// If it isn't, add it to the map
		data.Set(edge.child.name, []textEdge{})
	}
}

func (gp *graphProperties) parseString(line string) ([]textNode, error) {
	log.Debugf("Parsing line: %v", line)
	var lhs, rhs []textNode
	var err error
	// Patterns are matched in order
	patterns := []struct {
		regex   *regexp.Regexp
		handler func([]string) ([]textNode, error)
	}{
		{
			regex: regexp.MustCompile(`^\s*$`),
			handler: func(match []string) ([]textNode, error) {
				// Ignore empty lines
				return []textNode{}, nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+)\s+==>\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				if lhs, err = gp.parseString(match[0]); err != nil {
					lhs = []textNode{parseNode(match[0])}
				}
				if rhs, err = gp.parseString(match[1]); err != nil {
					rhs = []textNode{parseNode(match[1])}
				}
				return setEdgeWithLabel(lhs, rhs, "", graphEdgeLineStyleHeavy, true, gp), nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+)\s+==>\|(.+)\|\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				if lhs, err = gp.parseString(match[0]); err != nil {
					lhs = []textNode{parseNode(match[0])}
				}
				if rhs, err = gp.parseString(match[2]); err != nil {
					rhs = []textNode{parseNode(match[2])}
				}
				return setEdgeWithLabel(lhs, rhs, match[1], graphEdgeLineStyleHeavy, true, gp), nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+)\s+-\.->\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				if lhs, err = gp.parseString(match[0]); err != nil {
					lhs = []textNode{parseNode(match[0])}
				}
				if rhs, err = gp.parseString(match[1]); err != nil {
					rhs = []textNode{parseNode(match[1])}
				}
				return setEdgeWithLabel(lhs, rhs, "", graphEdgeLineStyleDashed, true, gp), nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+)\s+-\.->\|(.+)\|\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				if lhs, err = gp.parseString(match[0]); err != nil {
					lhs = []textNode{parseNode(match[0])}
				}
				if rhs, err = gp.parseString(match[2]); err != nil {
					rhs = []textNode{parseNode(match[2])}
				}
				return setEdgeWithLabel(lhs, rhs, match[1], graphEdgeLineStyleDashed, true, gp), nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+)\s+-\.-\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				if lhs, err = gp.parseString(match[0]); err != nil {
					lhs = []textNode{parseNode(match[0])}
				}
				if rhs, err = gp.parseString(match[1]); err != nil {
					rhs = []textNode{parseNode(match[1])}
				}
				return setEdgeWithLabel(lhs, rhs, "", graphEdgeLineStyleDashed, false, gp), nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+)\s+-->\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				if lhs, err = gp.parseString(match[0]); err != nil {
					lhs = []textNode{parseNode(match[0])}
				}
				if rhs, err = gp.parseString(match[1]); err != nil {
					rhs = []textNode{parseNode(match[1])}
				}
				return setArrow(lhs, rhs, gp), nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+)\s+-->\|(.+)\|\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				if lhs, err = gp.parseString(match[0]); err != nil {
					lhs = []textNode{parseNode(match[0])}
				}
				if rhs, err = gp.parseString(match[2]); err != nil {
					rhs = []textNode{parseNode(match[2])}
				}
				return setArrowWithLabel(lhs, rhs, match[1], gp), nil
			},
		},
		{
			regex: regexp.MustCompile(`^classDef\s+(\S+)\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				s := parseStyleClass(match)
				(*gp.styleClasses)[s.name] = s
				return []textNode{}, nil
			},
		},
		{
			regex: regexp.MustCompile(`^class\s+(.+)\s+(\S+)$`),
			handler: func(match []string) ([]textNode, error) {
				for _, nodeName := range strings.Split(match[0], ",") {
					nodeName = strings.TrimSpace(nodeName)
					if nodeName == "" {
						continue
					}
					spec := gp.nodeSpecs[nodeName]
					spec.styleClass = match[1]
					gp.nodeSpecs[nodeName] = spec
				}
				return []textNode{}, nil
			},
		},
		{
			regex: regexp.MustCompile(`^linkStyle\s+(\d+)\s+(.+)$`),
			handler: func(match []string) ([]textNode, error) {
				index, err := strconv.Atoi(match[0])
				if err != nil {
					return []textNode{}, err
				}
				gp.edgeStyles[index] = styleClass{
					name:   fmt.Sprintf("linkStyle-%d", index),
					styles: parseStyleMap(match[1]),
				}
				return []textNode{}, nil
			},
		},
		{
			regex: regexp.MustCompile(`(?s)^(.+) & (.+)$`),
			handler: func(match []string) ([]textNode, error) {
				log.Debugf("Found & pattern node %v to %v", match[0], match[1])
				var node textNode
				if lhs, err = gp.parseString(match[0]); err != nil {
					node = parseNode(match[0])
					lhs = []textNode{node}
				}
				if rhs, err = gp.parseString(match[1]); err != nil {
					node = parseNode(match[1])
					rhs = []textNode{node}
				}
				return append(lhs, rhs...), nil
			},
		},
	}
	for _, pattern := range patterns {
		if match := pattern.regex.FindStringSubmatch(line); match != nil {
			nodes, err := pattern.handler(match[1:])
			if err == nil {
				return nodes, nil
			}
		}
	}
	return []textNode{}, errors.New("Could not parse line: " + line)
}

func mermaidFileToMap(mermaid, styleType string) (*graphProperties, error) {
	rawLines := splitGraphLines(mermaid)

	// Process lines to remove comments
	lines := []string{}
	for _, line := range rawLines {
		// Stop processing at "---" separator (used in test files)
		if line == "---" {
			break
		}

		// Skip lines that start with %% (comment lines)
		if strings.HasPrefix(strings.TrimSpace(line), "%%") {
			continue
		}

		// Remove inline comments (anything after %%) and trim resulting whitespace
		if idx := strings.Index(line, "%%"); idx != -1 {
			line = strings.TrimSpace(line[:idx])
		}

		// Skip empty lines after comment removal
		if len(strings.TrimSpace(line)) > 0 {
			lines = append(lines, line)
		}
	}

	data := orderedmap.NewOrderedMap[string, []textEdge]()
	styleClasses := make(map[string]styleClass)
	properties := graphProperties{
		data:             data,
		nodeSpecs:        make(map[string]graphNodeSpec),
		edgeStyles:       make(map[int]styleClass),
		styleClasses:     &styleClasses,
		boxBorderPadding: boxBorderPadding,
		graphDirection:   "",
		styleType:        styleType,
		paddingX:         paddingBetweenX,
		paddingY:         paddingBetweenY,
		subgraphs:        []*textSubgraph{},
	}

	// Pick up optional padding directives before the graph definition
	paddingRegex := regexp.MustCompile(`^(?i)padding([xy])\s*=\s*(\d+)$`)
	for len(lines) > 0 {
		trimmed := strings.TrimSpace(lines[0])
		if trimmed == "" {
			lines = lines[1:]
			continue
		}
		if match := paddingRegex.FindStringSubmatch(trimmed); match != nil {
			paddingValue, err := strconv.Atoi(match[2])
			if err != nil {
				return &properties, err
			}
			if strings.EqualFold(match[1], "x") {
				properties.paddingX = paddingValue
			} else {
				properties.paddingY = paddingValue
			}
			lines = lines[1:]
			continue
		}
		break
	}

	if len(lines) == 0 {
		return &properties, errors.New("missing graph definition")
	}

	// First line should either say "graph TD" or "graph LR"
	switch lines[0] {
	case "graph LR", "flowchart LR":
		properties.graphDirection = "LR"
	case "graph TD", "flowchart TD", "graph TB", "flowchart TB":
		properties.graphDirection = "TD"
	default:
		return &properties, fmt.Errorf("unsupported graph type '%s'. Supported types: graph TD, graph TB, graph LR, flowchart TD, flowchart TB, flowchart LR", lines[0])
	}
	lines = lines[1:]

	// Track subgraph context using a stack
	subgraphStack := []*textSubgraph{}
	subgraphRegex := regexp.MustCompile(`^\s*subgraph\s+(.+)$`)
	endRegex := regexp.MustCompile(`^\s*end\s*$`)

	// Iterate over the lines
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		// Check for subgraph start
		if match := subgraphRegex.FindStringSubmatch(trimmedLine); match != nil {
			header := parseSubgraphHeader(match[1])
			newSubgraph := &textSubgraph{
				id:       header.id,
				name:     header.name,
				label:    header.label,
				nodes:    []string{},
				children: []*textSubgraph{},
			}

			// Set parent relationship if we're nested
			if len(subgraphStack) > 0 {
				parent := subgraphStack[len(subgraphStack)-1]
				newSubgraph.parent = parent
				parent.children = append(parent.children, newSubgraph)
			}

			subgraphStack = append(subgraphStack, newSubgraph)
			properties.subgraphs = append(properties.subgraphs, newSubgraph)
			log.Debugf("Started subgraph %s", newSubgraph.name)
			continue
		}

		// Check for subgraph end
		if endRegex.MatchString(trimmedLine) {
			if len(subgraphStack) > 0 {
				closedSubgraph := subgraphStack[len(subgraphStack)-1]
				subgraphStack = subgraphStack[:len(subgraphStack)-1]
				log.Debugf("Ended subgraph %s", closedSubgraph.name)
			}
			continue
		}

		// Remember nodes before parsing this line
		existingNodes := make(map[string]bool)
		for el := data.Front(); el != nil; el = el.Next() {
			existingNodes[el.Key] = true
		}

		// Parse nodes and edges normally
		nodes, err := properties.parseString(line)
		if err != nil {
			log.Debugf("Parsing remaining text to node %v", line)
			node := parseNode(line)
			addNode(node, properties.data, properties.nodeSpecs)
		} else {
			// Ensure all returned nodes are in the map
			for _, node := range nodes {
				addNode(node, properties.data, properties.nodeSpecs)
			}
		}

		// Add all new nodes to current subgraph(s)
		if len(subgraphStack) > 0 {
			for el := data.Front(); el != nil; el = el.Next() {
				nodeName := el.Key
				// If this is a new node (wasn't in existingNodes), add it to subgraph
				if !existingNodes[nodeName] {
					for _, sg := range subgraphStack {
						// Check if node is not already in the subgraph
						found := false
						for _, n := range sg.nodes {
							if n == nodeName {
								found = true
								break
							}
						}
						if !found {
							sg.nodes = append(sg.nodes, nodeName)
							log.Debugf("Added node %s to subgraph %s", nodeName, sg.name)
						}
					}
				}
			}
		}
	}
	return &properties, nil
}
