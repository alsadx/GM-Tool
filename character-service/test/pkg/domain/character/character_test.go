package character_test

import (
	"errors"
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
	"github.com/stretchr/testify/assert"
)

func TestNewCharacter(t *testing.T) {
	tests := []struct {
		name        string
		id          string
		ownerID     int
		charName    string
		className   string
		subclass    string
		race        string
		expectedErr error
	}{
		{
			name:        "Valid basic character",
			id:          "1",
			ownerID:     100,
			charName:    "TestHero",
			className:   "Warrior",
			subclass:    "Berserker",
			race:        "Human",
			expectedErr: nil,
		},
		{
			name:        "Empty character name",
			id:          "2",
			ownerID:     101,
			charName:    "",
			className:   "Mage",
			subclass:    "Necromancer",
			race:        "Elf",
			expectedErr: character.ErrInvalidCharacterName,
		},
		{
			name:        "Negative owner ID",
			id:          "3",
			ownerID:     -1,
			charName:    "TestVillain",
			className:   "Rogue",
			subclass:    "Thief",
			race:        "Halfling",
			expectedErr: character.ErrInvalidOwnerID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := character.New(tt.id, tt.ownerID, tt.charName, tt.className, tt.subclass, tt.race)

			if !errors.Is(err, tt.expectedErr) {
				t.Fatalf("Creation error: got %v, want %v", err, tt.expectedErr)
			}

			if tt.expectedErr != nil {
				return
			}

			assert.Equal(t, tt.id, c.ID, "Wrong id")
			assert.Equal(t, tt.ownerID, c.Owner, "Wrong owner")
			assert.Equal(t, tt.charName, c.Name, "Wrong name")
			assert.Equal(t, tt.className, c.Class, "Wrong class")
			assert.Equal(t, tt.subclass, c.Subclass, "Wrong subclass")
			assert.Equal(t, tt.race, c.Race, "Wrong race")
			assert.Equal(t, 1, c.GetLvl(), "Wrong lvl")
			assert.Equal(t, 0, c.GetCurrentExp(), "Wrong exp")

			abilities := []types.AbilityType{
				types.Strength,
				types.Dexterity,
				types.Constitution,
				types.Intelligence,
				types.Wisdom,
				types.Charisma,
			}

			for _, abilityType := range abilities {
				assert.Equal(t, 0, c.Ability(abilityType).Modifier(), "Wrong modifier")
			}

			assert.Greater(t, c.GetHP(), 0, "Wrong HP")
			assert.False(t, c.IsKnocked(), "The character was created in the knockout")
		})
	}
}

