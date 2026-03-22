package storage

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"mega-play/internal/domain/season"
)

type MongoSeasonRepository struct {
	collection *mongo.Collection
}

func NewMongoSeasonRepository(db *mongo.Database) *MongoSeasonRepository {
	return &MongoSeasonRepository{
		collection: db.Collection("seasons"),
	}
}

func (r *MongoSeasonRepository) Save(ctx context.Context, s season.Season) error {
	// Attempt to avoid duplicates
	count, err := r.collection.CountDocuments(ctx, bson.M{"name": s.Name})
	if err != nil {
		return err
	}
	if count > 0 {
		return nil // Already exists, acts idempotently
	}
	
	_, err = r.collection.InsertOne(ctx, s)
	return err
}

func (r *MongoSeasonRepository) Delete(ctx context.Context, name string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"name": name})
	return err
}

func (r *MongoSeasonRepository) GetAll(ctx context.Context) ([]season.Season, error) {
	var seasons []season.Season
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &seasons); err != nil {
		return nil, err
	}
	
	if seasons == nil {
		seasons = []season.Season{}
	}
	
	return seasons, nil
}
