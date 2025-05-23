package skill

import (
	"github.com/alsadx/GM-Tool/character-service/gen"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
)

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

func (s *Skill) ToProto() *gen.Skill {
	return &gen.Skill{
		Bonus: int32(s.Bonus),
	}
}

func FromProto(protoSkill *gen.Skill) *Skill {
	return &Skill{Bonus: int(protoSkill.Bonus)}
}
