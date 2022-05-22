package graph

// A Graph is a hierarchical undirected network of nodes and edges.
//
// A graph has at least one node, the root node.
// A node is something is connected to other nodes through edges. All nodes, except the root, have a parent.
// Each node may have an arbitrary number of children.
//
// An edge is a connection between nodes. It links at least two nodes that are siblings,
// meaning they have the same parent. Additionally an edge may also link descendents of the originally linked node.
// In each generation, only one child may be linked and no generations must be skipped. In other words, an edge linkes
// one unbroken line of parent-child related nodes with another.
//
// Both nodes and edges have Properties, which are maps with strings as keys and anything as values. They may be used
// arbitrarily by users of the graph.
//
// Even though graphs are undirected, an ordering of nodes can be inferred through the fact that Nodes() will
// always return a node as either the first or second slice.
type Graph interface {
	// Access
	Node(nidx NodeIndex) *Node
	Children(nidx NodeIndex) []NodeIndex
	Edge(eidx EdgeIndex) *Edge
	Nodes(eidx EdgeIndex) [2][]NodeIndex

	// Manipulate
	Add(parent NodeIndex, edges []EdgeIndex) (NodeIndex, error)
	Link(a, b NodeIndex) (EdgeIndex, error)
}

type NodeIndex [2]int
type EdgeIndex int
type Properties map[string]interface{}

type Node struct {
	Properties Properties
	Parent     NodeIndex
	Edges      []EdgeIndex
}

type Edge struct {
	Properties Properties
}
