package area_test

import (
	"testing"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/graph"
)

func TestGetOrientation(t *testing.T) {
	for _, c := range []struct {
		name        string
		graph       func() *graph.Graph
		nidx        graph.NodeIndex
		eidx        graph.EdgeIndex
		orientation area.Direction
	}{
		{
			"down center",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 10, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{2, 10, 12, 20}

				area.CreateDoor(g, n0, n1, .5)
				return g
			},
			graph.NodeIndex{1, 0},
			graph.EdgeIndex(0),
			area.Down,
		},
		{
			"up center",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 10, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{2, 10, 12, 20}

				area.CreateDoor(g, n0, n1, .5)
				return g
			},
			graph.NodeIndex{1, 1},
			graph.EdgeIndex(0),
			area.Up,
		},
		{
			"left off-center",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 10, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{10, 0, 20, 8}

				area.CreateDoor(g, n0, n1, .25)
				return g
			},
			graph.NodeIndex{1, 1},
			graph.EdgeIndex(0),
			area.Left,
		},
		{
			"right off-center",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 10, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{10, 0, 20, 8}

				area.CreateDoor(g, n0, n1, .25)
				return g
			},
			graph.NodeIndex{1, 0},
			graph.EdgeIndex(0),
			area.Right,
		},
		{
			"right off-center",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 10, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{10, 0, 20, 8}

				area.CreateDoor(g, n0, n1, .25)
				return g
			},
			graph.NodeIndex{1, 0},
			graph.EdgeIndex(0),
			area.Right,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			o := area.GetDirection(c.graph(), c.nidx, c.eidx)
			if o != c.orientation {
				t.Errorf("expected direction %v but got %v",
					c.orientation, o)
			}
		})
	}
}

func TestTurn(t *testing.T) {
	for _, c := range []struct {
		name  string
		in    area.Direction
		angle int
		out   area.Direction
	}{
		{
			"right, no turn",
			area.Right,
			0,
			area.Right,
		},
		{
			"down, no turn",
			area.Down,
			0,
			area.Down,
		},
		{
			"up, turn right",
			area.Up,
			90,
			area.Right,
		},
		{
			"up, turn around",
			area.Up,
			180,
			area.Down,
		},
		{
			"up, turn around",
			area.Up,
			180,
			area.Down,
		},
		{
			"left, turn around",
			area.Left,
			180,
			area.Right,
		},
		{
			"left, turn right",
			area.Left,
			90,
			area.Up,
		},
		{
			"up go spinny",
			area.Up,
			-1080,
			area.Up,
		},
		{
			"up, turn left thrice",
			area.Up,
			-270,
			area.Right,
		},
		{
			"up, turn just enough left",
			area.Up,
			-46,
			area.Left,
		},
		{
			"up and left, turn left",
			area.Up | area.Left,
			-90,
			area.Left | area.Down,
		},
		{
			"up, left and down, turn around",
			area.Up | area.Left | area.Down,
			180,
			area.Up | area.Right | area.Down,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			out := area.Turn(c.in, c.angle)
			if out != c.out {
				t.Errorf("expected direction %v but got %v",
					c.out, out)
			}
		})
	}
}
