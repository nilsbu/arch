package merge

import "github.com/nilsbu/arch/pkg/graph"

type Check interface {
	Match(graphs []*graph.Graph) (ok bool, matches []graph.NodeIndex, err error)
}

// TODO doc: incl expeced behaviour for zero graphs
