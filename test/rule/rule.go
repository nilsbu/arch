package rule

import (
	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
)

type RuleMock struct {
	Params []string
	Prep   func(
		g *graph.Graph,
		nidx graph.NodeIndex,
		children map[string][]graph.NodeIndex,
		bp *blueprint.Blueprint,
	) error
}

func (r *RuleMock) ChildParams() []string {
	return r.Params
}

func (r *RuleMock) PrepareGraph(
	g *graph.Graph, nidx graph.NodeIndex, children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {
	if r.Prep != nil {
		return r.Prep(g, nidx, children, bp)
	} else {
		return nil
	}
}
