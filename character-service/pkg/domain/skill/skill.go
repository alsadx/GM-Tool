package skill

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
)

type Skill struct {
	Ability *ability.Score
	Bonus   int
}

func (s *Skill) SetBonus(bonus int) {
	s.Bonus = bonus
}

func (s *Skill) Check() (diceRes, bonus, result int) {
	diceRes = dice.RollDice(dice.D20)
	bonus = s.Ability.Modifier() + bonus
	return diceRes, bonus, diceRes + bonus
}
