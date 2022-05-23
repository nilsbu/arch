package merge_test

import (
	"errors"
	"testing"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/pkg/merge"
	tg "github.com/nilsbu/arch/test/graph"
)

func TestBuild(t *testing.T) {
	for _, c := range []struct {
		name      string
		blueprint string
		graph     func() *graph.Graph
		err       error
	}{
		{
			"empty definition",
			"{}",
			func() *graph.Graph {
				return nil
			},
			merge.ErrInvalidBlueprint,
		},
		{
			"only a single rule",
			`{"Root":{"@":"R"}}`,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "R"
				return g
			},
			nil,
		},
		{
			"root references other property",
			`{"Root":"X","X":{"@":"R"}}`,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "R"
				return g
			},
			nil,
		},
		{
			"root has child",
			`{"Root":{"@":"1","a":{"@":"R"}}}`,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				nidx, _ := g.Add(graph.NodeIndex{}, nil)
				node = g.Node(nidx)
				node.Properties["name"] = "R"
				return g
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if bp, err := blueprint.Parse([]byte(c.blueprint)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else if graph, err := merge.Build(bp); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error")
			} else if !tg.AreEqual(c.graph(), graph) {
				t.Error("graph is wrong")
			}
		})
	}
}
