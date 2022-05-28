package rule

import (
	"encoding/json"

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
	a.Properties["render-walls"] = false
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
		} else {
			return area.CreateDoor(g, children["interior"][0], children["exterior"][0], .5)
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
	a := (*area.AreaNode)(g.Node(nidx))
	rect := a.GetRect()

	nidxs := []graph.NodeIndex{
		children["left"][0],
		children["corridor"][0],
		children["right"][0],
	}

	roomOrientation := area.Turn(area.GetDirection(g, nidx, a.Edges[0]), 180)
	var roomWidth int
	if roomOrientation == area.Up || roomOrientation == area.Down {
		roomWidth = rect.Y1 - rect.Y0 + 1
	} else {
		roomWidth = rect.X1 - rect.X0 + 1
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
		return nil
	}
}

type TwoRooms struct{}

func (r TwoRooms) ChildParams() []string {
	return []string{"rooms"}
}

func (r TwoRooms) PrepareGraph(
	g *graph.Graph,
	nidx graph.NodeIndex,
	children map[string][]graph.NodeIndex,
	bp *blueprint.Blueprint,
) error {
	roomOrientation := area.Turn(area.GetDirection(g, nidx, g.Node(nidx).Edges[0]), 180)
	if err := area.Split(g, nidx, children["rooms"], []float64{.5}, roomOrientation); err != nil {
		return err
	} else {
		return area.CreateDoor(g, children["rooms"][0], children["rooms"][1], .5)
	}
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
	return nil
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
	g.Node(nidx).Properties["render-walls"] = false
	return nil
}
