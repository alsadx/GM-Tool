package controller

import (
	"context"
	"fmt"

	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
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

func (ctrl *Controller) GainExp(ctx context.Context, id string, amount int) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.GainExp(amount)
	
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after exp gain",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) RemoveExp(ctx context.Context, id string, amount int) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.RemoveExp(amount)

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after remove exp",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) LevelUp(ctx context.Context, id string) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !char.LvlUp() {
		ctrl.logger.Warn(
			"cant lvl up for this character",
			zap.Int("current_exp", char.GetCurrentExp()),
			zap.Int("exp_to_next_lvl", char.ExpToNextLevel()),
		)
		return nil, fmt.Errorf("level up not possible")
	}
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after lvl up",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) LevelDown(ctx context.Context, id string) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !char.LvlDown() {
		ctrl.logger.Warn(
			"cant lvl down for this character",
			zap.Int("current_exp", char.GetCurrentExp()),
		)
		return nil, fmt.Errorf("level up not possible")
	}
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after lvl down",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
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

func (ctrl *Controller) CheckAbility(ctx context.Context, id string, abil types.AbilityType) (diseRes, bonus, result int, err error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	diseRes, bonus, result = char.CheckAbility(abil)
	return diseRes, bonus, result, nil
}

func (ctrl *Controller) CheckSkill(ctx context.Context, id string, skill types.SkillType) (diseRes, bonus, result int, err error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	diseRes, bonus, result = char.CheckSkill(skill)
	return diseRes, bonus, result, nil
}