package handler

import (
	"github.com/disdreamq/fantastic-telegram/services/post/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
)

func NewRouter(
	rdb *redis.Client,
	postCtrl *PostController,
	PublicRPM int,
	ProtectedPRM int,
	logger zerolog.Logger,
) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RecoveryMiddleware)
	r.Use(middleware.LoggingMiddleware(logger))

	// Public routes
	r.Group(func(r chi.Router) {
		r.Get("/posts/id/{postID}", postCtrl.GetByID)
		r.Get("/posts/title/{title}", postCtrl.GetByTitle)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		// обращение по GRPC TODO r.Use(middleware.NewAuthMiddleware(jwt.NewProvider(secret, expiry)).Authenticate)
		r.Use(middleware.NewRateLimitMiddleware(rdb, ProtectedPRM).Limit)
	})
	r.Route("/posts", func(r chi.Router) {
		r.Post("/", postCtrl.Create)
		r.Put("/{postID}", postCtrl.Update)
		r.Delete("/{postID}", postCtrl.Delete)
	})
	return r
}
