package handler_health

import (
	"context"
	"fmt"

	// character_health "github.com/alsadx/GM-Tool/character-service/gen/service/character-health"
	character_health "github.com/alsadx/gm-protos/gen/go/characterv1/character_health"
	health_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-health"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	character_health.UnimplementedCharacterHealthServiceServer
	ctrl   *health_controller.Controller
	logger *zap.Logger
}

func New(ctrl *health_controller.Controller) *Handler {
	return &Handler{
		ctrl:   ctrl,
		logger: zap.L().Named("health_handler"),
	}
}

func (h *Handler) SetMaxHP(ctx context.Context, req *character_health.SetHpRequest) (*character_health.HpState, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Value <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("value = %d must be greater than 0", req.Value))
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.SetMaxHP(ctx, req.Id, int(req.Value))
	if err != nil {
		h.logger.Error("failed to set max_hp", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_health.HpState{
		MaxHp:     int32(c.Health.MaxHP),
		CurrentHp: int32(c.Health.CurrentHP),
		TempHp:    int32(c.Health.TempHP),
		IsKnocked: c.IsKnocked,
	}, nil
}

func (h *Handler) SetCurrentHP(ctx context.Context, req *character_health.SetHpRequest) (*character_health.HpState, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Value < 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("value = %d must be greater than or equal 0", req.Value))
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.SetCurrentHP(ctx, req.Id, int(req.Value))
	if err != nil {
		h.logger.Error("failed to set current_hp", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_health.HpState{
		MaxHp:     int32(c.Health.MaxHP),
		CurrentHp: int32(c.Health.CurrentHP),
		TempHp:    int32(c.Health.TempHP),
		IsKnocked: c.IsKnocked,
	}, nil
}

func (h *Handler) SetTempHP(ctx context.Context, req *character_health.SetHpRequest) (*character_health.HpState, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Value < 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("value = %d must be greater than or equal 0", req.Value))
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.SetTempHP(ctx, req.Id, int(req.Value))
	if err != nil {
		h.logger.Error("failed to set temp_hp", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_health.HpState{
		MaxHp:     int32(c.Health.MaxHP),
		CurrentHp: int32(c.Health.CurrentHP),
		TempHp:    int32(c.Health.TempHP),
		IsKnocked: c.IsKnocked,
	}, nil
}

func (h *Handler) Heal(ctx context.Context, req *character_health.ModifyHpRequest) (*character_health.HpState, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Value <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("value = %d must be greater than 0", req.Value))
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.Heal(ctx, req.Id, int(req.Value))
	if err != nil {
		h.logger.Error("failed to heal", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_health.HpState{
		MaxHp:     int32(c.Health.MaxHP),
		CurrentHp: int32(c.Health.CurrentHP),
		TempHp:    int32(c.Health.TempHP),
		IsKnocked: c.IsKnocked,
	}, nil
}

func (h *Handler) TakeDamage(ctx context.Context, req *character_health.ModifyHpRequest) (*character_health.HpState, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Value <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("value = %d must be greater than 0", req.Value))
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.TakeDamage(ctx, req.Id, int(req.Value))
	if err != nil {
		h.logger.Error("failed to take damage", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_health.HpState{
		MaxHp:     int32(c.Health.MaxHP),
		CurrentHp: int32(c.Health.CurrentHP),
		TempHp:    int32(c.Health.TempHP),
		IsKnocked: c.IsKnocked,
	}, nil
}

func (h *Handler) IsKnocked(ctx context.Context, req *character_health.CharacterID) (*character_health.IsKnockedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	is_knocked, err := h.ctrl.IsKnocked(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get is_knocked", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_health.IsKnockedResponse{IsKnocked: is_knocked}, nil
}

func (h *Handler) GetHPState(ctx context.Context, req *character_health.CharacterID) (*character_health.HpState, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.GetHPState(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to heal", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_health.HpState{
		MaxHp:     int32(c.MaxHP),
		CurrentHp: int32(c.CurrentHP),
		TempHp:    int32(c.TempHP),
		IsKnocked: c.IsKnocked,
	}, nil
}
