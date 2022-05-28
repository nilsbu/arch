package area_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/graph"
)

func fillNodes(g *graph.Graph, base graph.Properties, n int) {
	g.Node(graph.NodeIndex{}).Properties = base
	for i := 0; i < n; i++ {
		g.Add(graph.NodeIndex{})
	}
}

func TestSplit(t *testing.T) {
	for _, c := range []struct {
		name      string
		base      graph.Properties
		at        []float64
		direction area.Direction
		out       []area.Rectangle
		err       error
	}{
		{
			"len(into) != len(at) + 1",
			graph.Properties{
				"rect": area.Rectangle{2, 3, 4, 5},
			},
			[]float64{0.2},
			area.Up,
			[]area.Rectangle{},
			area.ErrInvalidSplit,
		},
		// TODO with invalid Direction
		// TODO with at <0 or >1
		{
			"one split copies rect",
			graph.Properties{
				"rect": area.Rectangle{2, 3, 4, 5},
			},
			[]float64{},
			area.Up,
			[]area.Rectangle{{2, 3, 4, 5}},
			nil,
		},
		{
			"split upwards in the middle",
			graph.Properties{
				"rect": area.Rectangle{2, 0, 4, 6},
			},
			[]float64{.5},
			area.Up,
			[]area.Rectangle{{2, 3, 4, 6}, {2, 0, 4, 3}},
			nil,
		},
		{
			"split downwards twice",
			graph.Properties{
				"rect": area.Rectangle{2, 0, 4, 10},
			},
			[]float64{.3, .8},
			area.Down,
			[]area.Rectangle{{2, 0, 4, 3}, {2, 3, 4, 8}, {2, 8, 4, 10}},
			nil,
		},
		{
			"split right with one of size zero",
			graph.Properties{
				"rect": area.Rectangle{2, 0, 12, 77},
			},
			[]float64{.4, .4},
			area.Right,
			[]area.Rectangle{{2, 0, 6, 77}, {6, 0, 6, 77}, {6, 0, 12, 77}},
			nil,
		},
		{
			"split left thrice",
			graph.Properties{
				"rect": area.Rectangle{2, 0, 12, 77},
			},
			[]float64{.2, .4, .6},
			area.Left,
			[]area.Rectangle{{10, 0, 12, 77}, {8, 0, 10, 77}, {6, 0, 8, 77}, {2, 0, 6, 77}},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			g := graph.New(nil)
			fillNodes(g, c.base, len(c.out))
			if err := area.Split(g, graph.NodeIndex{}, g.Children(graph.NodeIndex{}), c.at, c.direction); err != nil && c.err == nil {
				t.Error("unexpected error:", err)
			} else if err == nil && c.err != nil {
				t.Error("expected error but none ocurred")
			} else if !errors.Is(err, c.err) {
				t.Error("wrong type of error")
			} else if err == nil {
				children := g.Children(graph.NodeIndex{})
				for i, out := range c.out {

					actual := (*area.AreaNode)(g.Node(children[i])).GetRect()
					if !reflect.DeepEqual(out, actual) {
						t.Errorf("rect %v, expect %v, actual %v", i, out, actual)
					}
				}
			}
		})
	}
}
