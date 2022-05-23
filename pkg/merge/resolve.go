package merge

import (
	"errors"
)

var ErrUnknownKey = errors.New("unknown key")

type resolver struct {
	name string
	keys map[string][]string
}
