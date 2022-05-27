package rule

import (
	"errors"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
)

var ErrInvalidGraph = errors.New("invalid graph")

var ErrPreparation = errors.New("error in preparation")

type Rule interface {
	ChildParams() []string
	PrepareGraph(
		g *graph.Graph,
		nidx graph.NodeIndex,
		children map[string][]graph.NodeIndex,
		bp *blueprint.Blueprint,
	) error
}
