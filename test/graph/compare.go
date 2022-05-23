package graph

import (
	"fmt"
	"reflect"

	g "github.com/nilsbu/arch/pkg/graph"
)

// AreEqual checks if two graphs are equal.
// Only the outputs are tested. If graphs have different parents but return the same values, true is returned.
func AreEqual(a, b *g.Graph) (equal bool, explanation string) {
	if a == nil && b != nil {
		return false, "first is nil"
	} else if a != nil && b == nil {
		return false, "second is nil"
	} else if a == nil && b == nil {
		return true, ""
	}

	queue := []g.NodeIndex{{}}
	testedEdges := map[g.EdgeIndex]bool{}

	for len(queue) > 0 {
		nidx := queue[0]
		queue = queue[1:]
		na, nb := a.Node(nidx), b.Node(nidx)

		if !reflect.DeepEqual(na.Properties, nb.Properties) {
			return false, fmt.Sprintf("properties of %v are different: %v vs. %v",
				nidx, na.Properties, nb.Properties)
		}

		dea, deb := getDisjunctEidxs(na.Edges, nb.Edges)
		if len(dea) != 0 || len(deb) != 0 {
			return false, fmt.Sprintf("edges of %v are disjunct in %v vs. %v",
				nidx, dea, deb)
		}

		for _, eidx := range na.Edges {
			testedEdges[eidx] = true
		}

		ca, cb := a.Children(nidx), b.Children(nidx)
		dna, dnb := getDisjunctNidxs(ca, cb)
		if len(dna) != 0 || len(dnb) != 0 {
			return false, fmt.Sprintf("children of %v are disjunct in %v vs. %v",
				nidx, dna, dnb)
		}

		queue = append(queue, ca...)
	}

	for eidx := range testedEdges {
		pa, pb := a.Edge(eidx).Properties, b.Edge(eidx).Properties
		if !reflect.DeepEqual(pa, pb) {
			return false, fmt.Sprintf("properties of edge %v don't match: %v vs. %v",
				eidx, pa, pb)
		}
		if !reflect.DeepEqual(a.Nodes(eidx), b.Nodes(eidx)) {
			return false, fmt.Sprintf("nodes of edge %v don't match: %v vs. %v",
				eidx, a.Nodes(eidx), b.Nodes(eidx))
		}
	}

	return true, ""
}

func getDisjunctNidxs(a, b []g.NodeIndex) (da, db []g.NodeIndex) {
	for _, nidx := range a {
		if !findNidx(nidx, b) {
			da = append(da, nidx)
		}
	}
	for _, nidx := range b {
		if !findNidx(nidx, a) {
			db = append(db, nidx)
		}
	}
	return
}

func findNidx(nidx g.NodeIndex, in []g.NodeIndex) bool {
	for _, oidx := range in {
		if nidx == oidx {
			return true
		}
	}
	return false
}

func getDisjunctEidxs(a, b []g.EdgeIndex) (da, db []g.EdgeIndex) {
	for _, eidx := range a {
		if !findEidx(eidx, b) {
			da = append(da, eidx)
		}
	}
	for _, eix := range b {
		if !findEidx(eix, a) {
			db = append(db, eix)
		}
	}
	return
}

func findEidx(eidx g.EdgeIndex, in []g.EdgeIndex) bool {
	for _, oidx := range in {
		if eidx == oidx {
			return true
		}
	}
	return false
}
