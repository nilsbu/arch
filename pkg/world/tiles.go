package world

import "github.com/nilsbu/async"

type TileType uint8

const (
	Free TileType = iota
	Wall
	Door
	Occupied
)

type Tile struct {
	Type    TileType
	Texture int16
}

type Tiles interface {
	Get(x, y int) Tile
	Set(x, y int, tile Tile)

	Width() int
	Height() int
}

type tiles struct {
	data          []Tile
	width, height int
}

func CreateTiles(width, height int, init Tile) Tiles {
	data := make([]Tile, width*height)
	ts := &tiles{
		data:  data,
		width: width, height: height,
	}

	async.Pi(ts.height, func(y int) {
		for x := 0; x < ts.width; x++ {
			ts.Set(x, y, init)
		}
	})
	return ts
}

func (ts *tiles) Get(x, y int) Tile {
	return ts.data[x+y*ts.width]
}

func (ts *tiles) Set(x, y int, tile Tile) {
	ts.data[x+y*ts.width] = tile
}

func (ts *tiles) Width() int {
	return ts.width
}

func (ts *tiles) Height() int {
	return ts.height
}

func DrawFrame(ts Tiles, x0, y0, x1, y1 int, tile Tile) {
	for xx := x0; xx <= x1; xx++ {
		ts.Set(xx, y0, tile)
	}
	for yy := y0 + 1; yy < y1; yy++ {
		ts.Set(x0, yy, tile)
		ts.Set(x1, yy, tile)
	}
	for xx := x0; xx <= x1; xx++ {
		ts.Set(xx, y1, tile)
	}
}

func DrawRect(ts Tiles, x0, y0, x1, y1 int, tile Tile) {
	for yy := y0; yy <= y1; yy++ {
		for xx := x0; xx <= x1; xx++ {
			ts.Set(xx, yy, tile)
		}
	}
}
