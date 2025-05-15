package types

//go:generate go run golang.org/x/tools/cmd/stringer -type=AbilityType,SkillType

type AbilityType int

const (
	Strength AbilityType = iota
	Dexterity
	Constitution
	Intelligence
	Wisdom
	Charisma
)

type SkillType int

const (
	Athletics SkillType = iota
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
