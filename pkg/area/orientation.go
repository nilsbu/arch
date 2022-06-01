package area

import (
	"errors"
	"fmt"

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

// Difference calculates the difference in angles between two Directions.
// The result is in range (-180, 180].
// It is only designed for single directions, not combined ones.
// If the input is a combined direction, the output is undefined.
func Difference(a, b Direction) int {
	i := 0
	for ; a > b; i++ {
		a >>= 1
	}
	for ; a < b; i-- {
		a <<= 1
	}
	if i <= -2 {
		i += 4
	} else if i >= 3 {
		i -= 4
	}
	return 90 * i
}

// ErrInvalidRotation is returned is a rectangle cannot be rotated.
var ErrInvalidRotation = errors.New("cannot rotate")

// RotateWithin rotates a rectangle within another.
func RotateWithin(rect, in Rectangle, from, to Direction, anchor Anchor) (Rectangle, error) {
	// TODO write more detailed documentation
	nFrom := Difference(Down, from)
	aFrom := Difference(from, Down)
	aD := Difference(to, from)

	pCenter := calcAnchorPoint(in, Center)

	rectNormalized, inNormalized := rotateAround(rect, pCenter, nFrom), rotateAround(in, pCenter, nFrom)

	pAnchorPre := calcAnchorPoint(inNormalized, anchor)
	pAnchorPost := calcAnchorPoint(inNormalized, rotateAnchor(anchor, aD))

	rectShifted := move(rectNormalized, pAnchorPost, pAnchorPre)
	rectRotated := rotateAround(rectShifted, pAnchorPost, aD)
	rectFinal := rotateAround(rectRotated, pCenter, aFrom)

	return rectFinal, assertInside(rectFinal, in)
}

func rotateAnchor(anchor Anchor, angle int) Anchor {
	if anchor == Center {
		return Center
	} else {
		switch angle {
		case -90:
			anchor = (2 - (anchor & 2)) | (anchor & 1)
			return ((anchor & 2) >> 1) | ((anchor & 1) << 1)
		case 0:
			return anchor
		case 90:
			anchor = (anchor & 2) | (1 - (anchor & 1))
			return ((anchor & 2) >> 1) | ((anchor & 1) << 1)
		default:
			return (2 - (anchor & 2)) | (1 - (anchor & 1))
		}
	}
}

func calcAnchorPoint(rect Rectangle, anchor Anchor) Point {
	if anchor == Center {
		return Point{
			X: (rect.X0 + rect.X1) / 2,
			Y: (rect.Y0 + rect.Y1) / 2,
		}
	} else {
		var res Point
		if anchor&1 > 0 {
			res.X = rect.X0
		} else {
			res.X = rect.X1
		}
		if anchor&2 > 0 {
			res.Y = rect.Y1
		} else {
			res.Y = rect.Y0
		}
		return res
	}
}

func move(rect Rectangle, fromAnchor, toAnchor Point) Rectangle {
	dx := fromAnchor.X - toAnchor.X
	dy := fromAnchor.Y - toAnchor.Y
	return Rectangle{
		X0: rect.X0 + dx,
		Y0: rect.Y0 + dy,
		X1: rect.X1 + dx,
		Y1: rect.Y1 + dy,
	}
}

func rotateAround(rect Rectangle, anchor Point, angle int) Rectangle {
	switch angle {
	case -90:
		return normalize(Rectangle{
			X0: anchor.X + (rect.Y0 - anchor.Y),
			Y0: anchor.Y - (rect.X0 - anchor.X),
			X1: anchor.X + (rect.Y1 - anchor.Y),
			Y1: anchor.Y - (rect.X1 - anchor.X),
		})
	case 0:
		return rect
	case 90:
		return normalize(Rectangle{
			X0: anchor.X - (rect.Y0 - anchor.Y),
			Y0: anchor.Y + (rect.X0 - anchor.X),
			X1: anchor.X - (rect.Y1 - anchor.Y),
			Y1: anchor.Y + (rect.X1 - anchor.X),
		})
	default:
		return normalize(Rectangle{
			X0: anchor.X - (rect.X0 - anchor.X),
			Y0: anchor.Y - (rect.Y0 - anchor.Y),
			X1: anchor.X - (rect.X1 - anchor.X),
			Y1: anchor.Y - (rect.Y1 - anchor.Y),
		})
	}
}

func normalize(rect Rectangle) Rectangle {
	if rect.X0 > rect.X1 {
		rect.X0, rect.X1 = rect.X1, rect.X0
	}
	if rect.Y0 > rect.Y1 {
		rect.Y0, rect.Y1 = rect.Y1, rect.Y0
	}
	return rect
}

func assertInside(rect, in Rectangle) error {
	if rect.X0 < in.X0 || rect.X1 > in.X1 || rect.Y0 < in.Y0 || rect.Y1 > in.Y1 {
		return fmt.Errorf("%w: rotation result %v not inside %v", ErrInvalidRotation,
			rect, in)
	} else {
		return nil
	}
}

// CalcAnchorPoint calculates the position of an anchor point.
func CalcAnchorPoint(rect Rectangle, anchor Anchor, direction Direction) Point {
	return calcAnchorPoint(rect, rotateAnchor(anchor, Difference(direction, Down)))
}
