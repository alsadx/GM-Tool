package skill_test

import (
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/skill"
)

func TestSkill_Check(t *testing.T) {
	// Сохраняем оригинальную функцию броска костей
	originalRoll := dice.RollDice
	defer func() { dice.RollDice = originalRoll }()

	tests := []struct {
		name          string
		abilityMod    int
		skillBonus    int
		mockDiceRoll  int
		wantDiceRes   int
		wantTotalBonus int
		wantResult    int
	}{
		{
			name:          "base case",
			abilityMod:    2,
			skillBonus:    3,
			mockDiceRoll:  15,
			wantDiceRes:   15,
			wantTotalBonus: 5,
			wantResult:    20,
		},
		{
			name:          "negative modifier",
			abilityMod:    -1,
			skillBonus:    2,
			mockDiceRoll:  10,
			wantDiceRes:   10,
			wantTotalBonus: 1,
			wantResult:    11,
		},
		{
			name:          "zero values",
			abilityMod:    0,
			skillBonus:    0,
			mockDiceRoll:  5,
			wantDiceRes:   5,
			wantTotalBonus: 0,
			wantResult:    5,
		},
		{
			name:          "max values",
			abilityMod:    5,
			skillBonus:    10,
			mockDiceRoll:  20,
			wantDiceRes:   20,
			wantTotalBonus: 15,
			wantResult:    35,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dice.RollDice = func(d dice.Dice) int { return tt.mockDiceRoll }

			abil := ability.NewScore(10)
			abil.SetBase(10 + tt.abilityMod*2)

			if abil.Modifier() != tt.abilityMod {
				t.Errorf("Modifier() = %v, want %v", abil.Modifier(), tt.abilityMod)
			}

			sk := &skill.Skill{
				Ability: abil,
				Bonus:   tt.skillBonus,
			}

			diceRes, bonus, result := sk.Check()

			if diceRes != tt.wantDiceRes {
				t.Errorf("Dice result = %d, want %d", diceRes, tt.wantDiceRes)
			}
			if bonus != tt.wantTotalBonus {
				t.Errorf("Total bonus = %d, want %d", bonus, tt.wantTotalBonus)
			}
			if result != tt.wantResult {
				t.Errorf("Final result = %d, want %d", result, tt.wantResult)
			}
		})
	}
}

func TestSkill_SetBonus(t *testing.T) {
	tests := []struct {
		name       string
		initial    int
		newBonus   int
		wantBonus  int
	}{
		{
			name:      "set positive bonus",
			initial:   0,
			newBonus:  5,
			wantBonus: 5,
		},
		{
			name:      "set negative bonus",
			initial:   3,
			newBonus:  -2,
			wantBonus: -2,
		},
		{
			name:      "set zero bonus",
			initial:   10,
			newBonus:  0,
			wantBonus: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := &skill.Skill{
				Ability: ability.NewScore(10),
				Bonus:   tt.initial,
			}

			sk.SetBonus(tt.newBonus)

			if sk.Bonus != tt.wantBonus {
				t.Errorf("Bonus = %d, want %d", sk.Bonus, tt.wantBonus)
			}
		})
	}
}