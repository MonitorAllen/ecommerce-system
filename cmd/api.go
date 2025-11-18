package main

import (
	db "ecommerce-system/internal/db/sqlc"
	"ecommerce-system/internal/orders"
	"ecommerce-system/internal/products"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"time"
)

type application struct {
	config config
	db     *pgx.Conn
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})
	repo := db.New(app.db)
	productService := products.NewProductService(repo)
	productHandler := products.NewHandler(productService)

	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.ListProducts)
		r.Get("/{id}", productHandler.FindProductByID)
	})

	orderService := orders.NewOrderService(repo, app.db)
	orderHandler := orders.NewHandler(orderService)
	r.Route("/orders", func(r chi.Router) {
		r.Post("/", orderHandler.PlaceOrder)
	})

	return r
}

func (app *application) run(h http.Handler) error {
	server := &http.Server{
		Addr:         app.config.addr,
		Handler:      h,
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  time.Minute,
	}

	log.Printf("server listening on %s", app.config.addr)

	return server.ListenAndServe()
}

type dbConfig struct {
	dsn string
}

type config struct {
	addr string
	db   dbConfig
}
