package controller

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/msjai/latency-ping/internal/config"
	"github.com/msjai/latency-ping/internal/entity"
)

type LatencyUseCaseProviderI interface {
	RefreshLatencyInfo()
	GetMinLatencyLogic() (*entity.WebSite, error)
	GetMaxLatencyLogic() (*entity.WebSite, error)
	GetLatencyLogic(siteName string) (*entity.WebSite, error)
}

// NewRouter -.
func NewRouter(router *chi.Mux, latencyUseCase LatencyUseCaseProviderI, cfg *config.Config) *chi.Mux {
	router.Use(middleware.Logger)

	// Routers
	router = newLatencyRoutes(router, latencyUseCase, cfg)

	return router
}
