package rule

import (
	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
)

type Path struct{}

func (r Path) ChildParams() []string {
	return []string{"steps"}
}

func (r Path) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	cnidxs := children["steps"]

	for i := 0; i < len(cnidxs)-1; i++ {
		if _, err := g.Link(cnidxs[i], cnidxs[i+1]); err != nil {
			return err
		}
	}

	return nil
}

type In struct{}

func (r In) ChildParams() []string {
	return []string{}
}

func (r In) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	g.Node(nidx).Properties["name"] = bp.Values("name")[0]
	return nil
}
