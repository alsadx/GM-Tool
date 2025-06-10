package stats_controller

import (
	"context"

	"github.com/alsadx/GM-Tool/character-service/internal/controller"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
	"go.uber.org/zap"
)

type Controller struct {
	repo   controller.CharacterRepository
	logger *zap.Logger
}

func New(repo controller.CharacterRepository) *Controller {
	return &Controller{
		repo:   repo,
		logger: zap.L().Named("stats_controller"),
	}
}

func (ctrl *Controller) CheckAbility(ctx context.Context, id string, abil types.AbilityType) (*dto.DiceResultDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	diseRes, bonus, result := char.CheckAbility(abil)
	return &dto.DiceResultDTO{DiceRes: diseRes, Bonus: bonus, Result: result}, nil
}

func (ctrl *Controller) CheckSkill(ctx context.Context, id string, skill types.SkillType) (*dto.DiceResultDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	diseRes, bonus, result := char.CheckSkill(skill)
	return &dto.DiceResultDTO{DiceRes: diseRes, Bonus: bonus, Result: result}, nil
}

func (ctrl *Controller) SetAbilityScore(ctx context.Context, id string, abil types.AbilityType, base int) (*dto.AbilityDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.Ability(abil).SetBase(base)
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after set ability score",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.AbilityDTOFromDomain(char.Ability(abil)), nil
}

func (ctrl *Controller) SetAbilityBonus(ctx context.Context, id string, abil types.AbilityType, bonus int) (*dto.AbilityDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.Ability(abil).SetBonus(bonus)
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after set ability bonus",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}

	return dto.AbilityDTOFromDomain(char.Ability(abil)), nil
}

func (ctrl *Controller) SetSkillBonus(ctx context.Context, id string, skill types.SkillType, bonus int) (*dto.SkillDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.Skill(skill).SetBonus(bonus)
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after set skill bonus",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.SkillDTOFromDomain(char.Skill(skill)), nil
}

func (ctrl *Controller) GetStats(ctx context.Context, id string) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) GetModifier(ctx context.Context, id string, abil types.AbilityType) (int, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return 0, err
	}
	return char.Ability(abil).Modifier(), nil
}

func (ctrl *Controller) GetAbility(ctx context.Context, id string, abil types.AbilityType) (*dto.AbilityDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.AbilityDTOFromDomain(char.Ability(abil)), nil
}