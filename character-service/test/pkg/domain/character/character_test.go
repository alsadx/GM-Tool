package character_test

import (
	"errors"
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
)

func TestNewCharacter(t *testing.T) {
    tests := []struct {
        name        string
        id          int
        ownerID     int
        charName    string
        className   string
        subclass    string
        race        string
        expectedErr error
    }{
        {
            name:        "Valid basic character",
            id:          1,
            ownerID:     100,
            charName:    "TestHero",
            className:   "Warrior",
            subclass:    "Berserker",
            race:        "Human",
            expectedErr: nil,
        },
        {
            name:        "Empty character name",
            id:          2,
            ownerID:     101,
            charName:    "",
            className:   "Mage",
            subclass:    "Necromancer",
            race:        "Elf",
            expectedErr: character.ErrInvalidCharacterName,
        },
        {
            name:        "Negative owner ID",
            id:          3,
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

            if c.ID != tt.id || c.Owner != tt.ownerID {
                t.Errorf("ID/Owner: got %d/%d, want %d/%d",
                    c.ID, c.Owner, tt.id, tt.ownerID)
            }

            if c.Name != tt.charName {
                t.Errorf("Character Name: got '%s', want '%s'", c.Name, tt.charName)
            }
            if c.Class != tt.className {
                t.Errorf("Class: got '%s', want '%s'", c.Class, tt.className)
            }
            if c.Subclass != tt.subclass {
                t.Errorf("Subclass: got '%s', want '%s'", c.Subclass, tt.subclass)
            }
            if c.Race != tt.race {
                t.Errorf("Race: got '%s', want '%s'", c.Race, tt.race)
            }

            if lvl := c.GetLvl(); lvl != 1 {
                t.Errorf("Level: got %d, want 1", lvl)
            }
            if exp := c.GetCurrentExp(); exp != 0 {
                t.Errorf("Experience: got %d, want 0", exp)
            }

            abilities := []types.AbilityType{
                types.Strength,
                types.Dexterity,
                types.Constitution,
                types.Intelligence,
                types.Wisdom,
                types.Charisma,
            }
            
            for _, abilityType := range abilities {
                if mod := c.Ability(abilityType).Modifier(); mod != 0 {
                    t.Errorf("Modifier %s: got %d, want 0", abilityType, mod)
                }
            }

            if hp := c.GetCurrentHP(); hp <= 0 {
                t.Errorf("Incorrect health: %d", hp)
            }
            if c.IsKnocked() {
                t.Error("The character was created in the knockout")
            }
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
                c.GainExp(1500) // Total 1500 (level 3)
            },
            expectedLvl:     1,
            expectedExp:     1500,
            expectedCanUp:   true,
            expectedCanDown: false,
            expToNext:       1200, // 2700 - 1500
        },
        {
            name: "Remove experience below previous tier",
            operation: func(c *character.Character) {
                c.GainExp(900)  // Level 3
                c.RemoveExp(601) // 900 - 601 = 299
            },
            expectedLvl:     1,
            expectedExp:     299,
            expectedCanUp:   false,
            expectedCanDown: false,
            expToNext:       1, // 300 - 299
        },
        {
            name: "Force level down",
            operation: func(c *character.Character) {
                c.SetLvl(3)
                c.RemoveExp(c.GetCurrentExp()) // Reset exp to 0
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
            expToNext:       -355000,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c, err := character.New(1, 100, "Test", "Wizard", "Necromancer", "Elf")
            if err != nil {
                t.Fatalf("Failed to create character: %v", err)
            }

            tt.operation(c)

            if lvl := c.GetLvl(); lvl != tt.expectedLvl {
                t.Errorf("Level mismatch: got %d, want %d", lvl, tt.expectedLvl)
            }
            if exp := c.GetCurrentExp(); exp != tt.expectedExp {
                t.Errorf("Experience mismatch: got %d, want %d", exp, tt.expectedExp)
            }
            if canUp := c.CanLvlUp(); canUp != tt.expectedCanUp {
                t.Errorf("CanLevelUp mismatch: got %t, want %t", canUp, tt.expectedCanUp)
            }
            if canDown := c.CanLvlDown(); canDown != tt.expectedCanDown {
                t.Errorf("CanLevelDown mismatch: got %t, want %t", canDown, tt.expectedCanDown)
            }
            if etn := c.ExpToNextLevel(); etn != tt.expToNext {
                t.Errorf("ExpToNextLevel mismatch: got %d, want %d", etn, tt.expToNext)
            }
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
            expectedMax: 50,
            isKnocked:    false,
        },
        {
            name: "Take damage with temp HP",
            setup: func(c *character.Character) {
                c.SetMaxHP(50)
                c.Heal(50)
                c.AddTempHp(10)
            },
            operation:    func(c *character.Character) { c.TakeDamage(35) },
            expectedHP:   25,
            expectedTemp: 0,
            expectedMax: 50,
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
            expectedMax: 50,
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
            expectedMax: 50,
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
            expectedMax: 50,
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
            expectedMax: 50,
            isKnocked:    false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c, _ := character.New(1, 100, "Test", "Cleric", "Life", "Dwarf")
            
            tt.setup(c)
            
            tt.operation(c)
            
            if currentHP := c.GetCurrentHP(); currentHP != tt.expectedHP {
                t.Errorf("Current HP: got %d, want %d", currentHP, tt.expectedHP)
            }
            
            if tempHP := c.GetTempHP(); tempHP != tt.expectedTemp {
                t.Errorf("Temp HP: got %d, want %d", tempHP, tt.expectedTemp)
            }
            
            if maxHP := c.GetMaxHP(); maxHP != tt.expectedMax {
                t.Errorf("Max HP: got %d, want %d", maxHP, tt.expectedMax)
            }
            
            if isKnocked := c.IsKnocked(); isKnocked != tt.isKnocked {
                t.Errorf("Knockout state: got %t, want %t", isKnocked, tt.isKnocked)
            }
        })
    }
}

