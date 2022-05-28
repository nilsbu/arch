package rule

import (
	"fmt"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/graph"
)

// RoomOrientation identifies the orientation of an area.
// The area is defined as the direction someone is looking in when they enter through the first door.
func RoomOrientation(g *graph.Graph, nidx graph.NodeIndex) area.Direction {
	return area.Turn(area.GetDirection(g, nidx, g.Node(nidx).Edges[0]), 180)
}

// InheritEdges passes on the edges of a parent to children depending on their position.
// ErrInvalidGraph is returned if a door cannot be assigned. This can be the case whe the door lies on the edge of
// multiple children.
func InheritEdges(g *graph.Graph, nidx graph.NodeIndex) error {
	nidxs := g.Children(nidx)
	if len(nidxs) == 0 {
		return nil
	}

	rects := make([]area.Rectangle, len(nidxs))
	for i, cnidx := range nidxs {
		rects[i] = (*area.AreaNode)(g.Node(cnidx)).GetRect()
	}

	for _, eidx := range g.Node(nidx).Edges {
		if err := passOn(g, eidx, nidx, nidxs, rects); err != nil {
			return err
		}
	}
	return nil
}

func passOn(
	g *graph.Graph,
	eidx graph.EdgeIndex,
	nidx graph.NodeIndex,
	nidxs []graph.NodeIndex,
	rects []area.Rectangle) error {
	door := (*area.DoorEdge)(g.Edge(eidx)).GetPos()

	for i, rect := range rects {
		if ((rect.X0 == door.X || rect.X1 == door.X) && rect.Y0 < door.Y && rect.Y1 > door.Y) ||
			((rect.Y0 == door.Y || rect.Y1 == door.Y) && rect.X0 < door.X && rect.X1 > door.X) {
			return g.InheritEdge(nidx, nidxs[i], []graph.EdgeIndex{eidx})
		}
	}

	return fmt.Errorf("%w: door %v at [%v, %v] cannot be assigned to child area", ErrInvalidGraph,
		eidx, door.X, door.Y)
}
