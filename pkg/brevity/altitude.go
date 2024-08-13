package brevity

import (
	"math"
	"slices"

	"github.com/martinlindhe/unit"
)

// Stack represents a single layer of an altitude STACK.
type Stack struct {
	Altitude unit.Length
	Count    int
}

// Stacks creates altitude STACKS from altitudes.
func Stacks(a ...unit.Length) []Stack {
	for i, alt := range a {
		a[i] = unit.Length(math.Round(alt.Feet()/1000)) * 1000 * unit.Foot
	}
	// reverse sort
	slices.SortFunc(a, func(i, j unit.Length) int {
		if i < j {
			return -1
		}
		if i > j {
			return 1
		}
		return 0
	})

	stacks := []Stack{}
	for i := len(a) - 1; i >= 0; i-- {
		if len(stacks) == 0 {
			stacks = append(stacks, Stack{Altitude: a[i], Count: 1})
		} else {
			j := len(stacks) - 1
			highest := stacks[j].Altitude
			if a[i] <= highest-9900*unit.Foot {
				stacks = append(stacks, Stack{Altitude: a[i], Count: 1})
			} else {
				stacks[j].Count++
			}
		}
	}

	return stacks
}