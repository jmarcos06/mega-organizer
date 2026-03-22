package storage

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"mega-play/internal/domain/bet"
)

type MongoRepository struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoRepository(uri string) (*MongoRepository, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	collection := client.Database("megahub").Collection("apostas")
	return &MongoRepository{
		client:     client,
		collection: collection,
	}, nil
}

func (r *MongoRepository) Database() *mongo.Database {
	return r.collection.Database()
}

func (r *MongoRepository) Load() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return r.client.Ping(ctx, nil)
}

func (r *MongoRepository) Save(b bet.Bet) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, b)
	return err
}

func (r *MongoRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("aposta não encontrada")
	}
	return nil
}

func (r *MongoRepository) GetAll() ([]bet.Bet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bets []bet.Bet
	if err = cursor.All(ctx, &bets); err != nil {
		return nil, err
	}

	if bets == nil {
		bets = []bet.Bet{}
	}
	return bets, nil
}

func (r *MongoRepository) GetByNickname(nickname string) ([]bet.Bet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"nickname": bson.M{"$regex": fmt.Sprintf("^%s$", nickname), "$options": "i"}}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var bets []bet.Bet
	if err = cursor.All(ctx, &bets); err != nil {
		return nil, err
	}

	if bets == nil {
		bets = []bet.Bet{}
	}
	return bets, nil
}

func (r *MongoRepository) CheckCollision(numeros []int) (bool, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var existing bet.Bet
	err := r.collection.FindOne(ctx, bson.M{"numeros": numeros}).Decode(&existing)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return false, "", nil
		}
		return false, "", err
	}
	return true, existing.Nickname, nil
}
