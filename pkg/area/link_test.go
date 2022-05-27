package area_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/graph"
)

func TestCreateDoor(t *testing.T) {
	for _, c := range []struct {
		name         string
		graph        func() *graph.Graph
		nidx0, nidx1 graph.NodeIndex
		position     float64
		linkPos      area.Point
		err          error
	}{
		{
			"pos out of range",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 2, 2}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{2, 0, 2, 2}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			-1,
			area.Point{},
			area.ErrInvalidDoor,
		},
		{
			"second to the right of the first",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 2, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{2, 0, 2, 10}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.3,
			area.Point{2, 3},
			nil,
		},
		{
			"second to the left of the first",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{2, 0, 2, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{0, 0, 2, 10}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.3,
			area.Point{2, 7},
			nil,
		},
		{
			"second to the bottom of the first",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{-2, 1, 2, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{0, 10, 2, 12}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.45,
			area.Point{1, 10},
			nil,
		},
		{
			"second to the top of the first",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{2, 4, 12, 60}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{4, 0, 14, 4}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.75,
			area.Point{6, 4},
			nil,
		},
		{
			"no intersection",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{2, 4, 12, 8}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{4, 0, 14, 3}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.75,
			area.Point{},
			area.ErrInvalidDoor,
		},
		{
			"meet at corner",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{2, 4, 12, 8}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{0, 0, 2, 4}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.75,
			area.Point{},
			area.ErrInvalidDoor,
		},
		{
			"intersect at more than one side",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{2, 4, 12, 8}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{0, 0, 3, 5}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.75,
			area.Point{},
			area.ErrInvalidDoor,
		},
		{
			"one rect inside the other",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{0, 0, 10, 10}
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{1, 1, 3, 5}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.75,
			area.Point{},
			area.ErrInvalidDoor,
		},
		{
			"invalid link",
			func() *graph.Graph {
				g := graph.New(nil)
				g.Node(graph.NodeIndex{}).Properties["rect"] = area.Rectangle{0, 0, 10, 10}
				n0, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["rect"] = area.Rectangle{1, 1, 3, 5}
				return g
			},
			graph.NodeIndex{0, 0}, graph.NodeIndex{1, 0},
			.75,
			area.Point{},
			graph.ErrIllegalAction,
		},
		{
			"one rect isn't set",
			func() *graph.Graph {
				g := graph.New(nil)
				g.Add(graph.NodeIndex{})
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{4, 0, 14, 4}
				return g
			},
			graph.NodeIndex{1, 0}, graph.NodeIndex{1, 1},
			.75,
			area.Point{},
			area.ErrInvalidDoor,
		},
		{
			"one rect isn't set (2)",
			func() *graph.Graph {
				g := graph.New(nil)
				g.Add(graph.NodeIndex{})
				n1, _ := g.Add(graph.NodeIndex{})
				g.Node(n1).Properties["rect"] = area.Rectangle{4, 0, 14, 4}
				return g
			},
			graph.NodeIndex{1, 1}, graph.NodeIndex{1, 0},
			.75,
			area.Point{},
			area.ErrInvalidDoor,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			g := c.graph()
			if err := area.CreateDoor(g, c.nidx0, c.nidx1, c.position); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error\nexpect: %v\nactual: %v", c.err, err)
			} else if err == nil {
				linkPos := (*area.DoorEdge)(g.Edge(g.Node(c.nidx0).Edges[0])).GetPos()
				if !reflect.DeepEqual(c.linkPos, linkPos) {
					t.Errorf("wrong position: expect %v, actual %v", c.linkPos, linkPos)
				}
			}
		})
	}
}

func TestGetUnsetPos(t *testing.T) {
	door := &area.DoorEdge{Properties: graph.Properties{}}
	if !reflect.DeepEqual(door.GetPos(), area.Point{}) {
		t.Errorf("door position was (unfathomably) set: %v", door.GetPos())
	}
}
