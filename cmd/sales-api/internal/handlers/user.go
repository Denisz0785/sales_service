package handlers

import (
	"context"

	"net/http"
	"sales_service/internal/platform/web"
	"sales_service/internal/user"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.opencensus.io/trace"

	"sales_service/internal/platform/auth"
)

// TODO: implement user handler
type Users struct {
	DB            *sqlx.DB
	authenticator *auth.Authenticator
}

// Token handles the authentication of a user and generates a JWT token.
func (u *Users) Token(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.user.Token")
	defer span.End()

	// Get the web values from the context.
	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return errors.New("web value missing from context")
	}

	// Get the user's email and password from the request's Basic Auth headers.
	email, pass, ok := r.BasicAuth()
	if !ok {
		err := errors.New("must provide email and password in basic auth")
		return web.NewRequestError(err, http.StatusUnauthorized)
	}

	// Authenticate the user with the provided email and password.
	claims, err := user.Authenticate(ctx, u.DB, v.Start, email, pass)
	if err != nil {
		switch err {
		case user.ErrAuthenticationFailure:
			return web.NewRequestError(err, http.StatusUnauthorized)
		default:
			return errors.Wrap(err, "authenticating user")
		}
	}

	// Generate a JWT token using the authenticator and the user's claims.
	var tkn struct {
		Token string `json:"token"`
	}

	tkn.Token, err = u.authenticator.GenerateToken(claims)
	if err != nil {
		return errors.Wrap(err, "generating token")
	}

	// Respond with the generated token.
	return web.Respond(ctx, w, tkn, http.StatusOK)
}
