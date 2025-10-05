package main

import (
	"awesomeProject/internal/application"
	"awesomeProject/internal/config"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/router"
	"awesomeProject/internal/router/handlers"
	"awesomeProject/internal/service"
	"awesomeProject/pkg/logger"
	"context"
	"fmt"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	cfg := config.MustLoad("config/.env")
	fmt.Printf("%#v\n", cfg)
	log, err := logger.NewLogger(cfg.LogLevel)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize logger: %v", err))
	}
	log.Info("initialize logger success")
	defer log.Sync()
	storage, err := repository.NewStorage(ctx, cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.SSLMode, log)
	if err != nil {
		log.Fatal("failed to initialize storage", zap.Error(err))
	}
	repo := storage.NewRepository()

	calendarService := service.NewCalendarService(repo, log)
	defer calendarService.CloseRepo()
	handler := handlers.NewCalendarHandler(calendarService)

	rout := router.NewRouter(handler, cfg.LogLevel, log)
	app := application.NewApp(rout, cfg.Addr, log)
	if err := app.Run(); err != nil {
		log.Fatal("failed to run app", zap.Error(err))
	}
}
