package api

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/arizz96/event-api/api/prometheus"
	requestId "github.com/arizz96/event-api/api/request-id"
	"github.com/gin-contrib/pprof"

	"github.com/arizz96/event-api/logging"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"

	"strings"
)

// Server represents this server
type Server struct {
	handlers    []RouteHandler
	server      *http.Server
	router      *gin.Engine
	adminServer *http.Server
	adminRouter *gin.Engine
	config      *ServerConfig
	healthy     bool
}

// ServerConfig holds configurations for this server
type ServerConfig struct {
	APIPort               string
	AdminPort             string
	APIInterface          string
	AdminInterface        string
	GracefulShutdownDelay int
	AllowedOrigins        string
}

// NewServer creates a new server
func NewServer(config *ServerConfig) *Server {
	// Create new server
	s := &Server{
		healthy:     false,
		router:      gin.Default(),
		adminRouter: gin.Default(),
	}

	// Set the config
	s.config = config

	// Create the actual http server
	s.server = &http.Server{
		Addr:    config.APIInterface + ":" + config.APIPort,
		Handler: s.router,
	}

	// Create the admin http server
	s.adminServer = &http.Server{
		Addr:    s.config.AdminInterface + ":" + s.config.AdminPort,
		Handler: s.adminRouter,
	}

	// Configure CORS on router
	s.router.Use(cors.New(cors.Config{
		AllowOrigins: strings.Split(config.AllowedOrigins, ","),
		AllowMethods: []string{"PUT", "PATCH", "GET", "POST"},
		AllowHeaders: []string{"Authorization", "Content-Length", "Content-Type", "Origin"},
	}))

	return s
}

// WithRouteHandler appends a new RouteHandler
func (s *Server) WithRouteHandler(rh RouteHandler) *Server {
	s.handlers = append(s.handlers, rh)
	return s
}

// WithPProf injects a middleware handler for pprof on the admin router
func (s *Server) WithPProf() *Server {
	pprof.Register(s.adminRouter, nil)
	return s
}

// WithPrometheus injects a middleware handler that will hook into the prometheus client
func (s *Server) WithPrometheus() *Server {
	prometheus.Register(s.adminRouter)
	s.router.Use(prometheus.Middleware())
	return s
}

// WithRequestID injects a request-id header
func (s *Server) WithRequestID() *Server {
	s.router.Use(requestId.Middleware())
	return s
}

// Serve starts the http server(s) and then listens for the shutdown signal
func (s *Server) Serve(shutdownChan <-chan struct{}) {
	log := logging.GetLogger(logrus.Fields{"package": "api"})

	if os.ExpandEnv("GIN_MODE") == gin.ReleaseMode {
		gin.DisableConsoleColor()
	}

	for _, handler := range s.handlers {
		handler.Register(s.router)
	}

	// Start admin server
	go func() {
		s.adminRouter.GET("/healthz", s.HealthCheckHandler)
		if err := s.adminServer.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	// Start events server
	go func() {
		s.SetHealthy()
		if err := s.server.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	<-shutdownChan
	log.Info("Webserver recieved shutdown signal")
}

// Close cleans up and shuts down the webservers
func (s *Server) Close() {
	log := logging.GetLogger(logrus.Fields{"package": "api"})
	s.SetUnhealthy()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		log.Errorf("Event server shutdown: %s", err)
	}

	if err := s.adminServer.Shutdown(ctx); err != nil {
		log.Errorf("Admin server shutdown: %s", err)
	}

	log.Info("Webserver has been shut down")
}
