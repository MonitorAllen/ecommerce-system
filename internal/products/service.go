package products

import (
	"context"
	db "ecommerce-system/internal/db/sqlc"
)

type Service interface {
	ListProducts(ctx context.Context) ([]db.Product, error)
	FindProductByID(ctx context.Context, id int64) (db.Product, error)
}

type ProductService struct {
	repo *db.Queries
}

func NewProductService(repo *db.Queries) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) ListProducts(ctx context.Context) ([]db.Product, error) {
	products, err := s.repo.ListProducts(ctx)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *ProductService) FindProductByID(ctx context.Context, id int64) (db.Product, error) {
	product, err := s.repo.FindProductByID(ctx, id)
	if err != nil {
		return db.Product{}, err
	}

	return product, nil
}