func TestLevelManagement(t *testing.T) {
	tests := []struct {
		name            string
		operation       func(c *character.Character)
		expectedLvl     int
		expectedExp     int
		expectedCanUp   bool
		expectedCanDown bool
		expToNext       int
	}{
		{
			name:            "Initial state",
			operation:       func(c *character.Character) {},
			expectedLvl:     1,
			expectedExp:     0,
			expectedCanUp:   false,
			expectedCanDown: false,
			expToNext:       300,
		},
		{
			name: "Add experience below threshold",
			operation: func(c *character.Character) {
				c.GainExp(299)
			},
			expectedLvl:     1,
			expectedExp:     299,
			expectedCanUp:   false,
			expectedCanDown: false,
			expToNext:       1, // 300 - 299
		},
		{
			name: "Reach level 2 threshold",
			operation: func(c *character.Character) {
				c.GainExp(300)
			},
			expectedLvl:     1,
			expectedExp:     300,
			expectedCanUp:   true,
			expectedCanDown: false,
			expToNext:       600, // 900 - 300
		},
		{
			name: "Level up to 2",
			operation: func(c *character.Character) {
				c.GainExp(300)
				c.LvlUp()
			},
			expectedLvl:     2,
			expectedExp:     300,
			expectedCanUp:   false,
			expectedCanDown: false,
			expToNext:       600,
		},
		{
			name: "Add excess experience",
			operation: func(c *character.Character) {
				c.GainExp(1500)
			},
			expectedLvl:     1,
			expectedExp:     1500,
			expectedCanUp:   true,
			expectedCanDown: false,
			expToNext:       1200,
		},
		{
			name: "Remove experience below previous tier",
			operation: func(c *character.Character) {
				c.GainExp(900)
				c.RemoveExp(601)
			},
			expectedLvl:     1,
			expectedExp:     299,
			expectedCanUp:   false,
			expectedCanDown: false,
			expToNext:       1,
		},
		{
			name: "Force level down",
			operation: func(c *character.Character) {
				c.SetLvl(3)
				c.RemoveExp(c.GetCurrentExp())
			},
			expectedLvl:     3,
			expectedExp:     0,
			expectedCanUp:   false,
			expectedCanDown: true,
			expToNext:       300,
		},
		{
			name: "Set invalid level recovery",
			operation: func(c *character.Character) {
				c.SetLvl(5)
				c.RemoveExp(9999999)
			},
			expectedLvl:     5,
			expectedExp:     0,
			expectedCanUp:   false,
			expectedCanDown: true,
			expToNext:       300,
		},
		{
			name: "Max level experience",
			operation: func(c *character.Character) {
				c.GainExp(355000)
				for c.CanLvlUp() {
					c.LvlUp()
				}
			},
			expectedLvl:     20,
			expectedExp:     355000,
			expectedCanUp:   false,
			expectedCanDown: false,
			expToNext:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, err := character.New("1", 100, "Test", "Wizard", "Necromancer", "Elf")
			if err != nil {
				t.Fatalf("Failed to create character: %v", err)
			}

			tt.operation(c)

			assert.Equal(t, tt.expectedLvl, c.GetLvl(), "Wrong lvl")
			assert.Equal(t, tt.expectedExp, c.GetCurrentExp(), "Wrong exp")
			assert.Equal(t, tt.expectedCanUp, c.CanLvlUp(), "CanLevelUp mismatch")
			assert.Equal(t, tt.expectedCanDown, c.CanLvlDown(), "CanLevelDown mismatch")
			assert.Equal(t, tt.expToNext, c.ExpToNextLevel(), "ExpToNextLevel mismatch")
		})
	}
}

func TestHealthManagement(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*character.Character)
		operation    func(*character.Character)
		expectedHP   int
		expectedTemp int
		expectedMax  int
		isKnocked    bool
	}{
		{
			name: "Take normal damage without temp HP",
			setup: func(c *character.Character) {
				c.SetMaxHP(50)
				c.Heal(50)
			},
			operation:    func(c *character.Character) { c.TakeDamage(30) },
			expectedHP:   20,
			expectedTemp: 0,
			expectedMax:  50,
			isKnocked:    false,
		},
		{
			name: "Take damage with temp HP",
			setup: func(c *character.Character) {
				c.SetMaxHP(50)
				c.Heal(50)
				c.SetTempHP(10)
			},
			operation:    func(c *character.Character) { c.TakeDamage(35) },
			expectedHP:   25,
			expectedTemp: 0,
			expectedMax:  50,
			isKnocked:    false,
		},
		{
			name: "Exact knockout",
			setup: func(c *character.Character) {
				c.SetMaxHP(50)
				c.Heal(50)
			},
			operation:    func(c *character.Character) { c.TakeDamage(50) },
			expectedHP:   0,
			expectedTemp: 0,
			expectedMax:  50,
			isKnocked:    true,
		},
		{
			name: "Overkill damage",
			setup: func(c *character.Character) {
				c.SetMaxHP(50)
				c.Heal(50)
			},
			operation:    func(c *character.Character) { c.TakeDamage(100) },
			expectedHP:   0,
			expectedTemp: 0,
			expectedMax:  50,
			isKnocked:    true,
		},
		{
			name: "Heal from knockout",
			setup: func(c *character.Character) {
				c.SetMaxHP(50)
				c.TakeDamage(50)
			},
			operation:    func(c *character.Character) { c.Heal(20) },
			expectedHP:   20,
			expectedTemp: 0,
			expectedMax:  50,
			isKnocked:    false,
		},
		{
			name: "Heal over max HP",
			setup: func(c *character.Character) {
				c.SetMaxHP(50)
				c.Heal(50)
			},
			operation:    func(c *character.Character) { c.Heal(100) },
			expectedHP:   50,
			expectedTemp: 0,
			expectedMax:  50,
			isKnocked:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := character.New("1", 100, "Test", "Cleric", "Life", "Dwarf")

			tt.setup(c)

			tt.operation(c)

			assert.Equal(t, tt.expectedHP, c.GetCurrentHP(), "Wrong current HP")
			assert.Equal(t, tt.expectedTemp, c.GetTempHP(), "Wrong temp HP")
			assert.Equal(t, tt.expectedMax, c.GetMaxHP(), "Wrong max HP")
			assert.Equal(t, tt.isKnocked, c.IsKnocked(), "Wrong knockout state")
		})
	}
}

