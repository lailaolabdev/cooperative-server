package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct{ collection *mongo.Collection }

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection("admins")}
}
func (r *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{{Key: "username", Value: 1}}, Options: options.Index().SetUnique(true)})
	return err
}
func (r *Repository) FindByUsername(ctx context.Context, username string) (Admin, error) {
	var admin Admin
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&admin)
	return admin, err
}
func (r *Repository) Create(ctx context.Context, admin Admin) error {
	_, err := r.collection.InsertOne(ctx, admin)
	return err
}
