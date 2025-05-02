package ability

import "github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"

type Ability int

const (
	Strength Ability = iota
	Dexterity
	Constitution
	Intelligence
	Wisdom
	Charisma
)

type Score struct {
	base int
	temp int
	mod  int
}

func (s *Score) UpdateModifier() {
	total := s.base + s.temp
	s.mod = (total - 10) / 2
}

func (s *Score) Modifier() int {
	return s.mod
}

func (s *Score) SetBase(base int) {
	s.base = base
	s.UpdateModifier()
}

func (s *Score) AddTemp(temp int) (removeTemp func()) {
	s.temp += temp
	return func() { s.temp -= temp }
}

func (s *Score) Check() (diceRes, modifier, result int) {
	diceRes = dice.RollDice(dice.K20)
	return diceRes, s.Modifier(), result
}
