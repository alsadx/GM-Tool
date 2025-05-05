package armor_test

import (
	"testing"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/armor"
)

func TestArmor_CalculateAC(t *testing.T) {
	tests := []struct {
		name     string
		armor    *armor.Armor
		dexMod   int
		expected int
	}{
		{
			name:     "No armor with positive dex",
			armor:    &armor.Armor{Type: armor.NoArmor},
			dexMod:   3,
			expected: 13,
		},
		{
			name:     "No armor with negative dex",
			armor:    &armor.Armor{Type: armor.NoArmor},
			dexMod:   -2,
			expected: 8,
		},

		{
			name:     "Light armor basic",
			armor:    &armor.Armor{Type: armor.LightArmor, BaseAC: 12},
			dexMod:   3,
			expected: 15,
		},
		{
			name:     "Light armor with negative dex",
			armor:    &armor.Armor{Type: armor.LightArmor, BaseAC: 11},
			dexMod:   -1,
			expected: 10,
		},

		{
			name:     "Medium armor below max dex",
			armor:    &armor.Armor{Type: armor.MediumArmor, BaseAC: 14, MaxDex: 2},
			dexMod:   3,
			expected: 16,
		},
		{
			name:     "Medium armor at max dex",
			armor:    &armor.Armor{Type: armor.MediumArmor, BaseAC: 15, MaxDex: 3},
			dexMod:   3,
			expected: 18,
		},
		{
			name:     "Medium armor with negative dex",
			armor:    &armor.Armor{Type: armor.MediumArmor, BaseAC: 13, MaxDex: 2},
			dexMod:   -2,
			expected: 11,
		},

		{
			name:     "Heavy armor with positive dex",
			armor:    &armor.Armor{Type: armor.HeavyArmor, BaseAC: 18},
			dexMod:   5,
			expected: 18,
		},
		{
			name:     "Heavy armor with negative dex",
			armor:    &armor.Armor{Type: armor.HeavyArmor, BaseAC: 16},
			dexMod:   -3,
			expected: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.armor.CalculateAC(tt.dexMod)
			if result != tt.expected {
				t.Errorf("Expected AC %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestArmor_CheckRequirements(t *testing.T) {
	tests := []struct {
		name     string
		armor    *armor.Armor
		strength int
		expected bool
	}{
		{
			name:     "Meets requirement",
			armor:    &armor.Armor{StrengthReq: 13},
			strength: 15,
			expected: false,
		},
		{
			name:     "Exactly meets requirement",
			armor:    &armor.Armor{StrengthReq: 13},
			strength: 13,
			expected: false,
		},
		{
			name:     "Below requirement",
			armor:    &armor.Armor{StrengthReq: 13},
			strength: 10,
			expected: true,
		},
		{
			name:     "No requirement",
			armor:    &armor.Armor{StrengthReq: 0},
			strength: 8,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.armor.CheckRequirements(tt.strength)
			if result != tt.expected {
				t.Errorf("Expected %t, got %t", tt.expected, result)
			}
		})
	}
}
