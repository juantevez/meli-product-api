package json

import (
	"context"
	"encoding/json"
	"meli-product-api/internal/domain/model"
	"os"
	"sync"
	"time"
)

type ReviewRepository struct {
	mu       sync.RWMutex
	reviews  []model.Review
	filePath string
}

func NewReviewRepository(filePath string) (*ReviewRepository, error) {
	repo := &ReviewRepository{
		filePath: filePath,
		reviews:  make([]model.Review, 0),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *ReviewRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.reviews)
}

func (r *ReviewRepository) GetByProductID(ctx context.Context, productID string) ([]model.Review, error) {
	// Simulate network latency
	time.Sleep(20 * time.Millisecond)

	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []model.Review
	for _, review := range r.reviews {
		if review.ProductID == productID {
			results = append(results, review)
		}
	}

	return results, nil
}

func (r *ReviewRepository) GetAverageRating(ctx context.Context, productID string) (float64, error) {
	reviews, err := r.GetByProductID(ctx, productID)
	if err != nil {
		return 0, err
	}

	if len(reviews) == 0 {
		return 0, nil
	}

	sum := 0
	for _, review := range reviews {
		sum += review.Rating
	}

	return float64(sum) / float64(len(reviews)), nil
}

func (r *ReviewRepository) GetTotalCount(ctx context.Context, productID string) (int, error) {
	reviews, err := r.GetByProductID(ctx, productID)
	if err != nil {
		return 0, err
	}

	return len(reviews), nil
}
