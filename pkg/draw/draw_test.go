package draw_test

import (
	"errors"
	"testing"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/draw"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/pkg/world"
)

func TestDraw(t *testing.T) {
	f := world.Tile{Type: world.Free}
	w := world.Tile{Type: world.Wall}
	d := world.Tile{Type: world.Door}
	o := world.Tile{Type: world.Occupied, Texture: 1}
	p := world.Tile{Type: world.Occupied, Texture: 2}

	for _, c := range []struct {
		name  string
		graph func() *graph.Graph
		tiles [][]world.Tile
		err   error
	}{
		{
			"no area set",
			func() *graph.Graph {
				return graph.New(nil)
			},
			[][]world.Tile{},
			draw.ErrInvalidGraph,
		},
		{
			"child has no rectangle",
			func() *graph.Graph {
				g := graph.New(nil)
				node := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 5, Y1: 3})
				n1, _ := g.Add(graph.NodeIndex{})
				node = (*area.AreaNode)(g.Node(n1))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 3, Y1: 3})
				g.Add(graph.NodeIndex{})

				return g
			},
			[][]world.Tile{
				{w, w, w, w, w, w},
				{w, f, f, d, f, w},
				{w, f, f, w, f, w},
				{w, w, w, w, w, w},
			},
			draw.ErrInvalidGraph,
		},
		{
			"only one room",
			func() *graph.Graph {
				g := graph.New(nil)
				node := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 3, Y1: 3})
				return g
			},
			[][]world.Tile{
				{w, w, w, w},
				{w, f, f, w},
				{w, f, f, w},
				{w, w, w, w},
			},
			nil,
		},
		{
			"two rooms with door",
			func() *graph.Graph {
				g := graph.New(nil)
				node := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 5, Y1: 3})
				n1, _ := g.Add(graph.NodeIndex{})
				node = (*area.AreaNode)(g.Node(n1))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 3, Y1: 3})
				n2, _ := g.Add(graph.NodeIndex{})
				node = (*area.AreaNode)(g.Node(n2))
				node.SetRect(area.Rectangle{X0: 3, Y0: 0, X1: 5, Y1: 3})
				e1, _ := g.Link(n1, n2)
				edge := (*area.DoorEdge)(g.Edge(e1))
				edge.SetPos(area.Point{X: 3, Y: 1})

				return g
			},
			[][]world.Tile{
				{w, w, w, w, w, w},
				{w, f, f, d, f, w},
				{w, f, f, w, f, w},
				{w, w, w, w, w, w},
			},
			nil,
		},
		{
			"render disabled",
			func() *graph.Graph {
				g := graph.New(nil)
				node := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 3, Y1: 3})
				node.Properties["render"] = false
				return g
			},
			[][]world.Tile{
				{f, f, f, f},
				{f, f, f, f},
				{f, f, f, f},
				{f, f, f, f},
			},
			nil,
		},
		{
			"two interior rooms disabled",
			func() *graph.Graph {
				g := graph.New(nil)
				node := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 5, Y1: 3})
				n1, _ := g.Add(graph.NodeIndex{})
				node = (*area.AreaNode)(g.Node(n1))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 3, Y1: 3})
				node.Properties["render"] = false
				n2, _ := g.Add(graph.NodeIndex{})
				node = (*area.AreaNode)(g.Node(n2))
				node.SetRect(area.Rectangle{X0: 3, Y0: 0, X1: 5, Y1: 3})
				node.Properties["render"] = false
				e1, _ := g.Link(n1, n2)
				edge := (*area.DoorEdge)(g.Edge(e1))
				edge.SetPos(area.Point{X: 3, Y: 1})
				edge.Properties["render"] = false

				return g
			},
			[][]world.Tile{
				{w, w, w, w, w, w},
				{w, f, f, f, f, w},
				{w, f, f, f, f, w},
				{w, w, w, w, w, w},
			},
			nil,
		},
		{
			"only one room",
			func() *graph.Graph {
				g := graph.New(nil)
				node := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
				node.SetRect(area.Rectangle{X0: 0, Y0: 0, X1: 3, Y1: 3})
				left, _ := g.Add(graph.NodeIndex{})
				(*area.AreaNode)(g.Node(left)).SetRect(area.Rectangle{X0: 1, Y0: 1, X1: 1, Y1: 2})
				g.Node(left).Properties["object"] = 1
				right, _ := g.Add(graph.NodeIndex{})
				(*area.AreaNode)(g.Node(right)).SetRect(area.Rectangle{X0: 2, Y0: 1, X1: 2, Y1: 2})
				g.Node(right).Properties["object"] = 2
				return g
			},
			[][]world.Tile{
				{w, w, w, w},
				{w, o, p, w},
				{w, o, p, w},
				{w, w, w, w},
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if data, err := draw.Draw(c.graph()); err != nil && c.err == nil {
				t.Error("unexpected error:", err)
			} else if err == nil && c.err != nil {
				t.Error("expected error but none ocurred")
			} else if !errors.Is(err, c.err) {
				t.Error("wrong type of error")
			} else if err == nil {
				if data.Height() != len(c.tiles) || data.Width() != len(c.tiles[0]) {
					t.Errorf("wrong size: expected (%v, %v) but got (%v, %v)",
						len(c.tiles[0]), len(c.tiles), data.Width(), data.Height())
				} else {
					for y, line := range c.tiles {
						for x, expect := range line {
							actual := data.Get(x, y)
							if expect.Type != actual.Type || expect.Texture != actual.Texture {
								t.Errorf("at (%v, %v): expect %v, actual %v",
									x, y, expect, actual)
							}
						}
					}
				}
			}
		})
	}
}
