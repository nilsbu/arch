package merge

import (
	"fmt"

	"github.com/nilsbu/arch/pkg/blueprint"
)

type choices struct {
	options   []group
	keys      []string
	blueprint blueprint.Blueprint
}

type group []*choices

func calcChoices(bp blueprint.Blueprint, property string, resolver *resolver) (*choices, error) {
	result := &choices{
		keys:      bp.Values(property),
		blueprint: bp,
	}

	if len(result.keys) == 0 {
		return nil, fmt.Errorf("%w: poperty '%v' has no values", ErrInvalidBlueprint, property)
	}

	result.options = make([]group, len(result.keys))
	for i, option := range result.keys {
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

func getGroup(bp blueprint.Blueprint, resolver *resolver) (group, error) {
	if name := bp.Values(resolver.name); len(name) != 1 {
		return nil, fmt.Errorf("%w: ", ErrInvalidBlueprint)
	} else {
		grp := []*choices{{
			blueprint: bp,
		}}
		for _, param := range resolver.keys[name[0]] {
			if choice, err := calcChoices(bp, param, resolver); err != nil {
				return nil, err
			} else {
				grp = append(grp, choice)
			}
		}
		return grp, nil
	}
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

func (c *choices) get(i int) []blueprint.Blueprint {
	for _, opt := range c.options {
		p := 1
		for _, s2 := range opt {
			p *= s2.n()
		}
		if i < p {
			bps := []blueprint.Blueprint{}
			for _, s2 := range opt {
				j := i % s2.n()
				i /= s2.n()
				bps = append(bps, s2.get(j)...)
			}
			return bps
		} else {
			i -= p
		}
	}
	return []blueprint.Blueprint{c.blueprint}
}
