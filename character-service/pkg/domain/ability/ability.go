package ability

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/skill"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
)

func floorDiv(a, b int) int {
    q := a / b
    if a % b != 0 && a < 0 {
        q--
    }
    return q
}

type Score struct {
	base int
	temp int
	mod  int
}

func NewScore(base int) *Score {
	score := Score{base: base}
	score.UpdateModifier()
	return &score
}

func (s *Score) UpdateModifier() {
	total := s.base + s.temp
	
	s.mod = floorDiv(total-10, 2)
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
	s.UpdateModifier()
	return func() {
		s.temp -= temp
		s.UpdateModifier()
	}
}

func (s *Score) Check() (diceRes, modifier, result int) {
	diceRes = dice.RollDice(dice.D20)
	result = diceRes + s.Modifier()
	return diceRes, s.Modifier(), result
}

func (s *Score) Temp() int {
	return s.temp
}

func (s *Score) Base() int {
	return s.base
}

type Ability struct {
	*Score
	Skills map[types.SkillType]skill.Skill
}

func New(abilityType types.AbilityType) *Ability {
	a := &Ability{
		Score:  NewScore(10),
		Skills: make(map[types.SkillType]skill.Skill),
	}
	for _, skillType := range types.AbilityToSkill[abilityType] {
		a.Skills[skillType] = *skill.NewSkill(0)
	}
	return a
}
