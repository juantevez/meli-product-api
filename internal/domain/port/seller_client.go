package port

import (
	"context"
	"meli-product-api/internal/domain/model"
)

// SellerClient simula llamada HTTP a microservicio de Sellers
type SellerClient interface {
	GetByID(ctx context.Context, sellerID string) (*model.Seller, error)
}
