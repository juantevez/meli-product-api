package json

import (
	"context"
	"encoding/json"
	"errors"
	"meli-product-api/internal/domain/model"
	"os"
	"strings"
	"sync"
)

type ProductRepository struct {
	mu       sync.RWMutex
	products []model.Product
	filePath string
}

func NewProductRepository(filePath string) (*ProductRepository, error) {
	repo := &ProductRepository{
		filePath: filePath,
		products: make([]model.Product, 0),
	}

	if err := repo.load(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *ProductRepository) load() error {
	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return err
	}

	return json.Unmarshal(data, &r.products)
}

func (r *ProductRepository) FindByID(ctx context.Context, id string) (*model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, p := range r.products {
		if p.ID == id {
			return &p, nil
		}
	}

	return nil, errors.New("product not found")
}

func (r *ProductRepository) Search(ctx context.Context, keyword string, limit, offset int) ([]model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keyword = strings.ToLower(keyword)
	var results []model.Product

	for _, p := range r.products {
		if r.matches(p, keyword) {
			results = append(results, p)
		}
	}

	// Pagination
	start := offset
	end := offset + limit

	if start > len(results) {
		return []model.Product{}, nil
	}
	if end > len(results) {
		end = len(results)
	}

	return results[start:end], nil
}

func (r *ProductRepository) Count(ctx context.Context, keyword string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	keyword = strings.ToLower(keyword)
	count := 0

	for _, p := range r.products {
		if r.matches(p, keyword) {
			count++
		}
	}

	return count, nil
}

func (r *ProductRepository) FindRelated(ctx context.Context, productID, category string, limit int) ([]model.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var results []model.Product

	for _, p := range r.products {
		if p.ID != productID && p.Category == category {
			results = append(results, p)
			if len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

func (r *ProductRepository) matches(p model.Product, keyword string) bool {
	title := strings.ToLower(p.Title)
	desc := strings.ToLower(p.Description)
	category := strings.ToLower(p.Category)
	brand := strings.ToLower(p.Brand)

	return strings.Contains(title, keyword) ||
		strings.Contains(desc, keyword) ||
		strings.Contains(category, keyword) ||
		strings.Contains(brand, keyword)
}
