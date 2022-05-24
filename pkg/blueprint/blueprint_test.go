package blueprint_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/nilsbu/arch/pkg/blueprint"
)

type ksvs struct {
	properties []string
	values     []string
}

func TestBlueprint(t *testing.T) {
	// access to children of parents is tested in TestBlueprintSibling
	for _, c := range []struct {
		name string
		json string

		ok         bool
		properties []string
		children   []string
		checks     []ksvs
	}{
		{
			"empty JSON",
			`{}`,
			true,
			[]string{},
			[]string{},
			[]ksvs{},
		},
		{
			"invalid JSON",
			`{"}`,
			false, nil, nil, nil,
		},
		{
			"bool not allowed as value",
			`{"k":true}`,
			false, nil, nil, nil,
		},
		{
			"int not allowed as value",
			`{"k":69}`,
			false, nil, nil, nil,
		},
		{
			"one value",
			`{"k":"v"}`,
			true,
			[]string{"k"},
			[]string{},
			[]ksvs{{[]string{"k"}, []string{"v"}}},
		},
		{
			"two properties",
			`{"k":"v","kk":"v"}`,
			true,
			[]string{"k", "kk"},
			[]string{},
			[]ksvs{{[]string{"k"}, []string{"v"}},
				{[]string{"kk"}, []string{"v"}}},
		},
		{
			"list of values",
			`{"k":["v", "w"]}`,
			true,
			[]string{"k"},
			[]string{},
			[]ksvs{{[]string{"k"}, []string{"v", "w"}}},
		},
		{
			"lists in lists of values",
			`{"k":["v", ["w"]]}`,
			true,
			[]string{"k"},
			[]string{},
			[]ksvs{{[]string{"k"}, []string{"v", "w"}}},
		},
		{
			"empty value list",
			`{"k":[]}`,
			true,
			[]string{"k"},
			[]string{},
			[]ksvs{{[]string{"k"}, []string{}}},
		},
		{
			"block as value",
			`{"k":{"kk":"v"}}`,
			true,
			[]string{"k"},
			[]string{"*k0"},
			[]ksvs{{[]string{"k"}, []string{"*k0"}},
				{[]string{"*k0", "kk"}, []string{"v"}}},
		},
		{
			"error in child",
			`{"k":{"kk":32}}`,
			false, nil, nil, nil,
		},
		{
			"error in list",
			`{"k":["kk",32]}`,
			false, nil, nil, nil,
		},
		// now let's test access semantics in more detail
		{
			"cannot access child properties",
			`{"k":{"o":"inner"}}`,
			true,
			[]string{"k"},
			[]string{"*k0"},
			[]ksvs{{[]string{"*k0"}, []string{}}},
		},
		{
			"local is preferred",
			`{"o":"outer","k":{"o":"inner"}}`,
			true,
			[]string{"k", "o"},
			[]string{"*k0"},
			[]ksvs{{[]string{"*k0", "o"}, []string{"inner"}}},
		},
		{
			"outer is chosen when local isn't available",
			`{"x":"outer","k":{}}`,
			true,
			[]string{"k", "x"},
			[]string{"*k0"},
			[]ksvs{{[]string{"*k0", "x"}, []string{"outer"}}},
		},
		{
			"choose closest outer",
			`{"x":"far out","k":{"x":"outer","m":{}}}`,
			true,
			[]string{"k", "x"},
			[]string{"*k0"},
			[]ksvs{{[]string{"*k0", "*m0", "x"}, []string{"outer"}}},
		},
		{
			"various blocks",
			`{"k":[{},{"x":"v"}],"asdf":{}}`,
			true,
			[]string{"k", "asdf"},
			[]string{"*k0", "*k1", "*asdf0"},
			[]ksvs{},
		},
		{
			"property with empty list is not omitted",
			`{"k":[]}`,
			true,
			[]string{"k"},
			[]string{},
			[]ksvs{{[]string{"k"}, []string{}}},
		},
		{
			"second block available",
			`{"k":[{},{"x":"v"}]}`,
			true,
			[]string{"k"},
			[]string{"*k0", "*k1"},
			[]ksvs{{[]string{"*k1", "x"}, []string{"v"}}},
		},
		{
			"choose closest when identical block names exist",
			`{"x":{},"k":{"x":{"v":"1"}}}`,
			true,
			[]string{"k", "x"},
			[]string{"*k0", "*x0"},
			[]ksvs{{[]string{"*k0", "*x0", "v"}, []string{"1"}}},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			bp, err := blueprint.Parse([]byte(c.json))

			if c.ok && err != nil {
				t.Error("didn't expect error:", err)
			} else if !c.ok && err == nil {
				t.Error("expected error but none occurred")
			} else if c.ok {
				checkBlueprint(t, bp, c.properties, c.children, c.checks)
			}
		})
	}
}

func checkBlueprint(
	t *testing.T,
	bp blueprint.Blueprint,
	properties []string,
	children []string,
	checks []ksvs) {

	{
		sort.Strings(properties)
		actual := bp.Properties()
		sort.Strings(actual)
		if !reflect.DeepEqual(properties, actual) {
			t.Fatalf("properties don't match expected: %v vs. %v",
				properties, actual)
		}
	}
	{
		sort.Strings(children)
		actual := bp.Children()
		sort.Strings(actual)
		if !reflect.DeepEqual(children, actual) {
			t.Fatalf("children don't match expected: %v vs. %v",
				children, actual)
		}
	}

	for _, check := range checks {
		bpTmp := bp
		for i, prop := range check.properties {
			if prop[0] == '*' {
				bpTmp = bpTmp.Child(prop)
				if bpTmp == nil {
					t.Fatalf("child %v doesn't exist", check.properties[:i+1])
				}
			} else {
				if !reflect.DeepEqual(check.values, bpTmp.Values(prop)) {
					t.Errorf("for checked properties %v: expected %v but got %v",
						check.properties, check.values, bpTmp.Values(prop))
				}
			}
		}
	}
}

func TestBlueprintSibling(t *testing.T) {
	bp, err := blueprint.Parse([]byte(`{"x":{"a":"b"},"k":{}}`))
	if err != nil {
		t.Fatal("unexptected error", err)
	}

	child := bp.Child("*k0")
	if child == nil {
		t.Fatal("child mustn't be nil")
	}

	sibling := child.Child("*x0")
	if sibling == nil {
		t.Fatal("sibling mustn't be nil")
	}

	value := sibling.Values("a")
	if !reflect.DeepEqual([]string{"b"}, value) {
		t.Errorf("value 'a' doesn't match: expect [b] but got %v", value)
	}
}

func TestBlueprintFalseProperty(t *testing.T) {
	bp, _ := blueprint.Parse([]byte("{}"))
	if bp.Values("IDontExist") != nil {
		t.Errorf("Values() shouldn't have returned %v",
			bp.Values("IDontExist"))
	}

}

func TestBlueprintFalseChild(t *testing.T) {
	bp, _ := blueprint.Parse([]byte("{}"))
	if bp.Child("*IDontExist") != nil {
		t.Errorf("Child() shouldn't have returned %v",
			bp.Child("*IDontExist"))
	}

}
