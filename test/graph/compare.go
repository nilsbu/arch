package graph

import (
	"reflect"

	g "github.com/nilsbu/arch/pkg/graph"
)

// AreEqual checks if two graphs are equal.
// Only the outputs are tested. If graphs have different parents but return the same values, true is returned.
func AreEqual(a, b *g.Graph) bool {
	if a == nil || b == nil {
		return a == b
	}

	queue := []g.NodeIndex{{}}
	testedEdges := map[g.EdgeIndex]bool{}

	for len(queue) > 0 {
		nidx := queue[0]
		queue = queue[1:]
		na, nb := a.Node(nidx), b.Node(nidx)

		if !reflect.DeepEqual(na, nb) {
			return false
		}

		for _, eidx := range na.Edges {
			testedEdges[eidx] = true
		}

		ca, cb := a.Children(nidx), b.Children(nidx)
		if !reflect.DeepEqual(ca, cb) {
			return false
		}

		queue = append(queue, ca...)
	}

	for eidx := range testedEdges {
		if !reflect.DeepEqual(a.Edge(eidx), b.Edge(eidx)) {
			return false
		}
		if !reflect.DeepEqual(a.Nodes(eidx), b.Nodes(eidx)) {
			return false
		}
	}

	return true
}
