package handlers

import (
	"log"
	"net/http"

	"sales_service/internal/platform/web"

	"github.com/jmoiron/sqlx"
)

// API creates a new web application with routes for handling Products.
func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	// Create a new web application with the logger
	app := web.NewApp(logger)

	// Create a new Product with the database connection and logger
	p := &Product{DB: db, Log: logger}

	// Register routes for retrieving all products
	app.Handle(http.MethodGet, "/v1/products", p.List)

	// Register route for retrieving a specific product
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	// Register route for creating a new product
	app.Handle(http.MethodPost, "/v1/products", p.Create)

	// Return the web application as an http.Handler
	return app
}
