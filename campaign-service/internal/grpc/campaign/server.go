package campaign

import (
	"campaigntool/internal/domain/models"
	"context"
	"errors"
	"log"
	"strconv"

	"github.com/alsadx/gm-protos/gen/go/campaignv1"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type Claims struct {
	UID   int64  `json:"uid"`
	Email string `json:"email"`

	jwt.RegisteredClaims
}

type contextKey string

const UserClaimsKey contextKey = "user_claims"

type CampaignTool interface {
	CreateCampaign(ctx context.Context, name, desc string, userId int64) (campaignId int64, err error)
	DeleteCampaign(ctx context.Context, campaignId int64, userId int64) (err error)
	UpdateCampaign(ctx context.Context, campaignId int64, name, desc string, userId int64) (err error)
	GenerateInviteCode(ctx context.Context, campaignId int64, userId int64) (inviteCode string, err error)
	JoinCampaign(ctx context.Context, inviteCode string, userId int64, charId int64) (err error)
	LeaveCampaign(ctx context.Context, campaignId int64, userId int64) (err error)
	RemovePlayer(ctx context.Context, campaignId int64, userId int64, masterId int64) (err error)
	AddCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) (err error)
	RemoveCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) (err error)
	GetCreatedCampaigns(ctx context.Context, userId int64) (campaigns []*models.Campaign, err error)
	GetCurrentCampaigns(ctx context.Context, userId int64) (campaigns []*models.CampaignForPlayer, err error)
	GetCampaignPlayers(ctx context.Context, campaignId int64) (players []int64, err error)
	GetCampaignCharacters(ctx context.Context, campaignId int64, userId int64) (characters []*models.CampaignCharacter, err error)
	GetPlayerCharacters(ctx context.Context, campaignId int64, playerId, userId int64) (characters []int64, err error)
}

type ServerAPI struct {
	campaignv1.UnimplementedCampaignToolServer
	CampaignTool CampaignTool
}

func RegisterServerAPI(gRPC *grpc.Server, campaignTool CampaignTool) {
	campaignv1.RegisterCampaignToolServer(gRPC, &ServerAPI{CampaignTool: campaignTool})
}

func (s *ServerAPI) CreateCampaign(ctx context.Context, req *campaignv1.CreateCampaignRequest) (*campaignv1.CreateCampaignResponse, error) {
	// TODO: validate
	if req.GetName() == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	campaignId, err := s.CampaignTool.CreateCampaign(ctx, req.GetName(), req.GetDescription(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignExists) {
			return nil, status.Errorf(codes.AlreadyExists, "campaign with this name already exists")
		}

		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.CreateCampaignResponse{CampaignId: campaignId}, nil
}

func (s *ServerAPI) DeleteCampaign(ctx context.Context, req *campaignv1.DeleteCampaignRequest) (*campaignv1.DeleteCampaignResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.DeleteCampaign(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &campaignv1.DeleteCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		} else if errors.Is(err, models.ErrNotMaster) {
			return &campaignv1.DeleteCampaignResponse{Success: false}, status.Errorf(codes.PermissionDenied, "user is not master")
		}

		return &campaignv1.DeleteCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.DeleteCampaignResponse{Success: true}, nil
}

func (s *ServerAPI) UpdateCampaign(ctx context.Context, req *campaignv1.UpdateCampaignRequest) (*campaignv1.UpdateCampaignResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.UpdateCampaign(ctx, req.GetCampaignId(), req.GetName(), req.GetDescription(), userId)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNotMaster):
			return &campaignv1.UpdateCampaignResponse{Success: false}, status.Errorf(codes.PermissionDenied, "user is not master")
		case errors.Is(err, models.ErrNoChanges):
			return &campaignv1.UpdateCampaignResponse{Success: false}, status.Errorf(codes.InvalidArgument, "no changes")
		case errors.Is(err, models.ErrCampaignNotFound):
			return &campaignv1.UpdateCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		default:
			return &campaignv1.UpdateCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
		}
	}

	return &campaignv1.UpdateCampaignResponse{Success: true}, nil
}

