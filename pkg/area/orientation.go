package area

import (
	"github.com/nilsbu/arch/pkg/graph"
)

// GetDirection returns the direction in which a door lies with respect to an area.
// The edge must belong to the node. The side of the rectangle on which the door lies, determines the result.
func GetDirection(g *graph.Graph, nidx graph.NodeIndex, eidx graph.EdgeIndex) Direction {
	rect := (*AreaNode)(g.Node(nidx)).GetRect()
	pos := (*DoorEdge)(g.Edge(eidx)).GetPos()
	switch {
	case rect.X0 == pos.X:
		return Left
	case rect.X1 == pos.X:
		return Right
	case rect.Y0 == pos.Y:
		return Up
	default:
		// Assume that door was set up properly and correct indices were used, no more to check.
		return Down
	}
}

// Turn rotates directions.
// Rotations are done in 90 degree intervals. If angle isn't divisible by 90, it is rounded to the closest valid value.
// Positive numbers denote clockwise turns, negative ones are counter-clockwise.
func Turn(direction Direction, angle int) Direction {
	angle = angle % 360
	if angle < 0 {
		angle += 360
	}
	a := ((angle/45 + 1) / 2) % 4

	direction <<= a
	direction |= direction >> 4
	return direction & 0b0000_1111
}
