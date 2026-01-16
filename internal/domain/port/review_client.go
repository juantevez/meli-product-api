package port

import (
	"context"
	"meli-product-api/internal/domain/model"
)

// ReviewClient simula llamada HTTP a microservicio de Reviews
type ReviewClient interface {
	GetByProductID(ctx context.Context, productID string) ([]model.Review, error)
	GetAverageRating(ctx context.Context, productID string) (float64, error)
	GetTotalCount(ctx context.Context, productID string) (int, error)
}
