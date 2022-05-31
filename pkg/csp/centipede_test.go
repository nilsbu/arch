package csp_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nilsbu/arch/pkg/csp"
	"github.com/nilsbu/arch/pkg/graph"
)

func TestCentipedeMatch(t *testing.T) {
	for _, c := range []struct {
		name    string
		graphs  []func() *graph.Graph
		ok      bool
		matches []graph.NodeIndex
		err     error
	}{
		{
			"zero graphs match",
			[]func() *graph.Graph{},
			true, nil, nil,
		},
		{
			"one graph always matches",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n0, _ := g.Add(graph.NodeIndex{})
					n1, _ := g.Add(graph.NodeIndex{})
					g.Link(n0, n1)
					return g
				},
			},
			true, nil, nil,
		},
		{
			"when parent or child match, parent is chosen",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					g.Add(graph.NodeIndex{})
					return g
				},
				func() *graph.Graph {
					leaves, _ := graph.New(nil).Leaves()
					return leaves
				},
			},
			true, []graph.NodeIndex{{0, 0}}, nil,
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
					g, _ = g.Leaves()
					return g
				},
			},
			false, nil, nil,
		},
		{
			"matching doesn't work when first is nameless",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "b"
					g, _ = g.Leaves()
					return g
				},
			},
			false, nil, nil,
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
					g, _ = g.Leaves()
					return g
				},
			},
			true, []graph.NodeIndex{{0, 0}}, nil,
		},
		{
			"matching names",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "a"
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "a"
					g, _ = g.Leaves()
					return g
				},
			},
			true, []graph.NodeIndex{{0, 0}}, nil,
		},
		{
			"names don't match",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "b"
					g.Node(graph.NodeIndex{}).Properties["names"] = []string{"y", "x"}
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					g.Node(graph.NodeIndex{}).Properties["name"] = "a"
					g, _ = g.Leaves()
					return g
				},
			},
			false, nil, nil,
		},
		{
			"more options than needed",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					nidx, _ := g.Add(graph.NodeIndex{})
					g.Node(nidx).Properties["name"] = "a"
					nidx, _ = g.Add(graph.NodeIndex{})
					g.Node(nidx).Properties["name"] = "a"
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					nidx, _ := g.Add(graph.NodeIndex{})
					g.Node(nidx).Properties["name"] = "a"
					return g
				},
			},
			true, []graph.NodeIndex{{1, 0}}, nil,
		},
		{
			"values are unique",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					nidx, _ := g.Add(graph.NodeIndex{})
					g.Node(nidx).Properties["name"] = "a"
					nidx, _ = g.Add(graph.NodeIndex{})
					g.Node(nidx).Properties["name"] = "a"
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					g.Add(graph.NodeIndex{})
					g.Add(graph.NodeIndex{})
					return g
				},
			},
			true, []graph.NodeIndex{{1, 0}, {1, 1}}, nil,
		},
		{
			"can't match because adjacency is invalid",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{})
					n11, _ := g.Add(graph.NodeIndex{})
					n12, _ := g.Add(graph.NodeIndex{})
					g.Node(n10).Properties["name"] = "a"
					g.Node(n11).Properties["name"] = "b"
					g.Node(n12).Properties["name"] = "c"
					g.Link(n10, n11)
					g.Link(n11, n12)
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{})
					n11, _ := g.Add(graph.NodeIndex{})
					g.Node(n10).Properties["name"] = "a"
					g.Node(n11).Properties["name"] = "c"
					g.Link(n10, n11)
					return g
				},
			},
			false, nil, nil,
		},
		{
			"complex valid match",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{})
					n11, _ := g.Add(graph.NodeIndex{})
					n12, _ := g.Add(graph.NodeIndex{})
					g.Add(graph.NodeIndex{})
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
					n10, _ := g.Add(graph.NodeIndex{})
					n11, _ := g.Add(graph.NodeIndex{})
					n12, _ := g.Add(graph.NodeIndex{})
					g.Node(n10).Properties["name"] = "b"
					g.Node(n11).Properties["name"] = "c"
					g.Node(n12).Properties["name"] = "a"
					g.Link(n10, n11)
					g.Link(n11, n12)
					g.Link(n12, n10)
					return g
				},
			},
			true, []graph.NodeIndex{{1, 1}, {1, 2}, {1, 0}}, nil,
		},
		{
			"never match both parent and child",
			[]func() *graph.Graph{
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{})
					g.Node(n10).Properties["name"] = "a"
					n20, _ := g.Add(n10)
					g.Node(n20).Properties["name"] = "c"
					n11, _ := g.Add(graph.NodeIndex{})
					g.Node(n11).Properties["name"] = "b"
					g.Link(n10, n11)
					return g
				},
				func() *graph.Graph {
					g := graph.New(nil)
					n10, _ := g.Add(graph.NodeIndex{})
					n11, _ := g.Add(graph.NodeIndex{})
					n12, _ := g.Add(graph.NodeIndex{})
					g.Node(n10).Properties["name"] = "b"
					g.Node(n11).Properties["name"] = "c"
					g.Node(n12).Properties["name"] = "a"
					return g
				},
			},
			false, nil, nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			graphs := make([]*graph.Graph, len(c.graphs))
			for i, f := range c.graphs {
				graphs[i] = f()
			}

			if ok, matches, err := (&csp.Centipede{}).Match(graphs); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error\nexpect: %v\nactual: %v", c.err, err)
			} else if c.ok && !ok {
				t.Error("should have matched")
			} else if !c.ok && ok {
				t.Error("shouldn't have matched")
			} else if ok {
				if !reflect.DeepEqual(c.matches, matches) {
					t.Errorf("matches don't match:\nexpect: %v\nactual: %v",
						c.matches, matches)
				}
			}
		})
	}
}

func TestResetNodes(t *testing.T) {
	c := &csp.Centipede{}
	g1 := graph.New(nil)
	g1.Add(graph.NodeIndex{})
	g1.Add(graph.NodeIndex{})

	g2 := graph.New(nil)
	g2.Add(graph.NodeIndex{})

	c.Match([]*graph.Graph{graph.New(nil), g1})
	c.Match([]*graph.Graph{graph.New(nil), g2})
	// nothing to check, if it doesn't crash, reset worked
}
