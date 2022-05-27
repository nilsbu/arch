package merge

import (
	"fmt"

	"github.com/nilsbu/arch/pkg/blueprint"
)

type choices struct {
	options   []group
	values    []string
	blueprint *blueprint.Blueprint
}

type group []*choices

func calcChoices(bp *blueprint.Blueprint, property string, resolver *Resolver) (*choices, error) {
	result := &choices{
		values:    bp.Values(property),
		blueprint: bp,
	}

	if len(result.values) == 0 {
		return nil, fmt.Errorf("%w: poperty '%v' has no values", ErrInvalidBlueprint, property)
	}

	result.options = make([]group, len(result.values))
	for i, option := range result.values {
		switch option[0] {
		case '*':
			if grp, err := getGroup(bp.Child(option), resolver); err != nil {
				return nil, err
			} else {
				result.options[i] = grp
			}
		default:
			if choice, err := calcChoices(bp, option, resolver); err != nil {
				return nil, err
			} else {
				result.options[i] = group{choice}
			}
		}
	}
	return result, nil
}

func getGroup(bp *blueprint.Blueprint, resolver *Resolver) (group, error) {
	if name := bp.Values(resolver.Name); len(name) != 1 {
		return nil, fmt.Errorf("%w: ", ErrInvalidBlueprint)
	} else {
		grp := []*choices{{
			blueprint: bp,
		}}
		rule := resolver.Keys[name[0]]
		if rule == nil {
			return nil, ErrUnknownKey
		}
		for _, param := range rule.ChildParams() {
			if choice, err := calcConjunction(bp, param, resolver); err != nil {
				return nil, err
			} else {
				grp = append(grp, choice)
			}
		}
		return grp, nil
	}
}

func calcConjunction(bp *blueprint.Blueprint, property string, resolver *Resolver) (*choices, error) {
	result := &choices{
		blueprint: bp,
	}

	values := bp.Values(property)
	if len(values) == 0 {
		return nil, fmt.Errorf("%w: poperty '%v' has no values", ErrInvalidBlueprint, property)
	}

	grp := make([]*choices, len(values))
	for i, option := range values {
		switch option[0] {
		case '*':
			if grp2, err := getGroup(bp.Child(option), resolver); err != nil {
				return nil, err
			} else {
				grp[i] = grp2[0]
			}
		default:
			if choice, err := calcChoices(bp, option, resolver); err != nil {
				return nil, err
			} else {
				grp[i] = choice
			}
		}
	}

	result.options = []group{grp}
	return result, nil
}

func (c *choices) n() int {
	if len(c.options) == 0 {
		return 1
	} else {
		n := 0
		for _, opt := range c.options {
			p := 1
			for _, s2 := range opt {
				p *= s2.n()
			}
			n += p
		}
		return n
	}
}

type bpNode struct {
	children []*bpNode
	bp       *blueprint.Blueprint
}

func (c *choices) get(i int) *bpNode {
	for _, opt := range c.options {
		p := 1
		for _, s2 := range opt {
			p *= s2.n()
		}
		if i < p {
			bps := &bpNode{}
			for _, s2 := range opt {
				j := i % s2.n()
				i /= s2.n()
				bps.children = append(bps.children, s2.get(j))
			}
			return bps
		} else {
			i -= p
		}
	}
	return &bpNode{bp: c.blueprint}
}
