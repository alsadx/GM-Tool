package dto_test

import (
	"testing"

	"github.com/alsadx/GM-Tool/character-service/gen"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCharacterDTO_Conversions(t *testing.T) {
	stats := map[string]*dto.AbilityDTO{
		"Strength": {
			ScoreDTO: &dto.ScoreDTO{Base: 16},
			Skills:   map[string]*dto.SkillDTO{"Athletics": {Bonus: 5}},
		},
	}

	healthDTO := &dto.HealthDTO{
		CurrentHP: 50,
		MaxHP:     60,
	}

	dtoObj := &dto.CharacterDTO{
		ID:       1,
		Owner:    100,
		Name:     "Test Hero",
		Class:    "Warrior",
		Subclass: "Berserker",
		Race:     "Human",
		Lvl:      dto.LevelDTO{CurrentExp: 900},
		Stats:    stats,
		Health:   healthDTO,
	}

	t.Run("ToProto", func(t *testing.T) {
		proto := dtoObj.ToProto()
		assert.Equal(t, int64(1), proto.Id)
		assert.Equal(t, "Test Hero", proto.Name)
		assert.Equal(t, "Warrior", proto.ClassName)
		assert.Equal(t, int32(900), proto.Lvl.CurrentExp)
		assert.Equal(t, int32(50), proto.Health.CurrentHp)
	})

	t.Run("ToDomain", func(t *testing.T) {
		domain, err := dtoObj.ToDomain()
		require.NoError(t, err)
		assert.Equal(t, "Test Hero", domain.Name)
		assert.Equal(t, "Berserker", domain.Subclass)
		assert.Equal(t, 50, domain.GetCurrentHP())
	})

	t.Run("FromProto", func(t *testing.T) {
		proto := &gen.Character{
			Id:        2,
			Owner:     200,
			Name:      "Proto Hero",
			ClassName: "Wizard",
			Subclass:  "Illusionist",
			Race:      "Elf",
			Lvl:       &gen.LevelSystem{CurrentExp: 1500},
			Health: &gen.HealthPoint{
				CurrentHp: 30,
				MaxHp:     35,
			},
			Stats: map[string]*gen.Ability{
				"Intelligence": {
					Score: &gen.Score{Base: 18},
					Skills: map[string]*gen.Skill{
						"Arcana": {Bonus: 7},
					},
				},
			},
		}

		dto := dto.FromProto(proto)
		assert.Equal(t, 2, dto.ID)
		assert.Equal(t, "Proto Hero", dto.Name)
		assert.Equal(t, "Wizard", dto.Class)
		assert.Equal(t, "Illusionist", dto.Subclass)
		assert.Equal(t, 30, dto.Health.CurrentHP)
		assert.Equal(t, 18, dto.Stats["Intelligence"].Base)
	})
}
