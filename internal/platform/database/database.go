package database

import (
	"context"
	"net/url"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Config struct {
	Host       string
	Name       string
	User       string
	Password   string
	DisableTLS string
}

// OpenDB create connection to DB
func OpenDB(cfg Config) (*sqlx.DB, error) {
	q := url.Values{}
	q.Set("sslmode", "require")
	if cfg.DisableTLS == "disable" {
		q.Set("sslmode", "disable")
	}
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	return sqlx.Open("postgres", u.String())
}

// StatusCheck checks if the database is reachable.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	// The query to send to the database.
	const q = `SELECT true`

	// Variable to store the result of the query.
	var status bool

	// Send the query to the database and store the result in the status variable.
	return db.QueryRowContext(ctx, q).Scan(&status)

}
