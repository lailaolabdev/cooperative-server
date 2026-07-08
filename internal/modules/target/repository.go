package target

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct{ collection *mongo.Collection }

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection("cooperative_targets")}
}

func (r *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "source", Value: 1}, {Key: "provinceCode", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	return err
}

func (r *Repository) List(ctx context.Context) ([]ProvinceTarget, error) {
	cursor, err := r.collection.Find(ctx, bson.M{}, options.Find().SetSort(bson.D{{Key: "targetCount", Value: -1}, {Key: "province", Value: 1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	items := make([]ProvinceTarget, 0)
	err = cursor.All(ctx, &items)
	return items, err
}

func (r *Repository) Upsert(ctx context.Context, item ProvinceTarget) error {
	_, err := r.collection.UpdateOne(
		ctx,
		bson.M{"source": item.Source, "provinceCode": item.ProvinceCode},
		bson.M{"$set": item},
		options.Update().SetUpsert(true),
	)
	return err
}
