package rest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	appContext "github.com/lstrgiang/cryptowatch-server/internal/app/rest/context"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest/controllers"
	"github.com/lstrgiang/cryptowatch-server/internal/app/rest/middleware"
	"github.com/lstrgiang/cryptowatch-server/internal/infra/cache"
	"github.com/lstrgiang/cryptowatch-server/internal/repositories"
)

type (
	Config struct {
		repositories.Config
		Host          string `envconfig:"SERVER_HOST"`
		Port          string `envconfig:"SERVER_PORT"`
		APIPath       string `envconfig:"SERVER_API_PATH"`
		Domain        string `envconfig:"SERVER_DOMAIN"`
		AuthSecretKey string `envconfig:"AUTH_SECRET_KEY"`
	}
	Server interface {
		Start(ctx context.Context, stop context.CancelFunc)
		Register(cache cache.Cache)
		Shutdown()
	}
	server struct {
		cfg        Config
		gin        *gin.Engine
		appContext appContext.Context
		httpServer *http.Server
	}
)

// create new server object
func NewServer(cfg Config) Server {
	return &server{
		cfg: cfg,
		gin: gin.Default(),
	}
}

// generate server address
func (cfg Config) ServerAddress() string {
	return fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
}

// Used for register controller or and given middleware
// also initialize app context that contains
// database connection and cache ref
func (s *server) Register(cache cache.Cache) {
	// initialize database connection
	db := repositories.InitDatabase(s.cfg.Config)
	// create new app context
	s.appContext = appContext.NewContext(
		db,
		cache,
		s.cfg.AuthSecretKey,
		s.cfg.Domain,
	)
	// recover middleware on panic
	s.gin.Use(middleware.ErrorHandlerMiddleware())
	// register controllers that need auth middleware
	api := s.gin.Group("/api")
	controllers.RegisterPrivateControllers(api, s.appContext, middleware.AuthMiddleware(), s.cfg.APIPath)
	// register controllers that does not need auth middleware
	controllers.RegisterPublicControllers(api, s.appContext, s.cfg.APIPath)
	s.gin.Use(static.Serve("/", static.LocalFile("./public", true)))
}

// start the server
func (s *server) Start(ctx context.Context, stop context.CancelFunc) {
	s.httpServer = &http.Server{
		Addr:    s.cfg.ServerAddress(),
		Handler: s.gin,
	}
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Listen for the interrupt signal.

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
}

func (s server) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// force shutdown
	_ = s.httpServer.Shutdown(ctx)
}
