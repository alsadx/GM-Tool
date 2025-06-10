package dto

import (
	gen "github.com/alsadx/GM-Tool/character-service/gen/character"
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
	Base  int `json:"base" bson:"base"`
	Bonus int `json:"bonus" bson:"bonus"`
	Mod   int `json:"mod" bson:"-"`
}

func (s *ScoreDTO) toDomain() *ability.Score {
	score := ability.NewScore(s.Base)
	score.SetBonus(s.Bonus)
	return score
}

func (s *ScoreDTO) toProto() *gen.Score {
	return &gen.Score{
		Base:  int32(s.Base),
		Bonus: int32(s.Bonus),
		Mod:   int32(s.Mod),
	}
}

func scoreDTOFromProto(protoScore *gen.Score) *ScoreDTO {
	return &ScoreDTO{
		Base:  int(protoScore.Base),
		Bonus: int(protoScore.Bonus),
		Mod:   int(protoScore.Mod),
	}
}

func scoreDTOFromDomain(scoreDomain *ability.Score) *ScoreDTO {
	return &ScoreDTO{
		Base:  scoreDomain.Base(),
		Bonus: scoreDomain.Bonus(),
		Mod:   scoreDomain.Modifier(),
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

func AbilityDTOFromDomain(abilityDomain *ability.Ability) *AbilityDTO {
	skills := make(map[string]*SkillDTO, len(abilityDomain.Skills))
	for skillType, skillDomain := range abilityDomain.Skills {
		skills[skillType.String()] = SkillDTOFromDomain(skillDomain)
	}
	return &AbilityDTO{
		ScoreDTO: scoreDTOFromDomain(abilityDomain.Score),
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

func HealthDTOFromDomain(healthDomain *health.Health) *HealthDTO {
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
	CurrentExp    int `json:"current_exp" bson:"current_exp"`
	EarnedLevel   int `json:"eanred_lvl" bson:"-"`
	CurrentLevel  int `json:"current_lvl" bson:"current_level"`
	NextThreshold int `json:"next_threshold" bson:"next_threshold"`
}

func (l *LevelDTO) ToProto() *gen.LevelSystem {
	return &gen.LevelSystem{
		CurrentExp:    int32(l.CurrentExp),
		CurrentLvl:    int32(l.CurrentLevel),
		EarnedLvl:     int32(l.EarnedLevel),
		NextThreshold: int32(l.NextThreshold),
	}
}

func (l *LevelDTO) ToDomain() *level.LevelSystem {
	lvl := level.NewLevelSystem().
		WithCurrentExp(l.CurrentExp).
		WithCurrentLvl(l.CurrentLevel).
		WithNextThreshold(l.NextThreshold)
	return lvl
}

func LevelDTOFromProto(genLevel *gen.LevelSystem) *LevelDTO {
	return &LevelDTO{
		CurrentExp:    int(genLevel.CurrentExp),
		EarnedLevel:   int(genLevel.EarnedLvl),
		CurrentLevel:  int(genLevel.CurrentLvl),
		NextThreshold: int(genLevel.NextThreshold),
	}
}

func LevelDTOFromDomain(domainLevel *level.LevelSystem) *LevelDTO {
	return &LevelDTO{
		CurrentExp:    domainLevel.CurrentExp(),
		EarnedLevel:   domainLevel.EarnedLevel(),
		CurrentLevel:  domainLevel.CurrentLevel(),
		NextThreshold: (domainLevel.ExpToNextLevel() + domainLevel.CurrentExp()),
	}
}

type CharacterDTO struct {
	ID    string `json:"id" bson:"_id,omitempty"`
	Owner int    `json:"owner_id" bson:"owner_id"`

	IsKnocked bool `json:"is_knocked" bson:"is_knocked"`

	Name     string `json:"name" bson:"name"`
	Class    string `json:"class" bson:"class"`
	Subclass string `json:"subclass" bson:"subclass"`
	Race     string `json:"race" bson:"race"`

	Lvl    *LevelDTO              `json:"level" bson:"level"`
	Stats  map[string]*AbilityDTO `json:"stats" bson:"stats"`
	Health *HealthDTO             `json:"health" bson:"health"`
}

func (c *CharacterDTO) ToProto() *gen.Character {
	statsProto := make(map[string]*gen.Ability, len(c.Stats))
	for ablType, abil := range c.Stats {
		statsProto[ablType] = abil.ToProto()
	}
	return &gen.Character{
		Id:        c.ID,
		Owner:     int64(c.Owner),
		Name:      c.Name,
		ClassName: c.Class,
		Subclass:  c.Subclass,
		Race:      c.Race,
		IsKnocked: c.IsKnocked,
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

func CharacterDTOFromProto(protoChar *gen.Character) *CharacterDTO {
	stats := make(map[string]*AbilityDTO, len(protoChar.Stats))
	for ablType, abil := range protoChar.Stats {
		stats[ablType] = AbilityDTOFromProto(abil)
	}
	return &CharacterDTO{
		ID:        protoChar.Id,
		Owner:     int(protoChar.Owner),
		Name:      protoChar.Name,
		Class:     protoChar.ClassName,
		Subclass:  protoChar.Subclass,
		Race:      protoChar.Race,
		IsKnocked: protoChar.IsKnocked,
		Lvl:       LevelDTOFromProto(protoChar.Lvl),
		Stats:     stats,
		Health:    HealthDTOFromProto(protoChar.Health),
	}
}

func CharacterDTOFromDomain(domainChar *character.Character) *CharacterDTO {
	stats := make(map[string]*AbilityDTO, len(domainChar.GetStats()))
	for ablType, abil := range domainChar.GetStats() {
		stats[ablType.String()] = AbilityDTOFromDomain(abil)
	}
	return &CharacterDTO{
		ID:        domainChar.ID,
		Owner:     domainChar.Owner,
		Name:      domainChar.Name,
		Class:     domainChar.Class,
		Subclass:  domainChar.Subclass,
		Race:      domainChar.Race,
		IsKnocked: domainChar.IsKnocked(),
		Lvl:       LevelDTOFromDomain(domainChar.GetLvlSystem()),
		Stats:     stats,
		Health:    HealthDTOFromDomain(domainChar.GetHealth()),
	}
}

type UpdateCharDTO struct {
	ID       string
	Name     string
	Class    string
	Subclass string
	Race     string
}

type CharacterInfoDTO struct {
	OwnerID  int
	Name     string
	Class    string
	Subclass string
	Race     string
}

type HPStateDTO struct {
	MaxHP     int
	CurrentHP int
	TempHP    int
	IsKnocked bool
}

type LvlStateDTO struct {
	CurrentLvl   int
	CurrentExp   int
	ExpToNextLvl int
}

type DiceResultDTO struct {
	DiceRes int
	Bonus   int
	result  int
}
