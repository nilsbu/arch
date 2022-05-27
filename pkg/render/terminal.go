package render

import (
	"errors"
	"fmt"
	"io"

	"github.com/nilsbu/arch/pkg/area"
	"github.com/nilsbu/arch/pkg/world"
)

// t0 is the unicode offset for simple box drawing
const t0 = 9472

// t1 is the unicode offset for double line box drawing
const t1 = 9552

// TODO Could use box drawing more extensively
// https://www.w3.org/TR/xml-entity-names/025.html

// ErrIllegalData is returned when false data was passed to a rendering function.
var ErrIllegalData = errors.New("illegal data")

// Terminal renders tiles by writing Unicode characters into a io.Writer.
// An error is returned when the input contains data that can't be rendered. This should not be the case when
// Terminal() is kept up-to-date.
func Terminal(w io.Writer, data *world.Tiles) error {
	w.Write([]byte(string(rune(t0 + 12))))

	for x := 0; x < data.Width(); x++ {
		w.Write([]byte(string(rune(t0 + 0))))
	}
	w.Write([]byte(string(rune(t0 + 16))))
	w.Write([]byte("\n"))

	for y := 0; y < data.Height(); y++ {
		w.Write([]byte(string(rune(t0 + 2))))

		for x := 0; x < data.Width(); x++ {
			if r, err := getChar(data, x, y); err != nil {
				return err
			} else {
				w.Write([]byte(string(r)))
			}
		}
		w.Write([]byte(string(rune(t0 + 2))))
		w.Write([]byte("\n"))
	}

	w.Write([]byte(string(rune(t0 + 20))))
	for x := 0; x < data.Width(); x++ {
		w.Write([]byte(string(rune(t0 + 0))))
	}
	w.Write([]byte(string(rune(t0 + 24))))
	w.Write([]byte("\n"))
	return nil
}

func getChar(data *world.Tiles, x, y int) (rune, error) {
	t := data.Get(x, y).Type

	switch t {
	case world.Free:
		return ' ', nil
	case world.Wall:
		return wall(data, x, y), nil
	case world.Door:
		return ' ', nil
	case world.Occupied:
		return '!', nil
	default:
		return 'x', fmt.Errorf("unexpected tile type: %v", t)
	}
}

func wall(data *world.Tiles, x, y int) rune {
	var o area.Direction
	if x > 0 && (data.Get(x-1, y).Type == world.Wall || data.Get(x-1, y).Type == world.Door) {
		o |= area.Left
	}
	if x+1 < data.Width() && (data.Get(x+1, y).Type == world.Wall || data.Get(x+1, y).Type == world.Door) {
		o |= area.Right
	}
	if y > 0 && (data.Get(x, y-1).Type == world.Wall || data.Get(x, y-1).Type == world.Door) {
		o |= area.Up
	}
	if y+1 < data.Height() && (data.Get(x, y+1).Type == world.Wall || data.Get(x, y+1).Type == world.Door) {
		o |= area.Down
	}

	switch o {
	case area.Left, area.Right, area.Left | area.Right:
		return t1
	case area.Up, area.Down, area.Up | area.Down:
		return t1 + 1
	case area.Down | area.Right:
		return t1 + 4
	case area.Down | area.Left:
		return t1 + 7
	case area.Up | area.Right:
		return t1 + 10
	case area.Up | area.Left:
		return t1 + 13
	case area.Up | area.Down | area.Right:
		return t1 + 16
	case area.Up | area.Down | area.Left:
		return t1 + 19
	case area.Left | area.Down | area.Right:
		return t1 + 22
	case area.Left | area.Up | area.Right:
		return t1 + 25
	case area.Left | area.Up | area.Right | area.Down:
		return t1 + 28
	default:
		return t1 + 91
	}
}
