package postgres

import (
	"campaigntool/internal/config"
	"campaigntool/internal/domain/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	dbPool *pgxpool.Pool
}

func New(dbConfig *config.DBConfig) (*Storage, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Database)

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to parse pool config: %s", err)
	}
	config.MaxConns = dbConfig.MaxConn

	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool: %s", err)
	}

	return &Storage{
		dbPool: pool,
	}, nil
}

func (s *Storage) SaveCampaign(ctx context.Context, name, desc string, userId int64) (int64, error) {
	op := "storage.postgres.SaveCampaign"
	var campaignId int64

	query := `
		INSERT INTO campaigns (name, description, master_id, created_at)
		VALUES ($1, $2, $3, now())
		RETURNING id
	`
	err := s.dbPool.QueryRow(ctx, query, name, desc, userId).Scan(&campaignId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, models.ErrCampaignExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return campaignId, nil
}

func (s *Storage) DeleteCampaign(ctx context.Context, campaignId int64, userId int64) error {
	op := "storage.postgres.DeleteCampaign"

	query := `
		DELETE FROM campaigns
		WHERE id = $1 AND master_id = $2
	`
	commandTag, err := s.dbPool.Exec(ctx, query, campaignId, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
	}

	return nil
}

func (s *Storage) UpdateCampaign(ctx context.Context, campaign *models.CampaignInfo) error {
	op := "storage.postgres.UpdateCampaign"

	query := `
		UPDATE campaigns
		SET name = $1, description = $2
		WHERE id = $3
	`
	commandTag, err := s.dbPool.Exec(ctx, query, campaign.Name, campaign.Description, campaign.Id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
	}

	return nil
}

func (s *Storage) GetCampaign(ctx context.Context, campaignId int64) (*models.CampaignInfo, error) {
	op := "storage.postgres.GetCampaign"

	query := `
		SELECT id, name, description, master_id, created_at
		FROM campaigns
		WHERE id = $1
	`
	var campaign models.CampaignInfo
	err := s.dbPool.QueryRow(ctx, query, campaignId).Scan(&campaign.Id, &campaign.Name, &campaign.Description, &campaign.MasterId, &campaign.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return nil, fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}

		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &campaign, nil
}

func (s *Storage) SetInviteCode(ctx context.Context, campaignId int64, masterId int64, inviteCode string) error {
	op := "storage.postgres.SetInviteCode"

	query := `
		UPDATE campaigns
		SET invite_code = $1
		WHERE id = $2 AND master_id = $3
	`
	result, err := s.dbPool.Exec(ctx, query, inviteCode, campaignId, masterId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	if result.RowsAffected() == 0 {
        return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
    }

	return nil
}

func (s *Storage) CheckInviteCode(ctx context.Context, inviteCode string) (int64, error) {
	op := "storage.postgres.CheckInviteCode"

	query := `
		SELECT id FROM campaigns
		WHERE invite_code = $1
	`
	var campaignId int64
	err := s.dbPool.QueryRow(ctx, query, inviteCode).Scan(&campaignId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return 0, fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return campaignId, nil
}

func (s *Storage) AddPlayer(ctx context.Context, campaignId int64, userId int64, charId int64) error {
	op := "storage.postgres.AddPlayer"

	query := `
		INSERT INTO campaign_characters (campaign_id, player_id, char_id)
		VALUES ($1, $2, $3);
	`
	_, err := s.dbPool.Exec(ctx, query, campaignId, userId, charId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			fmt.Println(pgErr)
			return fmt.Errorf("%s: %w", op, models.ErrPlayerInCampaign)
		} else if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RemovePlayer(ctx context.Context, campaignId int64, userId int64) error {
	op := "storage.postgres.RemovePlayer"

	query := `
		DELETE FROM campaign_characters
		WHERE campaign_id = $1 AND player_id = $2
	`
	_, err := s.dbPool.Exec(ctx, query, campaignId, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) AddCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) error {
	op := "storage.postgres.AddCharacter"

	query := `
		INSERT INTO campaign_characters (campaign_id, player_id, char_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (campaign_id, player_id, char_id) DO NOTHING;
	`
	_, err := s.dbPool.Exec(ctx, query, campaignId, userId, charId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			fmt.Println(pgErr)
			return fmt.Errorf("%s: %w", op, models.ErrCharacterInCampaign)
		} else if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, models.ErrCampaignNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) RemoveCharacter(ctx context.Context, campaignId int64, userId int64, charId int64) error {
	op := "storage.postgres.RemoveCharacter"

	query := `
		DELETE FROM campaign_characters
		WHERE campaign_id = $1 AND player_id = $2 AND char_id = $3
	`
	_, err := s.dbPool.Exec(ctx, query, campaignId, userId, charId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, models.ErrCharacterNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CreatedCampaigns(ctx context.Context, userId int64) ([]*models.Campaign, error) {
	op := "storage.postgres.CreatedCampaigns"

	query := `
		SELECT
            c.id,
            c.name,
            c.description,
			COALESCE(jsonb_agg(ch.player_id) FILTER (WHERE ch.player_id IS NOT NULL), '[]') AS player_ids,
            c.created_at
        FROM
            campaigns c
        LEFT JOIN
            campaign_characters ch ON c.id = ch.campaign_id
        WHERE
            c.master_id = $1
        GROUP BY
            c.id, c.name, c.created_at
        ORDER BY
            c.created_at DESC
	`
	rows, err := s.dbPool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var campaigns []*models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		var rawPlayersId []byte

		if err := rows.Scan(&campaign.Id, &campaign.Name, &campaign.Description, &rawPlayersId, &campaign.CreatedAt); err != nil {
			fmt.Println("rows.Scan error")
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		fmt.Println("Raw players_id:", string(rawPlayersId))
		var playersId []int64
		if err := json.Unmarshal(rawPlayersId, &playersId); err != nil {
			return nil, fmt.Errorf("failed to unmarshal players_id: %w", err)
		}

		campaign.PlayersId = &playersId
		campaign.PlayersCount = int64(len(playersId))

		campaigns = append(campaigns, &campaign)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(campaigns) == 0 {
		return nil, fmt.Errorf("%s: %w", op, models.ErrNoCampaigns)
	}

	return campaigns, nil
}

func (s *Storage) CurrentCampaigns(ctx context.Context, userId int64) ([]*models.CampaignForPlayer, error) {
	op := "storage.postgres.CurrentCampaigns"

	query := `
		SELECT c.id, c.name, c.master_id
		FROM campaigns c
		JOIN campaign_characters ch ON c.id = ch.campaign_id
		WHERE ch.player_id = $1
	`
	rows, err := s.dbPool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var campaigns []*models.CampaignForPlayer
	for rows.Next() {
		var campaign models.CampaignForPlayer
		if err := rows.Scan(&campaign.Id, &campaign.Name, &campaign.MasterId); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		campaigns = append(campaigns, &campaign)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if len(campaigns) == 0 {
		return nil, fmt.Errorf("%s: %w", op, models.ErrNoCampaigns)
	}

	return campaigns, nil
}

func (s *Storage) GetCampaignPlayers(ctx context.Context, campaignId int64) ([]int64, error) {
	op := "storage.postgres.GetCampaignPlayers"

	query := `
		SELECT DISTINCT player_id
		FROM campaign_characters
		WHERE campaign_id = $1;
	`
	rows, err := s.dbPool.Query(ctx, query, campaignId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var players []int64
	for rows.Next() {
		var player int64
		if err := rows.Scan(&player); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		players = append(players, player)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return players, nil
}

func (s *Storage) GetCampaignCharacters(ctx context.Context, campaignId int64) ([]*models.CampaignCharacter, error) {
	op := "storage.postgres.GetCampaignCharachters"

	query := `
		SELECT player_id, ARRAY_AGG(char_id) AS chars
        FROM campaign_characters
        WHERE campaign_id = $1
        GROUP BY player_id
	`
	rows, err := s.dbPool.Query(ctx, query, campaignId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	charactersMap := make(map[int64][]int64)
	for rows.Next() {
		var user int64
		var chars []int64
		if err := rows.Scan(&user, &chars); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		charactersMap[user] = chars
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var characters []*models.CampaignCharacter
	for user, chars := range charactersMap {
		characters = append(characters, &models.CampaignCharacter{
			UserId: user,
			CharIds:  chars,
		})
	}

	return characters, nil
}

func (s *Storage) GetPlayerCharacters(ctx context.Context, campaignId int64, userId int64) ([]int64, error) {
	op := "storage.postgres.GetPlayerCharachters"

	query := `
		SELECT char_id
		FROM campaign_characters
		WHERE campaign_id = $1 AND player_id = $2;
	`
	rows, err := s.dbPool.Query(ctx, query, campaignId, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var characters []int64
	for rows.Next() {
		var character int64
		if err := rows.Scan(&character); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		characters = append(characters, character)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return characters, nil
}

func (s *Storage) IsMaster(ctx context.Context, campaignId int64, userId int64) (bool, error) {
	op := "storage.postgres.IsMaster"

	query := `
		SELECT EXISTS
			(SELECT 1 
			FROM campaigns 
			WHERE id = $1 AND master_id = $2)
	`
	
	var isMaster bool
	err := s.dbPool.QueryRow(ctx, query, campaignId, userId).Scan(&isMaster)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isMaster, nil
}

func (s *Storage) IsPlayer(ctx context.Context, campaignId int64, userId int64) (bool, error) {
	op := "storage.postgres.IsPlayer"

	query := `
		SELECT EXISTS
			(SELECT 1 
			FROM campaign_characters 
			WHERE campaign_id = $1 AND player_id = $2)
	`
	
	var isPlayer bool
	err := s.dbPool.QueryRow(ctx, query, campaignId, userId).Scan(&isPlayer)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isPlayer, nil
}

func (s *Storage) IsCharacterOwner(ctx context.Context, campaignId int64, userId int64, charId int64) (bool, error) {
	op := "storage.postgres.IsCharacterOwner"

	query := `
		SELECT EXISTS
			(SELECT 1 
			FROM campaign_characters 
			WHERE campaign_id = $1 AND player_id = $2 AND char_id = $3)
	`
	
	var isOwner bool
	err := s.dbPool.QueryRow(ctx, query, campaignId, userId, charId).Scan(&isOwner)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isOwner, nil
}

func (s *Storage) Close() {
	s.dbPool.Close()
}
