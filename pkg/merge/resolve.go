package merge

import (
	"github.com/nilsbu/arch/pkg/rule"
)

type Resolver struct {
	Name string
	Keys map[string]rule.Rule
}
