package handlers

import (
	"log"
	"net/http"
	"os"

	mid "sales_service/internal/mid"
	"sales_service/internal/platform/auth"
	"sales_service/internal/platform/web"

	"github.com/jmoiron/sqlx"
)

// API creates a new web application with routes for handling Products.
func API(shutdown chan os.Signal, logger *log.Logger, db *sqlx.DB, authenticator *auth.Authenticator) http.Handler {
	// Create a new web application with the logger
	app := web.NewApp(shutdown, logger, mid.Logger(logger), mid.Errors(logger), mid.Metrics(), mid.Panics())

	// Create a new Product with the database connection and logger
	p := &Product{DB: db, Log: logger}
	c := &Check{DB: db}

	u := Users{DB: db, authenticator: authenticator}
	app.Handle(http.MethodGet, "/v1/users/token", u.Token)

	// Register routes for retrieving all products
	app.Handle(http.MethodGet, "/v1/products", p.List, mid.Authenticate(authenticator))

	// Register route for retrieving a specific product
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve, mid.Authenticate(authenticator))

	// Register route for creating a new product
	app.Handle(http.MethodPost, "/v1/products", p.Create, mid.Authenticate(authenticator))

	// Add a new sale to an existing product
	app.Handle(http.MethodPost, "/v1/products/{id}/sales", p.AddSale, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	// List all sales for an existing product
	app.Handle(http.MethodGet, "/v1/products/{id}/sales", p.ListSales, mid.Authenticate(authenticator))

	// Register route for updating an existing product
	app.Handle(http.MethodPut, "/v1/products/{id}", p.Update, mid.Authenticate(authenticator))

	// Register route for deleting an existing product
	app.Handle(http.MethodDelete, "/v1/products/{id}", p.Delete, mid.Authenticate(authenticator), mid.HasRole(auth.RoleAdmin))

	// Register route for checking status of database
	app.Handle(http.MethodGet, "/v1/health", c.Health)

	// Return the web application as an http.Handler
	return app

}
