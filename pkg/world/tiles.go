package world

import "github.com/nilsbu/async"

// TileType is the type of a tile.
// The type of a tile affects its role in the level.
type TileType uint8

const (
	Free TileType = iota
	Wall
	Door
	Occupied
)

// A Tile is the content of a slot in Tiles.
// It contains information about the TileType and about the appearance.
type Tile struct {
	Type    TileType
	Texture int
}

// Tiles is a field of tiles.
type Tiles struct {
	data          []Tile
	width, height int
}

// CreateTiles creates Tiles.
// The data is initializes with init in every slot.
func CreateTiles(width, height int, init Tile) *Tiles {
	data := make([]Tile, width*height)
	ts := &Tiles{
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

func (ts *Tiles) Get(x, y int) Tile {
	return ts.data[x+y*ts.width]
}

func (ts *Tiles) Set(x, y int, tile Tile) {
	ts.data[x+y*ts.width] = tile
}

func (ts *Tiles) Width() int {
	return ts.width
}

func (ts *Tiles) Height() int {
	return ts.height
}

// DrawFrame draws a non-filled rectangle.
// (x0, y0) is the top-left point, (x1, y1) is the bottom-right point.
func DrawFrame(ts *Tiles, x0, y0, x1, y1 int, tile Tile) {
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

// DrawRectangle fills a rectangle.
// (x0, y0) is the top-left point, (x1, y1) is the bottom-right point.
func DrawRectangle(ts *Tiles, x0, y0, x1, y1 int, tile Tile) {
	for yy := y0; yy <= y1; yy++ {
		for xx := x0; xx <= x1; xx++ {
			ts.Set(xx, yy, tile)
		}
	}
}
