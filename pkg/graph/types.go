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
type Graph struct {
	parent *Graph

	nodes     [][]*Node
	children  map[NodeIndex][]NodeIndex
	edges     []*Edge
	edgeNodes map[EdgeIndex]*edgeNodes
}

type edgeNodes struct {
	Nodes [2][]NodeIndex
}

// A NodeIndes is the index of a node in a graph.
// The first number refers to the layer on which the node exists, the second to the index of the node within the layer.
type NodeIndex [2]int

// An EdgeIndex is the index of an edge in a graph.
type EdgeIndex int

// Properties is keyed data associates with a node or edge.
type Properties map[string]interface{}

// A Node is a node (or vertex) in the fundamental unit of which graphs are formed.
type Node struct {
	// Properties containes user data.
	Properties Properties

	// Parent refers to the parent node in the hierarchy. The root node will have NoParent as its Parent.
	// Editing this value is not permitted.
	Parent NodeIndex

	// Edges are the edges that are connected to the node.
	// Editing this slice directly is not permitted.
	Edges []EdgeIndex
}

// An Edge is a link between two unbroken line of parent-child related nodes in a graph.
type Edge struct {
	// Properties containes user data.
	Properties Properties
}
