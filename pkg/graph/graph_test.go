package graph_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"

	"github.com/nilsbu/arch/pkg/graph"
)

func checkNode(t *testing.T, nidx graph.NodeIndex, expect, actual *graph.Node) {
	if actual == nil {
		t.Errorf("node %v doesn't exist", nidx)
		return
	}
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("node %v doesn't match\nexpect: %v\nactual: %v", nidx, expect, actual)
	}
}

func checkChildren(t *testing.T, nidx graph.NodeIndex, expect, actual []graph.NodeIndex) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("children %v don't match\nexpect: %v\nactual: %v", nidx, expect, actual)
	}
}

func checkEdge(t *testing.T, eidx graph.EdgeIndex, expect, actual *graph.Edge) {
	if actual == nil {
		t.Errorf("edge %v doesn't exist", eidx)
		return
	}
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("edge %v doesn't match\nexpect: %v\nactual: %v", eidx, expect, actual)
	}
}

func checkEdgeNodes(t *testing.T, eidx graph.EdgeIndex, expect, actual [2][]graph.NodeIndex) {
	if !reflect.DeepEqual(expect, actual) {
		t.Errorf("edge's nodes %v don't match\nexpect: %v\nactual: %v", eidx, expect, actual)
	}
}

func TestGraph(t *testing.T) {
	for _, c := range []struct {
		name      string
		setup     func() *graph.Graph
		nodes     map[graph.NodeIndex]*graph.Node
		children  map[graph.NodeIndex][]graph.NodeIndex
		edges     map[graph.EdgeIndex]*graph.Edge
		edgeNodes map[graph.EdgeIndex][2][]graph.NodeIndex
	}{
		{
			"only root",
			func() *graph.Graph {
				return graph.New(nil)
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent}},
			map[graph.NodeIndex][]graph.NodeIndex{},
			map[graph.EdgeIndex]*graph.Edge{},
			map[graph.EdgeIndex][2][]graph.NodeIndex{},
		},
		{
			"root has property",
			func() *graph.Graph {
				g := graph.New(nil)
				g.Node(graph.NodeIndex{}).Properties["A"] = "asdf"
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{"A": "asdf"}, Parent: graph.NoParent}},
			map[graph.NodeIndex][]graph.NodeIndex{},
			map[graph.EdgeIndex]*graph.Edge{},
			map[graph.EdgeIndex][2][]graph.NodeIndex{},
		},
		{
			"node with child",
			func() *graph.Graph {
				g := graph.New(nil)
				nidx, _ := g.Add(graph.NodeIndex{}, nil)
				g.Node(nidx).Properties["A"] = "asdf"
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{"A": "asdf"}, Parent: graph.NodeIndex{0, 0}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}},
			},
			map[graph.EdgeIndex]*graph.Edge{},
			map[graph.EdgeIndex][2][]graph.NodeIndex{},
		},
		{
			"link 2 children",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{}, nil)
				n1, _ := g.Add(graph.NodeIndex{}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties = graph.Properties{"X": 123}
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}, {1, 1}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{"X": 123}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}}},
			},
		},
		{
			"inherit edge",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{}, nil)
				n1, _ := g.Add(graph.NodeIndex{}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Add(n1, []graph.EdgeIndex{eidx})
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{2, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{1, 1}, Edges: []graph.EdgeIndex{0}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}, {1, 1}},
				{1, 1}: {{2, 0}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}, {2, 0}}},
			},
		},
		{
			"append empty graph",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{}, nil)
				n1, _ := g.Add(graph.NodeIndex{}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Add(n1, []graph.EdgeIndex{eidx})
				return graph.New(g)
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{2, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{1, 1}, Edges: []graph.EdgeIndex{0}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}, {1, 1}},
				{1, 1}: {{2, 0}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}, {2, 0}}},
			},
		},
		{
			"link inherited children",
			func() *graph.Graph {
				g := graph.New(graph.New(nil))
				n0, _ := g.Add(graph.NodeIndex{}, nil)
				n1, _ := g.Add(graph.NodeIndex{}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Edge(eidx).Properties = graph.Properties{"X": 123}
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}, {1, 1}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{"X": 123}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}}},
			},
		},
		{
			"link after inheritance",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{}, nil)
				n1, _ := g.Add(graph.NodeIndex{}, nil)
				eidx, _ := g.Link(n0, n1)
				g = graph.New(g)
				g.Add(n1, []graph.EdgeIndex{eidx})
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{2, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{1, 1}, Edges: []graph.EdgeIndex{0}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}, {1, 1}},
				{1, 1}: {{2, 0}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}, {2, 0}}},
			},
		},
		{
			"changes in later graphs don't affect the earlier ones",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{}, nil)
				n1, _ := g.Add(graph.NodeIndex{}, nil)
				eidx, _ := g.Link(n0, n1)
				g.Add(n1, []graph.EdgeIndex{eidx})
				g2 := graph.New(g)
				g2.Add(graph.NodeIndex{}, nil)
				n3, _ := g2.Add(n0, []graph.EdgeIndex{eidx})
				n4, _ := g2.Add(n1, nil)
				g2.Link(n3, n4)
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{2, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{1, 1}, Edges: []graph.EdgeIndex{0}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}, {1, 1}},
				{1, 1}: {{2, 0}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}, {2, 0}}},
			},
		},
		{
			"edge index continuity",
			func() *graph.Graph {
				g := graph.New(nil)
				n0, _ := g.Add(graph.NodeIndex{}, nil)
				n1, _ := g.Add(graph.NodeIndex{}, nil)
				g.Link(n0, n1)
				g = graph.New(g)
				n00, _ := g.Add(n0, nil)
				n01, _ := g.Add(n0, nil)
				g.Link(n00, n01)
				g = graph.New(g)
				n010, _ := g.Add(n01, nil)
				n011, _ := g.Add(n01, nil)
				g.Link(n010, n011)
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Parent: graph.NoParent},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{2, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{1, 0}, Edges: []graph.EdgeIndex{1}},
				{2, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{1, 0}, Edges: []graph.EdgeIndex{1}},
				{3, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{2, 1}, Edges: []graph.EdgeIndex{2}},
				{3, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{2, 1}, Edges: []graph.EdgeIndex{2}},
			},
			map[graph.NodeIndex][]graph.NodeIndex{
				{0, 0}: {{1, 0}, {1, 1}},
				{1, 0}: {{2, 0}, {2, 1}},
				{2, 1}: {{3, 0}, {3, 1}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{}},
				1: {Properties: graph.Properties{}},
				2: {Properties: graph.Properties{}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}}},
				1: {{{2, 0}}, {{2, 1}}},
				2: {{{3, 0}}, {{3, 1}}},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			g := c.setup()

			for nidx, node := range c.nodes {
				checkNode(t, nidx, node, g.Node(nidx))
			}

			for nidx, children := range c.children {
				checkChildren(t, nidx, children, g.Children(nidx))
			}

			for eidx, edge := range c.edges {
				checkEdge(t, eidx, edge, g.Edge(eidx))
			}

			for eidx, en := range c.edgeNodes {
				checkEdgeNodes(t, eidx, en, g.Nodes(eidx))
			}
		})
	}
}

