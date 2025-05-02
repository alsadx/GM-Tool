package ability

import "github.com/alsadx/GM-Tool/character-service/pkg/domain/utils"

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

func (s *Score) Check() (dice, modifier, result int) {
	dice = utils.RollDice(20)
	return dice, s.Modifier(), result
}
