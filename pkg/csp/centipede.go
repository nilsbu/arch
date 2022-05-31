package csp

import (
	"context"
	"fmt"

	"github.com/gnboorse/centipede"
	"github.com/nilsbu/arch/pkg/graph"
)

type Centipede struct {
	graphs []*graph.Graph
	// TODO make nodes actual nodes instead?
	nodes       [][]graph.NodeIndex
	vars        centipede.Variables[int]
	names       centipede.VariableNames
	constraints centipede.Constraints[int]
}

func (c *Centipede) Match(graphs []*graph.Graph) (ok bool, matches []graph.NodeIndex, err error) {
	if len(graphs) < 2 {
		return true, nil, nil
	}

	defer c.cleanup()

	if err := c.setup(graphs); err != nil {
		return false, nil, fmt.Errorf("cannot setup centipede matching: %w", err)
	}

	c.initVars()
	c.setAdjacencyConstraints()
	c.setHierarchyConstraints()

	solver := centipede.NewBackTrackingCSPSolver(c.vars, c.constraints)
	ok, err = solver.Solve(context.TODO())
	matches = c.getMatches()

	return
}

func (c *Centipede) cleanup() {
	c.graphs = nil
	c.nodes = nil
	c.vars = nil
	c.names = nil
	c.constraints = nil
}

func (c *Centipede) setup(graphs []*graph.Graph) error {
	for i, g := range graphs {
		c.graphs = append(c.graphs, g)
		if i == 0 {
			c.nodes = append(c.nodes, getAllChildren(g))
		} else {
			c.nodes = append(c.nodes, g.Children(graph.NodeIndex{}))
		}
	}
	return nil
}

func getAllChildren(g *graph.Graph) []graph.NodeIndex {
	nidxs := []graph.NodeIndex{{}}
	for i := 0; i < len(nidxs); i++ {
		nidxs = append(nidxs, g.Children(nidxs[i])...)
	}
	return nidxs
}

func (c *Centipede) initVars() {
	for i, nidx1 := range c.nodes[1] {
		name := centipede.VariableName(fmt.Sprint(i))
		domain := centipede.Domain[int]{}
		for j, nidx0 := range c.nodes[0] {
			if couldBe(c.graphs[0].Node(nidx0).Properties, c.graphs[1].Node(nidx1).Properties) {
				domain = append(domain, j)
			}
		}
		c.vars = append(c.vars, centipede.NewVariable(name, domain))
		c.names = append(c.names, name)
	}
}

func (c *Centipede) setAdjacencyConstraints() {
	adjacent := make([][][]bool, len(c.nodes))

	for i, nodes := range c.nodes {
		adjacent[i] = make([][]bool, len(nodes))
		lookup := map[graph.NodeIndex]int{}
		for i, nidx := range nodes {
			lookup[nidx] = i
		}

		for j, nidx := range nodes {
			adjacent[i][j] = make([]bool, len(nodes))
			for _, eidx := range c.graphs[i].Node(nidx).Edges {
				enidxs := c.graphs[i].Nodes(eidx)

				otherIdx := 0
				for _, nidx2 := range enidxs[0] {
					if nidx == nidx2 {
						otherIdx = 1
						break
					}
				}

				for _, onidx := range enidxs[otherIdx] {
					k := lookup[onidx]
					adjacent[i][j][k] = true
				}
			}
		}
	}

	for i, line := range adjacent[1] {
		for j, adj := range line {
			if i < j && adj {
				f := func(i, j int) centipede.VariablesConstraintFunction[int] {
					return func(vars *centipede.Variables[int]) bool {
						if vars.Find(c.names[i]).Empty || vars.Find(c.names[j]).Empty {
							return true
						}
						v0 := vars.Find(c.names[i]).Value
						v1 := vars.Find(c.names[j]).Value
						return adjacent[0][v0][v1]
					}
				}(i, j)

				c.constraints = append(c.constraints, centipede.Constraint[int]{
					Vars:               centipede.VariableNames{c.names[i], c.names[j]},
					ConstraintFunction: f,
				})
			}
		}
	}
}

func (c *Centipede) setHierarchyConstraints() {
	lookup := map[graph.NodeIndex]int{} // TODO don't duplicate lookup
	for i, nidx := range c.nodes[0] {
		lookup[nidx] = i
	}

	upAbove := make([][]bool, len(c.nodes[0]))
	for i := range c.nodes[0] {
		upAbove[i] = make([]bool, len(c.nodes[0]))
		for j := range c.nodes[0] {
			upAbove[i][j] = true
		}
	}

	for i := 0; i < len(c.nodes[0]); i++ {
		pre := c.nodes[0][i]
		idx0 := lookup[pre]
		upAbove[idx0][idx0] = false

		for parent := c.graphs[0].Node(pre).Parent; parent != graph.NoParent; parent = c.graphs[0].Node(pre).Parent {
			idx1 := lookup[parent]
			upAbove[idx0][idx1] = false
			upAbove[idx1][idx0] = false
			pre = parent
		}
	}

	for i := range c.nodes[1] {
		for j := range c.nodes[1] {
			if i != j {
				f := func(i, j int) centipede.VariablesConstraintFunction[int] {
					return func(vars *centipede.Variables[int]) bool {
						if vars.Find(c.names[i]).Empty || vars.Find(c.names[j]).Empty {
							return true
						}
						v0 := vars.Find(c.names[i]).Value
						v1 := vars.Find(c.names[j]).Value
						return upAbove[v0][v1]
					}
				}(i, j)

				c.constraints = append(c.constraints, centipede.Constraint[int]{
					Vars:               centipede.VariableNames{c.names[i], c.names[j]},
					ConstraintFunction: f,
				})
			}
		}
	}
}

func (c *Centipede) getMatches() []graph.NodeIndex {
	matches := make([]graph.NodeIndex, len(c.vars))
	for i, v := range c.vars {
		matches[i] = c.nodes[0][v.Value]
	}
	return matches
}

func couldBe(a, b graph.Properties) bool {
	// TODO find a better place for this
	if bname, ok := b["name"]; !ok {
		return true
	} else if aname, ok := a["name"]; !ok {
		return false
	} else {
		return aname == bname
	}
}
