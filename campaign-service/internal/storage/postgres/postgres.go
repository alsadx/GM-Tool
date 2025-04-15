package postgres

import (
	"campaigntool/internal/config"
	"campaigntool/internal/domain/models"
	"campaigntool/internal/storage"
	"context"
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

func (s *Storage) SaveCampaign(ctx context.Context, name, desc string, userId int) (int32, error) {
	op := "storage.postgres.SaveCampaign"
	var campaignId int32

	query := `
		INSERT INTO campaigns (name, description, master_id)
		VALUES ($1, $2, $3)
		RETURNING id
	`
	err := s.dbPool.QueryRow(ctx, query, name, desc, userId).Scan(&campaignId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrCampaignExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return campaignId, nil
}

func (s *Storage) DeleteCampaign(ctx context.Context, campaignId int32, userId int) error {
	op := "storage.postgres.DeleteCampaign"

	query := `
		DELETE FROM campaigns
		WHERE id = $1 AND master_id = $2
	`
	_, err := s.dbPool.Exec(ctx, query, campaignId, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, storage.ErrCampaignNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) SetInviteCode(ctx context.Context, campaignId int32, inviteCode string) error {
	op := "storage.postgres.SetInviteCode"

	query := `
		UPDATE campaigns
		SET invite_code = $1
		WHERE id = $2
	`
	_, err := s.dbPool.Exec(ctx, query, inviteCode, campaignId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, storage.ErrCampaignNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CheckInviteCode(ctx context.Context, inviteCode string) (int32, error) {
	op := "storage.postgres.CheckInviteCode"

	query := `
		SELECT id FROM campaigns
		WHERE invite_code = $1
	`
	var campaignId int32
	err := s.dbPool.QueryRow(ctx, query, inviteCode).Scan(&campaignId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrCampaignNotFound)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return campaignId, nil
}

func (s *Storage) AddPlayer(ctx context.Context, campaignId int32, userId int) error {
	op := "storage.postgres.AddPlayer"

	query := `
		INSERT INTO players (campaign_id, user_id)
		VALUES ($1, $2)
	`
	_, err := s.dbPool.Exec(ctx, query, campaignId, userId)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return fmt.Errorf("%s: %w", op, storage.ErrPlayerInCampaign)
		} else if errors.Is(err, pgx.ErrNoRows) || strings.Contains(err.Error(), "no rows in result set") {
			return fmt.Errorf("%s: %w", op, storage.ErrCampaignNotFound)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (s *Storage) CreatedCampaigns(ctx context.Context, userId int) ([]models.Campaign, error) {
	op := "storage.postgres.CreatedCampaigns"

	query := `
		SELECT
            c.id,
            c.name,
            c.description,
            COUNT(p.player_id) AS player_count,
            c.created_at
        FROM
            campaigns c
        LEFT JOIN
            players p ON c.id = p.campaign_id
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

	var campaigns []models.Campaign
	for rows.Next() {
		var campaign models.Campaign
		if err := rows.Scan(&campaign.Id, &campaign.Name, &campaign.Description, &campaign.PlayerCount, &campaign.CreatedAt); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		campaigns = append(campaigns, campaign)
	}
	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return campaigns, nil
}

func (s *Storage) CurrentCampaigns(ctx context.Context, userId int) ([]models.CampaignForPlayer, error) {
	op := "storage.postgres.CurrentCampaigns"

	query := `
		SELECT c.id, c.name
		FROM campaigns c
		JOIN players p ON c.id = p.campaign_id
		WHERE p.player_id = $1
	`
	rows, err := s.dbPool.Query(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
    defer rows.Close()

	var campaigns []models.CampaignForPlayer
	for rows.Next() {
		var campaign models.CampaignForPlayer
		if err := rows.Scan(&campaign.Id, &campaign.Name); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		campaigns = append(campaigns, campaign)
	}
	if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

	return campaigns, nil
}

func (s *Storage) Close() {
	s.dbPool.Close()
}
