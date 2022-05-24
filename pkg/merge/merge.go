package merge

import (
	"errors"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/check"
	"github.com/nilsbu/arch/pkg/graph"
)

// TODO find a package (name) that is better

var ErrInvalidBlueprint = errors.New("invalid blueprint")

var ErrNoSolution = errors.New("no solution found")

func Build(bp blueprint.Blueprint, check check.Check, r *Resolver) (*graph.Graph, error) {
	if choices, err := calcChoices(bp, "Root", r); err != nil {
		return nil, err
	} else {
		for i := 0; i < choices.n(); i++ {
			g := graph.New(nil)
			if err := parse(g, graph.NodeIndex{}, choices.get(i), r); err != nil {
				return nil, err
			} else if ok, err := check.Match([]*graph.Graph{g}); err != nil {
				return nil, err
			} else if ok {
				return g, nil
			}
		}

		return nil, ErrNoSolution
	}
}

func parse(g *graph.Graph, nidx graph.NodeIndex, choice *bpNode, r *Resolver) error {
	if choice.bp != nil {
		node := g.Node(nidx)
		node.Properties["name"] = choice.bp.Values(r.Name)[0]
		return nil
	} else if err := parse(g, nidx, choice.children[0], r); err != nil {
		return err
	} else {
		// each parameter has a child that may itself have multiple children
		for _, child := range choice.children[1:] {
			// child mustn't have a bp, so it's not checked here
			for _, grandchild := range child.children {
				if gcnidx, err := g.Add(nidx, nil); err != nil {
					return err
				} else if err := parse(g, gcnidx, grandchild, r); err != nil {
					return nil
				}
			}
		}
		return nil
	}
}
