package merge_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/nilsbu/arch/pkg/merge"
)

func TestInOrder(t *testing.T) {
	order := merge.InOrder([]int{4, 3, 1, 2})
	expect := [][]int{
		{0, 0, 0, 0},
		{1, 0, 0, 0},
		{2, 0, 0, 0},
		{3, 0, 0, 0},
		{0, 1, 0, 0},
		{1, 1, 0, 0},
		{2, 1, 0, 0},
		{3, 1, 0, 0},
		{0, 2, 0, 0},
		{1, 2, 0, 0},
		{2, 2, 0, 0},
		{3, 2, 0, 0},
		{0, 0, 0, 1},
		{1, 0, 0, 1},
		{2, 0, 0, 1},
		{3, 0, 0, 1},
		{0, 1, 0, 1},
		{1, 1, 0, 1},
		{2, 1, 0, 1},
		{3, 1, 0, 1},
		{0, 2, 0, 1},
		{1, 2, 0, 1},
		{2, 2, 0, 1},
		{3, 2, 0, 1},
	}
	if !reflect.DeepEqual(expect, order) {
		t.Errorf("order wrong:\nexpect: %v\nactual: %v",
			expect, order)
	}
}

func TestRandomOrder(t *testing.T) {
	rand.Seed(42)
	order := merge.RandomOrder([]int{4, 3, 1})
	expect := [][]int{
		{3, 1, 0},
		{1, 2, 0},
		{3, 0, 0},
		{1, 1, 0},
		{0, 2, 0},
		{3, 2, 0},
		{2, 0, 0},
		{2, 2, 0},
		{1, 0, 0},
		{2, 1, 0},
		{0, 0, 0},
		{0, 1, 0},
	}
	if !reflect.DeepEqual(expect, order) {
		t.Errorf("order wrong:\nexpect: %v\nactual: %v",
			expect, order)
	}
}