func TestHitDiceManagement(t *testing.T) {
    createCharacter := func() *character.Character {
        c, _ := character.New(1, 100, "Test", "Rogue", "Thief", "Halfling")
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
                c.AddHitDice(dice.D8)
                c.AddHitDice(dice.D8)
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
                c.AddHitDice(dice.D8)
            },
            operation: func(c *character.Character) error {
                return c.RemoveHidDice(dice.D8)
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
                return c.RemoveHidDice(dice.D10)
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
            
            if !errors.Is(err, tt.expectedErr) {
                t.Fatalf("Ошибка: получили %v, ожидалось %v", err, tt.expectedErr)
            }

            hitDice := c.GetHidDice()
            if len(hitDice) != len(tt.expectedDice) {
                t.Fatalf("Количество типов костей: получили %d, ожидалось %d", 
                    len(hitDice), len(tt.expectedDice))
            }

            for diceType, expectedCount := range tt.expectedDice {
                actual, ok := hitDice[diceType]
                if !ok {
                    t.Fatalf("Отсутствует ожидаемый тип кости: %v", diceType)
                }

                if actual.MaxAvailable != expectedCount {
                    t.Errorf("Количество костей %v: получили %d, ожидалось %d", 
                        diceType, actual.MaxAvailable, expectedCount)
                }
            }
        })
    }
}

func TestSkillChecks(t *testing.T) {
	originalRoll := dice.RollDice
	defer func() { dice.RollDice = originalRoll }()
	dice.RollDice = func(d dice.Dice) int { return 15 }

	c, _ := character.New(1, 100, "Test", "Bard", "Lore", "Half-Elf")
	ability := c.Ability(types.Charisma)
	ability.Skills[types.Persuasion].SetBonus(3)

	t.Run("Ability check", func(t *testing.T) {
		d, b, r := c.CheckAbility(types.Charisma)
		if d != 15 || b != 0 || r != 15 {
			t.Errorf("CheckAbility() = (%d, %d, %d), want (15, 0, 15)", d, b, r)
		}
	})

	t.Run("Skill check", func(t *testing.T) {
		d, b, r := c.CheckSkill(types.Persuasion)
		if d != 15 || b != 3 || r != 18 {
			t.Errorf("CheckSkill() = (%d, %d, %d), want (15, 3, 18)", d, b, r)
		}
	})
}
