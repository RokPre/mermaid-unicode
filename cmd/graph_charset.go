package cmd

import log "github.com/sirupsen/logrus"

type graphBoxChars struct {
	topLeft     string
	topRight    string
	bottomLeft  string
	bottomRight string
	horizontal  string
	vertical    string
}

type graphLineChars struct {
	horizontal        string
	vertical          string
	diagonalDownRight string
	diagonalUpRight   string
	topLeft           string
	topRight          string
	bottomLeft        string
	bottomRight       string
	teeRight          string
	teeLeft           string
	teeDown           string
	teeUp             string
	cross             string
}

type graphArrowHeadChars struct {
	up         string
	down       string
	left       string
	right      string
	upperRight string
	upperLeft  string
	lowerRight string
	lowerLeft  string
	fallback   string
}

type graphCharset struct {
	box         graphBoxChars
	line        graphLineChars
	heavyLine   graphLineChars
	dashedLine  graphLineChars
	arrowHeads  graphArrowHeadChars
	junctionMap map[string]map[string]string
	junctionSet map[string]struct{}
}

var asciiGraphCharset = graphCharset{
	box: graphBoxChars{
		topLeft:     "+",
		topRight:    "+",
		bottomLeft:  "+",
		bottomRight: "+",
		horizontal:  "-",
		vertical:    "|",
	},
	line: graphLineChars{
		horizontal:        "-",
		vertical:          "|",
		diagonalDownRight: "\\",
		diagonalUpRight:   "/",
		topLeft:           "+",
		topRight:          "+",
		bottomLeft:        "+",
		bottomRight:       "+",
		teeRight:          "+",
		teeLeft:           "+",
		teeDown:           "+",
		teeUp:             "+",
		cross:             "+",
	},
	heavyLine: graphLineChars{
		horizontal:        "=",
		vertical:          "|",
		diagonalDownRight: "\\",
		diagonalUpRight:   "/",
		topLeft:           "+",
		topRight:          "+",
		bottomLeft:        "+",
		bottomRight:       "+",
		teeRight:          "+",
		teeLeft:           "+",
		teeDown:           "+",
		teeUp:             "+",
		cross:             "+",
	},
	dashedLine: graphLineChars{
		horizontal:        ".",
		vertical:          ".",
		diagonalDownRight: "\\",
		diagonalUpRight:   "/",
		topLeft:           "+",
		topRight:          "+",
		bottomLeft:        "+",
		bottomRight:       "+",
		teeRight:          "+",
		teeLeft:           "+",
		teeDown:           "+",
		teeUp:             "+",
		cross:             "+",
	},
	arrowHeads: graphArrowHeadChars{
		up:       "^",
		down:     "v",
		left:     "<",
		right:    ">",
		fallback: "*",
	},
}

var lightGraphJunctionMap = map[string]map[string]string{
	"─": {"│": "┼", "┌": "┬", "┐": "┬", "└": "┴", "┘": "┴", "├": "┼", "┤": "┼", "┬": "┬", "┴": "┴"},
	"│": {"─": "┼", "┌": "├", "┐": "┤", "└": "├", "┘": "┤", "├": "├", "┤": "┤", "┬": "┼", "┴": "┼"},
	"┌": {"─": "┬", "│": "├", "┐": "┬", "└": "├", "┘": "┼", "├": "├", "┤": "┼", "┬": "┬", "┴": "┼"},
	"┐": {"─": "┬", "│": "┤", "┌": "┬", "└": "┼", "┘": "┤", "├": "┼", "┤": "┤", "┬": "┬", "┴": "┼"},
	"└": {"─": "┴", "│": "├", "┌": "├", "┐": "┼", "┘": "┴", "├": "├", "┤": "┼", "┬": "┼", "┴": "┴"},
	"┘": {"─": "┴", "│": "┤", "┌": "┼", "┐": "┤", "└": "┴", "├": "┼", "┤": "┤", "┬": "┼", "┴": "┴"},
	"├": {"─": "┼", "│": "├", "┌": "├", "┐": "┼", "└": "├", "┘": "┼", "┤": "┼", "┬": "┼", "┴": "┼"},
	"┤": {"─": "┼", "│": "┤", "┌": "┼", "┐": "┤", "└": "┼", "┘": "┤", "├": "┼", "┬": "┼", "┴": "┼"},
	"┬": {"─": "┬", "│": "┼", "┌": "┬", "┐": "┬", "└": "┼", "┘": "┼", "├": "┼", "┤": "┼", "┴": "┼"},
	"┴": {"─": "┴", "│": "┼", "┌": "┼", "┐": "┼", "└": "┴", "┘": "┴", "├": "┼", "┤": "┼", "┬": "┼"},
}

