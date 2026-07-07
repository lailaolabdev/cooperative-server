package cooperative

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Repository struct{ collection *mongo.Collection }

func NewRepository(db *mongo.Database) *Repository {
	return &Repository{collection: db.Collection("cooperatives")}
}
func (r *Repository) EnsureIndexes(ctx context.Context) error {
	_, err := r.collection.Indexes().CreateMany(ctx, []mongo.IndexModel{{Keys: bson.D{{Key: "type", Value: 1}}}, {Keys: bson.D{{Key: "provinceCode", Value: 1}}}, {Keys: bson.D{{Key: "name", Value: "text"}}}})
	return err
}

func (r *Repository) UpsertSeed(ctx context.Context, item Cooperative) error {
	filter := bson.M{"source": item.Source, "sourceNo": item.SourceNo}
	set := bson.M{"name": item.Name, "type": item.Type, "provinceCode": item.ProvinceCode, "province": item.Province, "district": item.District, "village": item.Village, "chairman": item.Chairman, "phone": item.Phone, "memberCount": item.MemberCount, "description": item.Description, "status": item.Status, "source": item.Source, "sourceNo": item.SourceNo, "productionAreaHa": item.ProductionAreaHa, "establishedYear": item.EstablishedYear, "updatedAt": item.UpdatedAt}
	update := bson.M{"$set": set, "$setOnInsert": bson.M{"createdAt": item.CreatedAt}}
	_, err := r.collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
	return err
}

func (r *Repository) DeleteByTypes(ctx context.Context, types []string) (int64, error) {
	result, err := r.collection.DeleteMany(ctx, bson.M{"type": bson.M{"$in": types}})
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}
func (r *Repository) List(ctx context.Context, coopType, province string) ([]Cooperative, error) {
	filter := bson.M{}
	if coopType != "" {
		filter["type"] = coopType
	}
	if province != "" {
		filter["provinceCode"] = province
	}
	cursor, err := r.collection.Find(ctx, filter, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	items := make([]Cooperative, 0)
	err = cursor.All(ctx, &items)
	return items, err
}
func (r *Repository) Get(ctx context.Context, id primitive.ObjectID) (Cooperative, error) {
	var item Cooperative
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&item)
	return item, err
}
func (r *Repository) Create(ctx context.Context, req UpsertRequest) (Cooperative, error) {
	now := time.Now().UTC()
	item := fromRequest(req)
	item.ID = primitive.NewObjectID()
	item.CreatedAt = now
	item.UpdatedAt = now
	_, err := r.collection.InsertOne(ctx, item)
	return item, err
}
func (r *Repository) Update(ctx context.Context, id primitive.ObjectID, req UpsertRequest) (Cooperative, error) {
	item := fromRequest(req)
	update := bson.M{"$set": bson.M{"name": item.Name, "type": item.Type, "provinceCode": item.ProvinceCode, "province": item.Province, "district": item.District, "village": item.Village, "chairman": item.Chairman, "phone": item.Phone, "memberCount": item.MemberCount, "description": item.Description, "status": item.Status, "updatedAt": time.Now().UTC()}}
	result, err := r.collection.UpdateByID(ctx, id, update)
	if err != nil {
		return Cooperative{}, err
	}
	if result.MatchedCount == 0 {
		return Cooperative{}, mongo.ErrNoDocuments
	}
	return r.Get(ctx, id)
}
func (r *Repository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	if err == nil && result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return err
}
func fromRequest(q UpsertRequest) Cooperative {
	status := q.Status
	if status == "" {
		status = "active"
	}
	return Cooperative{Name: q.Name, Type: q.Type, ProvinceCode: q.ProvinceCode, Province: q.Province, District: q.District, Village: q.Village, Chairman: q.Chairman, Phone: q.Phone, MemberCount: q.MemberCount, Description: q.Description, Status: status}
}
