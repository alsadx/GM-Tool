package skill_test

import (
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/skill"
)

func TestSkill_Check(t *testing.T) {
	// Сохраняем оригинальную функцию броска костей
	originalRoll := dice.RollDice
	defer func() { dice.RollDice = originalRoll }()

	tests := []struct {
		name           string
		skillBonus     int
		modifier       int
		mockDiceRoll   int
		expectedDice   int
		expectedBonus  int
		expectedResult int
	}{
		{
			name:           "Base case",
			skillBonus:     2,
			modifier:       3,
			mockDiceRoll:   15,
			expectedDice:   15,
			expectedBonus:  5,
			expectedResult: 20,
		},
		{
			name:           "Negative values",
			skillBonus:     -1,
			modifier:       -2,
			mockDiceRoll:   10,
			expectedDice:   10,
			expectedBonus:  -3,
			expectedResult: 7,
		},
		{
			name:           "Zero values",
			skillBonus:     0,
			modifier:       0,
			mockDiceRoll:   5,
			expectedDice:   5,
			expectedBonus:  0,
			expectedResult: 5,
		},
		{
			name:           "Max values",
			skillBonus:     10,
			modifier:       5,
			mockDiceRoll:   20,
			expectedDice:   20,
			expectedBonus:  15,
			expectedResult: 35,
		},
		{
			name:           "Mixed signs",
			skillBonus:     3,
			modifier:       -2,
			mockDiceRoll:   12,
			expectedDice:   12,
			expectedBonus:  1,
			expectedResult: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Мокаем бросок костей
			dice.RollDice = func(d dice.Dice) int { return tt.mockDiceRoll }

			// Создаем навык с заданным бонусом
			sk := skill.NewSkill(tt.skillBonus)

			// Выполняем проверку
			diceRes, bonus, result := sk.Check(tt.modifier)

			// Проверяем результаты
			if diceRes != tt.expectedDice {
				t.Errorf("Dice result = %d, want %d", diceRes, tt.expectedDice)
			}
			if bonus != tt.expectedBonus {
				t.Errorf("Total bonus = %d, want %d", bonus, tt.expectedBonus)
			}
			if result != tt.expectedResult {
				t.Errorf("Final result = %d, want %d", result, tt.expectedResult)
			}
		})
	}
}

func TestSkill_SetBonus(t *testing.T) {
	tests := []struct {
		name         string
		initialBonus int
		newBonus     int
	}{
		{
			name:         "Positive to positive",
			initialBonus: 2,
			newBonus:     4,
		},
		{
			name:         "Positive to negative",
			initialBonus: 3,
			newBonus:     -1,
		},
		{
			name:         "Zero to positive",
			initialBonus: 0,
			newBonus:     5,
		},
		{
			name:         "Negative to zero",
			initialBonus: -2,
			newBonus:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := skill.NewSkill(tt.initialBonus)
			sk.SetBonus(tt.newBonus)

			if sk.Bonus != tt.newBonus {
				t.Errorf("Bonus = %d, want %d", sk.Bonus, tt.newBonus)
			}
		})
	}
}

func TestNewSkill(t *testing.T) {
	tests := []struct {
		name     string
		bonus    int
		expected int
	}{
		{"Positive bonus", 5, 5},
		{"Negative bonus", -2, -2},
		{"Zero bonus", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sk := skill.NewSkill(tt.bonus)
			if sk.Bonus != tt.expected {
				t.Errorf("NewSkill() bonus = %d, want %d", sk.Bonus, tt.expected)
			}
		})
	}
}
