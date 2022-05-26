package merge_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/csp"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/pkg/merge"
	"github.com/nilsbu/arch/pkg/rule"
	tg "github.com/nilsbu/arch/test/graph"
	tr "github.com/nilsbu/arch/test/rule"
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

func with(r *merge.Resolver, add map[string]rule.Rule) *merge.Resolver {
	out := &merge.Resolver{
		Name: r.Name,
		Keys: map[string]rule.Rule{},
	}

	for k, v := range r.Keys {
		out.Keys[k] = v
	}
	for k, v := range add {
		out.Keys[k] = v
	}

	return out
}

func TestBuild(t *testing.T) {
	allOk := checker(func([]*graph.Graph) (bool, error) { return true, nil })

	resolver := &merge.Resolver{
		Name: "@",
		Keys: map[string]rule.Rule{
			"1": &tr.RuleMock{Params: []string{"a"}},
			"R": &tr.RuleMock{},
			"P": &tr.RuleMock{},
		},
	}

	for _, c := range []struct {
		name       string
		blueprints []string
		check      merge.Check
		resolver   *merge.Resolver
		graph      func() *graph.Graph
		err        error
	}{
		{
			"empty definition",
			[]string{"{}"},
			allOk, resolver,
			func() *graph.Graph {
				return nil
			},
			merge.ErrInvalidBlueprint,
		},
		{
			"only a single rule",
			[]string{`{"Root":{"@":"R"}}`},
			allOk, resolver,
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
			[]string{`{"Root":"X","X":{"@":"R"}}`},
			allOk, resolver,
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
			[]string{`{"Root":{"@":"1","a":{"@":"R"}}}`},
			allOk, resolver,
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
			[]string{`{"Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":{"@":"P"}}]}`},
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
			resolver,
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
		{
			"check with centipede",
			[]string{
				`{"Root":[
					{"@":"1","a":[{"@":"R"},{"@":"P"},{"@":"R"}]}, 
					{"@":"1","a":{"@":"P"}}]}`,
				`{"Root":[
						{"@":"1","a":{"@":"P"}}]}`,
			},
			&csp.Centipede{},
			resolver,
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
		{
			"set properties for self and child",
			[]string{`{"Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":{"@":"P"}}]}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"1": &tr.RuleMock{
					Params: []string{"a"},
					Prep: func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp blueprint.Blueprint) error {
						if bp == nil {
							return errors.New("bp not there")
						}

						node := g.Node(nidx)
						node.Properties["set"] = "meeee"
						child := g.Node(children["a"][0])
						child.Properties["asdf"] = "qwerty"
						return nil
					},
				},
			}),
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["set"] = "meeee"
				nidx, _ := g.Add(graph.NodeIndex{}, nil)
				node = g.Node(nidx)
				node.Properties["name"] = "R"
				node.Properties["asdf"] = "qwerty"
				return g
			},
			nil,
		},
		{
			"unrecoverable error in PrepareGraph() causes failure",
			[]string{`{"Root":{"@":"1","a":{"@":"R"}}}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"1": &tr.RuleMock{
					Params: []string{"a"},
					Prep: func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp blueprint.Blueprint) error {

						return fmt.Errorf("%w", rule.ErrPreparation)
					},
				},
			}),
			func() *graph.Graph {
				return nil
			},
			rule.ErrPreparation,
		},
		{
			"recoverable error in PrepareGraph() causes rejection of first option",
			[]string{`{"Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":{"@":"P"}}]}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"1": &tr.RuleMock{
					Params: []string{"a"},
					Prep: func() func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp blueprint.Blueprint) error {
						count := 0
						return func(
							g *graph.Graph, nidx graph.NodeIndex,
							children map[string][]graph.NodeIndex, bp blueprint.Blueprint) error {

							count++
							if count == 1 {
								return fmt.Errorf("%w", rule.ErrInvalidGraph)
							} else {
								return nil
							}
						}
					}(),
				},
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
			bps := make([]blueprint.Blueprint, len(c.blueprints))
			for i, bp := range c.blueprints {
				var err error
				if bps[i], err = blueprint.Parse([]byte(bp)); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if graph, err := merge.Build(bps, c.check, c.resolver); err != nil && c.err == nil {
				t.Errorf("unexpected error: %v", err)
			} else if err == nil && c.err != nil {
				t.Errorf("expected error but non ocurred")
			} else if !errors.Is(err, c.err) {
				t.Errorf("wrong type of error\nexpect: %v\nactual: %v", c.err, err)
			} else if eq, ex := tg.AreEqual(c.graph(), graph); !eq {
				t.Error("graph is wrong:", ex)
			}
		})
	}
}
