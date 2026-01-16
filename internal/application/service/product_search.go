package service

import (
	"context"
	"log/slog"
	"meli-product-api/internal/domain/model"
	"meli-product-api/internal/domain/port"
	"strings"
)

type ProductSearchService struct {
	productRepo port.ProductRepository
	logger      *slog.Logger
}

func NewProductSearchService(
	productRepo port.ProductRepository,
	logger *slog.Logger,
) *ProductSearchService {
	return &ProductSearchService{
		productRepo: productRepo,
		logger:      logger,
	}
}

func (s *ProductSearchService) Search(ctx context.Context, query string, limit, offset int) ([]model.Product, int, error) {
	s.logger.Info("Starting product search",
		"query", query,
		"limit", limit,
		"offset", offset,
	)

	// Validate and sanitize query
	query = strings.TrimSpace(query)
	if query == "" {
		s.logger.Warn("Empty search query provided")
		return []model.Product{}, 0, nil
	}

	query = strings.ToLower(query)

	// Count total results
	total, err := s.productRepo.Count(ctx, query)
	if err != nil {
		s.logger.Error("Failed to count results", "error", err)
		return nil, 0, err
	}

	if total == 0 {
		s.logger.Info("No products found", "query", query)
		return []model.Product{}, 0, nil
	}

	// Search products
	products, err := s.productRepo.Search(ctx, query, limit, offset)
	if err != nil {
		s.logger.Error("Search failed", "error", err)
		return nil, 0, err
	}

	s.logger.Info("Search completed",
		"query", query,
		"results", len(products),
		"total", total,
	)

	return products, total, nil
}
