package world_test

import (
	"testing"

	"github.com/nilsbu/arch/pkg/world"
)

func TestTiles(t *testing.T) {
	for _, c := range []struct {
		name    string
		create  func() world.Tiles
		types   []world.TileType
		texture []int16
	}{
		{
			"zero, zero",
			func() world.Tiles {
				return world.CreateTiles(0, 0, world.Tile{world.Free, 0})
			},
			[]world.TileType{},
			[]int16{},
		},
		{
			"single line free",
			func() world.Tiles {
				return world.CreateTiles(3, 1, world.Tile{world.Wall, 1})
			},
			[]world.TileType{world.Wall, world.Wall, world.Wall},
			[]int16{1, 1, 1, 1, 1},
		},
		{
			"2D all the same",
			func() world.Tiles {
				return world.CreateTiles(3, 3, world.Tile{world.Wall, 1})
			},
			[]world.TileType{
				1, 1, 1,
				1, 1, 1,
				1, 1, 1,
			},
			nil,
		},
		{
			"frame",
			func() world.Tiles {
				ts := world.CreateTiles(5, 5, world.Tile{world.Free, 1})
				world.DrawFrame(ts, 1, 0, 4, 2, world.Tile{world.Wall, 2})
				return ts
			},
			[]world.TileType{
				0, 1, 1, 1, 1,
				0, 1, 0, 0, 1,
				0, 1, 1, 1, 1,
				0, 0, 0, 0, 0,
				0, 0, 0, 0, 0,
			},
			nil,
		},
		{
			"rect",
			func() world.Tiles {
				ts := world.CreateTiles(6, 5, world.Tile{world.Free, 1})
				world.DrawRect(ts, 1, 0, 4, 2, world.Tile{world.Wall, 2})
				return ts
			},
			[]world.TileType{
				0, 1, 1, 1, 1, 0,
				0, 1, 1, 1, 1, 0,
				0, 1, 1, 1, 1, 0,
				0, 0, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0,
			},
			nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ts := c.create()

			if c.types != nil {
				for y := 0; y < ts.Height(); y++ {
					for x := 0; x < ts.Width(); x++ {
						if c.types[x+y*ts.Width()] != ts.Get(x, y).Type {
							t.Errorf("type @ (%v, %v): expected %v but got %v", x, y,
								c.types[x+y*ts.Width()], ts.Get(x, y).Type)
						}
					}
				}
			}

			if c.texture != nil {
				for y := 0; y < ts.Height(); y++ {
					for x := 0; x < ts.Width(); x++ {
						if c.texture[x+y*ts.Width()] != ts.Get(x, y).Texture {
							t.Errorf("texture @ (%v, %v): expected %v but got %v", x, y,
								c.texture[x+y*ts.Width()], ts.Get(x, y).Texture)
						}
					}
				}
			}
		})
	}
}
