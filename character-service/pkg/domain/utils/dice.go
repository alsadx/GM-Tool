package utils

import "math/rand"

func RollDice(dice int) int {
	return rand.Intn(dice) + 1
}
