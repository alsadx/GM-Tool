package dto

import (
	"github.com/alsadx/GM-Tool/character-service/gen"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/level"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/skill"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
)

type SkillDTO struct {
	Bonus int `json:"bonus" bson:"bonus"`
}

func (s *SkillDTO) ToDomain() *skill.Skill {
	return &skill.Skill{
		Bonus: s.Bonus,
	}
}

func (s *SkillDTO) ToProto() *gen.Skill {
	return &gen.Skill{
		Bonus: int32(s.Bonus),
	}
}

func SkillDTOFromDomain(s *skill.Skill) *SkillDTO {
	return &SkillDTO{
		Bonus: s.Bonus,
	}
}

func SkillDTOFromProto(protoSkill *gen.Skill) *SkillDTO {
	return &SkillDTO{Bonus: int(protoSkill.Bonus)}
}

type ScoreDTO struct {
	Base int `json:"base" bson:"base"`
	Temp int `json:"temp" bson:"temp"`
	Mod  int `json:"mod" bson:"mod"`
}

func (s *ScoreDTO) toDomain() *ability.Score {
	score := ability.NewScore(s.Base)
	score.AddTemp(s.Temp)
	return score
}

func (s *ScoreDTO) toProto() *gen.Score {
	return &gen.Score{
		Base: int32(s.Base),
		Temp: int32(s.Temp),
		Mod:  int32(s.Mod),
	}
}

func scoreDTOFromProto(protoScore *gen.Score) *ScoreDTO {
	return &ScoreDTO{
		Base: int(protoScore.Base),
		Temp: int(protoScore.Temp),
		Mod:  int(protoScore.Mod),
	}
}

type AbilityDTO struct {
	*ScoreDTO
	Skills map[string]*SkillDTO `json:"skills" bson:"skills"`
}

func (a *AbilityDTO) ToDomain() (*ability.Ability, error) {
	skills := make(map[types.SkillType]*skill.Skill, len(a.Skills))
	for skillName, skillDTO := range a.Skills {
		skillID, err := types.ParseSkillType(skillName)
		if err != nil {
			return nil, err
		}
		skills[skillID] = skillDTO.ToDomain()
	}
	return &ability.Ability{
		Score:  a.toDomain(),
		Skills: skills,
	}, nil
}

func (a *AbilityDTO) ToProto() *gen.Ability {
	protoSkills := make(map[string]*gen.Skill, len(a.Skills))
	for skillType, skill := range a.Skills {
		protoSkills[skillType] = skill.ToProto()
	}
	return &gen.Ability{
		Score:  a.toProto(),
		Skills: protoSkills,
	}
}

func AbilityDTOFromProto(abilityProto *gen.Ability) *AbilityDTO {
	skills := make(map[string]*SkillDTO, len(abilityProto.Skills))
	for skillType, skillProto := range abilityProto.Skills {
		skills[skillType] = SkillDTOFromProto(skillProto)
	}
	return &AbilityDTO{
		ScoreDTO: scoreDTOFromProto(abilityProto.Score),
		Skills:   skills,
	}
}

type AmountDTO struct {
	MaxAvailable int `json:"max_available" bson:"max_available"`
	Available    int `json:"available" bson:"available"`
}

func (a *AmountDTO) toProto() *gen.Amount {
	return &gen.Amount{
		MaxAvailable: int32(a.MaxAvailable),
		Available:    int32(a.Available),
	}
}

func (a *AmountDTO) toDomain() *health.Amount {
	return &health.Amount{
		MaxAvailable: a.MaxAvailable,
		Available:    a.Available,
	}
}

func amountDTOFromProto(amountProto *gen.Amount) *AmountDTO {
	return &AmountDTO{
		MaxAvailable: int(amountProto.MaxAvailable),
		Available:    int(amountProto.Available),
	}
}

func amountDTOFromDomain(amountDomain *health.Amount) *AmountDTO {
	return &AmountDTO{
		MaxAvailable: amountDomain.MaxAvailable,
		Available:    amountDomain.Available,
	}
}

type HealthDTO struct {
	CurrentHP int                      `json:"current_hp" bson:"current_hp"`
	MaxHP     int                      `json:"max_hp" bson:"max_hp"`
	TempHP    int                      `json:"temp_hp" bson:"temp_hp"`
	HitDice   map[dice.Dice]*AmountDTO `json:"hit_dice" bson:"hit_dice"`
}

func (h *HealthDTO) ToProto() *gen.HealthPoint {
	hitDiceProto := make(map[int32]*gen.Amount, len(h.HitDice))
	for diceType, amount := range h.HitDice {
		hitDiceProto[int32(diceType)] = amount.toProto()
	}
	return &gen.HealthPoint{
		CurrentHp: int32(h.CurrentHP),
		MaxHp:     int32(h.MaxHP),
		TempHp:    int32(h.TempHP),
		HitDice:   hitDiceProto,
	}
}

