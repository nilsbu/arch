package graph_test

import (
	"errors"
	"testing"

	"github.com/nilsbu/arch/pkg/graph"
	tg "github.com/nilsbu/arch/test/graph"
)

func TestLeaves(t *testing.T) {
	for _, c := range []struct {
		name   string
		graph  func() *graph.Graph
		expect func() *graph.Graph
		err    error
	}{
		{
			"root only",
			func() *graph.Graph { return graph.New(nil) },
			func() *graph.Graph {
				g := graph.New(nil)
				g.Add(graph.NodeIndex{})
				return g
			},
			nil,
		},
		{
			"single child with property",
			func() *graph.Graph {
				g := graph.New(nil)
				nidx, _ := g.Add(graph.NodeIndex{})
				g.Node(nidx).Properties["a"] = "b"
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				nidx, _ := g.Add(graph.NodeIndex{})
				g.Node(nidx).Properties["a"] = "b"
				return g
			},
			nil,
		},
		{
			"two linked nodes",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				n1, _ := g.Add(graph.NodeIndex{})
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties["a"] = "b"
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				n1, _ := g.Add(graph.NodeIndex{})
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties["a"] = "b"
				return g
			},
			nil,
		},
		{
			"not fully inherited edge",
			func() *graph.Graph {
				g := graph.New(nil)
				n10, _ := g.Add(graph.NodeIndex{})
				n11, _ := g.Add(graph.NodeIndex{})
				e0, _ := g.Link(n10, n11)
				n20, _ := g.Add(n10)
				n21, _ := g.Add(n10)
				g.InheritEdge(n10, n21, []graph.EdgeIndex{e0})
				g.Link(n20, n21)
				g.Add(n21)
				return g
			},
			func() *graph.Graph {
				return nil
			},
			graph.ErrNotLeafable,
		},
		{
			"node gets inherited",
			func() *graph.Graph {
				g := graph.New(nil)
				n10, _ := g.Add(graph.NodeIndex{})
				n11, _ := g.Add(graph.NodeIndex{})
				e0, _ := g.Link(n10, n11)
				g.Edge(e0).Properties["a"] = "b"
				n20, _ := g.Add(n10)
				g.Node(n20).Properties["n"] = 0
				n21, _ := g.Add(n10)
				g.InheritEdge(n10, n21, []graph.EdgeIndex{e0})
				g.Node(n21).Properties["n"] = 1
				e1, _ := g.Link(n20, n21)
				g.Edge(e1).Properties["a"] = "c"
				n22, _ := g.Add(n11)
				g.InheritEdge(n11, n22, []graph.EdgeIndex{e0})
				g.Node(n22).Properties["n"] = 2
				return g
			},
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{})
				n1, _ := g.Add(graph.NodeIndex{})
				n2, _ := g.Add(graph.NodeIndex{})
				g.Node(n0).Properties["n"] = 0
				g.Node(n1).Properties["n"] = 1
				g.Node(n2).Properties["n"] = 2
				e0, _ := g.Link(n1, n2)
				g.Edge(e0).Properties["a"] = "b"
				e1, _ := g.Link(n0, n1)
				g.Edge(e1).Properties["a"] = "c"
				return g
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if leaves, err := c.graph().Leaves(); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error\nexpect: %v\nactual: %v", c.err, err)
			} else if err == nil {
				if eq, explain := tg.AreEqual(c.expect(), leaves); !eq {
					t.Errorf("graphs don't match: %v", explain)
				}
			}
		})
	}
}
