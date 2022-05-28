package main

import (
	"fmt"
	"os"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/csp"
	"github.com/nilsbu/arch/pkg/draw"
	"github.com/nilsbu/arch/pkg/merge"
	"github.com/nilsbu/arch/pkg/render"
	"github.com/nilsbu/arch/pkg/rule"
)

func main() {
	bps := make([]*blueprint.Blueprint, len(os.Args)-1)
	for i := range bps {
		if file, err := os.ReadFile(os.Args[i+1]); err != nil {
			fmt.Println(err)
			return
		} else if bps[i], err = blueprint.Parse(file); err != nil {
			fmt.Println(err)
			return
		}
	}

	resolver := &merge.Resolver{
		Name: "@rule",
		Keys: map[string]rule.Rule{
			"House":    rule.House{},
			"Corridor": rule.Corridor{},
			"TwoRooms": rule.TwoRooms{},
			"Room":     rule.Room{},
			"NOP":      rule.NOP{},
		},
	}

	if g, err := merge.Build(bps, &csp.Centipede{}, resolver); err != nil {
		fmt.Println(err)
	} else if tiles, err := draw.Draw(g); err != nil {
		fmt.Println(err)
	} else {
		render.Terminal(os.Stdout, tiles)
	}
}
