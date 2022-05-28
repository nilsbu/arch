package merge

import "math/rand"

// TODO document

type Shuffle func(ns []int) (choiceIds [][]int)

func InOrder(ns []int) [][]int {
	total := 1
	for _, n := range ns {
		total *= n
	}

	orders := make([][]int, total)
	for i := range orders {
		order := make([]int, len(ns))
		exp := 1
		for j, n := range ns {
			order[j] = (i / exp) % n
			exp *= n
		}
		orders[i] = order
	}

	return orders
}

func RandomOrder(ns []int) [][]int {
	out := InOrder(ns)
	rand.Shuffle(len(out), func(i, j int) {
		out[i], out[j] = out[j], out[i]
	})
	return out
}
