package area_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/graph"
)

func TestSplit(t *testing.T) {
	for _, c := range []struct {
		name      string
		base      *area.AreaNode
		into      []*area.AreaNode
		at        []float64
		direction area.Direction
		out       []area.Rectangle
		err       error
	}{
		{
			"len(into) != len(at) + 1",
			&area.AreaNode{Properties: graph.Properties{
				"rect": area.Rectangle{2, 3, 4, 5},
			}},
			[]*area.AreaNode{{Properties: graph.Properties{}}},
			[]float64{0.2},
			area.Up,
			[]area.Rectangle{{2, 3, 4, 5}},
			area.ErrInvalidSplit,
		},
		// TODO with invalid Direction
		// TODO with at <0 or >1
		{
			"one split copies rect",
			&area.AreaNode{Properties: graph.Properties{
				"rect": area.Rectangle{2, 3, 4, 5},
			}},
			[]*area.AreaNode{{Properties: graph.Properties{}}},
			[]float64{},
			area.Up,
			[]area.Rectangle{{2, 3, 4, 5}},
			nil,
		},
		{
			"split upwards in the middle",
			&area.AreaNode{Properties: graph.Properties{
				"rect": area.Rectangle{2, 0, 4, 6},
			}},
			[]*area.AreaNode{
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
			},
			[]float64{.5},
			area.Up,
			[]area.Rectangle{{2, 3, 4, 6}, {2, 0, 4, 3}},
			nil,
		},
		{
			"split downwards twice",
			&area.AreaNode{Properties: graph.Properties{
				"rect": area.Rectangle{2, 0, 4, 10},
			}},
			[]*area.AreaNode{
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
			},
			[]float64{.3, .8},
			area.Down,
			[]area.Rectangle{{2, 0, 4, 3}, {2, 3, 4, 8}, {2, 8, 4, 10}},
			nil,
		},
		{
			"split right with one of size zero",
			&area.AreaNode{Properties: graph.Properties{
				"rect": area.Rectangle{2, 0, 12, 77},
			}},
			[]*area.AreaNode{
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
			},
			[]float64{.4, .4},
			area.Right,
			[]area.Rectangle{{2, 0, 6, 77}, {6, 0, 6, 77}, {6, 0, 12, 77}},
			nil,
		},
		{
			"split left thrice",
			&area.AreaNode{Properties: graph.Properties{
				"rect": area.Rectangle{2, 0, 12, 77},
			}},
			[]*area.AreaNode{
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
				{Properties: graph.Properties{}},
			},
			[]float64{.2, .4, .6},
			area.Left,
			[]area.Rectangle{{10, 0, 12, 77}, {8, 0, 10, 77}, {6, 0, 8, 77}, {2, 0, 6, 77}},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if err := area.Split(c.base, c.into, c.at, c.direction); err != nil && c.err == nil {
				t.Error("unexpected error:", err)
			} else if err == nil && c.err != nil {
				t.Error("expected error but none ocurred")
			} else if !errors.Is(err, c.err) {
				t.Error("wrong type of error")
			} else if err == nil {
				if len(c.out) != len(c.into) {
					t.Fatalf("expected rects = %v, received = %v", len(c.out), len(c.into))
				} else {
					for i, out := range c.out {
						actual := c.into[i].GetRect()
						if !reflect.DeepEqual(out, actual) {
							t.Errorf("rect %v, expect %v, actual %v", i, out, actual)
						}
					}

				}
			}
		})
	}
}
