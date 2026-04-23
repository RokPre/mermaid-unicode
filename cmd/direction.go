package cmd

type direction genericCoord

var (
	Up         = direction{1, 0}
	Down       = direction{1, 2}
	Left       = direction{0, 1}
	Right      = direction{2, 1}
	UpperRight = direction{2, 0}
	UpperLeft  = direction{0, 0}
	LowerRight = direction{2, 2}
	LowerLeft  = direction{0, 2}
	Middle     = direction{1, 1}
)

func (d direction) getOpposite() direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case Left:
		return Right
	case Right:
		return Left
	case UpperRight:
		return LowerLeft
	case UpperLeft:
		return LowerRight
	case LowerRight:
		return UpperLeft
	case LowerLeft:
		return UpperRight
	case Middle:
		return Middle
	}
	panic("Unknown direction")
}

func (c gridCoord) Direction(dir direction) gridCoord {
	return gridCoord{x: c.x + dir.x, y: c.y + dir.y}
}
func (c drawingCoord) Direction(dir direction) drawingCoord {
	return drawingCoord{x: c.x + dir.x, y: c.y + dir.y}
}

func (g graph) selfReferenceDirection() (direction, direction, direction, direction) {
	if g.isHorizontalLayout() {
		return Right, Down, Down, Right
	}
	return Down, Right, Right, Down
}

func (g graph) isBackwardFlowDirection(d direction) bool {
	if g.isHorizontalLayout() {
		return d == Left || d == UpperLeft || d == LowerLeft
	}
	return d == Up || d == UpperLeft || d == UpperRight
}

func (g graph) canonicalDirection(d direction) direction {
	switch g.graphDirection {
	case "RL":
		return mirrorHorizontalDirection(d)
	case "BT":
		return mirrorVerticalDirection(d)
	default:
		return d
	}
}

func (g graph) displayDirection(d direction) direction {
	return g.canonicalDirection(d)
}

func mirrorHorizontalDirection(d direction) direction {
	switch d {
	case Left:
		return Right
	case Right:
		return Left
	case UpperLeft:
		return UpperRight
	case UpperRight:
		return UpperLeft
	case LowerLeft:
		return LowerRight
	case LowerRight:
		return LowerLeft
	default:
		return d
	}
}

func mirrorVerticalDirection(d direction) direction {
	switch d {
	case Up:
		return Down
	case Down:
		return Up
	case UpperLeft:
		return LowerLeft
	case LowerLeft:
		return UpperLeft
	case UpperRight:
		return LowerRight
	case LowerRight:
		return UpperRight
	default:
		return d
	}
}

func (g graph) determineStartAndEndDir(e *edge) (direction, direction, direction, direction) {
	if e.from == e.to {
		return g.selfReferenceDirection()
	}
	d := g.canonicalDirection(determineDirection(genericCoord(*e.from.gridCoord), genericCoord(*e.to.gridCoord)))
	var preferredDir, preferredOppositeDir, alternativeDir, alternativeOppositeDir direction

	// Check if this is a backwards flowing edge
	isBackwards := g.isBackwardFlowDirection(d)

	// LR: prefer vertical over horizontal
	// TD: prefer horizontal over vertical
	// TODO: This causes some squirmy lines if the corner spot is already occupied.
	// For backwards edges, use special start positions: Down in LR mode, Right in TD mode
	switch d {
	case LowerRight:
		if g.isHorizontalLayout() {
			preferredDir = Down
			preferredOppositeDir = Left
			alternativeDir = Right
			alternativeOppositeDir = Up
		} else {
			preferredDir = Right
			preferredOppositeDir = Up
			alternativeDir = Down
			alternativeOppositeDir = Left
		}
	case UpperRight:
		if g.isHorizontalLayout() {
			preferredDir = Up
			preferredOppositeDir = Left
			alternativeDir = Right
			alternativeOppositeDir = Down
		} else {
			preferredDir = Right
			preferredOppositeDir = Down
			alternativeDir = Up
			alternativeOppositeDir = Left
		}
	case LowerLeft:
		if g.isHorizontalLayout() {
			// Backwards flow in LR mode - start from Down, arrive at Down
			preferredDir = Down
			preferredOppositeDir = Down // Edge goes to bottom of destination
			alternativeDir = Left
			alternativeOppositeDir = Up
		} else {
			preferredDir = Left
			preferredOppositeDir = Up
			alternativeDir = Down
			alternativeOppositeDir = Right
		}
	case UpperLeft:
		if g.isHorizontalLayout() {
			// Backwards flow in LR mode - start from Down, arrive at Down
			preferredDir = Down
			preferredOppositeDir = Down // Edge goes to bottom of destination
			alternativeDir = Left
			alternativeOppositeDir = Down
		} else {
			// Backwards flow in TD mode - start from Right, arrive at Right
			preferredDir = Right
			preferredOppositeDir = Right // Edge goes to right of destination
			alternativeDir = Up
			alternativeOppositeDir = Right
		}
	default:
		// Handle direct backwards flow cases
		if isBackwards {
			if g.isHorizontalLayout() && d == Left {
				// Direct left flow in LR mode - start from Down, arrive at Down
				preferredDir = Down
				preferredOppositeDir = Down // Edge goes to bottom of destination
				alternativeDir = Left
				alternativeOppositeDir = Right
			} else if g.isVerticalLayout() && d == Up {
				// Direct up flow in TD mode - start from Right, arrive at Right
				preferredDir = Right
				preferredOppositeDir = Right // Edge goes to right of destination
				alternativeDir = Up
				alternativeOppositeDir = Down
			} else {
				preferredDir = d
				preferredOppositeDir = preferredDir.getOpposite()
				alternativeDir = d
				alternativeOppositeDir = preferredOppositeDir
			}
		} else {
			preferredDir = d
			preferredOppositeDir = preferredDir.getOpposite()
			// TODO: just return null and don't calculate alternative path
			alternativeDir = d
			alternativeOppositeDir = preferredOppositeDir
		}
	}
	return g.displayDirection(preferredDir), g.displayDirection(preferredOppositeDir), g.displayDirection(alternativeDir), g.displayDirection(alternativeOppositeDir)
}
