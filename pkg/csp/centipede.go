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

func (c *Centipede) Match(graphs []*graph.Graph) (ok bool, err error) {
	if len(graphs) < 2 {
		return true, nil
	}

	defer c.cleanup()

	if err := c.setup(graphs); err != nil {
		return false, fmt.Errorf("cannot setup centipede matching: %w", err)
	}

	c.initVars()
	c.constraints = append(c.constraints, centipede.AllUnique[int](c.names...)...)
	c.setAdjacencyConstraints()

	solver := centipede.NewBackTrackingCSPSolver(c.vars, c.constraints)
	ok, err = solver.Solve(context.TODO())

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
	for _, g := range graphs {
		if flat, err := g.Leaves(); err != nil {
			return err
		} else {
			c.graphs = append(c.graphs, flat)
			c.nodes = append(c.nodes, flat.Children(graph.NodeIndex{}))
		}
	}
	return nil
}

func (c *Centipede) initVars() {
	for i, nidx1 := range c.nodes[1] {
		name := centipede.VariableName(fmt.Sprint(i))
		domain := centipede.Domain[int]{}
		for j, nidx0 := range c.nodes[0] {
			if couldBe(c.graphs[1].Node(nidx1).Properties, c.graphs[0].Node(nidx0).Properties) {
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
		for j, nidx0 := range nodes {
			adjacent[i][j] = make([]bool, len(nodes))
			for _, eidx := range c.graphs[i].Node(nidx0).Edges {
				enidxs := c.graphs[i].Nodes(eidx)
				var nidx1 graph.NodeIndex
				if enidxs[0][0] == nidx0 {
					nidx1 = enidxs[1][0]
				} else {
					nidx1 = enidxs[0][0]
				}
				adjacent[i][j][nidx1[1]] = true
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

func couldBe(a, b graph.Properties) bool {
	// TODO find a better place for this
	if aname, ok := a["name"]; !ok {
		return true
	} else if bname, ok := b["name"]; !ok {
		return true
	} else {
		return aname == bname
	}
}
