package dto_test

import (
	"testing"

	gen "github.com/alsadx/GM-Tool/character-service/gen/character"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/level"
	"github.com/stretchr/testify/assert"
)

func TestLevelDTO_Conversions(t *testing.T) {
	t.Run("ToProto", func(t *testing.T) {
		dto := &dto.LevelDTO{CurrentExp: 2700}
		proto := dto.ToProto()
		assert.Equal(t, int32(2700), proto.CurrentExp)
	})

	t.Run("ToDomain", func(t *testing.T) {
		dto := &dto.LevelDTO{CurrentExp: 300}
		domain := dto.ToDomain()
		assert.Equal(t, 300, domain.CurrentExp())
	})

	t.Run("FromProto", func(t *testing.T) {
		proto := &gen.LevelSystem{CurrentExp: 4500}
		dto := dto.LevelDTOFromProto(proto)
		assert.Equal(t, 4500, dto.CurrentExp)
	})

	t.Run("FromDomain", func(t *testing.T) {
		domain := level.NewLevelSystem()
		domain.AddExp(600)
		dto := dto.LevelDTOFromDomain(domain)
		assert.Equal(t, 600, dto.CurrentExp)
	})
}
