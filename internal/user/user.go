package user

import (
	"context"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

func Create(ctx context.Context, db *sqlx.DB, nu NewUser, now time.Time) (*User, error) {

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.Wrap(err, "generating password hash")
	}

	u := User{
		ID:           uuid.New().String(),
		Name:         nu.Name,
		Email:        nu.Email,
		Roles:        nu.Roles,
		PasswordHash: hash,
		DateCreated:  now.UTC(),
		DateUpdated:  now.UTC(),
	}

	const q = `INSERT INTO users (user_id, name, email, roles, password_hash, date_created, date_updated)
	VALUES($1, $2, $3, $4, $5, $6, $7)`
	if _, err := db.ExecContext(ctx, q, u.ID, u.Name, u.Email, u.Roles, u.PasswordHash, u.DateCreated, u.DateUpdated); err != nil {
		return nil, errors.Wrap(err, "inserting user" )
	}
	return &u, nil
}
