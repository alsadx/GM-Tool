package handler_stats

import (
	"context"
	"errors"
	"fmt"

	"github.com/alsadx/GM-Tool/character-service/gen/character"
	character_stats "github.com/alsadx/GM-Tool/character-service/gen/service/character-stats"
	stats_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-stats"
	"github.com/alsadx/GM-Tool/character-service/internal/repository"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	character_stats.UnimplementedCharacterStatsServiceServer
	ctrl   *stats_controller.Controller
	logger *zap.Logger
}

func New(ctrl *stats_controller.Controller) *Handler {
	return &Handler{
		ctrl:   ctrl,
		logger: zap.L().Named("stats_handler"),
	}
}

func (h *Handler) CheckAbility(ctx context.Context, req *character_stats.AbilityRequest) (*character_stats.CheckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	abilityType, err := types.ParseAbilityType(req.Ability)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect name for the ability")
	}

	h.logger.Info("got request", zap.Any("request", req))

	diceResult, err := h.ctrl.CheckAbility(ctx, req.Id, abilityType)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to check ability")
		return nil, err
	}

	return &character_stats.CheckResponse{
		DiceResult: int32(diceResult.DiceRes),
		Bonus:      int32(diceResult.Bonus),
		Total:      int32(diceResult.Result),
	}, nil
}

func (h *Handler) CheckSkill(ctx context.Context, req *character_stats.SkillRequest) (*character_stats.CheckResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	skillType, err := types.ParseSkillType(req.Skill)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect name for the skill")
	}

	h.logger.Info("got request", zap.Any("request", req))

	diceResult, err := h.ctrl.CheckSkill(ctx, req.Id, skillType)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to check skill")
		return nil, err
	}

	return &character_stats.CheckResponse{
		DiceResult: int32(diceResult.DiceRes),
		Bonus:      int32(diceResult.Bonus),
		Total:      int32(diceResult.Result),
	}, nil
}

func (h *Handler) SetAbilityScore(ctx context.Context, req *character_stats.SetAbilityRequest) (*character_stats.SetAbilityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("amount = %d must be greater than or equal 0", req.Amount))
	}

	abilityType, err := types.ParseAbilityType(req.Ability)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect name for the ability")
	}

	h.logger.Info("got request", zap.Any("request", req))

	abilDto, err := h.ctrl.SetAbilityScore(ctx, req.Id, abilityType, int(req.Amount))
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to set ability score", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_stats.SetAbilityResponse{
		Score: abilDto.ToProto().Score,
	}, nil
}

func (h *Handler) SetAbilityBonus(ctx context.Context, req *character_stats.SetAbilityRequest) (*character_stats.SetAbilityResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("amount = %d must be greater than or equal 0", req.Amount))
	}

	abilityType, err := types.ParseAbilityType(req.Ability)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect name for the ability")
	}

	h.logger.Info("got request", zap.Any("request", req))

	abilDto, err := h.ctrl.SetAbilityBonus(ctx, req.Id, abilityType, int(req.Amount))
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to set ability bonus", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_stats.SetAbilityResponse{
		Score: abilDto.ToProto().Score,
	}, nil
}

func (h *Handler) SetSkillBonus(ctx context.Context, req *character_stats.SetSkillBonusRequest) (*character_stats.SetSkillBonusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("amount = %d must be greater than or equal 0", req.Amount))
	}

	skillType, err := types.ParseSkillType(req.Skill)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect name for the skill")
	}

	h.logger.Info("got request", zap.Any("request", req))

	skillDto, err := h.ctrl.SetSkillBonus(ctx, req.Id, skillType, int(req.Amount))
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to set skill bonus", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_stats.SetSkillBonusResponse{
		Skill: req.Skill,
		Bonus: int32(skillDto.Bonus),
	}, nil
}

func (h *Handler) GetStats(ctx context.Context, req *character_stats.CharacterID) (*character_stats.GetStatsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}

	c, err := h.ctrl.GetStats(ctx, req.Id)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to set ability bonus", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &character_stats.GetStatsResponse{
		Stats: c.ToProto().Stats,
	}, nil
}

func (h *Handler) GetModifier(ctx context.Context, req *character_stats.GetModifierRequest) (*character_stats.GetModifierResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}

	abilityType, err := types.ParseAbilityType(req.Ability)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect name for the ability")
	}

	mod, err := h.ctrl.GetModifier(ctx, req.Id, abilityType)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to set ability bonus", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &character_stats.GetModifierResponse{
		Ability:  req.Ability,
		Modifier: int32(mod),
	}, nil
}

func (h *Handler) GetAbility(ctx context.Context, req *character_stats.AbilityRequest) (*character.Ability, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}

	abilityType, err := types.ParseAbilityType(req.Ability)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Incorrect name for the ability")
	}

	abilDto, err := h.ctrl.GetAbility(ctx, req.Id, abilityType)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to get ability", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return abilDto.ToProto(), nil
}