func TestHitDiceManagement(t *testing.T) {
	createCharacter := func() *character.Character {
		c, _ := character.New("1", 100, "Test", "Rogue", "Thief", "Halfling")
		c.WithDice(50, dice.D8)
		c.SetLvl(3)
		return c
	}

	tests := []struct {
		name         string
		setup        func(*character.Character)
		operation    func(*character.Character) error
		expectedErr  error
		expectedDice map[dice.Dice]int
	}{
		{
			name:  "Add first hit dice at level 3",
			setup: func(c *character.Character) {},
			operation: func(c *character.Character) error {
				return c.AddHitDice(dice.D8)
			},
			expectedErr: nil,
			expectedDice: map[dice.Dice]int{
				dice.D8: 2,
			},
		},
		{
			name: "Add beyond level limit",
			setup: func(c *character.Character) {
				_ = c.AddHitDice(dice.D8)
				_ = c.AddHitDice(dice.D8)
			},
			operation: func(c *character.Character) error {
				return c.AddHitDice(dice.D8)
			},
			expectedErr: character.ErrNotEnoughLvlForAddHidDice,
			expectedDice: map[dice.Dice]int{
				dice.D8: 3,
			},
		},
		{
			name: "Remove existing dice type",
			setup: func(c *character.Character) {
				_ = c.AddHitDice(dice.D8)
			},
			operation: func(c *character.Character) error {
				return c.RemoveHitDice(dice.D8)
			},
			expectedErr: nil,
			expectedDice: map[dice.Dice]int{
				dice.D8: 1,
			},
		},
		{
			name:  "Remove non-existent dice type",
			setup: func(c *character.Character) {},
			operation: func(c *character.Character) error {
				return c.RemoveHitDice(dice.D10)
			},
			expectedErr: health.ErrWrongTypeHitDice,
			expectedDice: map[dice.Dice]int{
				dice.D8: 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := createCharacter()
			tt.setup(c)

			err := tt.operation(c)
			if err != nil {
				assert.EqualError(t, err, tt.expectedErr.Error(), "Wrong error")
			}

			hitDice := c.GetHidDice()
			assert.Equal(t, len(tt.expectedDice), len(hitDice), "Wrong number of types")

			for diceType, expectedCount := range tt.expectedDice {
				actual, ok := hitDice[diceType]
				if !ok {
					t.Fatalf("The expected dice type is missing: %v", diceType)
				}
				assert.Equal(t, expectedCount, actual.MaxAvailable, "Wrong number of hit dice")
			}
		})
	}
}

func TestSkillChecks(t *testing.T) {
	originalRoll := dice.RollDice
	defer func() { dice.RollDice = originalRoll }()
	dice.RollDice = func(d dice.Dice) int { return 15 }

	c, _ := character.New("1", 100, "Test", "Bard", "Lore", "Half-Elf")
	ability := c.Ability(types.Charisma)
	ability.Skills[types.Persuasion].SetBonus(3)

	tests := []struct {
		name           string
		checkFunc      func() (int, int, int)
		expectedDice   int
		expectedBonus  int
		expectedResult int
	}{
		{
			name: "Ability check",
			checkFunc: func() (int, int, int) {
				return c.CheckAbility(types.Charisma)
			},
			expectedDice:   15,
			expectedBonus:  0,
			expectedResult: 15,
		},
		{
			name: "Skill check",
			checkFunc: func() (int, int, int) {
				return c.CheckSkill(types.Persuasion)
			},
			expectedDice:   15,
			expectedBonus:  3,
			expectedResult: 18,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, b, r := tt.checkFunc()
			assert.Equal(t, tt.expectedDice, d, "Dice roll mismatch")
			assert.Equal(t, tt.expectedBonus, b, "Bonus mismatch")
			assert.Equal(t, tt.expectedResult, r, "Result mismatch")
		})
	}
}
