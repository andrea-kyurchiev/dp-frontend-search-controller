package service

import (
	"context"
	"errors"

	"github.com/ONSdigital/dp-api-clients-go/renderer"
	search "github.com/ONSdigital/dp-api-clients-go/site-search"
	"github.com/ONSdigital/dp-frontend-search-controller/config"
	"github.com/ONSdigital/dp-frontend-search-controller/routes"
	"github.com/ONSdigital/log.go/log"
	"github.com/gorilla/mux"
)

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string
)

// Service contains the healthcheck, server and serviceList for the frontend search controller
type Service struct {
	Config      *config.Config
	HealthCheck HealthChecker
	Server      HTTPServer
	ServiceList *ExternalServiceList
}

// New creates a new service
func New() *Service {
	return &Service{}
}

// Init initialises all the service dependencies, including healthcheck with checkers, api and middleware
func (srv *Service) Init(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList) (err error) {
	log.Event(ctx, "initialising service", log.INFO)

	srv.Config = cfg
	srv.ServiceList = serviceList

	// Get health client for api router
	routerHealthClient := serviceList.GetHealthClient("api-router", cfg.APIRouterURL)

	// Initialise clients
	clients := routes.Clients{
		Renderer: renderer.New(cfg.RendererURL),
		Search:   search.NewWithHealthClient(routerHealthClient),
	}

	// Get healthcheck with checkers
	srv.HealthCheck, err = serviceList.GetHealthCheck(cfg, BuildTime, GitCommit, Version)
	if err != nil {
		log.Event(ctx, "failed to create health check", log.FATAL, log.Error(err))
		return err
	}
	if err = srv.registerCheckers(ctx, clients); err != nil {
		log.Event(ctx, "failed to register checkers", log.ERROR, log.Error(err))
		return err
	}
	clients.HealthCheckHandler = srv.HealthCheck.Handler

	// Initialise router
	r := mux.NewRouter()
	routes.Setup(ctx, r, cfg, clients)
	srv.Server = serviceList.GetHTTPServer(cfg.BindAddr, r)

	return nil
}

// Start starts an initialised service
func (srv *Service) Start(ctx context.Context, svcErrors chan error) {
	log.Event(ctx, "Starting service", log.Data{"config": srv.Config}, log.INFO)

	// Start healthcheck
	srv.HealthCheck.Start(ctx)

	// Start HTTP server
	log.Event(ctx, "Starting server", log.INFO)
	go func() {
		if err := srv.Server.ListenAndServe(); err != nil {
			log.Event(ctx, "failed to start http listen and serve", log.FATAL, log.Error(err))
			svcErrors <- err
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (srv *Service) Close(ctx context.Context) error {
	log.Event(ctx, "commencing graceful shutdown", log.INFO)
	ctx, cancel := context.WithTimeout(ctx, srv.Config.GracefulShutdownTimeout)
	hasShutdownError := false

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		log.Event(ctx, "stop health checkers", log.INFO)
		srv.HealthCheck.Stop()

		// stop any incoming requests
		if err := srv.Server.Shutdown(ctx); err != nil {
			log.Event(ctx, "failed to shutdown http server", log.Error(err), log.ERROR)
			hasShutdownError = true
		}
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Event(ctx, "shutdown timed out", log.ERROR, log.Error(ctx.Err()))
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		err := errors.New("failed to shutdown gracefully")
		log.Event(ctx, "failed to shutdown gracefully ", log.ERROR, log.Error(err))
		return err
	}

	log.Event(ctx, "graceful shutdown was successful", log.INFO)
	return nil
}

func (srv *Service) registerCheckers(ctx context.Context, c routes.Clients) (err error) {
	hasErrors := false

	if err = srv.HealthCheck.AddCheck("frontend renderer", c.Renderer.Checker); err != nil {
		hasErrors = true
		log.Event(ctx, "failed to add frontend renderer checker", log.ERROR, log.Error(err))
	}

	if err = srv.HealthCheck.AddCheck("Search API", c.Search.Checker); err != nil {
		hasErrors = true
		log.Event(ctx, "failed to add search API checker", log.ERROR, log.Error(err))
	}

	if hasErrors {
		return errors.New("Error(s) registering checkers for healthcheck")
	}

	return nil
}
