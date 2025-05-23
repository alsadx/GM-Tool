package ability

import (
	"github.com/alsadx/GM-Tool/character-service/gen"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/skill"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
)

func floorDiv(a, b int) int {
	q := a / b
	if a%b != 0 && a < 0 {
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

func (s *Score) toProto() *gen.Score {
	return &gen.Score{
		Base: int32(s.base),
		Temp: int32(s.temp),
		Mod:  int32(s.mod),
	}
}

func fromProto(protoScore *gen.Score) *Score {
	return &Score{
		base: int(protoScore.Base),
		temp: int(protoScore.Temp),
		mod:  int(protoScore.Mod),
	}
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
	Skills map[types.SkillType]*skill.Skill
}

func New(abilityType types.AbilityType) *Ability {
	a := &Ability{
		Score:  NewScore(10),
		Skills: make(map[types.SkillType]*skill.Skill),
	}
	for _, skillType := range types.AbilityToSkill[abilityType] {
		a.Skills[skillType] = skill.NewSkill(0)
	}
	return a
}

func (a *Ability) ToProto() *gen.Ability {
	protoSkills := make(map[int32]*gen.Skill, len(a.Skills))
	for skillType, skill := range a.Skills {
		protoSkills[int32(skillType)] = skill.ToProto()
	}
	return &gen.Ability{
		Score:  a.toProto(),
		Skills: protoSkills,
	}
}

func FromProto(abilityProto *gen.Ability) *Ability {
	skills := make(map[types.SkillType]*skill.Skill, len(abilityProto.Skills))
	for skillType, skillProto := range abilityProto.Skills {
		skills[types.SkillType(skillType)] = skill.FromProto(skillProto)
	}
	return &Ability{
		Score:  fromProto(abilityProto.Score),
		Skills: skills,
	}
}

func (a *Ability) CheckSkill(skillType types.SkillType) (diceRes, bonus, result int) {
	skill := a.Skills[skillType]
	diceRes, bonus, result = skill.Check(a.Modifier())
	return diceRes, bonus, result
}

func NewStats() map[types.AbilityType]*Ability {
	stats := make(map[types.AbilityType]*Ability, 6)
	abilityTypes := []types.AbilityType{types.Strength, types.Dexterity, types.Intelligence, types.Wisdom, types.Charisma, types.Constitution}
	for _, at := range abilityTypes {
		stats[at] = New(at)
	}
	return stats
}
