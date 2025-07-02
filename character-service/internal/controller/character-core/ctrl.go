package core_controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/alsadx/GM-Tool/character-service/internal/controller"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/internal/repository"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"go.uber.org/zap"
)

type Controller struct {
	repo   controller.CharacterRepository
	logger *zap.Logger
}

func New(repo controller.CharacterRepository) *Controller {
	return &Controller{
		repo:   repo,
		logger: zap.L().Named("core_controller"),
	}
}

func (ctrl *Controller) Create(ctx context.Context, ownerID int, name, class, subclass, race string) (*dto.CharacterDTO, error) {
	ctrl.logger.Debug("creating character",
		zap.Int("ownerID", ownerID),
		zap.String("name", name),
	)

	char, err := character.NewWitoutID(ownerID, name, class, subclass, race)
	if err != nil {
		ctrl.logger.Error("failed to create character", zap.Error(err))
		return nil, err
	}

	if err := ctrl.repo.Create(ctx, char); err != nil {
		ctrl.logger.Error("failed to save character", zap.Error(err))
		return nil, fmt.Errorf("failed to save character: %w", err)
	}

	ctrl.logger.Info("character created", zap.String("id", char.ID))
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) Get(ctx context.Context, id string) (*dto.CharacterDTO, error) {
	ctrl.logger.Debug("getting character", zap.String("id", id))

	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		ctrl.logger.Error("failed to get character", zap.Error(err))
		return nil, fmt.Errorf("failed to get character: %w", err)
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) Update(ctx context.Context, char *dto.UpdateCharDTO) error {
	ctrl.logger.Debug("updating character", zap.String("id", char.ID))

	existChar, err := ctrl.repo.Get(ctx, char.ID)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		ctrl.logger.Warn("character not found", zap.String("id", char.ID))
		return err
	} else if err != nil {
		ctrl.logger.Error("failed to get character", zap.Error(err))
		return err
	}
	if char.Name != "" {
		existChar.Name = char.Name
	}
	if char.Class != "" {
		existChar.Class = char.Class
	}
	if char.Race != "" {
		existChar.Race = char.Race
	}
	if char.Subclass != "" {
		existChar.Subclass = char.Subclass
	}

	if err := ctrl.repo.Update(ctx, existChar); err != nil {
		ctrl.logger.Error("failed to update character", zap.Error(err))
		return fmt.Errorf("failed to update character: %w", err)
	}
	return nil
}

func (ctrl *Controller) Delete(ctx context.Context, id string) error {
	ctrl.logger.Debug("deleting character", zap.String("id", id))

	if err := ctrl.repo.Delete(ctx, id); err != nil {
		ctrl.logger.Error("failed to delete character", zap.Error(err))
		return fmt.Errorf("failed to delete character: %w", err)
	}
	return nil
}

func (ctrl *Controller) GetByOwnerID(ctx context.Context, ownerID int) ([]*dto.CharacterDTO, error) {
	ctrl.logger.Debug("getting characters by owner", zap.Int("ownerID", ownerID))

	chars, err := ctrl.repo.GetByOwnerID(ctx, ownerID)
	if err != nil {
		ctrl.logger.Error("failed to get characters by owner", zap.Error(err))
		return nil, fmt.Errorf("failed to get characters: %w", err)
	}
	dtoChars := make([]*dto.CharacterDTO, len(chars))
	for i, char := range chars {
		dtoChar := dto.CharacterDTOFromDomain(char)
		dtoChars[i] = dtoChar
	}
	return dtoChars, nil
}

func (ctrl *Controller) GetInfoAboutCharacter(ctx context.Context, id string) (*dto.CharacterInfoDTO, error) {
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return &dto.CharacterInfoDTO{
		OwnerID:  char.Owner,
		Name:     char.Name,
		Class:    char.Class,
		Subclass: char.Subclass,
		Race:     char.Race,
	}, nil
}
