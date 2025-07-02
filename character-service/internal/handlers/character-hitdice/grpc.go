package handler_hitdice

import (
	"context"
	"errors"

	// character_hitdice "github.com/alsadx/GM-Tool/character-service/gen/service/character-hitdice"
	character_hitdice "github.com/alsadx/gm-protos/gen/go/characterv1/character_hitdice"
	hitdice_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-hitdice"
	"github.com/alsadx/GM-Tool/character-service/internal/repository"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/dice"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/health"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	character_hitdice.UnimplementedCharacterHitDiceServiceServer
	ctrl   *hitdice_controller.Controller
	logger *zap.Logger
}

func New(ctrl *hitdice_controller.Controller) *Handler {
	return &Handler{
		ctrl:   ctrl,
		logger: zap.L().Named("hitdice_handler"),
	}
}

func (h *Handler) GetHitDice(ctx context.Context, req *character_hitdice.CharacterID) (*character_hitdice.HitDiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	healthDto, err := h.ctrl.GetHealthDTO(ctx, req.Id)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to get hit dice")
		return nil, err
	}

	return &character_hitdice.HitDiceResponse{
		HitDice: healthDto.ToProto().HitDice,
	}, nil
}

func (h *Handler) AddHitDice(ctx context.Context, req *character_hitdice.HitDiceRequest) (*character_hitdice.HitDiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	healthDto, err := h.ctrl.AddHitDice(ctx, req.Id, dice.Dice(req.DiceType))
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to add hit dice")
		return nil, err
	}

	return &character_hitdice.HitDiceResponse{
		HitDice: healthDto.ToProto().HitDice,
	}, nil
}

func (h *Handler) RemoveHitDice(ctx context.Context, req *character_hitdice.HitDiceRequest) (*character_hitdice.HitDiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	healthDto, err := h.ctrl.RemoveHitDice(ctx, req.Id, dice.Dice(req.DiceType))
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, health.ErrWrongTypeHitDice) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		h.logger.Error("failed to remove hit dice")
		return nil, err
	}

	return &character_hitdice.HitDiceResponse{
		HitDice: healthDto.ToProto().HitDice,
	}, nil
}

func (h *Handler) RollHitDice(ctx context.Context, req *character_hitdice.HitDiceMapRequest) (*character_hitdice.RollHitDiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	rollingDice := make(map[dice.Dice]int, len(req.AmountDice))
	for diceType, amount := range req.AmountDice {
		rollingDice[dice.Dice(diceType)] = int(amount)
	}

	rollHitDiceDto, err := h.ctrl.RollHitDice(ctx, req.Id, rollingDice)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, health.ErrWrongTypeHitDice) || errors.Is(err, health.ErrNoHitDiceAvailable) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		h.logger.Error("failed to remove hit dice")
		return nil, err
	}
	return &character_hitdice.RollHitDiceResponse{
		Results: rollHitDiceDto.Result,
		Bonus:   rollHitDiceDto.Bonus,
		HitDice: rollHitDiceDto.Health.ToProto().HitDice,
	}, nil
}

func (h *Handler) ResetHitDice(ctx context.Context, req *character_hitdice.HitDiceMapRequest) (*character_hitdice.HitDiceResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	rollingDice := make(map[dice.Dice]int, len(req.AmountDice))
	for diceType, amount := range req.AmountDice {
		rollingDice[dice.Dice(diceType)] = int(amount)
	}

	healthDto, err := h.ctrl.ResetHitDice(ctx, req.Id, rollingDice)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, health.ErrWrongTypeHitDice) || errors.Is(err, health.ErrCantResetHitDice) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		h.logger.Error("failed to remove hit dice")
		return nil, err
	}

	return &character_hitdice.HitDiceResponse{
		HitDice: healthDto.ToProto().HitDice,
	}, nil
}
