package graph

type graph struct {
	// parent *graph

	nodes     [][]*Node
	edges     map[EdgeIndex]*Edge
	edgeNodes map[EdgeIndex]*edgeNodes
}
type edgeNodes struct {
	Nodes [2][]NodeIndex
}

func New() Graph {
	return &graph{
		nodes:     [][]*Node{{{Properties: Properties{}}}},
		edges:     map[EdgeIndex]*Edge{},
		edgeNodes: map[EdgeIndex]*edgeNodes{},
	}
}

func (g *graph) Node(nidx NodeIndex) *Node {
	if len(g.nodes) <= nidx[0] {
		return nil
	} else if len(g.nodes[nidx[0]]) <= nidx[1] {
		return nil
	} else {
		return g.nodes[nidx[0]][nidx[1]]
	}
}

func (g *graph) Edge(eidx EdgeIndex) *Edge {
	return g.edges[eidx]
}

func (g *graph) Nodes(eidx EdgeIndex) [2][]NodeIndex {
	return g.edgeNodes[eidx].Nodes
}

func (g *graph) Add(parent NodeIndex, edges []EdgeIndex) NodeIndex {
	if len(g.nodes) <= parent[0]+1 {
		g.nodes = append(g.nodes, []*Node{})
	}

	nidx := NodeIndex{parent[0] + 1, len(g.nodes[parent[0]+1])}
	node := &Node{
		Properties: Properties{},
		Parent:     parent,
	}
	node.Edges = append(node.Edges, edges...)

	g.nodes[parent[0]+1] = append(g.nodes[parent[0]+1], node)

	g.Node(parent).Children = append(g.Node(parent).Children, nidx)

	for _, eidx := range edges {
		s := g.findNodeInEdges(parent, eidx)
		// edge := g.edges[eidx]
		// edge.Nodes[s] = append(edge.Nodes[s], nidx)
		g.edgeNodes[eidx].Nodes[s] = append(g.edgeNodes[eidx].Nodes[s], nidx)
	}

	return nidx
}

func (g *graph) findNodeInEdges(nidx NodeIndex, eidx EdgeIndex) int {
	for _, idx := range g.edgeNodes[eidx].Nodes[0] {
		if nidx == idx {
			return 0
		}
	}
	return 1
}

func (g *graph) Link(a, b NodeIndex) EdgeIndex {
	eidx := EdgeIndex(len(g.edges))

	edge := &Edge{
		Properties: Properties{},
	}
	g.edges[eidx] = edge

	g.edgeNodes[eidx] = &edgeNodes{
		Nodes: [2][]NodeIndex{{a}, {b}},
	}

	g.Node(a).Edges = append(g.Node(a).Edges, eidx)
	g.Node(b).Edges = append(g.Node(b).Edges, eidx)

	return EdgeIndex(eidx)
}