type lineConnection uint8

const (
	connUp lineConnection = 1 << iota
	connRight
	connDown
	connLeft
)

const (
	graphLinePriorityDashed = iota
	graphLinePriorityLight
	graphLinePriorityHeavy
)

var graphCharConnections = map[string]lineConnection{
	"─": connLeft | connRight,
	"━": connLeft | connRight,
	"┄": connLeft | connRight,
	"═": connLeft | connRight,
	"│": connUp | connDown,
	"┃": connUp | connDown,
	"┆": connUp | connDown,
	"║": connUp | connDown,
	"┌": connRight | connDown,
	"┐": connLeft | connDown,
	"└": connRight | connUp,
	"┘": connLeft | connUp,
	"┏": connRight | connDown,
	"┓": connLeft | connDown,
	"┗": connRight | connUp,
	"┛": connLeft | connUp,
	"╭": connRight | connDown,
	"╮": connLeft | connDown,
	"╰": connRight | connUp,
	"╯": connLeft | connUp,
	"╔": connRight | connDown,
	"╗": connLeft | connDown,
	"╚": connRight | connUp,
	"╝": connLeft | connUp,
	"├": connUp | connRight | connDown,
	"┣": connUp | connRight | connDown,
	"┤": connUp | connDown | connLeft,
	"┫": connUp | connDown | connLeft,
	"┬": connRight | connDown | connLeft,
	"┳": connRight | connDown | connLeft,
	"┴": connUp | connRight | connLeft,
	"┻": connUp | connRight | connLeft,
	"┼": connUp | connRight | connDown | connLeft,
	"╋": connUp | connRight | connDown | connLeft,
}

var graphCharPriorities = map[string]int{
	"┄": graphLinePriorityDashed,
	"┆": graphLinePriorityDashed,
	"─": graphLinePriorityLight,
	"│": graphLinePriorityLight,
	"┌": graphLinePriorityLight,
	"┐": graphLinePriorityLight,
	"└": graphLinePriorityLight,
	"┘": graphLinePriorityLight,
	"╭": graphLinePriorityLight,
	"╮": graphLinePriorityLight,
	"╰": graphLinePriorityLight,
	"╯": graphLinePriorityLight,
	"├": graphLinePriorityLight,
	"┤": graphLinePriorityLight,
	"┬": graphLinePriorityLight,
	"┴": graphLinePriorityLight,
	"┼": graphLinePriorityLight,
	"━": graphLinePriorityHeavy,
	"┃": graphLinePriorityHeavy,
	"┏": graphLinePriorityHeavy,
	"┓": graphLinePriorityHeavy,
	"┗": graphLinePriorityHeavy,
	"┛": graphLinePriorityHeavy,
	"┣": graphLinePriorityHeavy,
	"┫": graphLinePriorityHeavy,
	"┳": graphLinePriorityHeavy,
	"┻": graphLinePriorityHeavy,
	"╋": graphLinePriorityHeavy,
	"═": graphLinePriorityHeavy,
	"║": graphLinePriorityHeavy,
	"╔": graphLinePriorityHeavy,
	"╗": graphLinePriorityHeavy,
	"╚": graphLinePriorityHeavy,
	"╝": graphLinePriorityHeavy,
}

var unicodeLightGraphCharset = graphCharset{
	box: graphBoxChars{
		topLeft:     "┌",
		topRight:    "┐",
		bottomLeft:  "└",
		bottomRight: "┘",
		horizontal:  "─",
		vertical:    "│",
	},
	line: graphLineChars{
		horizontal:        "─",
		vertical:          "│",
		diagonalDownRight: "╲",
		diagonalUpRight:   "╱",
		topLeft:           "┌",
		topRight:          "┐",
		bottomLeft:        "└",
		bottomRight:       "┘",
		teeRight:          "├",
		teeLeft:           "┤",
		teeDown:           "┬",
		teeUp:             "┴",
		cross:             "┼",
	},
	heavyLine: graphLineChars{
		horizontal:        "━",
		vertical:          "┃",
		diagonalDownRight: "╲",
		diagonalUpRight:   "╱",
		topLeft:           "┏",
		topRight:          "┓",
		bottomLeft:        "┗",
		bottomRight:       "┛",
		teeRight:          "┣",
		teeLeft:           "┫",
		teeDown:           "┳",
		teeUp:             "┻",
		cross:             "╋",
	},
	dashedLine: graphLineChars{
		horizontal:        "┄",
		vertical:          "┆",
		diagonalDownRight: "╲",
		diagonalUpRight:   "╱",
		topLeft:           "┌",
		topRight:          "┐",
		bottomLeft:        "└",
		bottomRight:       "┘",
		teeRight:          "├",
		teeLeft:           "┤",
		teeDown:           "┬",
		teeUp:             "┴",
		cross:             "┼",
	},
	arrowHeads: graphArrowHeadChars{
		up:         "▲",
		down:       "▼",
		left:       "◄",
		right:      "►",
		upperRight: "◥",
		upperLeft:  "◤",
		lowerRight: "◢",
		lowerLeft:  "◣",
		fallback:   "●",
	},
	junctionMap: lightGraphJunctionMap,
	junctionSet: map[string]struct{}{
		"─": {},
		"│": {},
		"┌": {},
		"┐": {},
		"└": {},
		"┘": {},
		"├": {},
		"┤": {},
		"┬": {},
		"┴": {},
		"┼": {},
		"╴": {},
		"╵": {},
		"╶": {},
		"╷": {},
	},
}

