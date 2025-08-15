package main

import (
	"context"
	"fmt"
	"net/http"

	"mini-ledger/internal/api"
	"mini-ledger/internal/config"
	"mini-ledger/internal/db"
	"mini-ledger/internal/repository"
	"mini-ledger/internal/service"

	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			config.New,
			db.New,
			repository.NewAccountRepository,
			repository.NewHoldingRepository,
			repository.NewOrderRepository,
			service.NewTradingService,
			api.NewHandler,
			api.NewRouter,
		),
		fx.Invoke(startServer),
	).Run()
}

func startServer(lc fx.Lifecycle, cfg *config.Config, router *chi.Mux) {
	server := &http.Server{
		Addr:    ":" + cfg.HTTPPort,
		Handler: router,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				fmt.Printf("Starting server on port %s\n", cfg.HTTPPort)
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					fmt.Printf("Server error: %v\n", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Println("Shutting down server...")
			return server.Shutdown(ctx)
		},
	})
}