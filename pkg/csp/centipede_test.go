package csp_test

import (
	"errors"
	"testing"

	"github.com/nilsbu/arch/pkg/csp"
	"github.com/nilsbu/arch/pkg/graph"
)

func TestCentipedeMatch(t *testing.T) {
	for _, c := range []struct {
		name   string
		graphs []func() *graph.Graph
		ok     bool
		err    error
	}{
		{
			"zero graphs match",
			[]func() *graph.Graph{},
			true, nil,
		},
		{
			"one graph always matches",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n0, _ := g.Add(graph.NodeIndex{}, nil)
					n1, _ := g.Add(graph.NodeIndex{}, nil)
					g.Link(n0, n1)
					return g
				},
			},
			true, nil,
		},
		{
			"both root only",
			[]func() *graph.Graph{
				func() *graph.Graph { return graph.New(nil) },
				func() *graph.Graph { return graph.New(nil) },
			},
			true, nil,
		},
		{
			"non-leafable graph",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{}, nil)
					n11, _ := g.Add(graph.NodeIndex{}, nil)
					g.Link(n10, n11)
					g.Add(n10, nil)
					return g
				},
				func() *graph.Graph { return graph.New(nil) },
			},
			false, graph.ErrNotLeafable,
		},
		{
			"one has a child",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					g.Add(graph.NodeIndex{}, nil)
					return g
				},
				func() *graph.Graph { return graph.New(nil) },
			},
			true, nil,
		},
		{
			"incompatible names",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "a"
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "b"
					return g
				},
			},
			false, nil,
		},
		{
			"matching works when first is nameless",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "b"
					return g
				},
			},
			true, nil,
		},
		{
			"matching works when second is nameless",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "b"
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					return g
				},
			},
			true, nil,
		},
		{
			"more options than needed",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					g.Add(graph.NodeIndex{}, nil)
					g.Add(graph.NodeIndex{}, nil)
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					g.Add(graph.NodeIndex{}, nil)
					return g
				},
			},
			true, nil,
		},
		{
			"can't match because adjacency is invalid",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{}, nil)
					n11, _ := g.Add(graph.NodeIndex{}, nil)
					n12, _ := g.Add(graph.NodeIndex{}, nil)
					g.Node(n10).Properties["name"] = "a"
					g.Node(n11).Properties["name"] = "b"
					g.Node(n12).Properties["name"] = "c"
					g.Link(n10, n11)
					g.Link(n11, n12)
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{}, nil)
					n11, _ := g.Add(graph.NodeIndex{}, nil)
					g.Node(n10).Properties["name"] = "a"
					g.Node(n11).Properties["name"] = "c"
					g.Link(n10, n11)
					return g
				},
			},
			false, nil,
		},
		{
			"complex valid match",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{}, nil)
					n11, _ := g.Add(graph.NodeIndex{}, nil)
					n12, _ := g.Add(graph.NodeIndex{}, nil)
					g.Add(graph.NodeIndex{}, nil)
					g.Node(n10).Properties["name"] = "a"
					g.Node(n11).Properties["name"] = "b"
					g.Node(n12).Properties["name"] = "c"
					g.Link(n10, n11)
					g.Link(n11, n12)
					g.Link(n12, n10)
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{}, nil)
					n11, _ := g.Add(graph.NodeIndex{}, nil)
					n12, _ := g.Add(graph.NodeIndex{}, nil)
					g.Node(n10).Properties["name"] = "a"
					g.Node(n11).Properties["name"] = "b"
					g.Node(n12).Properties["name"] = "c"
					g.Link(n10, n11)
					g.Link(n11, n12)
					g.Link(n12, n10)
					return g
				},
			},
			true, nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			graphs := make([]*graph.Graph, len(c.graphs))
			for i, f := range c.graphs {
				graphs[i] = f()
			}

			if ok, err := (&csp.Centipede{}).Match(graphs); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error\nexpect: %v\nactual: %v", c.err, err)
			} else if c.ok && !ok {
				t.Error("should have matched")
			} else if !c.ok && ok {
				t.Error("shouldn't have matched")
			}
			// t.Fail()
		})
	}
}
