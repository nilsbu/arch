package merge_test

import (
	"errors"
	"testing"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/check"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/pkg/merge"
	tg "github.com/nilsbu/arch/test/graph"
)

type checker func([]*graph.Graph) (bool, error)

func (fc checker) Match(graphs []*graph.Graph) (bool, error) {
	return fc(graphs)
}

func getNodes(g *graph.Graph) []graph.NodeIndex {
	out := []graph.NodeIndex{{}}
	for i := 0; i < len(out); i++ {
		children := g.Children(out[i])
		out = append(out, children...)
	}
	return out
}

func TestBuild(t *testing.T) {
	allOk := checker(func([]*graph.Graph) (bool, error) { return true, nil })

	for _, c := range []struct {
		name      string
		blueprint string
		check     check.Check
		graph     func() *graph.Graph
		err       error
	}{
		{
			"empty definition",
			"{}",
			allOk,
			func() *graph.Graph {
				return nil
			},
			merge.ErrInvalidBlueprint,
		},
		{
			"only a single rule",
			`{"Root":{"@":"R"}}`,
			allOk,
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
			allOk,
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
			allOk,
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
		{
			"reject R in child",
			`{"Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":{"@":"P"}}]}`,
			checker(func(graphs []*graph.Graph) (bool, error) {
				for _, g := range graphs {
					for _, nidx := range getNodes(g) {
						if g.Node(nidx).Properties["name"] == "R" {
							return false, nil
						}
					}
				}
				return true, nil
			}),
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				nidx, _ := g.Add(graph.NodeIndex{}, nil)
				node = g.Node(nidx)
				node.Properties["name"] = "P"
				return g
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			if bp, err := blueprint.Parse([]byte(c.blueprint)); err != nil {
				t.Fatalf("unexpected error: %v", err)
			} else if graph, err := merge.Build(bp, c.check); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error")
			} else if eq, ex := tg.AreEqual(c.graph(), graph); !eq {
				t.Error("graph is wrong:", ex)
			}
		})
	}
}
