package skill

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
)

type Name int

type Skill struct {
	Bonus int
}

func NewSkill(bonus int) *Skill {
	s := &Skill{Bonus: bonus}
	return s
}

func (s *Skill) SetBonus(bonus int) {
	s.Bonus = bonus
}

func (s *Skill) Check(modifier int) (diceRes, bonus, result int) {
	diceRes = dice.RollDice(dice.D20)
	bonus = modifier + s.Bonus
	return diceRes, bonus, diceRes + bonus
}
