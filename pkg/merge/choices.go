package merge

import (
	"fmt"

	"github.com/nilsbu/arch/pkg/blueprint"
)

// TODO document all of this
type block struct {
	params    []group
	blueprint *blueprint.Blueprint
}

func (b *block) n() int {
	value := 1
	for _, param := range b.params {
		for _, iob := range param {
			value *= iob.n()
		}
	}
	return value
}

type bpNode struct {
	children [][]*bpNode
	bp       *blueprint.Blueprint
}

func (b *block) get(i int) *bpNode {
	res := &bpNode{
		bp: b.blueprint,
	}
	res.children = make([][]*bpNode, len(b.params))
	for k, param := range b.params {
		pnodes := make([]*bpNode, len(param))
		for l, iob := range param {
			x := iob.n()
			j := i % x
			i /= x
			pnodes[l] = iob.get(j)
		}
		res.children[k] = pnodes
	}
	return res
}

type group []*groupOrBlock

func (g group) n() int {
	value := 0
	for _, iob := range g {
		value += iob.n()
	}
	return value
}
func (g group) get(i int) (bpn *bpNode) {
	for _, iob := range g {
		x := iob.n()
		if i < x {
			bpn = iob.get(i)
			break
		} else {
			i -= x
		}
	}
	// Slightly awkward construction due to my obsession with 100% coverage. There's no realistic way to produce an
	// error here, so I chose to break earlier and thus force to call the following line.
	return
}

type groupOrBlock struct {
	group group
	block *block
}

func (gob groupOrBlock) n() int {
	if gob.group != nil {
		return gob.group.n()
	} else {
		return gob.block.n()
	}
}

func (gob groupOrBlock) get(i int) *bpNode {
	if gob.group != nil {
		return gob.group.get(i)
	} else {
		return gob.block.get(i)
	}
}

func calcBlock(bp *blueprint.Blueprint, resolver *Resolver) (*block, error) {
	if name := bp.Values(resolver.Name); len(name) != 1 {
		return nil, fmt.Errorf("%w: ", ErrInvalidBlueprint)
	} else {
		blck := &block{
			blueprint: bp,
		}

		rule := resolver.Keys[name[0]]
		if rule == nil {
			return nil, fmt.Errorf("%w: key '%v' is not defined", ErrInvalidBlueprint, name[0])
		}
		for _, param := range rule.ChildParams() {
			if grp, err := calcGroup(bp, param, resolver); err != nil {
				return nil, err
			} else {
				blck.params = append(blck.params, grp)
			}
		}
		return blck, nil
	}
}

func calcGroup(bp *blueprint.Blueprint, property string, resolver *Resolver) (group, error) {
	values := bp.Values(property)
	group := make(group, len(values))
	for i, opt := range values {
		var err error
		if group[i], err = calcGroupOrBlock(bp, opt, resolver); err != nil {
			return nil, err
		}
	}
	return group, nil
}

func calcGroupOrBlock(bp *blueprint.Blueprint, property string, resolver *Resolver) (*groupOrBlock, error) {
	switch property[0] {
	case '*':
		blck, err := calcBlock(bp.Child(property), resolver)
		return &groupOrBlock{block: blck}, err
	default:
		choice, err := calcGroup(bp, property, resolver)
		return &groupOrBlock{group: choice}, err
	}
}
