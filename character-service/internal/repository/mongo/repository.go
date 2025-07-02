package mongo

import (
	"context"
	"errors"

	"github.com/alsadx/GM-Tool/character-service/internal/dto"
	"github.com/alsadx/GM-Tool/character-service/internal/repository"
	"github.com/alsadx/GM-Tool/character-service/pkg/domain/character"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.uber.org/zap"
)

type CharacterRepository struct {
	collection *mongo.Collection
	logger     *zap.Logger
}

func NewRepository(db *mongo.Database) *CharacterRepository {
	return &CharacterRepository{
		collection: db.Collection("characters"),
		logger:     zap.L().Named("character_repository"),
	}
}

func (r *CharacterRepository) Create(ctx context.Context, c *character.Character) error {
	log := r.logger.With(zap.String("character_id", c.ID), zap.String("name", c.Name))
	log.Info("creating character")

	dtoChar := dto.CharacterDTOFromDomain(c)
	res, err := r.collection.InsertOne(ctx, dtoChar)
	if err != nil {
		log.Error("failed to create character", zap.Error(err))
		return err
	}

	log.Info("character created successfully",
		zap.Any("inserted_id", res.InsertedID),
		zap.String("mongo_id", dtoChar.ID),
	)
	return nil
}

func (r *CharacterRepository) Get(ctx context.Context, id string) (*character.Character, error) {
	if id == "" {
		r.logger.Error("attempt to get character with empty ID")
		return nil, repository.ErrEmptyID
	}

	log := r.logger.With(zap.String("character_id", id))
	log.Info("getting character")

	dtoChar := new(dto.CharacterDTO)
	filter := bson.D{{Key: "_id", Value: id}}
	err := r.collection.FindOne(ctx, filter).Decode(dtoChar)

	switch {
	case errors.Is(err, mongo.ErrNoDocuments):
		log.Warn("character not found")
		return nil, repository.ErrCharacterNotFound
	case err != nil:
		log.Error("failed to get character", zap.Error(err))
		return nil, err
	}

	char, err := dtoChar.ToDomain()
	if err != nil {
		log.Error("failed to convert DTO to domain", zap.Error(err))
		return nil, err
	}

	log.Info("character retrieved successfully")
	return char, nil
}

func (r *CharacterRepository) GetByOwnerID(ctx context.Context, ownerID int) ([]*character.Character, error) {
	log := r.logger.With(zap.Int("owner_id", ownerID))
	log.Info("fetching characters by owner ID")

	filter := bson.D{{Key: "owner_id", Value: ownerID}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		log.Error("failed to execute find query", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			log.Warn("failed to close cursor", zap.Error(err))
		}
	}()

	var resultDto []dto.CharacterDTO
	if err := cursor.All(ctx, &resultDto); err != nil {
		log.Error("failed to decode results", zap.Error(err))
		return nil, err
	}

	log.Debug("characters found",
		zap.Int("dto_count", len(resultDto)),
		zap.Int("owner_id", ownerID),
	)

	var (
		result  = make([]*character.Character, 0, len(resultDto))
		errs    []error
		success int
	)

	for _, d := range resultDto {
		char, err := d.ToDomain()
		if err != nil {
			log.Error("failed to convert DTO to domain",
				zap.String("character_id", d.ID),
				zap.Error(err),
			)
			errs = append(errs, err)
			continue
		}
		result = append(result, char)
		success++
	}

	log.Info("characters conversion completed",
		zap.Int("total_documents", len(resultDto)),
		zap.Int("successful_conversions", success),
		zap.Int("failed_conversions", len(errs)),
	)

	if len(errs) > 0 {
		combinedErr := errors.Join(errs...)
		log.Warn("partial conversion errors occurred",
			zap.Int("error_count", len(errs)),
			zap.Error(combinedErr),
		)
		return result, combinedErr
	}

	return result, nil
}

func (r *CharacterRepository) Update(ctx context.Context, c *character.Character) error {
	if c.ID == "" {
		r.logger.Error("attempt to update character with empty ID")
		return repository.ErrEmptyID
	}

	log := r.logger.With(zap.String("character_id", c.ID))
	log.Info("updating character")

	dtoChar := dto.CharacterDTOFromDomain(c)
	res, err := r.collection.ReplaceOne(
		ctx,
		bson.D{{Key: "_id", Value: dtoChar.ID}},
		dtoChar,
	)
	if err != nil {
		log.Error("failed to update character", zap.Error(err))
		return err
	}
	if res.MatchedCount == 0 {
		log.Warn("character not found for update")
		return repository.ErrCharacterNotFound
	}

	log.Info("character updated successfully",
		zap.Int64("matched_count", res.MatchedCount),
		zap.Int64("modified_count", res.ModifiedCount),
	)
	return nil
}

func (r *CharacterRepository) Delete(ctx context.Context, id string) error {
	if id == "" {
		r.logger.Error("attempt to delete character with empty ID")
		return repository.ErrEmptyID
	}

	log := r.logger.With(zap.String("character_id", id))
	log.Info("deleting character")

	res, err := r.collection.DeleteOne(ctx, bson.D{{Key: "_id", Value: id}})
	if err != nil {
		log.Error("failed to delete character", zap.Error(err))
		return err
	}
	if res.DeletedCount == 0 {
		log.Warn("character not found for deletion")
		return repository.ErrCharacterNotFound
	}

	log.Info("character deleted successfully",
		zap.Int64("deleted_count", res.DeletedCount),
	)
	return nil
}
