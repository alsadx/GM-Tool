package dto_test

import (
	"testing"

	// gen "github.com/alsadx/GM-Tool/character-service/gen/character"
	gen "github.com/alsadx/gm-protos/gen/go/characterv1"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
	"github.com/stretchr/testify/assert"
)

func TestHealthDTO_Conversions(t *testing.T) {
	hitDice := map[dice.Dice]*dto.AmountDTO{
		dice.D6:  {MaxAvailable: 3, Available: 2},
		dice.D10: {MaxAvailable: 5, Available: 4},
	}

	dtoObj := &dto.HealthDTO{
		CurrentHP: 25,
		MaxHP:     30,
		TempHP:    5,
		HitDice:   hitDice,
	}

	t.Run("ToDomain", func(t *testing.T) {
		domain := dtoObj.ToDomain()
		assert.Equal(t, 25, domain.CurrentHP)
		assert.Equal(t, 30, domain.MaxHP)
		assert.Equal(t, 5, domain.TempHP)
		assert.Equal(t, 2, domain.HitDice[dice.D6].Available)
	})

	t.Run("ToProto", func(t *testing.T) {
		proto := dtoObj.ToProto()
		assert.Equal(t, int32(25), proto.CurrentHp)
		assert.Equal(t, int32(30), proto.MaxHp)
		assert.Equal(t, int32(5), proto.TempHp)
		assert.Equal(t, int32(2), proto.HitDice[int32(dice.D6)].Available)
	})

	t.Run("FromProto", func(t *testing.T) {
		proto := &gen.HealthPoint{
			CurrentHp: 18,
			MaxHp:     22,
			TempHp:    3,
			HitDice: map[int32]*gen.Amount{
				int32(dice.D8): {MaxAvailable: 4, Available: 3},
			},
		}

		dto := dto.HealthDTOFromProto(proto)
		assert.Equal(t, 18, dto.CurrentHP)
		assert.Equal(t, 22, dto.MaxHP)
		assert.Equal(t, 3, dto.TempHP)
		assert.Equal(t, 3, dto.HitDice[dice.D8].Available)
	})

	t.Run("FromDomain", func(t *testing.T) {
		domain := &health.Health{
			CurrentHP: 40,
			MaxHP:     45,
			TempHP:    2,
			HitDice: map[dice.Dice]*health.Amount{
				dice.D10: {MaxAvailable: 5, Available: 4},
			},
		}

		dto := dto.HealthDTOFromDomain(domain)
		assert.Equal(t, 40, dto.CurrentHP)
		assert.Equal(t, 45, dto.MaxHP)
		assert.Equal(t, 2, dto.TempHP)
		assert.Equal(t, 4, dto.HitDice[dice.D10].Available)
	})
}
