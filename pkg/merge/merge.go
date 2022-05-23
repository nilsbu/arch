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

		for _, ruleBp := range choices.get(0) {
			node := g.Node(graph.NodeIndex{})
			node.Properties["name"] = ruleBp.Values(r.name)[0]
		}
		return g, nil
	}
}
