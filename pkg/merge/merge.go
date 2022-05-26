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

func Build(bps []blueprint.Blueprint, check Check, r *Resolver) (*graph.Graph, error) {
	choicess := make([]*choices, len(bps))
	ns := &ns{}
	for i, bp := range bps {
		if choices, err := calcChoices(bp, "Root", r); err != nil {
			return nil, err
		} else {
			choicess[i] = choices
			ns.add(choicess[i].n())
		}
	}

	for i := 0; i < ns.total; i++ {
		is := ns.get(i)

		gs := make([]*graph.Graph, len(choicess))
		ok := true
		for j, choices := range choicess {

			gs[j] = graph.New(nil)
			if err := parse(gs[j], graph.NodeIndex{}, choices.get(is[j]), r); errors.Is(err, rule.ErrInvalidGraph) {
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

func parse(g *graph.Graph, nidx graph.NodeIndex, choice *bpNode, r *Resolver) error {
	if choice.bp != nil {
		node := g.Node(nidx)
		node.Properties["name"] = choice.bp.Values(r.Name)[0]
		return nil
	} else if choice.children[0].bp == nil {
		return parse(g, nidx, choice.children[0], r)
	} else {
		return parseBlock(g, nidx, choice, r)
	}
}

func parseBlock(g *graph.Graph, nidx graph.NodeIndex, choice *bpNode, r *Resolver) error {
	if err := parse(g, nidx, choice.children[0], r); err != nil {
		return err
	} else {
		nidxs := map[string][]graph.NodeIndex{}
		namedChoices := map[string]*bpNode{}
		name := g.Node(nidx).Properties["name"].(string)
		rule := r.Keys[name]
		names := rule.ChildParams()
		for i, child := range choice.children[1:] {
			for range child.children {
				if gcnidx, err := g.Add(nidx, nil); err != nil {
					return err
				} else {
					name := names[i]
					namedChoices[name] = child
					nidxs[name] = append(nidxs[name], gcnidx)
				}
			}
		}

		if err := rule.PrepareGraph(g, nidx, nidxs, nil); err != nil {
			return fmt.Errorf("couldn't create node of type '%v': %w", name, err)
		}

		for name, children := range nidxs {
			for _, child := range children {
				if err := parse(g, child, namedChoices[name], r); err != nil {
					return err
				}
			}
		}
		return nil
	}
}

// ns handles the iteration over multiple choices.
type ns struct {
	// TODO this is ugly, replace it by something better
	total int
	ns    []int
}

func (ns *ns) add(n int) {
	if ns.total == 0 {
		ns.total = n
	} else {
		ns.total *= n
	}
	ns.ns = append(ns.ns, n)
}

func (ns *ns) get(i int) []int {
	is := make([]int, len(ns.ns))
	exp := 1
	for j, n := range ns.ns {
		is[j] = (i / exp) % n
		exp *= n
	}

	return is
}
