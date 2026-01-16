package handler

import (
	"encoding/json"
	"log/slog"
	"meli-product-api/internal/application/service"
	"meli-product-api/internal/infrastructure/adapter/http/dto"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

type ProductHandler struct {
	aggregatorService *service.ProductAggregatorService
	searchService     *service.ProductSearchService
	logger            *slog.Logger
}

func NewProductHandler(
	aggregatorService *service.ProductAggregatorService,
	searchService *service.ProductSearchService,
	logger *slog.Logger,
) *ProductHandler {
	return &ProductHandler{
		aggregatorService: aggregatorService,
		searchService:     searchService,
		logger:            logger,
	}
}

// GetProductDetails godoc
// @Summary Get product details
// @Description Get complete product details with seller, reviews, questions, and related products
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dto.ProductDetailsResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products/{id} [get]
func (h *ProductHandler) GetProductDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	productID := vars["id"]

	h.logger.Info("HTTP GET /products/{id}",
		"product_id", productID,
		"method", r.Method,
		"remote_addr", r.RemoteAddr,
	)

	start := time.Now()

	// Call service
	details, err := h.aggregatorService.GetProductDetails(ctx, productID)
	if err != nil {
		if err == service.ErrProductNotFound {
			h.respondError(w, http.StatusNotFound, "Product not found with ID: "+productID, r.URL.Path)
			return
		}
		h.respondError(w, http.StatusInternalServerError, "Internal server error", r.URL.Path)
		return
	}

	// Map to DTO
	response := dto.ToProductDetailsResponse(details)

	duration := time.Since(start)
	h.logger.Info("HTTP 200 OK",
		"product_id", productID,
		"duration_ms", duration.Milliseconds(),
	)

	h.respondJSON(w, http.StatusOK, response)
}

// SearchProducts godoc
// @Summary Search products
// @Description Search products by keyword with pagination
// @Tags products
// @Accept json
// @Produce json
// @Param q query string true "Search keyword"
// @Param limit query int false "Limit" default(10) minimum(1) maximum(50)
// @Param offset query int false "Offset" default(0) minimum(0)
// @Success 200 {object} dto.ProductSearchResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/v1/products/search [get]
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Parse query parameters
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	h.logger.Info("HTTP GET /products/search",
		"query", query,
		"limit", limitStr,
		"offset", offsetStr,
	)

	// Validate query parameter
	if strings.TrimSpace(query) == "" {
		h.respondError(w, http.StatusBadRequest, "Required parameter 'q' is missing", r.URL.Path)
		return
	}

	// Parse and validate limit
	limit := 10
	if limitStr != "" {
		var err error
		limit, err = strconv.Atoi(limitStr)
		if err != nil || limit < 1 || limit > 50 {
			h.logger.Warn("Invalid limit, using default", "limit", limitStr)
			limit = 10
		}
	}

	// Parse and validate offset
	offset := 0
	if offsetStr != "" {
		var err error
		offset, err = strconv.Atoi(offsetStr)
		if err != nil || offset < 0 {
			h.logger.Warn("Invalid offset, using default", "offset", offsetStr)
			offset = 0
		}
	}

	start := time.Now()

	// Call service
	products, total, err := h.searchService.Search(ctx, query, limit, offset)
	if err != nil {
		h.respondError(w, http.StatusInternalServerError, "Error searching products", r.URL.Path)
		return
	}

	// Map to DTO
	response := dto.ToProductSearchResponse(query, products, total, limit, offset)

	duration := time.Since(start)
	h.logger.Info("HTTP 200 OK",
		"query", query,
		"results", len(products),
		"total", total,
		"duration_ms", duration.Milliseconds(),
	)

	h.respondJSON(w, http.StatusOK, response)
}

// HealthCheck godoc
// @Summary Health check
// @Description Check if the API is running
// @Tags health
// @Accept json
// @Produce plain
// @Success 200 {string} string "Product API is running"
// @Router /api/v1/products/health [get]
func (h *ProductHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Product API is running"))
}

// Helper methods
func (h *ProductHandler) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON response", "error", err)
	}
}

func (h *ProductHandler) respondError(w http.ResponseWriter, status int, message string, path string) {
	errorResponse := dto.ErrorResponse{
		Timestamp: time.Now(),
		Status:    status,
		Error:     http.StatusText(status),
		Message:   message,
		Path:      path,
	}

	h.logger.Warn("HTTP Error",
		"status", status,
		"message", message,
		"path", path,
	)

	h.respondJSON(w, status, errorResponse)
}
