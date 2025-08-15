package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(handler *Handler) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/accounts/{accountID}/balance", handler.GetAccountBalance)
		r.Get("/accounts/{accountID}/holdings", handler.GetAccountHoldings)
		r.Post("/orders", handler.CreateOrder)
		r.Delete("/orders/{orderID}", handler.CancelOrder)
	})

	return r
}