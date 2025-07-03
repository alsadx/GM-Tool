package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/alsadx/gm-protos/gen/go/campaignv1"
	"github.com/alsadx/gm-protos/gen/go/ssov1"
	"log"

	models "github.com/alsadx/GM-Tool/api/internal/domain"
	"github.com/alsadx/GM-Tool/api/internal/middleware"
)

func GetCreatedCampaignsHandler(campaignClient campaignv1.CampaignToolClient, userClient ssov1.UserInfoClient) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем user_id из токена
		claims, ok := r.Context().Value(middleware.UserClaimsKey).(*middleware.Claims)
		log.Printf("Claims in GetCreatedCampaignsHandler: %v", claims)
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		userID := claims.UID
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Шаг 1: Получаем кампании
		campReq := &campaignv1.GetCreatedCampaignsRequest{
			UserId: userID,
		}

		campRes, err := campaignClient.GetCreatedCampaigns(r.Context(), campReq)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to get campaigns: %v", err), http.StatusInternalServerError)
			return
		}

		// Шаг 2: Для каждого игрока в кампании — получить имя через UserInfo
		var result []models.CampaignResponse

		for _, campaign := range campRes.Campaigns {
			var players []models.Player

			for _, pid := range campaign.PlayersId {
				userRes, err := userClient.GetUserById(r.Context(), &ssov1.GetUserByIdRequest{UserId: pid})
				if err == nil && userRes.User != nil {
					players = append(players, models.Player{
						ID:   userRes.User.Id,
						Name: userRes.User.Name,
					})
				}
			}

			result = append(result, models.CampaignResponse{
				ID:          campaign.CampaignId,
				Name:        campaign.Name,
				Description: *campaign.Description,
				CreatedAt:   campaign.CreatedAt.AsTime().Format(time.RFC3339),
				Players:     players,
			})
		}

		// Отправляем JSON-ответ
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
