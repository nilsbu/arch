package rule

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
)

type House struct{}

func (r House) ChildParams() []string {
	return []string{"interior", "exterior"}
}

func (r House) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	a := (*area.AreaNode)(g.Node(nidx))
	a.Properties["render"] = false
	data := []int{}
	if err := json.Unmarshal([]byte(bp.Values("rect")[0]), &data); err != nil {
		return err
	} else {
		// Add one to have room at the bottom for the exterior
		data[3]++
		a.SetRect(area.Rectangle{X0: data[0], Y0: data[1], X1: data[2], Y1: data[3]})

		nidxs := []graph.NodeIndex{
			children["interior"][0],
			children["exterior"][0],
		}

		h := float64(data[3] - data[1])
		if err := area.Split(g, nidx, nidxs, []float64{(h - 1) / h}, area.Down); err != nil {
			return err
		} else if err := area.CreateDoor(g, children["interior"][0], children["exterior"][0], .5); err != nil {
			return err
		} else {
			return InheritEdges(g, nidx)
		}
	}
}

type Corridor struct{}

func (r Corridor) ChildParams() []string {
	return []string{"left", "corridor", "right"}
}

func (r Corridor) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	g.Node(nidx).Properties["render"] = false
	a := (*area.AreaNode)(g.Node(nidx))
	rect := a.GetRect()

	nidxs := []graph.NodeIndex{
		children["left"][0],
		children["corridor"][0],
		children["right"][0],
	}

	roomOrientation := RoomOrientation(g, nidx)
	var roomWidth int
	if roomOrientation == area.Up || roomOrientation == area.Down {
		roomWidth = rect.X1 - rect.X0
	} else {
		roomWidth = rect.Y1 - rect.Y0
	}

	corridorWidth := 3.
	at := []float64{
		.5 - corridorWidth/float64(roomWidth)/2,
		.5 + (corridorWidth+2)/float64(roomWidth)/2,
	}

	if err := area.Split(g, nidx, nidxs, at, area.Turn(roomOrientation, 90)); err != nil {
		return err
	} else {
		for _, side := range []string{"left", "right"} {
			at := make([]float64, len(children[side])-1)
			for i := range at {
				at[i] = float64(i+1) / float64(len(children[side]))
			}
			if err := area.Split(g, children[side][0], children[side], at, roomOrientation); err != nil {
				return err
			}
			for _, cnidx := range children[side] {
				if err := area.CreateDoor(g, children["corridor"][0], cnidx, .5); err != nil {
					return err
				}
			}
		}
		return InheritEdges(g, nidx)
	}
}

type RoomLine struct{}

func (r RoomLine) ChildParams() []string {
	return []string{"rooms"}
}

func (r RoomLine) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	g.Node(nidx).Properties["render"] = false
	cnidxs := children["rooms"]
	at := make([]float64, len(cnidxs)-1)
	for i := range at {
		at[i] = float64(i+1) / float64(len(cnidxs))
	}
	if err := area.Split(g, nidx, cnidxs, at, RoomOrientation(g, nidx)); err != nil {
		return err
	}
	for i := range at {
		if err := area.CreateDoor(g, cnidxs[i], cnidxs[i+1], .5); err != nil {
			return err
		}
	}

	return InheritEdges(g, nidx)
}

type Frame struct{}

func (r Frame) ChildParams() []string {
	return []string{"content"}
}

func (r Frame) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	cnidxs := children["content"]
	if len(cnidxs) != 1 {
		return ErrPreparation
	}

	if err := area.Split(g, nidx, cnidxs, []float64{}, RoomOrientation(g, nidx)); err != nil {
		return err
	}

	return InheritEdges(g, nidx)
}

type Room struct{}

func (r Room) ChildParams() []string {
	return []string{}
}

func (r Room) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	SetWall(g, nidx, true)
	return nil
}

type FurnishedRoom struct{}

func (r FurnishedRoom) ChildParams() []string {
	return []string{"furniture"}
}

func (r FurnishedRoom) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	SetWall(g, nidx, true)

	if furniture, ok := children["furniture"]; ok {
		interior := g.Node(furniture[0])
		a := (*area.AreaNode)(g.Node(nidx))
		rect := a.GetRect()
		(*area.AreaNode)(interior).SetRect(
			area.Rectangle{X0: rect.X0 + 1, Y0: rect.Y0 + 1, X1: rect.X1 - 1, Y1: rect.Y1 - 1})
	}

	return nil
}

type Furniture struct{}

func (r Furniture) ChildParams() []string {
	return []string{"elements"}
}

func (r Furniture) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	SetWall(g, nidx, false)

	elements := children["elements"]
	sizes := bp.Values("sizes")
	if len(elements) != len(sizes) {
		return fmt.Errorf("%w: have %v elements and %v sizes",
			ErrPreparation, len(elements), len(sizes))
	} else {
		a := (*area.AreaNode)(g.Node(nidx))
		rect := a.GetRect()
		for i := range sizes {
			size := []int{}
			if err := json.Unmarshal([]byte(bp.Values("sizes")[i]), &size); err != nil {
				return err
			} else {
				e := (*area.AreaNode)(g.Node(elements[i]))
				// TODO allow for positions other then top-left
				e.SetRect(area.Rectangle{
					X0: rect.X0, Y0: rect.Y0, X1: rect.X0 + size[0] - 1, Y1: rect.Y0 + size[1] - 1})
			}
		}
		return nil
	}
}

type NOP struct{}

func (r NOP) ChildParams() []string {
	return []string{}
}

func (r NOP) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	SetWall(g, nidx, false)
	return nil
}

type Occupy struct{}

func (r Occupy) ChildParams() []string {
	return []string{}
}

func (r Occupy) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	SetWall(g, nidx, false)

	if texture := bp.Values("texture"); len(texture) != 1 {
		return ErrPreparation
	} else if tex, err := strconv.Atoi(texture[0]); err != nil {
		return err
	} else {
		g.Node(nidx).Properties["object"] = tex
		return nil
	}
}
