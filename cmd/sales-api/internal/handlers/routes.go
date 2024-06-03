package handlers

import (
	"log"
	"net/http"

	"sales_service/internal/platform/web"

	"github.com/jmoiron/sqlx"
)

func API(logger *log.Logger, db *sqlx.DB) http.Handler {
	app := web.NewApp(logger)

	p := &Product{DB: db, Log: logger}

	app.Handle(http.MethodGet, "/v1/products", p.List)
	app.Handle(http.MethodGet, "/v1/products/{id}", p.Retrieve)

	return app
}
