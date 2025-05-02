package dice

import "math/rand"

type Dice int

const (
	K6 Dice = 6 + (iota*2)
	K8
	K10
	K12
	K16 Dice = iota*4
	K20
)

func RollDice(dice Dice) int {
	return rand.Intn(int(dice)) + 1
}