func (s *ServerAPI) GenerateInviteCode(ctx context.Context, req *campaignv1.GenerateInviteCodeRequest) (*campaignv1.GenerateInviteCodeResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	code, err := s.CampaignTool.GenerateInviteCode(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return nil, status.Errorf(codes.NotFound, "campaign not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.GenerateInviteCodeResponse{InviteCode: code}, nil
}

func (s *ServerAPI) JoinCampaign(ctx context.Context, req *campaignv1.JoinCampaignRequest) (*campaignv1.JoinCampaignResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.JoinCampaign(ctx, req.GetInviteCode(), userId, req.GetCharacterId())
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCampaignNotFound):
			return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		case errors.Is(err, models.ErrPlayerInCampaign):
			return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.AlreadyExists, "player is already in campaign")
		case errors.Is(err, models.ErrInvalidCode):
			return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.InvalidArgument, "invalid invite code")
		}

		return &campaignv1.JoinCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.JoinCampaignResponse{Success: true}, nil
}

func (s *ServerAPI) LeaveCampaign(ctx context.Context, req *campaignv1.LeaveCampaignRequest) (*campaignv1.LeaveCampaignResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.LeaveCampaign(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &campaignv1.LeaveCampaignResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		} else if errors.Is(err, models.ErrNotPlayer) {
			return &campaignv1.LeaveCampaignResponse{Success: false}, status.Errorf(codes.PermissionDenied, "user is not player")
		}

		return &campaignv1.LeaveCampaignResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.LeaveCampaignResponse{Success: true}, nil
}

func (s *ServerAPI) RemovePlayer(ctx context.Context, req *campaignv1.RemovePlayerRequest) (*campaignv1.RemovePlayerResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.RemovePlayer(ctx, req.GetCampaignId(), req.GetPlayerId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return &campaignv1.RemovePlayerResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		} else if errors.Is(err, models.ErrNotMaster) {
			return &campaignv1.RemovePlayerResponse{Success: false}, status.Errorf(codes.PermissionDenied, "user is not master")
		}
		return &campaignv1.RemovePlayerResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.RemovePlayerResponse{Success: true}, nil
}

func (s *ServerAPI) AddCharacter(ctx context.Context, req *campaignv1.AddCharacterRequest) (*campaignv1.AddCharacterResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.AddCharacter(ctx, req.GetCampaignId(), userId, req.GetCharacterId())
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCharacterInCampaign):
			return &campaignv1.AddCharacterResponse{Success: false}, status.Errorf(codes.AlreadyExists, "character already in campaign")
		case errors.Is(err, models.ErrCampaignNotFound):
			return &campaignv1.AddCharacterResponse{Success: false}, status.Errorf(codes.NotFound, "campaign not found")
		case errors.Is(err, models.ErrNotPlayer):
			return &campaignv1.AddCharacterResponse{Success: false}, status.Errorf(codes.PermissionDenied, "user is not player")
		default:
			return &campaignv1.AddCharacterResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
		}
	}

	return &campaignv1.AddCharacterResponse{Success: true}, nil
}

func (s *ServerAPI) RemoveCharacter(ctx context.Context, req *campaignv1.RemoveCharacterRequest) (*campaignv1.RemoveCharacterResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	err := s.CampaignTool.RemoveCharacter(ctx, req.GetCampaignId(), userId, req.GetCharacterId())
	if err != nil {
		switch {
		case errors.Is(err, models.ErrCharacterNotFound):
			return &campaignv1.RemoveCharacterResponse{Success: false}, status.Errorf(codes.NotFound, "character not found")
		case errors.Is(err, models.ErrNotCharacterOwner):
			return &campaignv1.RemoveCharacterResponse{Success: false}, status.Errorf(codes.PermissionDenied, "user is not character owner")
		default:
			return &campaignv1.RemoveCharacterResponse{Success: false}, status.Errorf(codes.Internal, "internal error")
		}
	}

	return &campaignv1.RemoveCharacterResponse{Success: true}, nil
}

