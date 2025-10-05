package application

import (
	"awesomeProject/internal/router"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"net/http"
)

type App struct {
	router     *router.Router
	httpServer *http.Server
	log        *zap.Logger
}

func NewApp(router *router.Router, addr string, log *zap.Logger) *App {
	return &App{router: router,
		httpServer: &http.Server{
			Addr:    addr,
			Handler: router.GetHTTPHandler(),
		},
		log: log}
}

func (a *App) Run() error {
	a.log.Info("Starting HTTP server", zap.String("address", a.httpServer.Addr))
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("HTTP server failed", zap.Error(err))
		return fmt.Errorf("HTTP server failed: %w", err)
	}
	return nil
}
