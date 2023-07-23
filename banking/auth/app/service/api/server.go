package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/SousaGLucas/apps/banking/auth/app/service/api/middlewares"
)

func NewServer(logger *zap.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(
		middlewares.Logger(logger),
		middleware.Recoverer,
	)

	return router
}
