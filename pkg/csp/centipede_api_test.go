package csp_test

import (
	"context"
	"testing"

	"github.com/gnboorse/centipede"
)

func TestCSP(t *testing.T) {
	// some integer variables
	vars := centipede.Variables[int]{
		centipede.NewVariable("A", centipede.IntRange(1, 10)),
		centipede.NewVariable("B", centipede.IntRange(1, 10)),
		centipede.NewVariable("C", centipede.IntRange(1, 10)),
		centipede.NewVariable("D", centipede.IntRange(1, 10)),
		centipede.NewVariable("E", centipede.IntRangeStep(0, 20, 2)), // even numbers < 20
	}

	// numeric constraints
	constraints := centipede.Constraints[int]{
		// using some constraint generators
		centipede.Equals[int]("A", "D"), // A = D
		// here we implement a custom constraint
		centipede.Constraint[int]{Vars: centipede.VariableNames{"A", "E"}, // E = A * 2
			ConstraintFunction: func(variables *centipede.Variables[int]) bool {
				// here we have to use type assertion for numeric methods since
				// Variable.Value is stored as interface{}
				if variables.Find("E").Empty || variables.Find("A").Empty {
					return true
				}
				return variables.Find("E").Value == variables.Find("A").Value*2
			}},
	}
	constraints = append(constraints, centipede.AllUnique[int]("A", "B", "C", "E")...) // A != B != C != E

	// solve the problem
	solver := centipede.NewBackTrackingCSPSolver(vars, constraints)
	// begin := time.Now()

	success, err := solver.Solve(context.TODO()) // run the solution
	if success == false && err != nil {
		t.Error(err)
	}
	// elapsed := time.Since(begin)

	// output results and time elapsed
	// if success {
	// 	fmt.Printf("Found solution in %s\n", elapsed)
	// 	for _, variable := range solver.State.Vars {
	// 		// print out values for each variable
	// 		fmt.Printf("Variable %v = %v\n", variable.Name, variable.Value)
	// 	}
	// } else {
	// 	fmt.Printf("Could not find solution in %s\n", elapsed)
	// }
}
