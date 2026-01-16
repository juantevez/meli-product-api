package router

import (
	"log/slog"
	"meli-product-api/internal/infrastructure/adapter/http/handler"
	"meli-product-api/internal/infrastructure/adapter/http/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func NewRouter(productHandler *handler.ProductHandler, logger *slog.Logger) *mux.Router {
	r := mux.NewRouter()

	// Global middlewares
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Logger(logger))
	r.Use(middleware.CORS)

	// API routes
	api := r.PathPrefix("/api/v1").Subrouter()

	// Product routes
	api.HandleFunc("/products/{id}", productHandler.GetProductDetails).Methods(http.MethodGet)
	api.HandleFunc("/products/search", productHandler.SearchProducts).Methods(http.MethodGet)
	api.HandleFunc("/products/health", productHandler.HealthCheck).Methods(http.MethodGet)

	// Root health check
	r.HandleFunc("/health", productHandler.HealthCheck).Methods(http.MethodGet)

	return r
}