var unicodeRoundedBoxChars = graphBoxChars{
	topLeft:     "╭",
	topRight:    "╮",
	bottomLeft:  "╰",
	bottomRight: "╯",
	horizontal:  "─",
	vertical:    "│",
}

var unicodeDoubleBoxChars = graphBoxChars{
	topLeft:     "╔",
	topRight:    "╗",
	bottomLeft:  "╚",
	bottomRight: "╝",
	horizontal:  "═",
	vertical:    "║",
}

var unicodeHeavyBoxChars = graphBoxChars{
	topLeft:     "┏",
	topRight:    "┓",
	bottomLeft:  "┗",
	bottomRight: "┛",
	horizontal:  "━",
	vertical:    "┃",
}

var unicodeDecisionBoxChars = graphBoxChars{
	topLeft:     "◇",
	topRight:    "◇",
	bottomLeft:  "◇",
	bottomRight: "◇",
	horizontal:  "─",
	vertical:    "│",
}

var unicodeHexagonBoxChars = graphBoxChars{
	topLeft:     "╱",
	topRight:    "╲",
	bottomLeft:  "╲",
	bottomRight: "╱",
	horizontal:  "─",
	vertical:    "│",
}

var unicodeParallelogramBoxChars = graphBoxChars{
	topLeft:     "╱",
	topRight:    "╱",
	bottomLeft:  "╲",
	bottomRight: "╲",
	horizontal:  "─",
	vertical:    "│",
}

func (g graph) charset() graphCharset {
	if g.useAscii {
		return asciiGraphCharset
	}
	return unicodeLightGraphCharset
}

func (g graph) boxCharsForNode(n *node) graphBoxChars {
	if g.useAscii {
		return asciiGraphCharset.box
	}

	shape := n.shape
	if !n.shapeIsExplicit {
		switch g.graphBoxStyle {
		case "rounded":
			shape = graphNodeShapeRounded
		case "double":
			shape = graphNodeShapeDouble
		case "heavy":
			return unicodeHeavyBoxChars
		}
	}

	switch shape {
	case graphNodeShapeRounded, graphNodeShapeStadium, graphNodeShapeCircle, graphNodeShapeDatabase:
		return unicodeRoundedBoxChars
	case graphNodeShapeDouble:
		return unicodeDoubleBoxChars
	case graphNodeShapeDecision:
		return unicodeDecisionBoxChars
	case graphNodeShapeHexagon:
		return unicodeHexagonBoxChars
	case graphNodeShapeParallelogram:
		return unicodeParallelogramBoxChars
	default:
		return unicodeLightGraphCharset.box
	}
}

func (g graph) boxCharsForSubgraph() graphBoxChars {
	if g.useAscii {
		return asciiGraphCharset.box
	}

	switch g.graphBoxStyle {
	case "rounded":
		return unicodeRoundedBoxChars
	case "double":
		return unicodeDoubleBoxChars
	case "heavy":
		return unicodeHeavyBoxChars
	default:
		return unicodeLightGraphCharset.box
	}
}

func (g graph) lineCharsForEdge(e *edge) graphLineChars {
	charset := g.charset()
	switch g.edgeLineStyle(e) {
	case graphEdgeLineStyleHeavy:
		return charset.heavyLine
	case graphEdgeLineStyleDashed:
		return charset.dashedLine
	default:
		return charset.line
	}
}

func (g graph) edgeLineStyle(e *edge) graphEdgeLineStyle {
	lineStyle := e.lineStyle
	if !e.lineStyleSet && g.graphEdgeStyle != "" {
		lineStyle = g.graphEdgeStyle
	}
	if lineStyle == "" {
		return graphEdgeLineStyleLight
	}
	return lineStyle
}

