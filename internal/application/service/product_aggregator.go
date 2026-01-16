package service

import (
	"context"
	"errors"
	"log/slog"
	"meli-product-api/internal/domain/model"
	"meli-product-api/internal/domain/port"
	"sync"
	"time"
)

var ErrProductNotFound = errors.New("product not found")

type ProductAggregatorService struct {
	productRepo    port.ProductRepository
	sellerClient   port.SellerClient
	reviewClient   port.ReviewClient
	questionClient port.QuestionClient
	logger         *slog.Logger
}

func NewProductAggregatorService(
	productRepo port.ProductRepository,
	sellerClient port.SellerClient,
	reviewClient port.ReviewClient,
	questionClient port.QuestionClient,
	logger *slog.Logger,
) *ProductAggregatorService {
	return &ProductAggregatorService{
		productRepo:    productRepo,
		sellerClient:   sellerClient,
		reviewClient:   reviewClient,
		questionClient: questionClient,
		logger:         logger,
	}
}

func (s *ProductAggregatorService) GetProductDetails(ctx context.Context, productID string) (*model.ProductDetails, error) {
	s.logger.Info("Starting product aggregation", "product_id", productID)
	start := time.Now()

	// PASO 1: Obtener producto principal
	product, err := s.productRepo.FindByID(ctx, productID)
	if err != nil {
		s.logger.Error("Product not found", "product_id", productID, "error", err)
		return nil, ErrProductNotFound
	}

	// PASO 2: Orquestar llamadas asíncronas a "microservicios"
	s.logger.Info("Orchestrating 4 parallel service calls")

	type result struct {
		seller    *model.Seller
		reviews   []model.Review
		avgRating float64
		totalRevs int
		questions []model.Question
		related   []model.Product
		err       error
	}

	resultChan := make(chan result, 1)

	go func() {
		var wg sync.WaitGroup
		var mu sync.Mutex
		res := result{}

		// Seller
		wg.Add(1)
		go func() {
			defer wg.Done()
			seller, err := s.fetchSeller(ctx, product.ID)
			mu.Lock()
			res.seller = seller
			if err != nil {
				s.logger.Warn("Seller fetch failed, using default", "error", err)
			}
			mu.Unlock()
		}()

		// Reviews
		wg.Add(1)
		go func() {
			defer wg.Done()
			reviews, avg, total := s.fetchReviews(ctx, productID)
			mu.Lock()
			res.reviews = reviews
			res.avgRating = avg
			res.totalRevs = total
			mu.Unlock()
		}()

		// Questions
		wg.Add(1)
		go func() {
			defer wg.Done()
			questions := s.fetchQuestions(ctx, productID)
			mu.Lock()
			res.questions = questions
			mu.Unlock()
		}()

		// Related Products
		wg.Add(1)
		go func() {
			defer wg.Done()
			related := s.fetchRelated(ctx, productID, product.Category)
			mu.Lock()
			res.related = related
			mu.Unlock()
		}()

		wg.Wait()
		resultChan <- res
	}()

	res := <-resultChan

	// PASO 3: Calcular shipping
	shipping := s.buildShipping(product)

	// PASO 4: Construir respuesta agregada
	details := &model.ProductDetails{
		Product:         *product,
		Seller:          *res.seller,
		Shipping:        shipping,
		Reviews:         res.reviews,
		AverageRating:   res.avgRating,
		TotalReviews:    res.totalRevs,
		Questions:       res.questions,
		RelatedProducts: res.related,
	}

	duration := time.Since(start)
	s.logger.Info("Aggregation completed",
		"product_id", productID,
		"duration_ms", duration.Milliseconds(),
		"reviews", len(res.reviews),
		"questions", len(res.questions),
		"related", len(res.related),
	)

	return details, nil
}

func (s *ProductAggregatorService) fetchSeller(ctx context.Context, sellerID string) (*model.Seller, error) {
	s.logger.Debug("Calling SellerService", "seller_id", sellerID)
	start := time.Now()

	seller, err := s.sellerClient.GetByID(ctx, sellerID)
	if err != nil {
		// Fallback: Default seller
		return &model.Seller{
			ID:              "default",
			Nickname:        "Vendedor",
			ReputationLevel: "green",
			TotalSales:      0,
			ReputationScore: 0.0,
			YearsActive:     0,
			IsOfficialStore: false,
		}, err
	}

	s.logger.Debug("SellerService responded",
		"duration_ms", time.Since(start).Milliseconds(),
		"nickname", seller.Nickname,
	)

	return seller, nil
}

func (s *ProductAggregatorService) fetchReviews(ctx context.Context, productID string) ([]model.Review, float64, int) {
	s.logger.Debug("Calling ReviewService", "product_id", productID)
	start := time.Now()

	reviews, _ := s.reviewClient.GetByProductID(ctx, productID)
	avgRating, _ := s.reviewClient.GetAverageRating(ctx, productID)
	total, _ := s.reviewClient.GetTotalCount(ctx, productID)

	s.logger.Debug("ReviewService responded",
		"duration_ms", time.Since(start).Milliseconds(),
		"count", len(reviews),
		"avg_rating", avgRating,
	)

	return reviews, avgRating, total
}

func (s *ProductAggregatorService) fetchQuestions(ctx context.Context, productID string) []model.Question {
	s.logger.Debug("Calling QuestionService", "product_id", productID)
	start := time.Now()

	questions, _ := s.questionClient.GetByProductID(ctx, productID, 10)

	s.logger.Debug("QuestionService responded",
		"duration_ms", time.Since(start).Milliseconds(),
		"count", len(questions),
	)

	return questions
}

func (s *ProductAggregatorService) fetchRelated(ctx context.Context, productID, category string) []model.Product {
	s.logger.Debug("Calling ProductService for related", "category", category)
	start := time.Now()

	related, _ := s.productRepo.FindRelated(ctx, productID, category, 4)

	s.logger.Debug("ProductService responded",
		"duration_ms", time.Since(start).Milliseconds(),
		"count", len(related),
	)

	return related
}

func (s *ProductAggregatorService) buildShipping(product *model.Product) model.Shipping {
	freeShipping := product.Price > 50000
	cost := 0.0
	if !freeShipping {
		cost = product.Price * 0.05
	}

	return model.Shipping{
		FreeShipping:      freeShipping,
		ShippingMode:      "standard",
		Cost:              cost,
		EstimatedDelivery: "Llega en 3-5 días",
		FullFulfillment:   false,
		PickupAvailable:   "Si",
	}
}
