package dice

import "math/rand"

type Dice int

const (
	D6 Dice = 6 + (iota * 2)
	D8
	D10
	D12
	D16 Dice = iota * 4
	D20
)

func RollDice(dice Dice) int {
	return rand.Intn(int(dice)) + 1
}

func MultiRollDice(dice Dice, amount int) []int {
	results := make([]int, amount)
	for i := range results {
		results[i] = RollDice(dice)
	}
	return results
}
