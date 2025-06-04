package dto_test

import (
	"testing"

	gen "github.com/alsadx/GM-Tool/character-service/gen/character"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/skill"
	"github.com/stretchr/testify/assert"
)

func TestSkillDTO_Conversions(t *testing.T) {
	t.Run("ToDomain", func(t *testing.T) {
		dto := &dto.SkillDTO{Bonus: 5}
		domain := dto.ToDomain()
		assert.Equal(t, 5, domain.Bonus)
	})

	t.Run("ToProto", func(t *testing.T) {
		dto := &dto.SkillDTO{Bonus: 3}
		proto := dto.ToProto()
		assert.Equal(t, int32(3), proto.Bonus)
	})

	t.Run("FromDomain", func(t *testing.T) {
		domain := &skill.Skill{Bonus: 7}
		dto := dto.SkillDTOFromDomain(domain)
		assert.Equal(t, 7, dto.Bonus)
	})

	t.Run("FromProto", func(t *testing.T) {
		proto := &gen.Skill{Bonus: 2}
		dto := dto.SkillDTOFromProto(proto)
		assert.Equal(t, 2, dto.Bonus)
	})
}
