package graph_test

import (
	"testing"

	"github.com/nilsbu/arch/pkg/graph"
	tg "github.com/nilsbu/arch/test/graph"
)

func TestCompare(t *testing.T) {
	for _, c := range []struct {
		name        string
		a, b        func() *graph.Graph
		equal       bool
		explanation string
	}{
		{
			"both are nil",
			func() *graph.Graph { return nil },
			func() *graph.Graph { return nil },
			true, "",
		},
		{
			"a is nil",
			func() *graph.Graph { return nil },
			func() *graph.Graph { return graph.New(nil) },
			false, "first is nil",
		},
		{
			"b is nil",
			func() *graph.Graph { return graph.New(nil) },
			func() *graph.Graph { return nil },
			false, "second is nil",
		},
		{
			"both only root",
			func() *graph.Graph { return graph.New(nil) },
			func() *graph.Graph { return graph.New(nil) },
			true, "",
		},
		{
			"a has additional child",
			func() *graph.Graph {
				g := graph.New(nil)
				g.Add(graph.NodeIndex{0, 0}, nil)
				return g
			},
			func() *graph.Graph { return graph.New(nil) },
			false, "children of [0 0] are disjunct in [[1 0]] vs. []",
		},
		{
			"b has additional child",
			func() *graph.Graph { return graph.New(nil) },
			func() *graph.Graph {
				g := graph.New(nil)
				g.Add(graph.NodeIndex{0, 0}, nil)
				return g
			},
			false, "children of [0 0] are disjunct in [] vs. [[1 0]]",
		},
		{
			"different properties",
			func() *graph.Graph {
				g := graph.New(nil)
				nidx, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				g.Node(nidx).Properties["A"] = 3
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				nidx, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				g.Node(nidx).Properties["A"] = 2
				return g
			},
			false, "properties of [1 0] are different: map[A:3] vs. map[A:2]",
		},
		{
			"same link",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				g.Link(n0, n1)
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				g.Link(n0, n1)
				return g
			},
			true, "",
		},
		{
			"only first has link",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				g.Link(n0, n1)
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				g.Add(graph.NodeIndex{0, 0}, nil)
				g.Add(graph.NodeIndex{0, 0}, nil)
				return g
			},
			false, "edges of [1 0] are disjunct in [0] vs. []",
		},
		{
			"link in different direction",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				g.Link(n0, n1)
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				g.Link(n1, n0)
				return g
			},
			false, "nodes of edge 0 don't match: [[[1 0]] [[1 1]]] vs. [[[1 1]] [[1 0]]]",
		},
		{
			"different edge properties",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties["A"] = 123
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties["A"] = 124
				return g
			},
			false, "properties of edge 0 don't match: map[A:123] vs. map[A:124]",
		},
		{
			"children added in different graph instance",
			func() *graph.Graph {
				g := graph.New(nil)
				g = graph.New(g)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties["A"] = 123
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				n1, _ := g.Add(graph.NodeIndex{0, 0}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties["A"] = 123
				return g
			},
			true, "",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			equal, explanation := tg.AreEqual(c.a(), c.b())
			if equal && !c.equal {
				t.Error("graphs aren't equal but true was returned")
			} else if !equal && c.equal {
				t.Error("graphs are equal but false was returned")
			}

			if c.explanation != explanation {
				t.Errorf("explanations doesn't match:\nexpect: '%v'\nactual: '%v'",
					c.explanation, explanation)
			}
		})
	}
}
