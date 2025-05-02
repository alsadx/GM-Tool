package domain

import (
	"errors"
	"sync"
)

// Ability represents a characteristic of the character
type Ability int

const (
	Strength Ability = iota
	Dexterity
	Constitution
	Intelligence
	Wisdom
	Charisma
)

type AbilityScore struct {
	base  int
	temp  int
	mod   int
	mutex sync.RWMutex
}

func (a *AbilityScore) Modifier() int {
	a.mutex.RLock()
	defer a.mutex.RUnlock()
	return a.mod
}

func (a *AbilityScore) updateModifier() {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	total := a.base + a.temp
	a.mod = (total - 10) / 2
}

// AC system
type ArmorType int

const (
	NoArmor ArmorType = iota
	LightArmor
	MediumArmor
	HeavyArmor
	Shield
)

type Armor struct {
	Name           string
	Type           ArmorType
	BaseAC         int
	MaxDex         int
	StrengthReq    int
	StealthPenalty bool
	Disadvantage   []Ability
}

type SavingThrow struct {
	Ability     Ability
	Proficiency bool
	Bonus       int
}

type Character struct {
	abilities map[Ability]*AbilityScore
	armor     *Armor
	shield    *Armor

	ac int

	skills       map[string]*Skill
	savingThrows map[Ability]*SavingThrow

	proficiencyBonus int
	level            int

	movementPenalty int
	activeEffects   map[string]Effect

	maxHP     int
	currentHP int
	tempHP    int
	hitDice   string

	mutex sync.RWMutex
}

func (c *Character) SetMaxHP(hp int) {
	c.maxHP = max(hp, 1)
	c.currentHP = min(c.currentHP, c.maxHP)
}

func (c *Character) TakeDamage(amount int) {
	amount = max(amount, 0)

	if c.tempHP > 0 {
		damageToTemp := min(amount, c.tempHP)
		c.tempHP -= damageToTemp
		amount -= damageToTemp
	}

	if amount > 0 {
		c.currentHP = max(c.currentHP-amount, 0)
	}
}

// Геттеры с блокировкой
func (c *Character) AbilityMod(a Ability) int {
	return c.abilities[a].Modifier()
}

func (c *Character) SetBaseAbility(a Ability, value int) {
	c.abilities[a].base = clamp(value, 1, 30)
	c.abilities[a].updateModifier()
	c.checkArmorRequirements()
	c.recalculateAC()
}

func (c *Character) EquipArmor(armor Armor) error {
	if armor.Type == Shield {
		return errors.New("use EquipShield for shields")
	}

	c.armor = &armor
	c.checkArmorRequirements()
	c.recalculateAC()
	return nil
}

func (c *Character) checkArmorRequirements() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if c.armor == nil {
		c.movementPenalty = 0
		return
	}

	if c.armor.Type == HeavyArmor {
		strValue := c.abilities[Strength].base + c.abilities[Strength].temp
		if strValue < c.armor.StrengthReq {
			c.movementPenalty = 10
		} else {
			c.movementPenalty = 0
		}
	} else {
		c.movementPenalty = 0
	}
}

func (c *Character) SetTempAbility(a Ability, value int) {
	c.abilities[a].mutex.Lock()
	c.abilities[a].temp = clamp(value, 0, 30)
	c.abilities[a].mutex.Unlock()

	c.abilities[a].updateModifier()
	c.checkArmorRequirements()
	c.recalculateAC()
	c.updateSkillsForAbility(a)
	c.UpdateSavingThrows()
}

func (c *Character) updateSkillsForAbility(a Ability) {
	for _, skill := range c.skills {
		if skill.ability == a {
			skill.calculateBonus(c.AbilityMod(a), c.proficiencyBonus)
		}
	}
}

func (c *Character) recalculateAC() {
	base := 10
	var dexMod int

	c.mutex.Lock()
	defer c.mutex.Unlock()

	dexMod = c.abilities[Dexterity].mod

	if c.armor != nil {
		switch c.armor.Type {
		case LightArmor:
			base = c.armor.BaseAC + dexMod
		case MediumArmor:
			base = c.armor.BaseAC + min(dexMod, c.armor.MaxDex)
		case HeavyArmor:
			base = c.armor.BaseAC
		}
	}

	if c.shield != nil {
		base += c.shield.BaseAC
	}

	for _, eff := range c.activeEffects {
		base += eff.ACBonus
	}

	c.ac = max(base, 0)
}

func (c *Character) SetLevel(level int) {
	c.level = clamp(level, 1, 20)
	c.proficiencyBonus = 2 + (level-1)/4
	c.UpdateAllSkills()
	c.UpdateSavingThrows()
	c.recalculateAC()
}

func (c *Character) UpdateAllSkills() {
	for _, skill := range c.skills {
		skill.calculateBonus(
			c.AbilityMod(skill.ability),
			c.proficiencyBonus,
		)
	}
}

func (c *Character) UpdateSavingThrows() {
	for _, st := range c.savingThrows {
		st.Bonus = c.AbilityMod(st.Ability)
		if st.Proficiency {
			st.Bonus += c.proficiencyBonus
		}
	}
}

// Skills system
type Skill struct {
	ability    Ability
	proficient bool
	expertise  bool
	bonus      int
}

func (s *Skill) calculateBonus(abilityMod, proficiencyBonus int) {
	s.bonus = abilityMod
	if s.proficient {
		s.bonus += proficiencyBonus
	}
	if s.expertise {
		s.bonus += proficiencyBonus
	}
}

func (c *Character) UpdateSkill(skillName string, proficient, expertise bool) {
	skill := c.skills[skillName]
	skill.proficient = proficient
	skill.expertise = expertise
	skill.calculateBonus(
		c.AbilityMod(skill.ability),
		c.proficiencyBonus,
	)
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

type Effect struct {
	Name          string
	ACBonus       int
	AbilityScores map[Ability]int
	SpeedMod      int
	Duration      int
}

func (c *Character) AddEffect(e Effect) {
	for ab, val := range e.AbilityScores {
		c.abilities[ab].temp += val
		c.abilities[ab].updateModifier()
	}

	c.activeEffects[e.Name] = e
	c.recalculateAC()
	c.checkArmorRequirements()
	c.UpdateAllSkills()
	c.UpdateSavingThrows()
}

func (c *Character) RemoveEffect(name string) {
	delete(c.activeEffects, name)
	c.recalculateAC()
}

func (c *Character) AC() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.ac
}

func (c *Character) MovementSpeed(base int) int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	total := base - c.movementPenalty
	for _, eff := range c.activeEffects {
		total += eff.SpeedMod
	}
	return max(total, 0)
}

func (c *Character) Initiative() int {
	return c.AbilityMod(Dexterity)
}

func (c *Character) EquipShield(shield Armor) error {
	if shield.Type != Shield {
		return errors.New("not a shield")
	}

	c.shield = &shield
	c.recalculateAC()
	return nil
}
