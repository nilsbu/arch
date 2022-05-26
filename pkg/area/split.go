package area

import (
	"errors"
	"fmt"
)

var ErrInvalidSplit = errors.New("invalid split")

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
