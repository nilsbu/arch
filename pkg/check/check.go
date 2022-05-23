package check

import "github.com/nilsbu/arch/pkg/graph"

type Check interface {
	Match(graphs []*graph.Graph) (bool, error)
}
