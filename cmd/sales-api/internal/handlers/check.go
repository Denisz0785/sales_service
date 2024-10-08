package handlers

import (
	"context"
	"net/http"
	"sales_service/internal/platform/database"
	"sales_service/internal/platform/web"

	"github.com/jmoiron/sqlx"
)

type Check struct {
	DB *sqlx.DB
}

func (c *Check) Health(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var health struct {
		Status string `json:"status"`
	}

	if err := database.StatusCheck(ctx, c.DB); err != nil {
		health.Status = "database is not ready"
		return web.Respond(ctx, w, health, http.StatusServiceUnavailable)
	}

	health.Status = "database is ready"
	return web.Respond(ctx, w, health, http.StatusOK)

}
