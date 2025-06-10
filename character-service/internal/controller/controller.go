package controller

import (
	"context"

	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"go.uber.org/zap"
)

type CharacterRepository interface {
	Get(ctx context.Context, id string) (*character.Character, error)
	Create(ctx context.Context, c *character.Character) error
	GetByOwnerID(ctx context.Context, ownerID int) ([]*character.Character, error)
	Update(ctx context.Context, c *character.Character) error
	Delete(ctx context.Context, id string) error
}

type Controller struct {
	repo   CharacterRepository
	logger *zap.Logger
}

func New(repo CharacterRepository) *Controller {
	return &Controller{
		repo:   repo,
		logger: zap.L().Named("character_controller"),
	}
}

func (ctrl *Controller) AddHitDice(ctx context.Context, id string, diceType dice.Dice) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := char.AddHitDice(diceType); err != nil {
		ctrl.logger.Warn("domain error", zap.Error(err))
		return nil, err
	}
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after add hit dice",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) RemoveHitDice(ctx context.Context, id string, diceType dice.Dice) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := char.RemoveHitDice(diceType); err != nil {
		ctrl.logger.Warn("domain error", zap.Error(err))
		return nil, err
	}
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after remove hit dice",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return dto.CharacterDTOFromDomain(char), nil
}
