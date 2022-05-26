package render_test

import (
	"strings"
	"testing"

	r "github.com/nilsbu/arch/pkg/render"
	"github.com/nilsbu/arch/pkg/world"
)

func toStr(codes []int) string {
	runes := make([]rune, len(codes))
	for i, c := range codes {
		runes[i] = rune(c)
	}
	return string(runes)
}

const t0 = 9472

func TestTerminal(t *testing.T) {
	for _, c := range []struct {
		name   string
		tiles  func() world.Tiles
		ok     bool
		expect string
	}{
		{
			"empty",
			func() world.Tiles { return world.CreateTiles(0, 0, world.Tile{}) },
			true,
			toStr([]int{
				t0 + 12, t0 + 16, 10,
				t0 + 20, t0 + 24, 10}),
		},
		{
			"all free",
			func() world.Tiles { return world.CreateTiles(3, 3, world.Tile{Type: world.Free}) },
			true,
			toStr([]int{
				t0 + 12, t0, t0, t0, t0 + 16, 10,
				t0 + 2, 32, 32, 32, t0 + 2, 10,
				t0 + 2, 32, 32, 32, t0 + 2, 10,
				t0 + 2, 32, 32, 32, t0 + 2, 10,
				t0 + 20, t0, t0, t0, t0 + 24, 10}),
		},
		{
			"all types",
			func() world.Tiles {
				data := world.CreateTiles(3, 3, world.Tile{Type: world.Free})
				world.DrawRect(data, 0, 0, 0, 1, world.Tile{Type: world.Wall})
				world.DrawRect(data, 1, 1, 1, 1, world.Tile{Type: world.Door})
				world.DrawRect(data, 2, 1, 2, 1, world.Tile{Type: world.Wall})
				world.DrawRect(data, 2, 2, 2, 2, world.Tile{Type: world.Occupied})
				return data
			},
			true,
			toStr([]int{
				t0 + 12, t0, t0, t0, t0 + 16, 10,
				t0 + 2, 9553, int(' '), int(' '), t0 + 2, 10,
				t0 + 2, 9562, int(' '), 9552, t0 + 2, 10,
				t0 + 2, int(' '), int(' '), int('!'), t0 + 2, 10,
				t0 + 20, t0, t0, t0, t0 + 24, 10}),
		},
		// TODO test all that could go wrong
	} {
		t.Run(c.name, func(t *testing.T) {
			b := &strings.Builder{}
			err := r.Terminal(b, c.tiles())
			if err != nil && c.ok {
				t.Fatal("unexpected error:", err)
			} else if err == nil && !c.ok {
				t.Fatal("expected error but none occured")
			}
			if err == nil {
				actual := b.String()
				if c.expect != actual {
					t.Errorf("expect:\n%v\nactual:\n%v",
						c.expect, actual)
				}
			}
		})
	}
}