func (h *HealthDTO) ToDomain() *health.Health {
	hitDiceDomain := make(map[dice.Dice]*health.Amount, len(h.HitDice))
	for diceType, amount := range h.HitDice {
		hitDiceDomain[diceType] = amount.toDomain()
	}
	return &health.Health{
		CurrentHP: h.CurrentHP,
		MaxHP:     h.MaxHP,
		TempHP:    h.TempHP,
		HitDice:   hitDiceDomain,
	}
}

func HealthDTOFromProto(healthProto *gen.HealthPoint) *HealthDTO {
	hitDice := make(map[dice.Dice]*AmountDTO, len(healthProto.HitDice))
	for diceType, amount := range healthProto.HitDice {
		hitDice[dice.Dice(diceType)] = amountDTOFromProto(amount)
	}
	return &HealthDTO{
		CurrentHP: int(healthProto.CurrentHp),
		MaxHP:     int(healthProto.MaxHp),
		TempHP:    int(healthProto.TempHp),
		HitDice:   hitDice,
	}
}

func HealthDTOFromDomain(healthDomain health.Health) *HealthDTO {
	hitDice := make(map[dice.Dice]*AmountDTO, len(healthDomain.HitDice))
	for diceType, amount := range healthDomain.HitDice {
		hitDice[dice.Dice(diceType)] = amountDTOFromDomain(amount)
	}
	return &HealthDTO{
		CurrentHP: int(healthDomain.CurrentHP),
		MaxHP:     int(healthDomain.MaxHP),
		TempHP:    int(healthDomain.TempHP),
		HitDice:   hitDice,
	}
}

type LevelDTO struct {
	CurrentExp int `json:"current_exp" bson:"current_exp"`
}

func (l *LevelDTO) ToProto() *gen.LevelSystem {
	return &gen.LevelSystem{
		CurrentExp: int32(l.CurrentExp),
	}
}

func (l *LevelDTO) ToDomain() *level.LevelSystem {
	lvl := level.NewLevelSystem()
	lvl.AddExp(l.CurrentExp)
	return lvl
}

func LevelDTOFromProto(genLevel *gen.LevelSystem) *LevelDTO {
	return &LevelDTO{
		CurrentExp: int(genLevel.CurrentExp),
	}
}

func LevelDTOFromDomain(domainLevel *level.LevelSystem) *LevelDTO {
	return &LevelDTO{
		CurrentExp: domainLevel.CurrentExp(),
	}
}

type CharacterDTO struct {
	ID    int `json:"id" bson:"id"`
	Owner int `json:"id_owner" bson:"id_owner"`

	IsKnocked bool `json:"is_knocked" bson:"is_knocked"`

	Name     string `json:"name" bson:"name"`
	Class    string `json:"class" bson:"class"`
	Subclass string `json:"subclass" bson:"subclass"`
	Race     string `json:"race" bson:"race"`

	Lvl    LevelDTO               `json:"level" bson:"level"`
	Stats  map[string]*AbilityDTO `json:"stats" bson:"stats"`
	Health *HealthDTO             `json:"health" bson:"health"`
}

func (c *CharacterDTO) ToProto() *gen.Character {
	statsProto := make(map[string]*gen.Ability, len(c.Stats))
	for ablType, abil := range c.Stats {
		statsProto[ablType] = abil.ToProto()
	}
	return &gen.Character{
		Id:        int64(c.ID),
		Owner:     int64(c.Owner),
		Name:      c.Name,
		ClassName: c.Class,
		Subclass:  c.Subclass,
		Race:      c.Race,
		Lvl:       c.Lvl.ToProto(),
		Stats:     statsProto,
		Health:    c.Health.ToProto(),
	}
}

func (c *CharacterDTO) ToDomain() (*character.Character, error) {
	statsDomain := make(map[types.AbilityType]*ability.Ability, len(c.Stats))
	for ablType, abil := range c.Stats {
		ablID, err := types.ParseAbilityType(ablType)
		if err != nil {
			return nil, err
		}
		abilDomain, err := abil.ToDomain()
		if err != nil {
			return nil, err
		}
		statsDomain[ablID] = abilDomain
	}
	char := &character.Character{
		ID:       c.ID,
		Owner:    c.Owner,
		Name:     c.Name,
		Class:    c.Class,
		Subclass: c.Subclass,
		Race:     c.Race,
	}
	char.WithLvl(c.Lvl.ToDomain()).WithStats(statsDomain).WithHealth(*c.Health.ToDomain())
	return char, nil
}

func FromProto(protoChar *gen.Character) *CharacterDTO {
	stats := make(map[string]*AbilityDTO, len(protoChar.Stats))
	for ablType, abil := range protoChar.Stats {
		stats[ablType] = AbilityDTOFromProto(abil)
	}
	return &CharacterDTO{
		ID:       int(protoChar.Id),
		Owner:    int(protoChar.Owner),
		Name:     protoChar.Name,
		Class:    protoChar.ClassName,
		Subclass: protoChar.Subclass,
		Race:     protoChar.Race,
		Lvl:      *LevelDTOFromProto(protoChar.Lvl),
		Stats:    stats,
		Health:   HealthDTOFromProto(protoChar.Health),
	}
}
