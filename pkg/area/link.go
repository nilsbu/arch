package area

import (
	"errors"
	"fmt"
	"math"

	"github.com/nilsbu/arch/pkg/graph"
)

// ErrInvalidDoor is returned when linking two areas through a door isn't possible.
var ErrInvalidDoor = errors.New("door invalid")

// CreateDoor creates DoorEdge between two AreaNodes.
// The position parameter specifies where the door is positioned. First a line, where the two areas touch, is
// established. From the perspective of the first area, facing the second one, 0 denotes the left end of the line, 1
// deontes the right end. Any value in between can be used to denote other points of the line.
//
// The areas must touch in a line of non-zero length. If they aren't connected, overlap or touch in only a point, an
// error is returned.
func CreateDoor(g *graph.Graph, nidx0, nidx1 graph.NodeIndex, position float64) error {
	if position < 0 || position > 1 {
		return fmt.Errorf("%w: position must be in range [0, 1] but was %v",
			ErrInvalidDoor, position)
	}

	node0, node1 := (*AreaNode)(g.Node(nidx0)), (*AreaNode)(g.Node(nidx1))
	rect0, rect1 := node0.GetRect(), node1.GetRect()
	if rect0.X0 == rect0.X1 && rect0.Y0 == rect0.Y1 {
		return fmt.Errorf("%w: first rectangle isn't set",
			ErrInvalidDoor)
	}
	if rect1.X0 == rect1.X1 && rect1.Y0 == rect1.Y1 {
		return fmt.Errorf("%w: second rectangle isn't set",
			ErrInvalidDoor)
	}

	if eidx, err := g.Link(nidx0, nidx1); err != nil {
		return err
	} else if inter, err := intersect(rect0, rect1); err != nil {
		return err
	} else {
		if inter.X1 == rect0.X0 || inter.Y1 == rect0.Y0 {
			position = 1 - position
		}
		pt := Point{
			X: inter.X0 + int(math.Round(float64(inter.X1-inter.X0)*position)),
			Y: inter.Y0 + int(math.Round(float64(inter.Y1-inter.Y0)*position)),
		}
		(*DoorEdge)(g.Edge(eidx)).SetPos(pt)
		return nil
	}
}

func intersect(ar, br Rectangle) (Rectangle, error) {
	x0 := max(ar.X0, br.X0)
	y0 := max(ar.Y0, br.Y0)
	x1 := min(ar.X1, br.X1)
	y1 := min(ar.Y1, br.Y1)
	inter := Rectangle{x0, y0, x1, y1}
	if (x0 > x1 || y0 > y1) || (x0 == x1 && y0 == y1) {
		return Rectangle{}, fmt.Errorf("%w, rectangles %v and %v don't intersect",
			ErrInvalidDoor, ar, br)
	} else if x1 > x0 && y1 > y0 {
		return Rectangle{}, fmt.Errorf("%w, rectangles %v and %v intersect at more than one side",
			ErrInvalidDoor, ar, br)
	} else {
		return inter, nil
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}
