package merge

import (
	"errors"

	"github.com/nilsbu/arch/pkg/rule"
)

var ErrUnknownKey = errors.New("unknown key")

type Resolver struct {
	Name string
	Keys map[string]rule.Rule
}
