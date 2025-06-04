package character

import (
	"context"
	"errors"
	"fmt"

	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/internal/repository"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
	"go.uber.org/zap"
)

type characterRepository interface {
	Get(ctx context.Context, id string) (*character.Character, error)
	Create(ctx context.Context, c *character.Character) error
	GetByOwnerID(ctx context.Context, ownerID int) ([]*character.Character, error)
	Update(ctx context.Context, c *character.Character) error
	Delete(ctx context.Context, id string) error
}

type Controller struct {
	repo   characterRepository
	logger *zap.Logger
}

func New(repo characterRepository) *Controller {
	return &Controller{
		repo:   repo,
		logger: zap.L().Named("character_controller"),
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

func (ctrl *Controller) SetMaxHP(ctx context.Context, id string, maxHP int) (*dto.CharacterDTO, error) {
	if maxHP <= 0 {
		return nil, fmt.Errorf("invalid maxHP value: %d", maxHP)
	}
	
	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.SetMaxHP(maxHP)
	
	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after setting max HP",
			zap.String("id", id), 
			zap.Int("maxHP", maxHP),
			zap.Error(err))
		return nil, err
	}
	return dto.CharacterDTOFromDomain(char), nil
}

func (ctrl *Controller) SetTempHP(ctx context.Context, id string, tempHP int) (*dto.CharacterDTO, error) {
	if tempHP <= 0 {
		return nil, fmt.Errorf("invalid tempHP value: %d", tempHP)
	}

	char, err := ctrl.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	char.SetTempHP(tempHP)

	if err := ctrl.repo.Update(ctx, char); err != nil {
		ctrl.logger.Error("failed to update character after setting temp HP",
			zap.String("id", id), 
			zap.Int("tempHP", tempHP),
			zap.Error(err))
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