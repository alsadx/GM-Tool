package campaign

import (
	"campaigntool/internal/domain/models"
	"campaigntool/internal/lib/jwt"
	"campaigntool/internal/services/campaign"
	"context"
	"errors"

	campaignv1 "github.com/alsadx/protos/gen/go/campaign"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type CampaignTool interface {
	CreateCampaign(ctx context.Context, name, desc string, userId int) (campaignId int32, err error)
	DeleteCampaign(ctx context.Context, campaignId int32, userId int) (err error)
	GenerateInviteCode(ctx context.Context, campaignId int32, userId int) (inviteCode string, err error)
	JoinCampaign(ctx context.Context, inviteCode string, userId int) (err error)
	LeaveCampaign(ctx context.Context, campaignId int32, userId int) (err error)
	GetCreatedCampaigns(ctx context.Context, userId int) (campaigns []models.Campaign, err error)
	GetCurrentCampaigns(ctx context.Context, userId int) (campaigns []models.CampaignForPlayer, err error)
}

type serverAPI struct {
	campaignv1.UnimplementedCampaignToolServer
	campaignTool CampaignTool
}

func RegisterServerAPI(gRPC *grpc.Server, campaignTool CampaignTool) {
	campaignv1.RegisterCampaignToolServer(gRPC, &serverAPI{campaignTool: campaignTool})
}

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("missing metadata")
	}

	authHeader, ok := md["authorization"]
	if !ok || len(authHeader) == 0 {
		return nil, errors.New("missing authorization header")
	}

	token := authHeader[0]
	if token == "" {
		return nil, errors.New("empty token")
	}

	userId, err := jwt.ValidateToken(token, "secret")
	if err != nil {
		return nil, errors.New("invalid token")
	}

	ctx = context.WithValue(ctx, "user_id", userId)

	return handler(ctx, req)
}

func (s *serverAPI) CreateCampaign(ctx context.Context, req *campaignv1.CreateCampaignRequest) (*campaignv1.CreateCampaignResponse, error) {
	// TODO: validate
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, errors.New("missing user_id in context")
	}

	campaignId, err := s.campaignTool.CreateCampaign(ctx, req.GetName(), req.GetDescription(), userId)
	if err != nil {
		if errors.Is(err, campaign.ErrCampaignExists) {
			return nil, status.Errorf(codes.AlreadyExists, "campaign with this name already exists")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.CreateCampaignResponse{CampaignId: campaignId}, nil
}

func (s *serverAPI) DeleteCampaign(ctx context.Context, req *campaignv1.DeleteCampaignRequest) (*campaignv1.DeleteCampaignResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, errors.New("missing user_id in context")
	}

	err := s.campaignTool.DeleteCampaign(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, campaign.ErrCampaignNotFound) {
			return &campaignv1.DeleteCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		}

		return &campaignv1.DeleteCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.DeleteCampaignResponse{Success: true}, nil
}

func (s *serverAPI) GenerateInviteCode(ctx context.Context, req *campaignv1.GenerateInviteCodeRequest) (*campaignv1.GenerateInviteCodeResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, errors.New("missing user_id in context")
	}

	code, err := s.campaignTool.GenerateInviteCode(ctx, req.GetCampaignId(), userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.GenerateInviteCodeResponse{InviteCode: code}, nil
}

func (s *serverAPI) JoinGame(ctx context.Context, req *campaignv1.JoinCampaignRequest) (*campaignv1.JoinCampaignResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, errors.New("missing user_id in context")
	}

	err := s.campaignTool.JoinCampaign(ctx, req.GetInviteCode(), userId)
	if err != nil {
		if errors.Is(err, campaign.ErrInvalidCode) {
			return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.InvalidArgument, "invalid invite code")
		}
		return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.JoinCampaignResponse{Success: true}, nil
}

func (s *serverAPI) LeaveGame(ctx context.Context, req *campaignv1.LeaveCampaignRequest) (*campaignv1.LeaveCampaignResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, errors.New("missing user_id in context")
	}

	err := s.campaignTool.LeaveCampaign(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, campaign.ErrCampaignNotFound) {
			return &campaignv1.LeaveCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		}
		return &campaignv1.LeaveCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.LeaveCampaignResponse{Success: true}, nil
}

func (s *serverAPI) GetCreatedCampaigns(ctx context.Context, req *campaignv1.GetCreatedCampaignsRequest) (*campaignv1.GetCreatedCampaignsResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, errors.New("missing user_id in context")
	}

	campaigns, err := s.campaignTool.GetCreatedCampaigns(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var res []*campaignv1.Campaign
    for _, campaign := range campaigns {
        res = append(res, &campaignv1.Campaign{
            CampaignId:       campaign.Id,
            Name:         campaign.Name,
            Description:  &campaign.Description,
            PlayerCount:  campaign.PlayerCount,
        })
    }

	return &campaignv1.GetCreatedCampaignsResponse{Campaigns: res}, nil
}

func (s *serverAPI) GetCurrentCampaigns(ctx context.Context, req *campaignv1.GetCurrentCampaignsRequest) (*campaignv1.GetCurrentCampaignsResponse, error) {
	userId, ok := ctx.Value("user_id").(int)
	if !ok {
		return nil, errors.New("missing user_id in context")
	}

	campaigns, err := s.campaignTool.GetCurrentCampaigns(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var res []*campaignv1.CampaignForPlayer
    for _, campaign := range campaigns {
        res = append(res, &campaignv1.CampaignForPlayer{
            CampaignId:       campaign.Id,
            Name:         campaign.Name,
        })
    }

	return &campaignv1.GetCurrentCampaignsResponse{Campaigns: res}, nil
}
