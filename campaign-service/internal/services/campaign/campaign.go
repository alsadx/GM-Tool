package campaign

import (
	"campaigntool/internal/domain/models"
	"campaigntool/internal/lib/code"
	"context"
	"errors"
	"fmt"
	"log/slog"
)

type GameSaver interface {
	SaveCampaign(ctx context.Context, name, desc string, userId int64) (int64, error)
	DeleteCampaign(ctx context.Context, campaignId int64, userId int64) error
	UpdateCampaign(ctx context.Context, campaign *models.CampaignInfo) error
	AddPlayer(ctx context.Context, campaignId int64, userId int64, charId int64) error
	RemovePlayer(ctx context.Context, campaignId int64, userId int64) error
	SetInviteCode(ctx context.Context, campaignId int64, masterId int64, inviteCode string) error
	AddCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) error
	RemoveCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) error
}

type GameProvider interface {
	GetCampaign(ctx context.Context, campaignId int64) (*models.CampaignInfo, error)
	CheckInviteCode(ctx context.Context, inviteCode string) (int64, error)
	CreatedCampaigns(ctx context.Context, userId int64) ([]*models.Campaign, error)
	CurrentCampaigns(ctx context.Context, userId int64) ([]*models.CampaignForPlayer, error)
	GetCampaignPlayers(ctx context.Context, campaignId int64) ([]int64, error)
	GetCampaignCharacters(ctx context.Context, campaignId int64) ([]*models.CampaignCharacter, error)
	GetPlayerCharacters(ctx context.Context, campaignId int64, userId int64) ([]int64, error)
	IsMaster(ctx context.Context, campaignId int64, userId int64) (bool, error)
	IsPlayer(ctx context.Context, campaignId int64, userId int64) (bool, error)
	IsCharacterOwner(ctx context.Context, campaignId int64, userId int64, charId int64) (bool, error)
}

type CampaignTool struct {
	Log          *slog.Logger
	GameSaver    GameSaver
	GameProvider GameProvider
}

func New(log *slog.Logger, gameSaver GameSaver, gameProvider GameProvider) *CampaignTool {
	return &CampaignTool{
		Log:          log,
		GameSaver:    gameSaver,
		GameProvider: gameProvider,
	}
}

func (s *CampaignTool) CreateCampaign(ctx context.Context, name, desc string, userId int64) (campaignId int64, err error) {
	op := "campaign.CreateCampaign"

	log := s.Log.With(slog.String("op", op), slog.String("name", name))

	log.Info("creating new campaign")

	campaignId, err = s.GameSaver.SaveCampaign(ctx, name, desc, userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignExists) {
			s.Log.Error("campaign already exists", slog.String("error", err.Error()))

			return 0, fmt.Errorf("%s: %w", op, models.ErrCampaignExists)
		}
		s.Log.Error("failed to save campaign", slog.String("error", err.Error()))

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("campaign created", slog.Int("campaignId", int(campaignId)))

	return campaignId, nil
}

