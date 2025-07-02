package handler_progress

import (
	"context"
	"errors"
	"fmt"

	// character_progress "github.com/alsadx/GM-Tool/character-service/gen/service/character-progress"
	character_progress "github.com/alsadx/gm-protos/gen/go/characterv1/character_progress"
	progress_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-progress"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	character_progress.UnimplementedCharacterProgressServiceServer
	ctrl   *progress_controller.Controller
	logger *zap.Logger
}

func New(ctrl *progress_controller.Controller) *Handler {
	return &Handler{
		ctrl:   ctrl,
		logger: zap.L().Named("progress_handler"),
	}
}

func (h *Handler) GainExp(ctx context.Context, req *character_progress.ModifyExpRequest) (*character_progress.ExpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("amount = %d must be greater than 0", req.Amount))
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.GainExp(ctx, req.Id, int(req.Amount))
	if err != nil {
		h.logger.Error("failed to gain exp", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.ExpResponse{
		CurrentLvl:   int32(c.Lvl.CurrentLevel),
		CurrentExp:   int32(c.Lvl.CurrentExp),
		ExpToNextLvl: int32(max(c.Lvl.NextThreshold-c.Lvl.CurrentExp, 0)),
	}, nil
}

func (h *Handler) RemoveExp(ctx context.Context, req *character_progress.ModifyExpRequest) (*character_progress.ExpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	if req.Amount <= 0 {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("amount = %d must be greater than 0", req.Amount))
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.RemoveExp(ctx, req.Id, int(req.Amount))
	if err != nil {
		h.logger.Error("failed to remove exp", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.ExpResponse{
		CurrentLvl:   int32(c.Lvl.CurrentLevel),
		CurrentExp:   int32(c.Lvl.CurrentExp),
		ExpToNextLvl: int32(max(c.Lvl.NextThreshold-c.Lvl.CurrentExp, 0)),
	}, nil
}

func (h *Handler) LvlUp(ctx context.Context, req *character_progress.CharacterID) (*character_progress.ExpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.LvlUp(ctx, req.Id)
	if errors.Is(err, progress_controller.ErrCantLvlUpDown) {
		return nil, status.Error(
			codes.InvalidArgument,
			fmt.Sprintf(
				"You cannot raise the level of this character. You need experience for the next level: %d",
				c.Lvl.NextThreshold-c.Lvl.CurrentExp,
			),
		)
	} else if err != nil {
		h.logger.Error("failed to lvl up", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.ExpResponse{
		CurrentLvl:   int32(c.Lvl.CurrentLevel),
		CurrentExp:   int32(c.Lvl.CurrentExp),
		ExpToNextLvl: int32(max(c.Lvl.NextThreshold-c.Lvl.CurrentExp, 0)),
	}, nil
}

func (h *Handler) LvlDown(ctx context.Context, req *character_progress.CharacterID) (*character_progress.ExpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.LvlDonw(ctx, req.Id)
	if errors.Is(err, progress_controller.ErrCantLvlUpDown) {
		return nil, status.Error(
			codes.InvalidArgument,
			"You cannot downgrade this character",
		)
	} else if err != nil {
		h.logger.Error("failed to lvl down", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.ExpResponse{
		CurrentLvl:   int32(c.Lvl.CurrentLevel),
		CurrentExp:   int32(c.Lvl.CurrentExp),
		ExpToNextLvl: int32(max(c.Lvl.NextThreshold-c.Lvl.CurrentExp, 0)),
	}, nil
}

func (h *Handler) CanLvlUp(ctx context.Context, req *character_progress.CharacterID) (*character_progress.CanLvlResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	can, c, err := h.ctrl.CanLvlUp(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get CanLvlUp", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.CanLvlResponse{
		Can: can,
		ExpInfo: &character_progress.ExpResponse{
			CurrentLvl:   int32(c.Lvl.CurrentLevel),
			CurrentExp:   int32(c.Lvl.CurrentExp),
			ExpToNextLvl: int32(max(c.Lvl.NextThreshold-c.Lvl.CurrentExp, 0)),
		},
	}, nil
}

func (h *Handler) CanLvlDown(ctx context.Context, req *character_progress.CharacterID) (*character_progress.CanLvlResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	can, c, err := h.ctrl.CanLvlDown(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get CanLvlDonw", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.CanLvlResponse{
		Can: can,
		ExpInfo: &character_progress.ExpResponse{
			CurrentLvl:   int32(c.Lvl.CurrentLevel),
			CurrentExp:   int32(c.Lvl.CurrentExp),
			ExpToNextLvl: int32(max(c.Lvl.NextThreshold-c.Lvl.CurrentExp, 0)),
		},
	}, nil
}

func (h *Handler) CurrentLvl(ctx context.Context, req *character_progress.CharacterID) (*character_progress.CurrentLvlResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	lvl, err := h.ctrl.CurrentLvl(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get current lvl", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.CurrentLvlResponse{CurrentLvl: int32(lvl)}, nil
}

func (h *Handler) ExpToNextLvl(ctx context.Context, req *character_progress.CharacterID) (*character_progress.ExpToNextLvlResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	exp, err := h.ctrl.ExpToNextLvl(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get exp to next lvl", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.ExpToNextLvlResponse{ExpToNextLvl: int32(exp)}, nil
}

func (h *Handler) SetLvl(ctx context.Context, req *character_progress.SetLevelRequest) (*character_progress.ExpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.SetLvl(ctx, req.Id, int(req.Level))
	if err != nil {
		h.logger.Error("failed to set lvl", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.ExpResponse{
		CurrentLvl:   int32(c.Lvl.CurrentLevel),
		CurrentExp:   int32(c.Lvl.CurrentExp),
		ExpToNextLvl: int32(max(c.Lvl.NextThreshold-c.Lvl.CurrentExp, 0)),
	}, nil
}

func (h *Handler) GetLvlState(ctx context.Context, req *character_progress.CharacterID) (*character_progress.ExpResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	lvlState, err := h.ctrl.GetLvlState(ctx, req.Id)
	if err != nil {
		h.logger.Error("failed to get lvl state", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_progress.ExpResponse{
		CurrentLvl:   int32(lvlState.CurrentLvl),
		CurrentExp:   int32(lvlState.CurrentExp),
		ExpToNextLvl: int32(lvlState.ExpToNextLvl),
	}, nil
}
