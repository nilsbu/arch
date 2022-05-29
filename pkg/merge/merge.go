package merge

import (
	"errors"
	"fmt"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/pkg/rule"
)

// TODO find a package (name) that is better

var ErrInvalidBlueprint = errors.New("invalid blueprint")

var ErrNoSolution = errors.New("no solution found")

func Build(bps []*blueprint.Blueprint, check Check, resolver *Resolver, shuffle Shuffle) (*graph.Graph, error) {
	choices := make([]*block, len(bps))
	ns := make([]int, len(bps))
	for i, bp := range bps {
		if block, err := calcBlock(bp, resolver); err != nil {
			return nil, err
		} else {
			choices[i] = block
			ns[i] = choices[i].n()
		}
	}

	for _, is := range shuffle(ns) {

		gs := make([]*graph.Graph, len(choices))
		ok := true
		for j, choices := range choices {

			gs[j] = graph.New(nil)
			if err := parse(gs[j], graph.NodeIndex{}, choices.get(is[j]), resolver); errors.Is(err, rule.ErrInvalidGraph) {
				ok = false
				break
			} else if err != nil {
				return nil, err
			}
		}
		if !ok {
			continue
		}

		if ok, err := check.Match(gs); err != nil {
			return nil, err
		} else if ok {
			return gs[0], nil
		}
	}

	return nil, ErrNoSolution
}

func parse(g *graph.Graph, nidx graph.NodeIndex, choice *bpNode, resolver *Resolver) error {
	node := g.Node(nidx)
	node.Properties["name"] = choice.bp.Values(resolver.Name)[0]

	nidxs := map[string][]graph.NodeIndex{}
	namedChoices := map[string][]*bpNode{}
	name := g.Node(nidx).Properties["name"].(string)
	r := resolver.Keys[name]
	names := r.ChildParams()

	for i, params := range choice.children {
		name := names[i]
		for _, grandchild := range params {
			gcnidx, _ := g.Add(nidx)
			namedChoices[name] = append(namedChoices[name], grandchild)
			nidxs[name] = append(nidxs[name], gcnidx)
		}
	}

	if err := r.PrepareGraph(g, nidx, nidxs, choice.bp); err != nil {
		return fmt.Errorf("couldn't create node of type '%v': %w", name, err)
	}

	for name, children := range nidxs {
		for i, child := range children {
			if err := parse(g, child, namedChoices[name][i], resolver); err != nil {
				return err
			}
		}
	}
	return nil
}
