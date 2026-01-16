package port

import (
	"context"
	"meli-product-api/internal/domain/model"
)

// QuestionClient simula llamada HTTP a microservicio de Questions
type QuestionClient interface {
	GetByProductID(ctx context.Context, productID string, limit int) ([]model.Question, error)
}
