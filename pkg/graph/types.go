package graph

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
