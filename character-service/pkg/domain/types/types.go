package types

import "errors"

//go:generate go run golang.org/x/tools/cmd/stringer -type=AbilityType,SkillType

var (
	ErrInvalidAbilityType = errors.New("invalid AbilityType")
	ErrInvalidSkillType   = errors.New("invalid SkillType")
)

type AbilityType int

const (
	Strength AbilityType = iota + 1
	Dexterity
	Constitution
	Intelligence
	Wisdom
	Charisma
)

type SkillType int

const (
	Athletics SkillType = iota + 1
	Acrobatics
	SleightOfHand
	Stealth
	Arcana
	History
	Investigation
	Nature
	Religion
	AnimalHandling
	Insight
	Medicine
	Perception
	Survival
	Deception
	Intimidation
	Performance
	Persuasion
)

var SkillToAbility = map[SkillType]AbilityType{
	Athletics:      Strength,
	Acrobatics:     Dexterity,
	SleightOfHand:  Dexterity,
	Stealth:        Dexterity,
	Arcana:         Intelligence,
	History:        Intelligence,
	Investigation:  Intelligence,
	Nature:         Intelligence,
	Religion:       Intelligence,
	AnimalHandling: Wisdom,
	Insight:        Wisdom,
	Medicine:       Wisdom,
	Perception:     Wisdom,
	Survival:       Wisdom,
	Deception:      Charisma,
	Intimidation:   Charisma,
	Performance:    Charisma,
	Persuasion:     Charisma,
}

var AbilityToSkill = map[AbilityType][]SkillType{
	Strength:     {Athletics},
	Dexterity:    {Acrobatics, SleightOfHand, Stealth},
	Intelligence: {Arcana, History, Investigation, Nature, Religion},
	Wisdom:       {AnimalHandling, Insight, Medicine, Perception, Survival},
	Charisma:     {Deception, Intimidation, Performance, Persuasion},
	Constitution: {},
}

var (
	abilityTypes = []AbilityType{
		Strength,
		Dexterity,
		Constitution,
		Intelligence,
		Wisdom,
		Charisma,
	}
	abilityTypeMap = make(map[string]AbilityType)

	skillTypes = []SkillType{
		Athletics,
		Acrobatics,
		SleightOfHand,
		Stealth,
		Arcana,
		History,
		Investigation,
		Nature,
		Religion,
		AnimalHandling,
		Insight,
		Medicine,
		Perception,
		Survival,
		Deception,
		Intimidation,
		Performance,
		Persuasion,
	}
	skillTypeMap = make(map[string]SkillType)
)

func init() {
	for _, at := range abilityTypes {
		abilityTypeMap[at.String()] = at
	}

	for _, st := range skillTypes {
		skillTypeMap[st.String()] = st
	}
}

func ParseAbilityType(s string) (AbilityType, error) {
	if at, ok := abilityTypeMap[s]; ok {
		return at, nil
	}
	return 0, ErrInvalidAbilityType
}

func ParseSkillType(s string) (SkillType, error) {
	if st, ok := skillTypeMap[s]; ok {
		return st, nil
	}
	return 0, ErrInvalidSkillType
}
