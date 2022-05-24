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
		grandchildren := map[string][]graph.NodeIndex{}
		// each parameter has a child that may itself have multiple children
		names := r.Keys[g.Node(nidx).Properties["name"].(string)].ChildParams()
		for i, child := range choice.children[1:] {
			for range child.children {
				if gcnidx, err := g.Add(nidx, nil); err != nil {
					return err
				} else {
					grandchildren[names[i]] = append(grandchildren[names[i]], gcnidx)
				}
			}

			rule := r.Keys[g.Node(nidx).Properties["name"].(string)]
			rule.PrepareGraph(g, nidx, grandchildren, nil)

			i := 0
			for _, children := range grandchildren {
				for _, child := range children {
					if err := parse(g, child, choice.children[i+1], r); err != nil {
						return nil
					}
				}
				i++
			}
		}
		return nil
	}
}
