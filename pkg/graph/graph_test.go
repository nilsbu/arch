package graph_test

import (
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
		setup     func() graph.Graph
		nodes     map[graph.NodeIndex]*graph.Node
		edges     map[graph.EdgeIndex]*graph.Edge
		edgeNodes map[graph.EdgeIndex][2][]graph.NodeIndex
	}{
		{
			"only root",
			func() graph.Graph {
				return graph.New()
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}}},
			map[graph.EdgeIndex]*graph.Edge{},
			map[graph.EdgeIndex][2][]graph.NodeIndex{},
		},
		{
			"root has property",
			func() graph.Graph {
				g := graph.New()
				g.Node(graph.NodeIndex{}).Properties["A"] = "asdf"
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{"A": "asdf"}}},
			map[graph.EdgeIndex]*graph.Edge{},
			map[graph.EdgeIndex][2][]graph.NodeIndex{},
		},
		{
			"node with child",
			func() graph.Graph {
				g := graph.New()
				g.Node(g.Add(graph.NodeIndex{}, nil)).Properties["A"] = "asdf"
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Children: []graph.NodeIndex{{1, 0}}},
				{1, 0}: {Properties: graph.Properties{"A": "asdf"}, Parent: graph.NodeIndex{0, 0}},
			},
			map[graph.EdgeIndex]*graph.Edge{},
			map[graph.EdgeIndex][2][]graph.NodeIndex{},
		},
		{
			"link 2 children",
			func() graph.Graph {
				g := graph.New()
				n0, n1 := g.Add(graph.NodeIndex{}, nil), g.Add(graph.NodeIndex{}, nil)
				g.Edge(g.Link(n0, n1)).Properties = graph.Properties{"X": 123}
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Children: []graph.NodeIndex{{1, 0}, {1, 1}}},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
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
			func() graph.Graph {
				g := graph.New()
				n0, n1 := g.Add(graph.NodeIndex{}, nil), g.Add(graph.NodeIndex{}, nil)
				eidx := g.Link(n0, n1)
				g.Add(n1, []graph.EdgeIndex{eidx})
				return g
			},
			map[graph.NodeIndex]*graph.Node{
				{0, 0}: {Properties: graph.Properties{}, Children: []graph.NodeIndex{{1, 0}, {1, 1}}},
				{1, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Edges: []graph.EdgeIndex{0}},
				{1, 1}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{0, 0}, Children: []graph.NodeIndex{{2, 0}}, Edges: []graph.EdgeIndex{0}},
				{2, 0}: {Properties: graph.Properties{}, Parent: graph.NodeIndex{1, 1}, Edges: []graph.EdgeIndex{0}},
			},
			map[graph.EdgeIndex]*graph.Edge{
				0: {Properties: graph.Properties{}},
			},
			map[graph.EdgeIndex][2][]graph.NodeIndex{
				0: {{{1, 0}}, {{1, 1}, {2, 0}}},
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			g := c.setup()

			for nidx, node := range c.nodes {
				checkNode(t, nidx, node, g.Node(nidx))
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
