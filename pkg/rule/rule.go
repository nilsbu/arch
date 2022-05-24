package rule

import (
	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
)

type Rule interface {
	ChildParams() []string
	PrepareGraph(
		g *graph.Graph,
		nidx graph.NodeIndex,
		children map[string][]graph.NodeIndex,
		bp blueprint.Blueprint,
	) error
}
