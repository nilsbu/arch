package merge

import (
	"errors"
	"testing"

	"github.com/nilsbu/arch/pkg/blueprint"
	"github.com/nilsbu/arch/pkg/rule"
	tr "github.com/nilsbu/arch/test/rule"
)

func TestChoices(t *testing.T) {
	type out struct {
		c []*out
		n string
		// keys []string
	}

	stdResolver := &Resolver{"@", map[string]rule.Rule{
		"a": &tr.RuleMock{Params: []string{}},
		"b": &tr.RuleMock{Params: []string{}},
		"T": &tr.RuleMock{Params: []string{"1", "2"}},
		"Q": &tr.RuleMock{Params: []string{"3"}},
	}}

	for _, c := range []struct {
		name      string
		blueprint string
		root      string
		resolver  *Resolver
		outs      []*out
		err       error
	}{
		{
			"no root",
			`{}`,
			"Root",
			stdResolver,
			nil,
			ErrInvalidBlueprint,
		},
		{
			"only root",
			`{"Root":{"@":"a"}}`,
			"Root",
			stdResolver,
			[]*out{{c: []*out{{n: "a"}}}},
			nil,
		},
		{
			"root has two options",
			`{"Root":[{"@":"a"}, {"@":"b"}]}`,
			"Root",
			stdResolver,
			[]*out{
				{c: []*out{{n: "a"}}},
				{c: []*out{{n: "b"}}},
			},
			nil,
		},
		{
			"root refers to other",
			`{"Root":"Other","Other":{"@":"a"}}`,
			"Root",
			stdResolver,
			[]*out{{c: []*out{{c: []*out{{n: "a"}}}}}},
			nil,
		},
		{
			"block has parameters",
			`{"Root":{"@":"T","1":{"@":"a"},"2":[{"@":"a"},{"@":"b"}]}}`,
			"Root",
			stdResolver,
			[]*out{
				{c: []*out{{n: "T"}, {c: []*out{{n: "a"}}}, {c: []*out{{n: "a"}, {n: "b"}}}}},
			},
			nil,
		},
		{
			"more children",
			`{"R":["O","P"],"O":{"@":"Q","3":{"@":"a"}},"P":[{"@":"a"},{"@":"T","1":[{"@":"a"},{"@":"b"}],"2":[{"@":"a"},{"@":"b"}]}]}`,
			"R",
			stdResolver,
			[]*out{
				{c: []*out{{c: []*out{{n: "Q"}, {c: []*out{{n: "a"}}}}}}},
				{c: []*out{{c: []*out{{n: "a"}}}}},
				{c: []*out{{c: []*out{{n: "T"}, {c: []*out{{n: "a"}, {n: "b"}}}, {c: []*out{{n: "a"}, {n: "b"}}}}}}},
			},
			nil,
		},
		{
			"options in parameters",
			`{"R":["O","P"],"O":{"@":"Q","3":{"@":"a"}},"P":[{"@":"a"},{"@":"T","1":"P1","P1":[{"@":"a"},{"@":"b"}],"2":"P2","P2":[{"@":"a"},{"@":"b"}]}]}`,
			"R",
			stdResolver,
			[]*out{
				{c: []*out{{c: []*out{{n: "Q"}, {c: []*out{{n: "a"}}}}}}},
				{c: []*out{{c: []*out{{n: "a"}}}}},
				{c: []*out{{c: []*out{{n: "T"}, {c: []*out{{c: []*out{{n: "a"}}}}}, {c: []*out{{c: []*out{{n: "a"}}}}}}}}},
				{c: []*out{{c: []*out{{n: "T"}, {c: []*out{{c: []*out{{n: "b"}}}}}, {c: []*out{{c: []*out{{n: "a"}}}}}}}}},
				{c: []*out{{c: []*out{{n: "T"}, {c: []*out{{c: []*out{{n: "a"}}}}}, {c: []*out{{c: []*out{{n: "b"}}}}}}}}},
				{c: []*out{{c: []*out{{n: "T"}, {c: []*out{{c: []*out{{n: "b"}}}}}, {c: []*out{{c: []*out{{n: "b"}}}}}}}}},
			},
			nil,
		},
		// TODO test errors
	} {
		t.Run(c.name, func(t *testing.T) {
			if bp, err := blueprint.Parse([]byte(c.blueprint)); err != nil {
				t.Fatal("cannot parse blueprint:", err)
			} else {
				if choices, err := calcChoices(bp, c.root, c.resolver); err != nil && c.err == nil {
					t.Error("unexpected error:", err)
				} else if err == nil && c.err != nil {
					t.Error("expected error but none ocurred")
				} else if !errors.Is(err, c.err) {
					t.Error("wrong type of error")
				} else if err == nil {
					n := choices.n()
					if len(c.outs) < n {
						t.Errorf("outputs < choices.n(): %v vs. %v", len(c.outs), n)
					} else if len(c.outs) > n {
						t.Fatalf("outputs > choices.n(): %v vs. %v", len(c.outs), n)
					}

					for i, o := range c.outs {
						outBps := []*bpNode{choices.get(i)}
						outs := []*out{o}
						i := 0
						for len(outBps) > 0 {
							outBp := outBps[0]
							o := outs[0]
							if o.n != "" && outBp.bp == nil {
								t.Error("expected block but none ocurred")
							} else if o.n == "" && outBp.bp != nil {
								t.Error("expected no block but one ocurred")
							} else if o.n != "" && o.n != outBp.bp.Values("@")[0] {
								t.Errorf("value %v doesn't match: expect '%v', actual '%v'",
									i, o.n, outBp.bp.Values("@")[0])
							}

							outBps = append(outBps, outBp.children...)
							outs = append(outs, o.c...)

							if len(outBps) != len(outs) {
								t.Fatalf("step %v, '%v': expected %v children but got %v",
									i, o.n, len(o.c), len(outBp.children))
							}
							outBps = outBps[1:]
							outs = outs[1:]
							i++
						}
					}
				}
			}
		})
	}
}
