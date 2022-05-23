package merge

import (
	"errors"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
)

// TODO find a package (name) that is better

var ErrInvalidBlueprint = errors.New("invalid blueprint")

func Build(bp blueprint.Blueprint) (*graph.Graph, error) {
	r := &resolver{
		name: "@",
		keys: map[string][]string{
			"R": {},
		}}

	if choices, err := calcChoices(bp, "Root", r); err != nil {
		return nil, err
	} else {
		g := graph.New(nil)
		if err := parse(g, graph.NodeIndex{}, choices.get(0), r); err != nil {
			return nil, err
		} else {
			return g, nil
		}
	}
}

func parse(g *graph.Graph, nidx graph.NodeIndex, choice *bpNode, r *resolver) error {
	if choice.bp != nil {
		node := g.Node(graph.NodeIndex{})
		node.Properties["name"] = choice.bp.Values(r.name)[0]
		return nil
	} else if err := parse(g, nidx, choice.children[0], r); err != nil {
		return err
	} else {
		// node :=
		return nil
	}
}
