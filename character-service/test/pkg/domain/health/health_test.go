package health_test

import (
	"errors"
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
)

func TestHealthPoint(t *testing.T) {
	t.Run("SetMaxHP", func(t *testing.T) {
		tests := []struct {
			name            string
			initialHP       int
			newMaxHP        int
			expectedHP      int
			expectedCurrent int
		}{
			{"Increase max HP", 50, 60, 60, 50},
			{"Decrease max HP below current", 50, 30, 30, 30},
			{"Same max HP", 50, 50, 50, 50},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				hp := &health.HealthPoint{CurrentHP: tt.initialHP, MaxHP: tt.initialHP}
				hp.SetMaxHP(tt.newMaxHP)
				if hp.MaxHP != tt.expectedHP || hp.CurrentHP != tt.expectedCurrent {
					t.Errorf("SetMaxHP() = %d/%d, want %d/%d", hp.MaxHP, hp.CurrentHP, tt.expectedHP, tt.expectedCurrent)
				}
			})
		}
	})

	t.Run("Add/RemoveHitDice", func(t *testing.T) {
		d6 := dice.D6
		hp := health.New(100, d6)

		t.Run("Add existing hit dice", func(t *testing.T) {
			hp.AddHitDice(d6)
			if hd := hp.HitDice[d6]; hd.MaxAvailable != 2 || hd.Available != 2 {
				t.Error("AddHitDice() failed to increment existing entry")
			}
		})

		t.Run("Remove hit dice", func(t *testing.T) {
			err := hp.RemoveHitDice(d6)
			if err != nil || hp.HitDice[d6].MaxAvailable != 1 || hp.HitDice[d6].Available != 1 {
				t.Error("RemoveHitDice() failed")
			}
		})

		t.Run("Remove non-existent dice", func(t *testing.T) {
			err := hp.RemoveHitDice(dice.D8)
			if !errors.Is(err, health.ErrWrongTypeHitDice) {
				t.Error("RemoveHitDice() should return error for wrong type")
			}
		})
	})

	t.Run("RollHitDiceRest", func(t *testing.T) {
		originalRoll := dice.MultiRollDice
		defer func() { dice.MultiRollDice = originalRoll }()
		dice.MultiRollDice = func(d dice.Dice, n int) []int { return []int{4, 5} }

		hp := &health.HealthPoint{
			HitDice: map[dice.Dice]*health.Amount{
				dice.D8: {MaxAvailable: 3, Available: 2},
			},
		}

		tests := []struct {
			name        string
			rolling     map[dice.Dice]int
			expectedErr error
			expectedLen int
		}{
			{"Valid roll", map[dice.Dice]int{dice.D8: 2}, nil, 2},
			{"Not enough dice", map[dice.Dice]int{dice.D8: 3}, health.ErrNoHitDiceAvailable, 0},
			{"Wrong dice type", map[dice.Dice]int{dice.D6: 1}, health.ErrWrongTypeHitDice, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				result, err := hp.RollHitDiceRest(tt.rolling)
				if !errors.Is(err, tt.expectedErr) {
					t.Errorf("RollHitDiceRest() error = %v, want %v", err, tt.expectedErr)
				}
				if len(result) != tt.expectedLen {
					t.Errorf("RollHitDiceRest() result length = %d, want %d", len(result), tt.expectedLen)
				}
			})
		}
	})

	t.Run("TakeDamage/Heal", func(t *testing.T) {
		hp := &health.HealthPoint{
			CurrentHP: 30,
			MaxHP:     50,
			TempHP:    10,
		}

		tests := []struct {
			name         string
			operation    func()
			expectedHP   int
			expectedTemp int
		}{
			{"Damage less than temp", func() { hp.TakeDamage(5) }, 30, 5},
			{"Damage through temp", func() { hp.TakeDamage(15) }, 25, 0},
			{"Overkill damage", func() { hp.TakeDamage(40) }, 0, 0},
			{"Heal normal", func() { hp.Heal(15) }, 45, 10},
			{"Overheal", func() { hp.Heal(60) }, 50, 10},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				hp.CurrentHP = 30
				hp.MaxHP = 50
				hp.TempHP = 10
				tt.operation()
				if hp.CurrentHP != tt.expectedHP || hp.TempHP != tt.expectedTemp {
					t.Errorf("After operation: HP = %d/%d, Temp = %d, want %d/%d",
						hp.CurrentHP, hp.MaxHP, hp.TempHP, tt.expectedHP, tt.expectedTemp)
				}
			})
		}
	})

	t.Run("AddTempHP", func(t *testing.T) {
		hp := &health.HealthPoint{TempHP: 5}
		hp.AddTempHP(3)
		if hp.TempHP != 5 {
			t.Error("AddTempHP() should keep higher value")
		}
		hp.AddTempHP(8)
		if hp.TempHP != 8 {
			t.Error("AddTempHP() should update to higher value")
		}
	})
}
