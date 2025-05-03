package ability_test

import (
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
)

func TestScore_UpdateModifier(t *testing.T) {
	tests := []struct {
		name    string
		base    int
		temp    int
		wantMod int
	}{
		{
			name:    "average score",
			base:    10,
			temp:    0,
			wantMod: 0,
		},
		{
			name:    "positive modifier",
			base:    15,
			temp:    0,
			wantMod: 2,
		},
		{
			name:    "temp adjustment to average",
			base:    8,
			temp:    2,
			wantMod: 0,
		},
		{
			name:    "negative temp adjustment",
			base:    12,
			temp:    -2,
			wantMod: 0,
		},
		{
			name:    "high score with temp bonus",
			base:    18,
			temp:    3,
			wantMod: 5,
		},
		{
			name:    "negative modifier",
			base:    9,
			temp:    -1,
			wantMod: -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ability.NewScore(tt.base)

			if s.Base() != tt.base {
				t.Errorf("Base() = %v, want %v", s.Base(), tt.base)
			}

			if tt.temp != 0 {
				remove := s.AddTemp(tt.temp)
				defer remove()
			}

			if s.Temp() != tt.temp {
				t.Errorf("Temp() = %v, want %v", s.Temp(), tt.temp)
			}

			if got := s.Modifier(); got != tt.wantMod {
				t.Errorf("Modifier() = %v, want %v", got, tt.wantMod)
			}
		})
	}
}

func TestScore_SetBase(t *testing.T) {
	tests := []struct {
		name    string
		initial int
		newBase int
		wantMod int
	}{
		{
			name:    "increase base",
			initial: 10,
			newBase: 14,
			wantMod: 2,
		},
		{
			name:    "decrease base",
			initial: 16,
			newBase: 12,
			wantMod: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := ability.NewScore(tt.initial)
			s.SetBase(tt.newBase)

			if s.Modifier() != tt.wantMod {
				t.Errorf("Modifier = %v, want %v", s.Modifier(), tt.wantMod)
			}
		})
	}
}

func TestScore_AddTemp(t *testing.T) {
	s := ability.NewScore(10)

	t.Run("add temporary bonus", func(t *testing.T) {
		remove := s.AddTemp(3)
		if s.Modifier() != 1 { // (10+3-10)/2 = 1.5 → 1
			t.Errorf("Modifier = %v, want 1", s.Modifier())
		}

		if s.Temp() != 3 {
			t.Errorf("Temp = %v, want 3", s.Temp())
		}

		remove()
		if s.Modifier() != 0 {
			t.Errorf("Modifier after remove = %v, want 0", s.Modifier())
		}

		if s.Temp() != 0 {
			t.Errorf("Temp = %v, want 0", s.Temp())
		}
	})
}

func TestScore_Check(t *testing.T) {
	originalRoll := dice.RollDice
	defer func() { dice.RollDice = originalRoll }()
	dice.RollDice = func(d dice.Dice) int { return 10 }

	s := ability.NewScore(14)

	t.Run("check roll", func(t *testing.T) {
		diceRes, modifier, result := s.Check()

		if diceRes != 10 {
			t.Errorf("Dice result = %v, want 10", diceRes)
		}
		if modifier != 2 {
			t.Errorf("Modifier = %v, want 2", modifier)
		}
		if result != 12 {
			t.Errorf("Final result = %v, want 12", result)
		}
	})
}
