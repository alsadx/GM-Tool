package handler_core

import (
	"context"
	"errors"

	gen "github.com/alsadx/GM-Tool/character-service/gen/character"
	character_core "github.com/alsadx/GM-Tool/character-service/gen/service/character-core"
	core_controller "github.com/alsadx/GM-Tool/character-service/internal/controller/character-core"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Handler struct {
	character_core.UnimplementedCharacterCoreServiceServer
	ctrl *core_controller.Controller
	logger *zap.Logger
}

func New(ctrl *core_controller.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
		logger: zap.L().Named("core_handler"),
	}
}

func (h *Handler) CreateCharacter(ctx context.Context, req *character_core.CreateCharRequest) (*character_core.CharacterID, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.Create(ctx, int(req.OwnerId), req.Name, req.ClassName, req.Subclass, req.Race)
	if err != nil {
		h.logger.Error("failed to create character", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &character_core.CharacterID{Id: c.ID}, nil
}

func (h *Handler) GetCharacter(ctx context.Context, req *character_core.CharacterID) (*character_core.GetCharResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.Get(ctx, req.Id)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		h.logger.Error("failed to get character", zap.Error(err))
	}
	return &character_core.GetCharResponse{Character: c.ToProto()}, nil
}

func (h *Handler) UpdateCharacter(ctx context.Context, req *character_core.UpdateCharRequest) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	charUpdate := dto.UpdateCharDTO{ID: req.Id, Name: *req.Name, Class: *req.Class, Subclass: *req.Subclass, Race: *req.Race}
	err := h.ctrl.Update(ctx, &charUpdate)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) DeleteCharacter(ctx context.Context, req *character_core.CharacterID) (*emptypb.Empty, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))
	err := h.ctrl.Delete(ctx, req.Id)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if errors.Is(err, repository.ErrEmptyID) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	} else if err != nil {
		h.logger.Error("error when processing a deletion request", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) GetCharactersByUserID(ctx context.Context, req *character_core.ListCharRequest) (*character_core.ListCharResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))
	charList, err := h.ctrl.GetByOwnerID(ctx, int(req.OwnerId))
	if err != nil {
		h.logger.Error("error when processing a request to get user characters", zap.Error(err))
	}
	protoCharList := make([]*gen.Character, len(charList))
	for i, char := range charList {
		protoCharList[i] = char.ToProto()
	}
	return &character_core.ListCharResponse{Characters: protoCharList}, nil
}

func (h *Handler) GetInfoAboutCharacter(ctx context.Context, req *character_core.CharacterID) (*character_core.GetInfoResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))
	charInfo, err := h.ctrl.GetInfoAboutCharacter(ctx, req.Id)
	if err != nil {
		h.logger.Error("error when processing a request to get character info", zap.Error(err))
	}
	return &character_core.GetInfoResponse{
		OwnerId: int64(charInfo.OwnerID),
		Name: charInfo.Name,
		Class: charInfo.Class,
		Subclass: charInfo.Subclass,
		Race: charInfo.Race,
	}, nil
}