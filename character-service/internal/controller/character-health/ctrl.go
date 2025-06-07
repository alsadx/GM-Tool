package health_controller

import (
	"context"
	"fmt"

	"github.com/alsadx/GM-Tool/character-service/internal/controller"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"go.uber.org/zap"
)

type Controller struct {
	repo   controller.CharacterRepository
	logger *zap.Logger
}

func New(repo controller.CharacterRepository) *Controller {
	return &Controller{
		repo:   repo,
		logger: zap.L().Named("health_controller"),
	}
}

func (ctrl *Controller) SetMaxHP(ctx context.Context, id string, max_hp int) (*dto.CharacterDTO, error) {
	if max_hp <= 0 {
		return nil, fmt.Errorf("invalid max_hp value: %d", max_hp)
	}

	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.SetMaxHP(max_hp)

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after setting max HP",
			zap.String("id", id),
			zap.Int("max_hp", max_hp),
			zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) SetCurrentHP(ctx context.Context, id string, current_hp int) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.SetCurrentHP(current_hp)

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after setting current HP",
			zap.String("id", id),
			zap.Int("current_hp", current_hp))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) SetTempHP(ctx context.Context, id string, temp_hp int) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.SetTempHP(temp_hp)

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after setting temp HP",
			zap.String("id", id),
			zap.Int("temp_hp", temp_hp),
			zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) Heal(ctx context.Context, id string, amount int) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.Heal(amount)

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after heal",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) TakeDamage(ctx context.Context, id string, damage int) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.TakeDamage(damage)

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after take damage",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) IsKnocked(ctx context.Context, id string) (bool, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return false, err
	}

	return char.IsKnocked(), nil
}

func (ctrl *Controller) GetHPState(ctx context.Context, id string) (*dto.CharacterHPState, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	return &dto.CharacterHPState{
		MaxHP:     char.GetMaxHP(),
		CurrentHP: char.GetCurrentHP(),
		TempHP:    char.GetTempHP(),
		IsKnocked: char.IsKnocked(),
	}, nil
}
