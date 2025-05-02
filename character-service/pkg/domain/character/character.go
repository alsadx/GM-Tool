package character

import (
	"sync"

	"github.com/alsadx/GM-Tool/character-service/pkg/domain/ability"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/armor"
)

type Character struct {
	mu        sync.RWMutex
	abilities map[ability.Ability]*ability.Score
	armor     *armor.Armor
	shield    bool
}
