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
	choicess := make([]*choices, len(bps))
	ns := make([]int, len(bps))
	for i, bp := range bps {
		if choices, err := calcChoices(bp, "Root", resolver); err != nil {
			return nil, err
		} else {
			choicess[i] = choices
			ns[i] = choicess[i].n()
		}
	}

	for _, is := range shuffle(ns) {

		gs := make([]*graph.Graph, len(choicess))
		ok := true
		for j, choices := range choicess {

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
	if choice.bp != nil {
		setName(g, nidx, choice, resolver)
		rule := resolver.Keys[choice.bp.Values(resolver.Name)[0]]
		return rule.PrepareGraph(g, nidx, map[string][]graph.NodeIndex{}, choice.bp)
	} else if choice.children[0].bp == nil {
		return parse(g, nidx, choice.children[0], resolver)
	} else {
		return parseBlock(g, nidx, choice, resolver)
	}
}

func setName(g *graph.Graph, nidx graph.NodeIndex, choice *bpNode, r *Resolver) {
	node := g.Node(nidx)
	node.Properties["name"] = choice.bp.Values(r.Name)[0]
}

func parseBlock(g *graph.Graph, nidx graph.NodeIndex, choice *bpNode, resolver *Resolver) error {
	setName(g, nidx, choice.children[0], resolver)
	nidxs := map[string][]graph.NodeIndex{}
	namedChoices := map[string][]*bpNode{}
	name := g.Node(nidx).Properties["name"].(string)
	r := resolver.Keys[name]
	names := r.ChildParams()
	for i, child := range choice.children[1:] {
		for _, grandchild := range child.children {
			gcnidx, _ := g.Add(nidx)
			name := names[i]
			namedChoices[name] = append(namedChoices[name], grandchild)
			nidxs[name] = append(nidxs[name], gcnidx)
		}
	}

	if err := r.PrepareGraph(g, nidx, nidxs, choice.children[0].bp); err != nil {
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
