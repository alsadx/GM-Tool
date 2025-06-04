package grpc

import (
	"context"
	"errors"

	"github.com/alsadx/GM-Tool/character-service/gen/service"
	gen "github.com/alsadx/GM-Tool/character-service/gen/character"
	"github.com/alsadx/GM-Tool/character-service/internal/controller/character"
	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/internal/repository"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	service.UnimplementedCharacterServiceServer
	ctrl *character.Controller
	logger *zap.Logger
}

func New(ctrl *character.Controller) *Handler {
	return &Handler{
		ctrl: ctrl,
		logger: zap.L().Named("handler"),
	}
}

func (h *Handler) CreateCharacter(ctx context.Context, req *service.CreateCharRequest) (*service.CreateCharResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	c, err := h.ctrl.Create(ctx, int(req.OwnerId), req.Name, req.ClassName, req.Subclass, req.Race)
	if err != nil {
		h.logger.Error("failed to create character", zap.Error(err))
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &service.CreateCharResponse{Id: c.ID}, nil
}

func (h *Handler) GetCharacter(ctx context.Context, req *service.GetCharRequest) (*service.GetCharResponse, error) {
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
	return &service.GetCharResponse{Character: c.ToProto()}, nil
}

func (h *Handler) UpdateCharacter(ctx context.Context, req *service.UpdateCharRequest) (*service.UpdateCharResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "nil req")
	}
	h.logger.Info("got request", zap.Any("request", req))

	charUpdate := dto.UpdateCharDTO{ID: req.Id, Name: req.Name, Class: req.Class, Subclass: req.Subclass, Race: req.Race}
	err := h.ctrl.Update(ctx, &charUpdate)
	if errors.Is(err, repository.ErrCharacterNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &service.UpdateCharResponse{Success: true}, nil
}

func (h *Handler) DeleteCharacter(ctx context.Context, req *service.DeleteCharRequest) (*service.DeleteCharResponse, error) {
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
	return &service.DeleteCharResponse{Success: true}, nil
}

func (h *Handler) GetCharactersByUserID(ctx context.Context, req *service.ListCharRequest) (*service.ListCharResponse, error) {
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
	return &service.ListCharResponse{Characters: protoCharList}, nil
}