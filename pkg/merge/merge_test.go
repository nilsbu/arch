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

type checker func([]*graph.Graph) (bool, []graph.NodeIndex, error)

func (fc checker) Match(graphs []*graph.Graph) (bool, []graph.NodeIndex, error) {
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
	allOk := checker(func([]*graph.Graph) (bool, []graph.NodeIndex, error) { return true, nil, nil })

	resolver := &merge.Resolver{
		Name: "@",
		Keys: map[string]rule.Rule{
			"1": &tr.RuleMock{Params: []string{"a"}},
			"R": &tr.RuleMock{},
			"P": &tr.RuleMock{},
			"Error": &tr.RuleMock{Prep: func(
				g *graph.Graph,
				nidx graph.NodeIndex,
				children map[string][]graph.NodeIndex,
				bp *blueprint.Blueprint) error {
				return merge.ErrInvalidBlueprint
			},
			},
			"Leaf": &tr.RuleMock{Prep: func(
				g *graph.Graph,
				nidx graph.NodeIndex,
				children map[string][]graph.NodeIndex,
				bp *blueprint.Blueprint) error {
				g.Node(nidx).Properties["leaf"] = true
				return nil
			},
			},
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
			"no @",
			[]string{`{"@":""}`},
			allOk, resolver,
			func() *graph.Graph { return nil },
			merge.ErrInvalidBlueprint,
		},
		{
			"too many @",
			[]string{`{"@":["R","B"]}`},
			allOk, resolver,
			func() *graph.Graph { return nil },
			merge.ErrInvalidBlueprint,
		},
		{
			"@ is block",
			[]string{`{"@":{"@":"R"}}`},
			allOk, resolver,
			func() *graph.Graph { return nil },
			merge.ErrInvalidBlueprint,
		},
		{
			"@ has unknown name",
			[]string{`{"@":"Whoami"}`},
			allOk, resolver,
			func() *graph.Graph { return nil },
			merge.ErrInvalidBlueprint,
		},
		{
			"@ has unknown name in some depth",
			[]string{`{"@":"1","a":"V","V":{"@":"1","a":{"@":"Whoami"}}}`},
			allOk, resolver,
			func() *graph.Graph { return nil },
			merge.ErrInvalidBlueprint,
		},
		{
			"only a single rule",
			[]string{`{"@":"R"}`},
			allOk, resolver,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "R"
				node.Properties["names"] = []string{}
				return g
			},
			nil,
		},
		{
			"root has child",
			[]string{`{"@":"1","a":{"@":"R"}}`},
			allOk, resolver,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{}
				nidx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nidx)
				node.Properties["name"] = "R"
				node.Properties["names"] = []string{"1", "a"}
				return g
			},
			nil,
		},
		{
			"property 'a' required but not defined",
			[]string{`{"@":"1"}`},
			allOk, resolver,
			func() *graph.Graph { return nil },
			merge.ErrInvalidBlueprint,
		},
		{
			"reject R in child",
			[]string{`{"@":"1","a":"Root","Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":{"@":"P"}}]}`},
			checker(func(graphs []*graph.Graph) (bool, []graph.NodeIndex, error) {
				for _, g := range graphs {
					for _, nidx := range getNodes(g) {
						if g.Node(nidx).Properties["name"] == "R" {
							return false, nil, nil
						}
					}
				}
				return true, nil, nil
			}),
			resolver,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{}
				nidx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nidx)
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{"1", "a"}
				nidx, _ = g.Add(nidx)
				node = g.Node(nidx)
				node.Properties["name"] = "P"
				node.Properties["names"] = []string{"1", "a", "1", "a"}
				return g
			},
			nil,
		},
		{
			"check with centipede",
			[]string{
				`{"@":"1","a":"Root","Root":[
					{"@":"1","a":[{"@":"R"},{"@":"R"},{"@":"R"}]},
					{"@":"1","a":{"@":"P"}}]}`,
				`{"@":"1","a":{"@":"P"}}`,
			},
			&csp.Centipede{},
			resolver,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{}
				nidx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nidx)
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{"1", "a"}
				nidx, _ = g.Add(nidx)
				node = g.Node(nidx)
				node.Properties["name"] = "P"
				node.Properties["names"] = []string{"1", "a", "1", "a"}
				return g
			},
			nil,
		},
		{
			"error in matching",
			[]string{
				`{"@":"1","a":"Root","Root":[{"@":"1","a":{"@":"P"}}]}`,
				`{"@":"1","a":"Root","Root":[{"@":"1","a":{"@":"P"}}]}`,
			},
			checker(func(graphs []*graph.Graph) (bool, []graph.NodeIndex, error) {
				return false, nil, merge.ErrInvalidBlueprint // TODO this is not the correct error
			}),
			resolver,
			func() *graph.Graph { return nil },
			merge.ErrInvalidBlueprint,
		},
		{
			"different types of children",
			[]string{
				`{"@":"1","a":"Root","Root":[
					{"@":"1","a":[{"@":"R"},{"@":"P"}]}]}`,
			},
			&csp.Centipede{},
			resolver,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{}
				nx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nx)
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{"1", "a"}
				n0, _ := g.Add(nx)
				g.Node(n0).Properties["name"] = "R"
				g.Node(n0).Properties["names"] = []string{"1", "a", "1", "a"}
				n1, _ := g.Add(nx)
				g.Node(n1).Properties["name"] = "P"
				g.Node(n1).Properties["names"] = []string{"1", "a", "1", "a"}
				return g
			},
			nil,
		},
		{
			"set properties for self and child",
			[]string{`{"@":"1","a":"Root","Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":"Param","Param":{"@":"P"}}]}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"1": &tr.RuleMock{
					Params: []string{"a"},
					Prep: func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {
						if bp == nil {
							return errors.New("bp not there")
						}
						node := g.Node(nidx)
						node.Properties["set"] = "self"
						if as, ok := children["a"]; ok {
							child := g.Node(as[0])
							child.Properties["asdf"] = "child"
						}
						return nil
					},
				},
			}),
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["set"] = "self"
				node.Properties["names"] = []string{}
				nidx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nidx)
				node.Properties["name"] = "1"
				node.Properties["set"] = "self"
				node.Properties["asdf"] = "child"
				node.Properties["names"] = []string{"1", "a"}
				nidx, _ = g.Add(nidx)
				node = g.Node(nidx)
				node.Properties["name"] = "R"
				node.Properties["asdf"] = "child"
				node.Properties["names"] = []string{"1", "a", "1", "a"}
				return g
			},
			nil,
		},
		{
			"unrecoverable error in PrepareGraph() causes failure",
			[]string{`{"@":"1","a":"Root","Root":{"@":"1","a":{"@":"R"}}}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"1": &tr.RuleMock{
					Params: []string{"a"},
					Prep: func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {

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
			[]string{`{"@":"1","a":"Root","Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":{"@":"P"}}]}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"1": &tr.RuleMock{
					Params: []string{"a"},
					Prep: func() func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {
						count := 0
						return func(
							g *graph.Graph, nidx graph.NodeIndex,
							children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {

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
				node.Properties["names"] = []string{}
				nidx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nidx)
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{"1", "a"}
				nidx, _ = g.Add(nidx)
				node = g.Node(nidx)
				node.Properties["name"] = "P"
				node.Properties["names"] = []string{"1", "a", "1", "a"}
				return g
			},
			nil,
		},
		{
			"reject all options",
			[]string{`{"@":"1","a":"Root","Root":[{"@":"1","a":{"@":"R"}}, {"@":"1","a":{"@":"P"}}]}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"1": &tr.RuleMock{
					Params: []string{"a"},
					Prep: func() func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {
						return func(
							g *graph.Graph, nidx graph.NodeIndex,
							children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {
							return rule.ErrInvalidGraph
						}
					}(),
				},
			}),
			func() *graph.Graph { return nil },
			merge.ErrNoSolution,
		},
		{
			// this test was introduced because of a bug where PrepareGraph wasn't called on leafs
			"element that always fails",
			[]string{`{"@":"1","a":"Root","Root":{"@":"1","a":{"@":"X"}}}`},
			allOk,
			with(resolver, map[string]rule.Rule{
				"X": &tr.RuleMock{
					Prep: func() func(
						g *graph.Graph, nidx graph.NodeIndex,
						children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {
						return func(
							g *graph.Graph, nidx graph.NodeIndex,
							children map[string][]graph.NodeIndex, bp *blueprint.Blueprint) error {

							return fmt.Errorf("%w", rule.ErrPreparation)
						}
					}(),
				},
			}),
			func() *graph.Graph { return nil },
			rule.ErrPreparation,
		},
		{
			"set property on leaf",
			[]string{`{"@":"1","a":"Root","Root":{"@":"1","a":"X"},"X":{"@":"Leaf"}}`},
			allOk,
			resolver,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{}
				nidx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nidx)
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{"1", "a"}
				nidx, _ = g.Add(nidx)
				node = g.Node(nidx)
				node.Properties["name"] = "Leaf"
				node.Properties["names"] = []string{"1", "a", "1", "a"}
				node.Properties["leaf"] = true
				return g
			},
			nil,
		},
		{
			"block in block",
			[]string{`{"@":"1","a":"Root","Root":{"@":"1","a":{"@":"1","a":"X"}},"X":{"@":"Leaf"}}`},
			allOk,
			resolver,
			func() *graph.Graph {
				g := graph.New(nil)
				node := g.Node(graph.NodeIndex{})
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{}
				nidx, _ := g.Add(graph.NodeIndex{})
				node = g.Node(nidx)
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{"1", "a"}
				nidx, _ = g.Add(nidx)
				node = g.Node(nidx)
				node.Properties["name"] = "1"
				node.Properties["names"] = []string{"1", "a", "1", "a"}
				nidx, _ = g.Add(nidx)
				node = g.Node(nidx)
				node.Properties["name"] = "Leaf"
				node.Properties["names"] = []string{"1", "a", "1", "a", "1", "a"}
				node.Properties["leaf"] = true
				return g
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			bps := make([]*blueprint.Blueprint, len(c.blueprints))
			for i, bp := range c.blueprints {
				var err error
				if bps[i], err = blueprint.Parse([]byte(bp)); err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}

			if graph, err := merge.Build(bps, c.check, c.resolver, merge.InOrder); err != nil && c.err == nil {
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