func (s *ServerAPI) GetCreatedCampaigns(ctx context.Context, req *campaignv1.GetCreatedCampaignsRequest) (*campaignv1.GetCreatedCampaignsResponse, error) {
	// md, ok := metadata.FromIncomingContext(ctx)
	// if !ok {
	// 	return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	// }

	// log.Printf("Received metadata: %+v", md)

	// userIDStrs := md["user_id"]
	// if len(userIDStrs) == 0 {
	// 	return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	// }

	// userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	userId := req.GetUserId()
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	campaigns, err := s.CampaignTool.GetCreatedCampaigns(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrNoCampaigns) {
			return nil, status.Errorf(codes.NotFound, "campaigns not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var res []*campaignv1.Campaign
	for _, campaign := range campaigns {
		res = append(res, &campaignv1.Campaign{
			CampaignId:   campaign.Id,
			Name:         campaign.Name,
			Description:  &campaign.Description,
			PlayersCount: campaign.PlayersCount,
			PlayersId:    *campaign.PlayersId,
		})
	}

	return &campaignv1.GetCreatedCampaignsResponse{Campaigns: res}, nil
}

func (s *ServerAPI) GetCurrentCampaigns(ctx context.Context, req *campaignv1.GetCurrentCampaignsRequest) (*campaignv1.GetCurrentCampaignsResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	campaigns, err := s.CampaignTool.GetCurrentCampaigns(ctx, userId)
	if err != nil {
		if errors.Is(err, models.ErrNoCampaigns) {
			return nil, status.Errorf(codes.NotFound, "campaigns not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var res []*campaignv1.CampaignForPlayer
	for _, campaign := range campaigns {
		res = append(res, &campaignv1.CampaignForPlayer{
			CampaignId: campaign.Id,
			Name:       campaign.Name,
			MasterId:   campaign.MasterId,
		})
	}

	return &campaignv1.GetCurrentCampaignsResponse{Campaigns: res}, nil
}

func (s *ServerAPI) GetCampaignPlayers(ctx context.Context, req *campaignv1.GetCampaignPlayersRequest) (*campaignv1.GetCampaignPlayersResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	players, err := s.CampaignTool.GetCampaignPlayers(ctx, req.GetCampaignId())
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return nil, status.Errorf(codes.NotFound, "campaign not found")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.GetCampaignPlayersResponse{PlayersId: players}, nil
}

func (s *ServerAPI) GetCampaignCharacters(ctx context.Context, req *campaignv1.GetCampaignCharactersRequest) (*campaignv1.GetCampaignCharactersResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	characters, err := s.CampaignTool.GetCampaignCharacters(ctx, req.GetCampaignId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return nil, status.Errorf(codes.NotFound, "campaign not found")
		} else if errors.Is(err, models.ErrNotMaster) {
			return nil, status.Errorf(codes.PermissionDenied, "user is not master")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	var res []*campaignv1.CampaignCharacters
	for _, character := range characters {
		res = append(res, &campaignv1.CampaignCharacters{
			UserId:       character.UserId,
			CharacterIds: character.CharIds,
		})
	}

	return &campaignv1.GetCampaignCharactersResponse{Characters: res}, nil
}

func (s *ServerAPI) GetPlayerCharacters(ctx context.Context, req *campaignv1.GetPlayerCharactersRequest) (*campaignv1.GetPlayerCharactersResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	log.Printf("Received metadata: %+v", md)

	userIDStrs := md["user_id"]
	if len(userIDStrs) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "missing user_id in metadata")
	}

	userId, _ := strconv.ParseInt(userIDStrs[0], 10, 64)
	if userId == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "unauthenticated")
	}

	characters, err := s.CampaignTool.GetPlayerCharacters(ctx, req.GetCampaignId(), req.GetPlayerId(), userId)
	if err != nil {
		if errors.Is(err, models.ErrCampaignNotFound) {
			return nil, status.Errorf(codes.NotFound, "campaign not found")
		} else if errors.Is(err, models.ErrNotCharacterOwner) {
			return nil, status.Errorf(codes.PermissionDenied, "user is not character owner")
		}
		return nil, status.Errorf(codes.Internal, "internal error")
	}

	return &campaignv1.GetPlayerCharactersResponse{CharactersId: characters}, nil
}
