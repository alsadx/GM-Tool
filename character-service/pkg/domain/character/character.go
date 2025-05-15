package character

import (
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/armor"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/level"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
)

type Character struct {
	ID    int
	Owner int

	Name  string
	Class string
	Race  string
	Lvl   *level.LevelSystem
	Armor *armor.Armor

	Stats map[types.AbilityType][]ability.Ability

	Health *health.HealthPoint
}
