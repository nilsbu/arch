package merge

import (
	"errors"
	"testing"

	"github.com/nilsbu/arch/pkg/blueprint"
)

func TestChoices(t *testing.T) {
	type out struct {
		// keys []string
		out []string
	}

	stdResolver := &resolver{"@", map[string][]string{
		"a": {},
		"b": {},
		"T": {"1", "2"},
		"Q": {"3"},
	}}

	for _, c := range []struct {
		name      string
		blueprint string
		root      string
		resolver  *resolver
		outs      []out
		err       error
	}{
		{
			"no root",
			`{}`,
			"Root",
			stdResolver,
			[]out{{
				[]string{"a"},
			}},
			ErrInvalidBlueprint,
		},
		{
			"only root",
			`{"Root":{"@":"a"}}`,
			"Root",
			stdResolver,
			[]out{{
				[]string{"a"},
			}},
			nil,
		},
		{
			"root has two options",
			`{"Root":[{"@":"a"}, {"@":"b"}]}`,
			"Root",
			stdResolver,
			[]out{{
				[]string{"a"},
			}, {
				[]string{"b"},
			}},
			nil,
		},
		{
			"root refers to other",
			`{"Root":"Other","Other":{"@":"a"}}`,
			"Root",
			stdResolver,
			[]out{{
				[]string{"a"},
			}},
			nil,
		},
		{
			"block has parameters",
			`{"Root":{"@":"T","1":{"@":"a"},"2":[{"@":"a"},{"@":"b"}]}}`,
			"Root",
			stdResolver,
			[]out{{
				[]string{"T", "a", "a"},
			}, {
				[]string{"T", "a", "b"},
			}},
			nil,
		},
		{
			"more complex scenario",
			`{"R":["O","P"],"O":{"@":"Q","3":{"@":"a"}},"P":[{"@":"a"},{"@":"T","1":[{"@":"a"},{"@":"b"}],"2":[{"@":"a"},{"@":"b"}]}]}`,
			"R",
			stdResolver,
			[]out{{
				[]string{"Q", "a"},
			}, {
				[]string{"a"},
			}, {
				[]string{"T", "a", "a"},
			}, {
				[]string{"T", "b", "a"},
			}, {
				[]string{"T", "a", "b"},
			}, {
				[]string{"T", "b", "b"},
			}},
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
					for i, out := range c.outs {
						outBp := choices.get(i)
						if len(out.out) != len(outBp) {
							t.Fatalf("expected %v blueprints but got %v", len(out.out), len(outBp))
						}
						for j := range out.out {
							if out.out[j] != outBp[j].Values("@")[0] {
								t.Errorf("value %v, %v doesn't match: expect '%v', actual '%v'",
									i, j, out.out[j], outBp[j].Values("@")[0])
							}
						}
					}
				}
			}
		})
	}
}
