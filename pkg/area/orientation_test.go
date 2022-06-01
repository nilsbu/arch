package area_test

import (
	"errors"
	"reflect"
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

func TestDifferenct(t *testing.T) {
	for _, c := range []struct {
		name string
		a, b area.Direction
		diff int
	}{
		{
			"both up",
			area.Up,
			area.Up,
			0,
		},
		{
			"both left",
			area.Left,
			area.Left,
			0,
		},
		{
			"up & right",
			area.Up,
			area.Right,
			-90,
		},
		{
			"down & up",
			area.Down,
			area.Up,
			180,
		},
		{
			"left & right",
			area.Left,
			area.Right,
			180,
		},
		{
			"right & left",
			area.Right,
			area.Left,
			180,
		},
		{
			"up & left",
			area.Up,
			area.Left,
			90,
		},
		{
			"left & up",
			area.Left,
			area.Up,
			-90,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			diff := area.Difference(c.a, c.b)
			if c.diff != diff {
				t.Errorf("expected diff to be %v but was %v", c.diff, diff)
			}
		})
	}
}

func TestRotateWithin(t *testing.T) {
	for _, c := range []struct {
		name     string
		rect, in area.Rectangle
		from, to area.Direction
		anchor   area.Anchor
		result   area.Rectangle
		err      error
	}{
		{
			"size one, don't rotate",
			area.Rectangle{0, 0, 1, 1},
			area.Rectangle{0, 0, 10, 10},
			area.Down,
			area.Down,
			area.FarLeft,
			area.Rectangle{0, 0, 1, 1},
			nil,
		},
		{
			"size one, don't rotate, orientation up",
			area.Rectangle{0, 0, 1, 1},
			area.Rectangle{0, 0, 10, 10},
			area.Up,
			area.Up,
			area.FarLeft,
			area.Rectangle{0, 0, 1, 1},
			nil,
		},
		{
			"size one, rotated 180 deg",
			area.Rectangle{0, 0, 1, 1},
			area.Rectangle{0, 0, 10, 10},
			area.Up,
			area.Down,
			area.FarLeft,
			area.Rectangle{9, 9, 10, 10},
			nil,
		},
		{
			"size 2x1, rotated 90 deg",
			area.Rectangle{0, 0, 1, 2},
			area.Rectangle{0, 0, 10, 10},
			area.Up,
			area.Right,
			area.FarLeft,
			area.Rectangle{8, 0, 10, 1},
			nil,
		},
		{
			"size 2x1, rotated -90 deg, non-rectangular room",
			area.Rectangle{17, 7, 18, 9},
			area.Rectangle{0, 0, 20, 10},
			area.Right,
			area.Up,
			area.FarRight,
			area.Rectangle{17, 2, 19, 3},
			nil,
		},
		{
			"size 3x1, rotate 90 around center",
			area.Rectangle{8, 4, 9, 6},
			area.Rectangle{0, 0, 20, 10},
			area.Right,
			area.Left,
			area.Center,
			area.Rectangle{11, 4, 12, 6},
			nil,
		},
		{
			"non-zero top-left",
			area.Rectangle{2, 3, 3, 4},
			area.Rectangle{2, 2, 10, 10},
			area.Left,
			area.Up,
			area.NearLeft,
			area.Rectangle{8, 2, 9, 3},
			nil,
		},
		{
			"near right as anchor",
			area.Rectangle{9, 3, 10, 7},
			area.Rectangle{0, 2, 10, 10},
			area.Left,
			area.Down,
			area.NearRight,
			area.Rectangle{1, 2, 5, 3},
			nil,
		},
		{
			"doesn't fit after rotation",
			area.Rectangle{2, 3, 3, 7},
			area.Rectangle{2, 0, 5, 10},
			area.Left,
			area.Down,
			area.Center,
			area.Rectangle{2, 6, 6, 7},
			area.ErrInvalidRotation,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if result, err := area.RotateWithin(c.rect, c.in, c.from, c.to, c.anchor); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error\nexpect: %v\nactual: %v", c.err, err)
			} else if err == nil {
				if !reflect.DeepEqual(c.result, result) {
					t.Errorf("expected %v but got %v", c.result, result)
				}
			}
		})
	}
}

func TestCalcAnchorPoint(t *testing.T) {
	for _, c := range []struct {
		name      string
		rect      area.Rectangle
		anchor    area.Anchor
		direction area.Direction
		point     area.Point
	}{
		{
			"down, near-right",
			area.Rectangle{3, 3, 20, 24},
			area.NearRight,
			area.Down,
			area.Point{3, 3},
		},
		{
			"right, far-left",
			area.Rectangle{3, 3, 20, 24},
			area.FarLeft,
			area.Right,
			area.Point{20, 3},
		},
		{
			"center, up",
			area.Rectangle{10, 10, 20, 24},
			area.Center,
			area.Up,
			area.Point{15, 17},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			point := area.CalcAnchorPoint(c.rect, c.anchor, c.direction)
			if c.point != point {
				t.Errorf("expected %v but got %v", c.point, point)
			}
		})
	}
}
