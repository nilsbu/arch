package draw

import (
	"errors"
	"fmt"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/pkg/world"
)

var ErrInvalidGraph = errors.New("graph cannot be drawn")

func Draw(g *graph.Graph) (*world.Tiles, error) {
	root := (*area.AreaNode)(g.Node(graph.NodeIndex{}))
	rect := root.GetRect()
	if rect.X0 >= rect.X1 && rect.Y0 >= rect.Y1 {
		return nil, fmt.Errorf("%w: ", ErrInvalidGraph)
	}

	data := world.CreateTiles(rect.X1+1, rect.Y1+1, world.Tile{Type: world.Free})
	return data, draw(g, graph.NodeIndex{}, data)
}

func draw(g *graph.Graph, nidx graph.NodeIndex, tiles *world.Tiles) error {
	a := (*area.AreaNode)(g.Node(nidx))
	rect := a.GetRect()
	if rect.X1 == 0 || rect.Y1 == 0 {
		return fmt.Errorf("%w: rect for %v not set", ErrInvalidGraph, nidx)
	} else if object, ok := a.Properties["object"]; ok {
		for y := rect.Y0; y <= rect.Y1; y++ {
			for x := rect.X0; x <= rect.X1; x++ {
				tiles.Set(x, y, world.Tile{Type: world.Occupied, Texture: object.(int)})
			}
		}
	} else if render, ok := a.Properties["render"]; !ok || render.(bool) {
		for x := rect.X0; x <= rect.X1; x++ {
			tiles.Set(x, rect.Y0, world.Tile{Type: world.Wall})
			tiles.Set(x, rect.Y1, world.Tile{Type: world.Wall})
		}
		for y := rect.Y0; y <= rect.Y1; y++ {
			tiles.Set(rect.X0, y, world.Tile{Type: world.Wall})
			tiles.Set(rect.X1, y, world.Tile{Type: world.Wall})
		}
	}

	for _, eidx := range a.Edges {
		door := (*area.DoorEdge)(g.Edge(eidx))
		pos := door.GetPos()
		if render, ok := door.Properties["render"]; !ok || render.(bool) {
			tiles.Set(pos.X, pos.Y, world.Tile{Type: world.Door})
		} else {
			tiles.Set(pos.X, pos.Y, world.Tile{Type: world.Free})
		}
	}

	for _, cnidx := range g.Children(nidx) {
		if err := draw(g, cnidx, tiles); err != nil {
			return err
		}
	}
	return nil
}
