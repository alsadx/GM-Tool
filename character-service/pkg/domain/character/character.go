package character

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/level"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
)

type Character struct {
	ID    int
	Owner int

	isKnocked bool

	Name     string
	Class    string
	Subclass string
	Race     string

	lvl    *level.LevelSystem
	stats  map[types.AbilityType]*ability.Ability
	health *health.HealthPoint
}

func New(id int, ownerID int, name, class, subclass, race string) (*Character, error) {
	if name == "" {
		return nil, ErrInvalidCharacterName
	}
	if ownerID < 0 {
		return nil, ErrInvalidOwnerID
	}
	return &Character{
		ID:       id,
		Owner:    ownerID,
		Name:     name,
		Class:    class,
		Subclass: subclass,
		Race:     race,
		lvl:      level.NewLevelSystem(),
		stats:    ability.NewStats(),
		health:   health.New(9, dice.D6),
	}, nil
}

func (c *Character) WithDice(maxHp int, hitDice dice.Dice) *Character {
	c.health = health.New(maxHp, hitDice)
	return c
}

func (c *Character) GainExp(amountExp int) {
	c.lvl.AddExp(amountExp)
}

func (c *Character) RemoveExp(amountExp int) {
	c.lvl.RemoveExp(amountExp)
}

func (c *Character) GetLvl() int { return c.lvl.CurrentLevel() }

func (c *Character) SetLvl(lvl int) { c.lvl.SetLevel(lvl) }

func (c *Character) GetCurrentExp() int { return c.lvl.CurrentExp() }

func (c *Character) ExpToNextLevel() int { return c.lvl.ExpToNextLevel() }

func (c *Character) CanLvlUp() bool { return c.lvl.CanLevelUp() }

func (c *Character) CanLvlDown() bool { return c.lvl.CanLevelDown() }

func (c *Character) LvlUp() bool { return c.lvl.LevelUp() }

func (c *Character) LvlDown() bool { return c.lvl.LevelDown() }

func (c *Character) TakeDamage(dmg int) {
	c.health.TakeDamage(dmg)
	if c.health.CurrentHP == 0 {
		c.isKnocked = true
	}
}

func (c *Character) Heal(healthAmount int) {
	c.health.Heal(healthAmount)
	if healthAmount > 0 {
		c.isKnocked = false
	}
}

func (c *Character) GetCurrentHP() int { return c.health.CurrentHP }

func (c *Character) GetTempHP() int { return c.health.TempHP }

func (c *Character) GetHP() int { return c.health.CurrentHP + c.health.TempHP }

func (c *Character) GetMaxHP() int { return c.health.MaxHP }

func (c *Character) IsKnocked() bool { return c.isKnocked }

func (c *Character) SetMaxHP(maxHp int) { c.health.SetMaxHP(maxHp) }

func (c *Character) AddHitDice(hidDiceType dice.Dice) error {
	if (c.health.GetHitDiceCount() + 1) > c.lvl.CurrentLevel() {
		return ErrNotEnoughLvlForAddHidDice
	}
	c.health.AddHitDice(hidDiceType)
	return nil
}

func (c *Character) AddTempHp(tempHp int) { c.health.AddTempHP(tempHp) }

func (c *Character) RemoveHidDice(hitDiceType dice.Dice) error {
	return c.health.RemoveHitDice(hitDiceType)
}

func (c *Character) GetHidDice() map[dice.Dice]*health.Amount {
	return c.health.HitDice
}

func (c *Character) CheckAbility(abilityType types.AbilityType) (diceRes, bonus, result int) {
	return c.stats[abilityType].Check()
}

func (c *Character) CheckSkill(skillType types.SkillType) (diceRes, bonus, result int) {
	abil := types.SkillToAbility[skillType]
	return c.stats[abil].CheckSkill(skillType)
}

func (c *Character) Ability(abilityType types.AbilityType) *ability.Ability {
	return c.stats[abilityType]
}
