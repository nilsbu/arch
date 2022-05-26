package graph

import (
	"errors"
)

var ErrNotLeafable = errors.New("cannot build valid leaf graph when nodes aren't fully inherited")

func (g *Graph) Leaves() (*Graph, error) {
	out := New(nil)

	mapping := make(map[NodeIndex]NodeIndex)
	g.leaves(NodeIndex{}, out, mapping)

	return out, g.leafEdges(out, mapping)
}

func (g *Graph) leaves(nidx NodeIndex, out *Graph, mapping map[NodeIndex]NodeIndex) {
	node := g.Node(nidx)
	children := g.Children(nidx)
	if len(children) == 0 {
		cidx, _ := out.Add(NodeIndex{})
		outNode := out.Node(cidx)
		for k, v := range node.Properties {
			outNode.Properties[k] = v
		}
		mapping[nidx] = cidx
	} else {
		for _, cidx := range children {
			g.leaves(cidx, out, mapping)
		}
	}
}

func (g *Graph) leafEdges(out *Graph, mapping map[NodeIndex]NodeIndex) error {
	for eidx := EdgeIndex(0); int(eidx) < g.countEdges(); eidx++ {
		edgeNodes := g.Nodes(eidx)

		var oidxs [2]NodeIndex
		for j := 0; j < 2; j++ {
			if oidx, ok := mapping[edgeNodes[j][len(edgeNodes[j])-1]]; ok {
				oidxs[j] = oidx
			} else {
				return ErrNotLeafable
			}
		}

		oeidx, _ := out.Link(oidxs[0], oidxs[1])
		new := out.Edge(oeidx)
		for k, v := range g.Edge(eidx).Properties {
			new.Properties[k] = v
		}
	}
	return nil
}
