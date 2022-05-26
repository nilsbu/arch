package area

import "github.com/nilsbu/arch/pkg/graph"

type AreaNode graph.Node

func (n *AreaNode) GetRect() Rectangle {
	if rect, ok := n.Properties["rect"]; ok {
		return rect.(Rectangle)
	} else {
		return Rectangle{}
	}
}

func (n *AreaNode) SetRect(rect Rectangle) {
	n.Properties["rect"] = rect
}

// Rectangle describes an axis-aligned rectangle.
// It fills the area of all points (x, y) that fulfill X0 <= x <= X1 && Y0 <= y <= Y1.
// In other words, (X0, Y0) is the minimal point, (X1, Y1) is the maximal point.
type Rectangle struct {
	X0, Y0, X1, Y1 int
}

type DoorEdge graph.Edge

func (e *DoorEdge) GetRect() Point {
	if pos, ok := e.Properties["pos"]; ok {
		return pos.(Point)
	} else {
		return Point{}
	}
}

func (e *DoorEdge) SetPos(pos Point) {
	e.Properties["pos"] = pos
}

type Point struct {
	X, Y int
}

type Direction byte

const (
	Up    Direction = 1
	Down  Direction = 2
	Left  Direction = 4
	Right Direction = 8
)
