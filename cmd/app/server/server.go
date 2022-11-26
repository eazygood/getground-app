package server

import (
	"context"
	"net/http"

	"github.com/eazygood/getground-app/internal/config"
	"github.com/gin-gonic/gin"
	logger "github.com/sirupsen/logrus"
)

func Start(ctx context.Context, cfg config.App) {
	dependencies, err := initDependencies(&cfg)
	if err != nil {
		panic(err)
	}

	router := gin.New()
	router.GET("/health", func(request *gin.Context) { request.String(http.StatusOK, "pong") })

	initRoutes(router, dependencies)

	run(ctx, router, cfg.Server)
}

func run(ctx context.Context, router *gin.Engine, cfg config.Server) {
	logger.Info(cfg.Http.Host + ":" + cfg.Http.Port)
	srv := &http.Server{
		Addr:    cfg.Http.Host + ":" + cfg.Http.Port,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("failed to start server: %s\n", err)
		}
	}()

	<-ctx.Done()

	shutdownCtx, cancelFn := context.WithTimeout(context.Background(), cfg.Http.ShutdownTimeout)
	defer cancelFn()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Fatal("server shutdown: ", err)
	}

	logger.Info("server exiting")
}
