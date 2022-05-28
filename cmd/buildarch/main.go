package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/csp"
	"github.com/nilsbu/arch/pkg/draw"
	"github.com/nilsbu/arch/pkg/merge"
	"github.com/nilsbu/arch/pkg/render"
	"github.com/nilsbu/arch/pkg/rule"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	if err := buildArchitecture(); err != nil {
		fmt.Println(err)
	}
}

func buildArchitecture() error {
	bps := make([]*blueprint.Blueprint, len(os.Args)-1)
	for i := range bps {
		if file, err := os.ReadFile(os.Args[i+1]); err != nil {
			return err
		} else if bps[i], err = blueprint.Parse(file); err != nil {
			return err
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

	if g, err := merge.Build(bps, &csp.Centipede{}, resolver, merge.RandomOrder); err != nil {
		return err
	} else if tiles, err := draw.Draw(g); err != nil {
		return err
	} else {
		render.Terminal(os.Stdout, tiles)
		return nil
	}
}
