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
	SaveCampaign(ctx context.Context, name, desc string, userId int) (int32, error)
	DeleteCampaign(ctx context.Context, campaignId int32, userId int) error
	AddPlayer(ctx context.Context, campaignId int32, userId int) error
	RemovePlayer(ctx context.Context, campaignId int32, userId int) error
	SetInviteCode(ctx context.Context, campaignId int32, inviteCode string) error
}

type GameProvider interface {
	CheckInviteCode(ctx context.Context, inviteCode string) (int32, error)
	CreatedCampaigns(ctx context.Context, userId int) ([]models.Campaign, error)
	CurrentCampaigns(ctx context.Context, userId int) ([]models.CampaignForPlayer, error)
}

type CampaignTool struct {
	Log          *slog.Logger
	GameSaver    GameSaver
	GameProvider GameProvider
}

// var (
// 	ErrInvalidCredentials = errors.New("invalid credentials")
// 	ErrCampaignNotFound   = errors.New("campaign not found")
// 	ErrCampaignExists     = errors.New("campaign with this name already exists")
// 	ErrInvalidCode        = errors.New("invalid invite code")
// 	ErrPlayerInCampaign   = errors.New("player already in campaign")
// )

func New(log *slog.Logger, gameSaver GameSaver, gameProvider GameProvider) *CampaignTool {
	return &CampaignTool{
		Log:          log,
		GameSaver:    gameSaver,
		GameProvider: gameProvider,
	}
}

func (s *CampaignTool) CreateCampaign(ctx context.Context, name, desc string, userId int) (campaignId int32, err error) {
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

func (s *CampaignTool) DeleteCampaign(ctx context.Context, campaignId int32, userId int) (err error) {
	op := "campaign.DeleteCampaign"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("deleting campaign")

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

func (s *CampaignTool) GenerateInviteCode(ctx context.Context, campaignId int32, userId int) (inviteCode string, err error) {
	op := "campaign.GenerateInviteCode"

	log := s.Log.With(slog.String("op", op), slog.Int("campaignId", int(campaignId)))

	log.Info("generating invite code")

	inviteCode, err = code.GenerateCode()
	if err != nil {
		s.Log.Error("failed to generate invite code", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	err = s.GameSaver.SetInviteCode(ctx, campaignId, inviteCode)
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

func (s *CampaignTool) JoinCampaign(ctx context.Context, inviteCode string, userId int) (err error) {
	op := "campaign.JoinCampaign"

	log := s.Log.With(slog.String("op", op), slog.String("inviteCode", inviteCode))

	log.Info("attempting to join campaign")

	campaignId, err := s.GameProvider.CheckInviteCode(ctx, inviteCode)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			s.Log.Error("campaign not found", slog.String("error", err.Error()))

			return fmt.Errorf("%s: %w", op, models.ErrInvalidCode)
		}
		s.Log.Error("failed to check invite code", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.GameSaver.AddPlayer(ctx, campaignId, userId)
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

	log.Info("joined campaign")

	return nil
}

func (s *CampaignTool) LeaveCampaign(ctx context.Context, campaignId int32, userId int) (err error) {
	op := "campaign.LeaveCampaign"

	log := s.Log.With(slog.String("op", op))

	log.Info("attempting to leave campaign")

	err = s.GameSaver.RemovePlayer(ctx, campaignId, userId)
	if err != nil {
		log.Error("failed to leave campaign", slog.String("error", err.Error()))

		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("left campaign")

	return nil
}

func (s *CampaignTool) GetCreatedCampaigns(ctx context.Context, userId int) (campaigns []models.Campaign, err error) {
	op := "campaign.GetCreatedCampaigns"

	log := s.Log.With(slog.String("op", op), slog.Int("id", int(userId)))

	log.Info("getting created campaigns")

	campaigns, err = s.GameProvider.CreatedCampaigns(ctx, userId)
	if err != nil {
		log.Error("failed to get created campaigns", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got created campaigns")

	return campaigns, nil
}

func (s *CampaignTool) GetCurrentCampaigns(ctx context.Context, userId int) (campaigns []models.CampaignForPlayer, err error) {
	op := "campaign.GetCreatedCampaigns"

	log := s.Log.With(slog.String("op", op), slog.Int("id", int(userId)))

	log.Info("getting current campaigns")

	campaigns, err = s.GameProvider.CurrentCampaigns(ctx, userId)
	if err != nil {
		log.Error("failed to get current campaigns", slog.String("error", err.Error()))

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("got current campaigns")

	return campaigns, nil
}
