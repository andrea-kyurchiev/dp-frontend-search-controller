package routes

import (
	"context"

	"github.com/ONSdigital/dp-api-clients-go/renderer"
	search "github.com/ONSdigital/dp-api-clients-go/site-search"
	"github.com/ONSdigital/dp-frontend-search-controller/config"
	"github.com/ONSdigital/dp-frontend-search-controller/handlers"

	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

// Setup registers routes for the service
func Setup(ctx context.Context, r *mux.Router, cfg *config.Config, hc health.HealthCheck) {
	log.Event(ctx, "adding routes")

	rendC := renderer.New(cfg.RendererURL)
	searchC := search.NewClient(cfg.SearchQueryURL)

	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
	r.StrictSlash(true).Path("/search").Methods("GET").HandlerFunc(handlers.Read(rendC, searchC))
}
