package area

import (
	"errors"
	"fmt"
)

// ErrInvalidSplit is returned when an AreaNode cannot be split.
var ErrInvalidSplit = errors.New("invalid split")

// Split splits an area into smaller areas.
// The resulting areas are strung together either horizontally or vertically.
// Their nodes must already exist and are passed as "into". The original node, named "base", may occur in the resulting
// ares, thus shrinking it.
// "direction" is the direction along which the areas are aligned. E.g. if Down is chosen, the first resulting area is
// the hightest (smallest y value), and the other ones follow below it.
// "at" is a sequence of numbers in range [0, 1] determining where along the splitting axis the borders between the
// areas lie.
func Split(base *AreaNode, into []*AreaNode, at []float64, direction Direction) error {
	if len(into) != len(at)+1 {
		return fmt.Errorf("%w: tried to split into %v nodes with %v dividers", ErrInvalidSplit, len(into), len(at))
	}
	ats := []float64{0}
	ats = append(ats, at...)
	ats = append(ats, 1)

	flipped := flip(base.GetRect(), direction)
	for i := 0; i < len(ats)-1; i++ {
		into[i].SetRect(flip(crop(flipped, ats[i], ats[i+1]), direction))
	}

	return nil
}

func flip(rect Rectangle, direction Direction) Rectangle {
	if direction == Up {
		return Rectangle{rect.X0, rect.Y1, rect.X1, rect.Y0}
	} else if direction == Down {
		return rect
	} else if direction == Left {
		return Rectangle{rect.Y1, rect.X1, rect.Y0, rect.X0}
	} else {
		return Rectangle{rect.Y0, rect.X0, rect.Y1, rect.X1}
	}
}

func crop(rect Rectangle, from, to float64) Rectangle {
	return Rectangle{
		rect.X0,
		rect.Y0 + int(float64(rect.Y1-rect.Y0)*from),
		rect.X1,
		rect.Y0 + int(float64(rect.Y1-rect.Y0)*to),
	}
}
