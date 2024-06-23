package user

import (
	"context"
	"database/sql"
	"sales_service/internal/platform/auth"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrAuthenticationFailure = errors.New("authentication failed")
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
		return nil, errors.Wrap(err, "inserting user")
	}
	return &u, nil
}

// Authenticate authenticates a user by their email and password.

func Authenticate(
	ctx context.Context,
	db *sqlx.DB,
	now time.Time,
	email,
	password string,
) (auth.Claims, error) {

	// Query the database for the user with the given email.
	const q = `SELECT * FROM users WHERE email = $1`
	var u User
	if err := db.GetContext(ctx, &u, q, email); err != nil {

		// If no rows were found, return an authentication failure error.
		if err == sql.ErrNoRows {
			return auth.Claims{}, ErrAuthenticationFailure
		}

		// If any other error occurred, wrap it with a descriptive message.
		return auth.Claims{}, errors.Wrap(err, "selecting single user")
	}

	// Compare the provided password with the hashed password in the database.
	err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(password))

	// If the passwords do not match, return an authentication failure error.
	if err != nil {
		return auth.Claims{}, ErrAuthenticationFailure
	}

	// Generate the user's claims with the user's ID, roles, and an expiration time.
	claims := auth.NewClaims(u.ID, u.Roles, now, time.Hour)

	// Return the user's claims.
	return claims, nil
}
