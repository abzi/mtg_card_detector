package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/abzi/mtg_card_detector/internal/auth"
	"github.com/abzi/mtg_card_detector/internal/middleware"
)

func NewRouter(handler *Handler, authService *auth.Service) *chi.Mux {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.LoggingMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Public routes
	r.Get("/health", handler.HandleHealthCheck)
	r.Post("/api/v1/auth/anonymous", handler.HandleAnonymousAuth)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(authService))

		r.Post("/api/v1/cards/scan", handler.HandleSingleScan)
		r.Post("/api/v1/cards/scan/bulk", handler.HandleBulkScan)
		r.Get("/api/v1/inventory", handler.HandleGetInventory)
		r.Get("/api/v1/cards", handler.HandleGetCard)
	})

	return r
}
