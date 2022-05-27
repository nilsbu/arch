package area

import "github.com/nilsbu/arch/pkg/graph"

// AreaNode is a node that has a rectangular area.
// The rect will default to {0, 0, 0, 0} is not otherwise specified.
// It uses the property "rect".
type AreaNode graph.Node

// GetRect returns the area of the node.
func (n *AreaNode) GetRect() Rectangle {
	if rect, ok := n.Properties["rect"]; ok {
		return rect.(Rectangle)
	} else {
		return Rectangle{}
	}
}

// SetRect sets the area of the node.
func (n *AreaNode) SetRect(rect Rectangle) {
	n.Properties["rect"] = rect
}

// Rectangle describes an axis-aligned rectangle.
// It fills the area of all points (x, y) that fulfill X0 <= x <= X1 && Y0 <= y <= Y1.
// In other words, (X0, Y0) is the minimal point, (X1, Y1) is the maximal point.
type Rectangle struct {
	X0, Y0, X1, Y1 int
}

// DoorEdge that represents a door.
// It uses the property "pos" to define the position of the door. The orientation is defined implicitely by the nodes
// that it connects. The position defaults to [0, 0] when not set.
type DoorEdge graph.Edge

// GetPos returns the position of the door.
func (e *DoorEdge) GetPos() Point {
	if pos, ok := e.Properties["pos"]; ok {
		return pos.(Point)
	} else {
		return Point{}
	}
}

// SetPos sets the position of the door.
func (e *DoorEdge) SetPos(pos Point) {
	e.Properties["pos"] = pos
}

// Point specifies a position.
type Point struct {
	X, Y int
}

// Direction speficies a direction.
type Direction byte

const (
	Up    Direction = 1
	Right Direction = 2
	Down  Direction = 4
	Left  Direction = 8
)
