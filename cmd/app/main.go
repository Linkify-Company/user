package main

import (
	"context"
	"fmt"
	auth "github.com/Linkify-Company/auth-client"
	"github.com/Linkify-Company/common_utils/errify"
	"github.com/Linkify-Company/common_utils/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"user/internal/config"
	"user/internal/handler"
	v1 "user/internal/handler/v1"
	"user/internal/repository"
	"user/internal/repository/postgres"
	"user/internal/service"
)

func main() {
	cfg := config.MustLoad()

	log := logger.GetLogger(cfg.Application.Env)

	pool, err := postgres.New(context.Background(), log, cfg.Application.Debug)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	userService := service.NewService(
		log,
		pool,
		repository.NewRepository(),
		newAuthService(cfg.Auth, log),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = userService.RebootSystem(ctx)
	if err != nil {
		log.Error(err)
		panic(err)
	}

	router := handler.Run(
		log,
		v1.NewHandler(&cfg.Handler, log, userService),
	)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port), router); err != nil {
			log.Error(errify.NewInternalServerError(err.Error(), "main/ListenAndServe"))
			panic(err)
		}
	}()

	log.Infof("server listening on port %d", cfg.Server.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Infof("application stopped")
}

func newAuthService(cfg config.AuthConfig, log logger.Logger) *auth.Service {
	authService := auth.NewClient(
		cfg.Host+":",
		cfg.Port,
		cfg.Timeout,
	)
	pong, err := authService.Ping(log)
	if err != nil {
		log.Warn(err)
	} else {
		log.Infof(pong)
	}

	return authService
}
