package progress_controller

import (
	"context"

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
		logger: zap.L().Named("progress_controller"),
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

func (ctrl *Controller) LvlUp(ctx context.Context, id string) (*dto.CharacterDTO, error) {
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
		return dto.CharacterDTOFromDomain(char), ErrCantLvlUpDown
	}
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after lvl up",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) LvlDonw(ctx context.Context, id string) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if !char.LvlDown() {
		ctrl.logger.Warn(
			"cant lvl down for this character",
			zap.Int("current_exp", char.GetCurrentExp()),
		)
		return nil, ErrCantLvlUpDown
	}
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after lvl down",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) CanLvlUp(ctx context.Context, id string) (bool, *dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return false, nil, err
	}
	return char.CanLvlUp(), dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) CanLvlDown(ctx context.Context, id string) (bool, *dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return false, nil, err
	}
	return char.CanLvlDown(), dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) CurrentLvl(ctx context.Context, id string) (int, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return 0, err
	}
	return char.GetLvl(), nil
}

func (ctrl *Controller) ExpToNextLvl(ctx context.Context, id string) (int, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return 0, err
	}
	return char.ExpToNextLevel(), nil
}

func (ctrl *Controller) SetLvl(ctx context.Context, id string, lvl int) (*dto.CharacterDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.SetLvl(lvl)
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after set lvl",
			zap.String("id", id), zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) GetLvlState(ctx context.Context, id string) (*dto.LvlStateDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.LvlStateDTO{
		CurrentLvl:   char.GetLvl(),
		CurrentExp:   char.GetCurrentExp(),
		ExpToNextLvl: char.ExpToNextLevel(),
	}, nil
}
