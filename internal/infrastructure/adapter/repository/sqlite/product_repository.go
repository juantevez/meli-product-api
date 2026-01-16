package port

import (
	"context"
	"meli-product-api/internal/domain/model"
)

type ProductRepository interface {
	FindByID(ctx context.Context, id string) (*model.Product, error)
	Search(ctx context.Context, keyword string, limit, offset int) ([]model.Product, error)
	Count(ctx context.Context, keyword string) (int, error)
	FindRelated(ctx context.Context, productID, category string, limit int) ([]model.Product, error)
}
