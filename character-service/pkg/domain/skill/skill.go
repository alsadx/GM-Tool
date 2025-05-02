package skill

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/utils"
)

type Skill struct {
	Ability *ability.Score
	Bonus   int
}

func (s *Skill) SetBonus(bonus int) {
	s.Bonus = bonus
}

func (s *Skill) Check() (dice, bonus, result int) {
	dice = utils.RollDice(20)
	bonus = s.Ability.Modifier() + bonus
	return dice, bonus, dice + bonus
}
