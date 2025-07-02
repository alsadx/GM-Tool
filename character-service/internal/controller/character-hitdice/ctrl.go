package hitdice_controller

import (
	"context"

	"github.com/alsadx/GM-Tool/character-service/internal/controller"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"go.uber.org/zap"
)

type Controller struct {
	repo   controller.CharacterRepository
	logger *zap.Logger
}

func New(repo controller.CharacterRepository) *Controller {
	return &Controller{
		repo:   repo,
		logger: zap.L().Named("hitdice_controller"),
	}
}

func (ctrl *Controller) GetHealthDTO(ctx context.Context, id string) (*dto.HealthDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char).Health, nil
}

func (ctrl *Controller) AddHitDice(ctx context.Context, id string, hitDice dice.Dice) (*dto.HealthDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = char.AddHitDice(hitDice)
	if err != nil {
		return nil, err
	}

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after add hit dice",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char).Health, nil
}

func (ctrl *Controller) RemoveHitDice(ctx context.Context, id string, hitDice dice.Dice) (*dto.HealthDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = char.RemoveHitDice(hitDice)
	if err != nil {
		return nil, err
	}

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after remove hit dice",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char).Health, nil
}

func (ctrl *Controller) RollHitDice(ctx context.Context, id string, rollingDice map[dice.Dice]int) (*dto.RollHitDiceResultDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	res, bonus, err := char.RollHitDice(rollingDice)
	if err != nil {
		return nil, err
	}

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after roll hit dice",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}

	resInt32 := make([]int32, len(res))

	for diceType, amount := range res {
		resInt32[diceType] = int32(amount)
	}

	rollDto := &dto.RollHitDiceResultDTO{
		Result: resInt32,
		Bonus:  int32(bonus),
		Health: dto.HealthDTOFromDomain(char.GetHealth()),
	}

	return rollDto, nil
}

func (ctrl *Controller) ResetHitDice(ctx context.Context, id string, resetDice map[dice.Dice]int) (*dto.HealthDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	err = char.ResetHitDice(resetDice)
	if err != nil {
		return nil, err
	}

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after reset hit dice",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char).Health, nil
}
