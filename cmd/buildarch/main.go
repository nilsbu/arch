package main

import (
	"fmt"
	"os"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/csp"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/pkg/merge"
	"github.com/nilsbu/arch/pkg/render"
	"github.com/nilsbu/arch/pkg/rule"
	"github.com/nilsbu/arch/pkg/world"
)

func main() {
	bps := make([]blueprint.Blueprint, len(os.Args)-1)
	for i := range bps {
		if file, err := os.ReadFile(os.Args[i+1]); err != nil {
			fmt.Println(err)
			return
		} else if bps[i], err = blueprint.Parse(file); err != nil {
			fmt.Println(err)
			return
		}
	}

	resolver := &merge.Resolver{
		Name: "@rule",
		Keys: map[string]rule.Rule{
			"House":    rule.House{},
			"Corridor": rule.Corridor{},
			"TwoRooms": rule.TwoRooms{},
			"Room":     rule.Room{},
			"NOP":      rule.NOP{},
		},
	}

	if g, err := merge.Build(bps, &csp.Centipede{}, resolver); err != nil {
		fmt.Println(err)
	} else {
		root := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
		rect := root.GetRect()
		tiles := world.CreateTiles(rect.X1+1, rect.Y1+1, world.Tile{})
		if err := draw(g, graph.NodeIndex{}, tiles); err != nil {
			fmt.Println(err)
		} else {
			render.Terminal(os.Stdout, tiles)
		}
	}
}

func draw(g *graph.Graph, nidx graph.NodeIndex, tiles world.Tiles) error {
	// TODO this is ugly, right a correct function for this
	a := (*area.AreaNode)(g.Node(nidx))
	rect := a.GetRect()
	if rect.X1 == 0 || rect.Y1 == 0 {
		return fmt.Errorf("rect for %v not set", nidx)
	} else if render, ok := a.Properties["render-walls"]; !ok || render.(bool) {
		for x := rect.X0; x <= rect.X1; x++ {
			tiles.Set(x, rect.Y0, world.Tile{Type: world.Wall})
			tiles.Set(x, rect.Y1, world.Tile{Type: world.Wall})
		}
		for y := rect.Y0; y <= rect.Y1; y++ {
			tiles.Set(rect.X0, y, world.Tile{Type: world.Wall})
			tiles.Set(rect.X1, y, world.Tile{Type: world.Wall})
		}
	}

	for _, cnidx := range g.Children(nidx) {
		if err := draw(g, cnidx, tiles); err != nil {
			return err
		}
	}

	// TODO If this is moved before the recursive call, it doesn't work. Why?
	for _, eidx := range a.Edges {
		pos := (*area.DoorEdge)(g.Edge(eidx)).GetPos()
		tiles.Set(pos.X, pos.Y, world.Tile{Type: world.Free})
	}
	return nil
}
