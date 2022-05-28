package rule

import (
	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/graph"
)

// RoomOrientation identifies the orientation of an area.
// The area is defined as the direction someone is looking in when they enter through the first door.
func RoomOrientation(g *graph.Graph, nidx graph.NodeIndex) area.Direction {
	return area.Turn(area.GetDirection(g, nidx, g.Node(nidx).Edges[0]), 180)
}
