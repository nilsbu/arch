package rule_test

import (
	"errors"
	"reflect"
	"testing"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/graph"
	"github.com/nilsbu/arch/test/rule"
)

func TestMockRuleParams(t *testing.T) {
	expect := []string{"a", "b"}
	rm := &rule.RuleMock{
		Params: expect,
	}

	if !reflect.DeepEqual(expect, rm.ChildParams()) {
		t.Errorf("wrong params")
	}
}

func TestMockRuleNoPrep(t *testing.T) {
	data := ""
	for _, c := range []struct {
		name string
		f    func(
			g *graph.Graph,
			nidx graph.NodeIndex,
			children map[string][]graph.NodeIndex,
			bp blueprint.Blueprint,
		) error
		ok   bool
		data string
	}{
		{
			"no f",
			nil,
			true,
			"",
		},
		{
			"error",
			func(
				g *graph.Graph,
				nidx graph.NodeIndex,
				children map[string][]graph.NodeIndex,
				bp blueprint.Blueprint,
			) error {
				return errors.New("bah")
			},
			false,
			"",
		},
		{
			"data changed",
			func(
				g *graph.Graph,
				nidx graph.NodeIndex,
				children map[string][]graph.NodeIndex,
				bp blueprint.Blueprint,
			) error {
				data = "booo"
				return nil
			},
			true,
			"booo",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			data = ""
			rm := &rule.RuleMock{
				Params: []string{"s"},
				Prep:   c.f,
			}

			bp, _ := blueprint.Parse([]byte("{}"))
			err := rm.PrepareGraph(
				graph.New(nil),
				graph.NodeIndex{},
				map[string][]graph.NodeIndex{},
				bp)

			if c.ok && err != nil {
				t.Error("unexpected error:", err)
			} else if !c.ok && err == nil {
				t.Error("expected error but none ocurred")
			} else if err == nil {
				if c.data != data {
					t.Errorf("expected '%v' but got '%v'",
						c.data, data)
				}
			}
		})
	}

}
