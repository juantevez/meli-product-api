package json

import (
	"context"
	"encoding/json"
	"errors"
	"meli-product-api/internal/domain/model"
	"os"
	"sync"
	"time"
)

type SellerRepository struct {
	mu       sync.RWMutex
	sellers  map[string]model.Seller
	filePath string
}

func NewSellerRepository(filePath string) (*SellerRepository, error) {
	repo := &SellerRepository{
		filePath: filePath,
		sellers:  make(map[string]model.Seller),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *SellerRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	var sellers []model.Seller
	if err := json.Unmarshal(data, &sellers); err != nil {
		return err
	}

	for _, s := range sellers {
		r.sellers[s.ID] = s
	}

	return nil
}

func (r *SellerRepository) GetByID(ctx context.Context, sellerID string) (*model.Seller, error) {
	// Simulate network latency
	time.Sleep(15 * time.Millisecond)

	r.mu.RLock()
	defer r.mu.RUnlock()

	seller, exists := r.sellers[sellerID]
	if !exists {
		return nil, errors.New("seller not found")
	}

	return &seller, nil
}
