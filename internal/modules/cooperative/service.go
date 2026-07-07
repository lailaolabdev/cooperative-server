package cooperative

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var ErrInvalidType = errors.New("type must be agriculture or scu")

type Service struct{ repository *Repository }

func NewService(r *Repository) *Service { return &Service{repository: r} }
func (s *Service) List(ctx context.Context, t, p string) ([]Cooperative, error) {
	return s.repository.List(ctx, t, p)
}
func (s *Service) Get(ctx context.Context, id string) (Cooperative, error) {
	oid, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return Cooperative{}, e
	}
	return s.repository.Get(ctx, oid)
}
func (s *Service) Create(ctx context.Context, q UpsertRequest) (Cooperative, error) {
	if !q.ValidType() {
		return Cooperative{}, ErrInvalidType
	}
	return s.repository.Create(ctx, q)
}
func (s *Service) Update(ctx context.Context, id string, q UpsertRequest) (Cooperative, error) {
	if !q.ValidType() {
		return Cooperative{}, ErrInvalidType
	}
	oid, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return Cooperative{}, e
	}
	return s.repository.Update(ctx, oid, q)
}
func (s *Service) Delete(ctx context.Context, id string) error {
	oid, e := primitive.ObjectIDFromHex(id)
	if e != nil {
		return e
	}
	return s.repository.Delete(ctx, oid)
}