func (g graph) edgeDrawPriority(e *edge) int {
	return edgeLineStylePriority(g.edgeLineStyle(e))
}

func edgeLineStylePriority(style graphEdgeLineStyle) int {
	switch style {
	case graphEdgeLineStyleHeavy:
		return graphLinePriorityHeavy
	case graphEdgeLineStyleDashed:
		return graphLinePriorityDashed
	default:
		return graphLinePriorityLight
	}
}

func (c graphCharset) arrowHead(dir, fallback direction) string {
	if dir == Middle {
		dir = fallback
	}
	switch dir {
	case Up:
		return c.arrowHeads.up
	case Down:
		return c.arrowHeads.down
	case Left:
		return c.arrowHeads.left
	case Right:
		return c.arrowHeads.right
	case UpperRight:
		if c.arrowHeads.upperRight != "" {
			return c.arrowHeads.upperRight
		}
	case UpperLeft:
		if c.arrowHeads.upperLeft != "" {
			return c.arrowHeads.upperLeft
		}
	case LowerRight:
		if c.arrowHeads.lowerRight != "" {
			return c.arrowHeads.lowerRight
		}
	case LowerLeft:
		if c.arrowHeads.lowerLeft != "" {
			return c.arrowHeads.lowerLeft
		}
	}

	if fallback != dir {
		return c.arrowHead(fallback, Middle)
	}
	return c.arrowHeads.fallback
}

func (c graphLineChars) corner(prevDir, nextDir direction) string {
	switch {
	case (prevDir == Right && nextDir == Down) || (prevDir == Up && nextDir == Left):
		return c.topRight
	case (prevDir == Right && nextDir == Up) || (prevDir == Down && nextDir == Left):
		return c.bottomRight
	case (prevDir == Left && nextDir == Down) || (prevDir == Up && nextDir == Right):
		return c.topLeft
	case (prevDir == Left && nextDir == Up) || (prevDir == Down && nextDir == Right):
		return c.bottomLeft
	default:
		return c.cross
	}
}

func (c graphCharset) isJunctionChar(char string) bool {
	if _, ok := c.junctionSet[char]; ok {
		return true
	}
	_, ok := graphCharConnections[char]
	return ok
}

func (c graphCharset) mergeJunctions(c1, c2 string) string {
	if merged, ok := c.junctionMap[c1][c2]; ok {
		log.Debugf("Merging %s and %s to %s", c1, c2, merged)
		return merged
	}
	conn1, ok1 := graphCharConnections[c1]
	conn2, ok2 := graphCharConnections[c2]
	if ok1 && ok2 {
		merged := junctionForConnections(conn1|conn2, maxLinePriority(graphCharPriority(c1), graphCharPriority(c2)))
		log.Debugf("Merging %s and %s to %s", c1, c2, merged)
		return merged
	}
	return c1
}

func graphCharPriority(char string) int {
	if priority, ok := graphCharPriorities[char]; ok {
		return priority
	}
	return graphLinePriorityLight
}

func maxLinePriority(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func junctionForConnections(conn lineConnection, priority int) string {
	if priority >= graphLinePriorityHeavy {
		return heavyJunctionForConnections(conn)
	}
	return lightJunctionForConnections(conn)
}

func heavyJunctionForConnections(conn lineConnection) string {
	switch conn {
	case connLeft | connRight:
		return "━"
	case connUp | connDown:
		return "┃"
	case connRight | connDown:
		return "┏"
	case connDown | connLeft:
		return "┓"
	case connUp | connRight:
		return "┗"
	case connUp | connLeft:
		return "┛"
	case connUp | connRight | connDown:
		return "┣"
	case connUp | connDown | connLeft:
		return "┫"
	case connRight | connDown | connLeft:
		return "┳"
	case connUp | connRight | connLeft:
		return "┻"
	default:
		return "╋"
	}
}

func lightJunctionForConnections(conn lineConnection) string {
	switch conn {
	case connLeft | connRight:
		return "─"
	case connUp | connDown:
		return "│"
	case connRight | connDown:
		return "┌"
	case connDown | connLeft:
		return "┐"
	case connUp | connRight:
		return "└"
	case connUp | connLeft:
		return "┘"
	case connUp | connRight | connDown:
		return "├"
	case connUp | connDown | connLeft:
		return "┤"
	case connRight | connDown | connLeft:
		return "┬"
	case connUp | connRight | connLeft:
		return "┴"
	default:
		return "┼"
	}
}
