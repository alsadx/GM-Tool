package ability_test

import (
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/skill"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
)

func TestScore(t *testing.T) {
	t.Run("Modifier calculations", func(t *testing.T) {
		tests := []struct {
			name     string
			base     int
			temp     int
			expected int
		}{
			{"Average score", 10, 0, 0},
			{"Positive modifier", 15, 0, 2},
			{"Temp adjustment", 8, 2, 0},
			{"Negative temp", 12, -2, 0},
			{"High score", 18, 3, 5},
			{"Negative modifier", 9, -1, -1},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := ability.NewScore(tt.base)
				if tt.temp != 0 {
					remove := s.AddTemp(tt.temp)
					defer remove()
				}
				if got := s.Modifier(); got != tt.expected {
					t.Errorf("Modifier() = %d, want %d", got, tt.expected)
				}
			})
		}
	})

	t.Run("SetBase", func(t *testing.T) {
		tests := []struct {
			name     string
			initial  int
			newBase  int
			expected int
		}{
			{"Increase base", 10, 14, 2},
			{"Decrease base", 16, 12, 1},
			{"Edge case 9", 9, 9, -1},
			{"Edge case 11", 11, 11, 0},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				s := ability.NewScore(tt.initial)
				s.SetBase(tt.newBase)
				if s.Modifier() != tt.expected {
					t.Errorf("Modifier after SetBase() = %d, want %d", s.Modifier(), tt.expected)
				}
			})
		}
	})

	t.Run("Temp operations", func(t *testing.T) {
		s := ability.NewScore(10)
		remove := s.AddTemp(3)

		t.Run("Add temp", func(t *testing.T) {
			if s.Modifier() != 1 {
				t.Error("Temp modifier not applied correctly")
			}
		})

		t.Run("Remove temp", func(t *testing.T) {
			remove()
			if s.Modifier() != 0 {
				t.Error("Temp modifier not removed correctly")
			}
		})
	})

	t.Run("Check", func(t *testing.T) {
		originalRoll := dice.RollDice
		defer func() { dice.RollDice = originalRoll }()
		dice.RollDice = func(d dice.Dice) int { return 15 }

		s := ability.NewScore(16) // Modifier 3
		diceRes, mod, result := s.Check()

		if diceRes != 15 || mod != 3 || result != 18 {
			t.Errorf("Check() = (%d, %d, %d), want (15, 3, 18)", diceRes, mod, result)
		}
	})
}

func TestAbility(t *testing.T) {
	t.Run("New Ability initialization", func(t *testing.T) {
		abilityTypes := []types.AbilityType{
			types.Strength,
			types.Dexterity,
			types.Intelligence,
		}

		for _, at := range abilityTypes {
			t.Run(at.String(), func(t *testing.T) {
				a := ability.New(at)
				expectedSkills := types.AbilityToSkill[at]

				if len(a.Skills) != len(expectedSkills) {
					t.Fatalf("Expected %d skills, got %d", len(expectedSkills), len(a.Skills))
				}

				for _, st := range expectedSkills {
					if _, ok := a.Skills[st]; !ok {
						t.Errorf("Missing skill: %s", st.String())
					}
				}
			})
		}
	})

	t.Run("Skill management", func(t *testing.T) {
		a := ability.New(types.Wisdom)
		skillType := types.SkillType(types.AnimalHandling)

		t.Run("Initial skill", func(t *testing.T) {
			if _, ok := a.Skills[skillType]; !ok {
				t.Error("Default skill not present")
			}
		})

		t.Run("Modify skill", func(t *testing.T) {
			newSkill := &skill.Skill{Bonus: 2}
			a.Skills[skillType] = newSkill
			if a.Skills[skillType].Bonus != 2 {
				t.Error("Skill modification failed")
			}
		})
	})
}
