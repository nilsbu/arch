package graph

type graph struct {
	parent Graph

	nodes     [][]*Node
	children  map[NodeIndex][]NodeIndex
	edges     map[EdgeIndex]*Edge
	edgeNodes map[EdgeIndex]*edgeNodes
}
type edgeNodes struct {
	Nodes [2][]NodeIndex
}

func New(parent Graph) Graph {
	if parent == nil {
		return &graph{
			nodes:     [][]*Node{{{Properties: Properties{}}}},
			children:  map[NodeIndex][]NodeIndex{},
			edges:     map[EdgeIndex]*Edge{},
			edgeNodes: map[EdgeIndex]*edgeNodes{},
		}
	} else {
		return &graph{
			parent:    parent,
			nodes:     [][]*Node{},
			children:  map[NodeIndex][]NodeIndex{},
			edges:     map[EdgeIndex]*Edge{},
			edgeNodes: map[EdgeIndex]*edgeNodes{},
		}
	}
}

func (g *graph) Node(nidx NodeIndex) *Node {
	if len(g.nodes) <= nidx[0] || len(g.nodes[nidx[0]]) <= nidx[1] {
		if g.parent == nil {
			return nil
		} else {
			return g.parent.Node(nidx)
		}
	} else {
		return g.nodes[nidx[0]][nidx[1]]
	}
}

func (g *graph) Children(nidx NodeIndex) []NodeIndex {
	var children []NodeIndex
	if g.parent != nil {
		children = g.parent.Children(nidx)
	}
	if en, ok := g.children[nidx]; ok {
		children = append(children, en...)
	}
	return children
}

func (g *graph) Edge(eidx EdgeIndex) *Edge {
	if edge, ok := g.edges[eidx]; ok {
		return edge
	} else {
		return g.parent.Edge(eidx)
	}
}

func (g *graph) Nodes(eidx EdgeIndex) [2][]NodeIndex {
	var nodes [2][]NodeIndex
	if g.parent != nil {
		nodes = g.parent.Nodes(eidx)
	}
	if en, ok := g.edgeNodes[eidx]; ok {
		nodes[0] = append(nodes[0], en.Nodes[0]...)
		nodes[1] = append(nodes[1], en.Nodes[1]...)
	}
	return nodes
}

func (g *graph) Add(parent NodeIndex, edges []EdgeIndex) NodeIndex {
	nidx := NodeIndex{parent[0] + 1, g.nodesInLayer(parent[0] + 1)}
	node := g.createNodeAt(nidx)
	node.Parent = parent
	node.Edges = append(node.Edges, edges...)

	// g.Node(parent).Children = append(g.Node(parent).Children, nidx)
	if children, ok := g.children[parent]; ok {
		g.children[parent] = append(children, nidx)
	} else {
		g.children[parent] = []NodeIndex{nidx}
	}

	for _, eidx := range edges {
		s := g.findNodeInEdges(parent, eidx)
		if nodes, ok := g.edgeNodes[eidx]; ok {
			nodes.Nodes[s] = append(nodes.Nodes[s], nidx)
		} else {
			n := &edgeNodes{}
			n.Nodes[s] = []NodeIndex{nidx}
			g.edgeNodes[eidx] = n
		}
	}

	return nidx
}

func (g *graph) nodesInLayer(l int) int {
	n := 0
	if g.parent != nil {
		n = g.parent.(*graph).nodesInLayer(l)
	}
	if len(g.nodes) > l {
		n += len(g.nodes[l])
	}
	return n
}

func (g *graph) createNodeAt(nidx NodeIndex) *Node {
	for len(g.nodes) <= nidx[0] {
		g.nodes = append(g.nodes, []*Node{})
	}
	for len(g.nodes[nidx[0]]) <= nidx[1] {
		g.nodes[nidx[0]] = append(g.nodes[nidx[0]], nil)
	}

	node := &Node{
		Properties: Properties{},
	}
	g.nodes[nidx[0]][nidx[1]] = node

	return node
}

func (g *graph) findNodeInEdges(nidx NodeIndex, eidx EdgeIndex) int {
	if nodes, ok := g.edgeNodes[eidx]; ok {
		for i, side := range nodes.Nodes {
			for _, idx := range side {
				if nidx == idx {
					return i
				}
			}
		}
	}
	if g.parent != nil {
		return g.parent.(*graph).findNodeInEdges(nidx, eidx)
	} else {
		// TODO: untestable
		return 0
	}
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
