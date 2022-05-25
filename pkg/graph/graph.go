package graph

import (
	"errors"
	"fmt"
)

// ErrIllegalAction is returned when an illegal action was attempted to be performed on a graph.
var ErrIllegalAction = errors.New("error during graph manipulation")

// NoParent is the NodeIndex used for a node's Parent that doesn't exist.
var NoParent = NodeIndex{-1, -1}

// New creates a new Graph.
// If parent is not nil, the new graph will inherit nodes and edges from the parent. Adding nodes and connections to the
// new graph will not affect the parent. The existing nodes and edges, however, are shared. This means that alterations
// of Properties will be shared. Making changes to parents using Add() and Link() after creating a new graph from it
// may not be safe and should generally not be done.
func New(parent *Graph) *Graph {
	if parent == nil {
		return &Graph{
			nodes:     [][]*Node{{{Properties: Properties{}, Parent: NoParent}}},
			children:  map[NodeIndex][]NodeIndex{},
			edges:     []*Edge{},
			edgeNodes: map[EdgeIndex]*edgeNodes{},
		}
	} else {
		return &Graph{
			parent:    parent,
			nodes:     [][]*Node{},
			children:  map[NodeIndex][]NodeIndex{},
			edges:     []*Edge{},
			edgeNodes: map[EdgeIndex]*edgeNodes{},
		}
	}
}

// Node returns the Node associated with a NodeIndex.
func (g *Graph) Node(nidx NodeIndex) *Node {
	if node := g.nodeSameInstance(nidx); node != nil {
		return node
	} else if g.parent != nil {
		return g.parent.Node(nidx)
	} else {
		return nil
	}
}

func (g *Graph) nodeSameInstance(nidx NodeIndex) *Node {
	if len(g.nodes) <= nidx[0] || len(g.nodes[nidx[0]]) <= nidx[1] {
		return nil
	} else {
		return g.nodes[nidx[0]][nidx[1]]
	}
}

// Children returns the children of a node.
func (g *Graph) Children(nidx NodeIndex) []NodeIndex {
	var children []NodeIndex
	if g.parent != nil {
		children = g.parent.Children(nidx)
	}
	if en, ok := g.children[nidx]; ok {
		children = append(children, en...)
	}
	return children
}

// Edge returns the Edge associated with an EdgeIndex.
func (g *Graph) Edge(eidx EdgeIndex) *Edge {
	if g.parent == nil {
		return g.edges[eidx]
	} else {
		offset := g.parent.countEdges()
		if int(eidx) < offset {
			return g.parent.Edge(eidx)
		} else {
			return g.edges[int(eidx)-offset]
		}
	}
}

func (g *Graph) countEdges() int {
	if g.parent == nil {
		return len(g.edges)
	} else {
		return len(g.edges) + g.parent.countEdges()
	}
}

// Nodes returns the nodes linked by an edge.
func (g *Graph) Nodes(eidx EdgeIndex) [2][]NodeIndex {
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

// Add creates a new node.
// The node has to have a parent and it may inherit edges that the parent has. Each edge may only be inherited once.
// The node will lie on the layer directly under the parent.
func (g *Graph) Add(parent NodeIndex, edges []EdgeIndex) (NodeIndex, error) {
	if parent[0] < 0 || parent[1] < 0 {
		return NodeIndex{}, fmt.Errorf("%w: cannot add node without parent", ErrIllegalAction)
	} else if g.Node(parent) == nil {
		return NodeIndex{}, fmt.Errorf("%w: parent node %v doesn't exist", ErrIllegalAction, parent)
	}

	nidx := NodeIndex{parent[0] + 1, g.nodesInLayer(parent[0] + 1)}
	node := g.createNodeAt(nidx)
	node.Parent = parent
	node.Edges = append(node.Edges, edges...)

	if children, ok := g.children[parent]; ok {
		g.children[parent] = append(children, nidx)
	} else {
		g.children[parent] = []NodeIndex{nidx}
	}

	for _, eidx := range edges {
		s := g.findNodeInEdges(parent, eidx)
		if s == -1 {
			return NodeIndex{}, fmt.Errorf("%w: edge %v doesn't belong to parent", ErrIllegalAction, eidx)
		}

		if nodes, ok := g.edgeNodes[eidx]; ok {
			nodes.Nodes[s] = append(nodes.Nodes[s], nidx)
		} else {
			n := &edgeNodes{}
			n.Nodes[s] = []NodeIndex{nidx}
			g.edgeNodes[eidx] = n
		}
	}

	return nidx, nil
}

func (g *Graph) nodesInLayer(l int) int {
	n := 0
	if g.parent != nil {
		n = g.parent.nodesInLayer(l)
	}
	if len(g.nodes) > l {
		n += len(g.nodes[l])
	}
	return n
}

func (g *Graph) createNodeAt(nidx NodeIndex) *Node {
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

func (g *Graph) findNodeInEdges(nidx NodeIndex, eidx EdgeIndex) int {
	if nodes, ok := g.edgeNodes[eidx]; ok {
		for i, side := range nodes.Nodes {
			if len(side) > 0 && side[len(side)-1] == nidx {
				return i
			}
		}
	}
	if g.parent != nil {
		return g.parent.findNodeInEdges(nidx, eidx)
	} else {
		return -1
	}
}

// Link creates an edge between two nodes.
// They must have the same parent and not be linked, already.
func (g *Graph) Link(a, b NodeIndex) (EdgeIndex, error) {
	if nodeA, nodeB := g.nodeSameInstance(a), g.nodeSameInstance(b); nodeA == nil || nodeB == nil {
		return -1, fmt.Errorf("%w: nodes must be created in the same graph as the edge", ErrIllegalAction)
	} else if nodeA.Parent != nodeB.Parent {
		return -1, fmt.Errorf("%w: nodes must have the same parent", ErrIllegalAction)
	} else if g.linkExists(a, b) {
		return -1, fmt.Errorf("%w: nodes are already linked", ErrIllegalAction)
	} else {
		eidx := EdgeIndex(g.countEdges())

		edge := &Edge{
			Properties: Properties{},
		}
		g.edges = append(g.edges, edge)

		g.edgeNodes[eidx] = &edgeNodes{
			Nodes: [2][]NodeIndex{{a}, {b}},
		}

		nodeA.Edges = append(nodeA.Edges, eidx)
		nodeB.Edges = append(nodeB.Edges, eidx)

		return EdgeIndex(eidx), nil
	}
}

// linkExists checks if two nodes that have the same parent are already linked. It will not look for edges between nodes
// of different parents.
func (g *Graph) linkExists(a, b NodeIndex) bool {
	for _, en := range g.edgeNodes {
		if (en.Nodes[0][0] == a && en.Nodes[1][0] == b) || (en.Nodes[0][0] == b && en.Nodes[1][0] == a) {
			return true
		}
	}

	// g.parent isn't tested since linkExists() is used in Link() and a link across graph instances isn't permitted
	return false
}