func TestGetNonExistingNode(t *testing.T) {
	g := graph.New(nil)
	if g.Node(graph.NodeIndex{1, 0}) != nil {
		t.Error("no node should be returned")
	}
}

func TestAddToFalseParent(t *testing.T) {
	g := graph.New(nil)
	if _, err := g.Add(graph.NodeIndex{1, 0}, nil); err == nil {
		t.Error("must return error")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		t.Error("error must be 'ErrIllegalAction'")
	}
}

func TestAddWithNonexistentEdge(t *testing.T) {
	g := graph.New(nil)
	if _, err := g.Add(graph.NodeIndex{0, 0}, []graph.EdgeIndex{1}); err == nil {
		t.Error("must return error")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		t.Error("error must be 'ErrIllegalAction'")
	}
}

func TestAddWithFalseEdge(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	n1, _ := g.Add(graph.NodeIndex{}, nil)
	eidx, _ := g.Link(n0, n1)
	if _, err := g.Add(graph.NodeIndex{0, 0}, []graph.EdgeIndex{eidx}); err == nil {
		t.Error("must return error")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		t.Error("error must be 'ErrIllegalAction'")
	}
}

func TestInheritEdgeTwice(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	n1, _ := g.Add(graph.NodeIndex{}, nil)
	eidx, _ := g.Link(n0, n1)
	g.Add(n1, []graph.EdgeIndex{eidx})
	if _, err := g.Add(n1, []graph.EdgeIndex{eidx}); err == nil {
		t.Error("must return error")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		t.Error("error must be 'ErrIllegalAction'")
	}
}

func TestAddWithoutParent(t *testing.T) {
	g := graph.New(nil)
	if _, err := g.Add(graph.NoParent, nil); err == nil {
		t.Error("must return error")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		t.Error("error must be 'ErrIllegalAction'")
	}
}

func TestLinkeNodesTwice(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	n1, _ := g.Add(graph.NodeIndex{}, nil)
	g.Link(n0, n1)
	if _, err := g.Link(n0, n1); err == nil {
		t.Error("must return error")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		t.Error("error must be 'ErrIllegalAction'")
	}
}

func TestLinkANonExistent(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	if _, err := g.Link(n0, graph.NodeIndex{1, 1}); err == nil {
		t.Error("expected error but none ocurred")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		fmt.Println(err)
		t.Error("link error must be an 'ErrIllegalAction'")
	}
}

func TestLinkBNonExistent(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	if _, err := g.Link(graph.NodeIndex{1, 1}, n0); err == nil {
		t.Error("expected error but none ocurred")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		fmt.Println(err)
		t.Error("link error must be an 'ErrIllegalAction'")
	}
}

func TestLinkWithOwnParent(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	if _, err := g.Link(graph.NodeIndex{}, n0); err == nil {
		t.Error("expected error but none ocurred")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		fmt.Println(err)
		t.Error("link error must be an 'ErrIllegalAction'")
	}
}

func TestLinkAfterInheritance(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	n1, _ := g.Add(graph.NodeIndex{}, nil)
	g = graph.New(g)

	if _, err := g.Link(n0, n1); err == nil {
		t.Error("expected error but none ocurred")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		fmt.Println(err)
		t.Error("link error must be an 'ErrIllegalAction'")
	}
}

func TestOneLinkedNodeBeforeInheritance(t *testing.T) {
	g := graph.New(nil)
	n0, _ := g.Add(graph.NodeIndex{}, nil)
	g = graph.New(g)
	n1, _ := g.Add(graph.NodeIndex{}, nil)
	if _, err := g.Link(n0, n1); err == nil {
		t.Error("expected error but none ocurred")
	} else if !errors.Is(err, graph.ErrIllegalAction) {
		t.Error("link error must be an 'ErrIllegalAction'")
	}
}
