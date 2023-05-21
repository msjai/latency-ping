package controller

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/msjai/latency-ping/internal/config"
	"github.com/msjai/latency-ping/internal/controller/middleware"
)

const (
	TextPlain       = "text/plain"
	ApplicationJSON = "application/json"
)

// latencyRoutes -
type latencyRoutes struct {
	latencyUseCase LatencyUseCaseProviderI
	cfg            *config.Config
}

// newLoyaltyRoutes -.
func newLatencyRoutes(router *chi.Mux, latencyUseCase LatencyUseCaseProviderI, cfg *config.Config) *chi.Mux {
	routes := &latencyRoutes{
		latencyUseCase: latencyUseCase,
		cfg:            cfg,
	}

	// Public Routes
	// // Only text/plain request type accepted
	router.Group(func(router chi.Router) {
		router.Use(middleware.AllowContentType(TextPlain))
		router.Get("/api/minlatency", routes.GetMinLatency)
		router.Get("/api/maxlatency", routes.GetMaxLatency)
		router.Get("/api/latency/{name}", routes.GetLatency)
	})

	go routes.latencyUseCase.RefreshLatencyInfo()

	return router
}

func (routes *latencyRoutes) GetMinLatency(w http.ResponseWriter, r *http.Request) {
	website, err := routes.latencyUseCase.GetMinLatencyLogic()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(website)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response) //nolint:errcheck
}

func (routes *latencyRoutes) GetMaxLatency(w http.ResponseWriter, r *http.Request) {
	website, err := routes.latencyUseCase.GetMaxLatencyLogic()

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(website)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response) //nolint:errcheck
}

func (routes *latencyRoutes) GetLatency(w http.ResponseWriter, r *http.Request) {
	siteName := chi.URLParam(r, "name")

	if siteName == "/favicon.ico" {
		return
	}

	if siteName == "" {
		http.Error(w, "The query parameter is missing", http.StatusBadRequest)
		return
	}

	website, err := routes.latencyUseCase.GetLatencyLogic(siteName)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(website)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", ApplicationJSON)
	w.WriteHeader(http.StatusOK)
	w.Write(response) //nolint:errcheck
}