func (s *CampaignTool) DeleteCampaign(ctx context.Context, campaignId int64, userId int64) (err error) {
	op := "campaign.DeleteCampaign"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("deleting campaign")

	ok, err := s.GameProvider.IsMaster(ctx, campaignId, userId) 
	if err != nil {
		s.Log.Error("failed to check if user is master", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if !ok {
		s.Log.Error("user is not master", slog.String("error", models.ErrNotMaster.Error()))

		return fmt.Errorf("%s: %w", op, models.ErrNotMaster)
	}

	players, err := s.GameProvider.GetCampaignPlayers(ctx, campaignId)
	if err != nil{
		s.Log.Error("failed to get campaign players", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if len(players) > 0 {
		log.Info("removing players from campaign")
		for _, player := range players {
			err = s.GameSaver.RemovePlayer(ctx, campaignId, player)
			if err != nil {
				s.Log.Error("failed to remove player", slog.String("error", err.Error()))

				return fmt.Errorf("%s: %w", op, err)
			}
		}
		log.Info("removed players from campaign")
	}

	err = s.GameSaver.DeleteCampaign(ctx, campaignId, userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}
		s.Log.Error("failed to delete campaign", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("campaign deleted")

	return nil
}

func (s *CampaignTool) UpdateCampaign(ctx context.Context, campaignId int64, name, desc string, userId int64) (err error) {
	op := "campaign.UpdateCampaign"
	change := false

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("updating campaign")

	ok, err := s.GameProvider.IsMaster(ctx, campaignId, userId)
	if err != nil {
		s.Log.Error("failed to check if user is master", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if !ok {
		s.Log.Error("user is not master", slog.String("error", models.ErrNotMaster.Error()))

		return fmt.Errorf("%s: %w", op, models.ErrNotMaster)
	}

	updatedCampaign, err := s.GameProvider.GetCampaign(ctx, campaignId)
	if err != nil {
		s.Log.Error("failed to get campaign", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if name != "" && name != updatedCampaign.Name {
		updatedCampaign.Name = name
		change = true
	}

	if desc != "" && desc != updatedCampaign.Description {
		updatedCampaign.Description = desc
		change = true
	}

	if !change {
		log.Info("no changes detected")

		return fmt.Errorf("%s: %w", op, models.ErrNoChanges)
	}

	err = s.GameSaver.UpdateCampaign(ctx, updatedCampaign)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}
		s.Log.Error("failed to update campaign", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("campaign updated")

	return nil
}

func (s *CampaignTool) GenerateInviteCode(ctx context.Context, campaignId int64, userId int64) (inviteCode string, err error) {
	op := "campaign.GenerateInviteCode"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("generating invite code", slog.Int("userId", int(userId)))

	inviteCode, err = code.GenerateCode()
	if err != nil {
		s.Log.Error("failed to generate invite code", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.GameSaver.SetInviteCode(ctx, campaignId, userId, inviteCode)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return "", fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}
		s.Log.Error("failed to set invite code", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("invite code generated", slog.String("inviteCode", inviteCode))

	return inviteCode, nil
}

func (s *CampaignTool) JoinCampaign(ctx context.Context, inviteCode string, userId int64, charId int64) (err error) {
	op := "campaign.JoinCampaign"

	log := s.Log.With(slog.String("op", op), slog.String("inviteCode", inviteCode))

	log.Info("attempting to join campaign", slog.Int("userId", int(userId)))

	campaignId, err := s.GameProvider.CheckInviteCode(ctx, inviteCode)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrInvalidCode)
		}
		s.Log.Error("failed to check invite code", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.GameSaver.AddPlayer(ctx, campaignId, userId, charId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		} else if errors.Is(err, models.ErrPlayerInCampaign) {
			s.Log.Error("player already in campaign", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrPlayerInCampaign)
		}

		s.Log.Error("failed to add player", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("joined campaign", slog.Int("userId", int(userId)))

	return nil
}

func (s *CampaignTool) LeaveCampaign(ctx context.Context, campaignId int64, userId int64) (err error) {
	op := "campaign.LeaveCampaign"

	log := s.Log.With(slog.String("op", op))

	log.Info("attempting to leave campaign")

	isPlayer, err := s.GameProvider.IsPlayer(ctx, campaignId, userId)
	if err != nil {
		s.Log.Error("failed to check if user is player", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if !isPlayer {
		s.Log.Error("user is not player", slog.String("error", models.ErrNotPlayer.Error()))

		return fmt.Errorf("%s: %w", op, models.ErrNotPlayer)
	}

	err = s.GameSaver.RemovePlayer(ctx, campaignId, userId)
	if err != nil {
		log.Error("failed to leave campaign", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("left campaign")

	return nil
}

func (s *CampaignTool) RemovePlayer(ctx context.Context, campaignId int64, userId int64, masterId int64) (err error) {
	op := "campaign.RemovePlayer"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("removing players from campaign")

	ok, err := s.GameProvider.IsMaster(ctx, campaignId, masterId)
	if err != nil {
		s.Log.Error("failed to check if user is master", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if !ok {
		s.Log.Error("user is not master", slog.String("error", models.ErrNotMaster.Error()))

		return fmt.Errorf("%s: %w", op, models.ErrNotMaster)
	}

	err = s.GameSaver.RemovePlayer(ctx, campaignId, userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}
		s.Log.Error("failed to remove players", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("removed players from campaign")

	return nil
}	

func (s *CampaignTool) AddCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) (err error) {
	op := "campaign.AddCharacter"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)), slog.Int("userId", int(userId)))

	log.Info("adding character to campaign")

	isMaster, err := s.GameProvider.IsMaster(ctx, campaignId, userId)
	if err != nil {
		s.Log.Error("failed to check if user is master", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	isPlayer, err := s.GameProvider.IsPlayer(ctx, campaignId, userId)
	if err != nil {
		s.Log.Error("failed to check if user is player", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if !isMaster && !isPlayer {
		s.Log.Error("user is not master or player", slog.String("error", models.ErrNotPlayer.Error()))

		return fmt.Errorf("%s: %w", op, models.ErrNotPlayer)
	}

	err = s.GameSaver.AddCharacter(ctx, campaignId, userId, charId)
	if err != nil {

		if errors.Is(err, models.ErrCharacterInCampaign) {
			s.Log.Error("player already in campaign", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrCharacterInCampaign)
		} else if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}

		s.Log.Error("failed to add character", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("added character to campaign")

	return nil
}

func (s *CampaignTool) RemoveCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) (err error) {
	op := "campaign.RemoveCharacter"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)), slog.Int("userId", int(userId)))

	log.Info("removing character from campaign")

	isMaster, err := s.GameProvider.IsMaster(ctx, campaignId, userId)
	if err != nil {
		s.Log.Error("failed to check if user is master", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	isOwner, err := s.GameProvider.IsCharacterOwner(ctx, campaignId, userId, charId)
	if err != nil {
		s.Log.Error("failed to check if user is owner", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	if !isMaster && !isOwner {
		s.Log.Error("user is not master or owner", slog.String("error", models.ErrNotMaster.Error()))

		return fmt.Errorf("%s: %w", op, models.ErrNotCharacterOwner)
	}

	err = s.GameSaver.RemoveCharacter(ctx, campaignId, userId, charId)
	if err != nil {
		if errors.Is(err, models.ErrCharacterNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrCharacterNotFound)
		}
		s.Log.Error("failed to remove character", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("removed character from campaign")

	return nil
}

func (s *CampaignTool) GetCampaignPlayers(ctx context.Context, campaignId int64) (players []int64, err error) {
	op := "campaign.GetCampaignPlayers"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("getting campaign players")

	players, err = s.GameProvider.GetCampaignPlayers(ctx, campaignId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return nil, fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}
		s.Log.Error("failed to get campaign players", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got campaign players")

	return players, nil
}

func (s *CampaignTool) GetCreatedCampaigns(ctx context.Context, userId int64) (campaigns []*models.Campaign, err error) {
	op := "campaign.GetCreatedCampaigns"

	log := s.Log.With(slog.String("op", op), slog.Int("user_id", int(userId)))

	log.Info("getting created campaigns")

	campaigns, err = s.GameProvider.CreatedCampaigns(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrNoCampaigns) {
			s.Log.Error("created campaigns not found", slog.String("error", err.Error()))

			return nil, fmt.Errorf("%s: %w", op, models.ErrNoCampaigns)
		}
		log.Error("failed to get created campaigns", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got created campaigns")

	return campaigns, nil
}

func (s *CampaignTool) GetCurrentCampaigns(ctx context.Context, userId int64) (campaigns []*models.CampaignForPlayer, err error) {
	op := "campaign.GetCreatedCampaigns"

	log := s.Log.With(slog.String("op", op), slog.Int("user_id", int(userId)))

	log.Info("getting current campaigns")

	campaigns, err = s.GameProvider.CurrentCampaigns(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrNoCampaigns) {
			s.Log.Error("current campaigns not found", slog.String("error", err.Error()))

			return nil, fmt.Errorf("%s: %w", op, models.ErrNoCampaigns)
		}
		log.Error("failed to get current campaigns", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got current campaigns")

	return campaigns, nil
}

func (s *CampaignTool) GetCampaignCharacters(ctx context.Context, campaignId int64, userId int64) (characters []*models.CampaignCharacter, err error) {
	op := "campaign.GetCampaignCharacters"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("getting campaign characters")

	isMaster, err := s.GameProvider.IsMaster(ctx, campaignId, userId)
	if err != nil {
		s.Log.Error("failed to check if user is master", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !isMaster {
		s.Log.Error("user is not master", slog.String("error", models.ErrNotMaster.Error()))

		return nil, fmt.Errorf("%s: %w", op, models.ErrNotMaster)
	}

	characters, err = s.GameProvider.GetCampaignCharacters(ctx, campaignId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return nil, fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}
		s.Log.Error("failed to get campaign characters", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got campaign characters")

	return characters, nil
}

func (s *CampaignTool) GetPlayerCharacters(ctx context.Context, campaignId int64, playerId, userId int64) (characters []int64, err error) {
	op := "campaign.GetPlayerCharacters"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)), slog.Int("userId", int(userId)))

	log.Info("getting player characters")

	ok, err := s.GameProvider.IsMaster(ctx, campaignId, userId)
	if err != nil {
		s.Log.Error("failed to check if user is master", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if !ok && playerId != userId {
		s.Log.Error("user is not master", slog.String("error", models.ErrNotCharacterOwner.Error()))

		return nil, fmt.Errorf("%s: %w", op, models.ErrNotCharacterOwner)
	}

	characters, err = s.GameProvider.GetPlayerCharacters(ctx, campaignId, playerId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return nil, fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}
		s.Log.Error("failed to get player characters", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got player characters")

	return characters, nil
}
