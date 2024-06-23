package mid

import (
	"context"
	"errors"
	"net/http"
	"sales_service/internal/platform/auth"
	"sales_service/internal/platform/web"
	"strings"

	"go.opencensus.io/trace"
)

// ErrForbidden is an error that indicates that the request is forbidden.
var ErrForbidden = web.NewRequestError(errors.New("request is forbidden"), http.StatusForbidden)

// Authenticate is a middleware function that authenticates the request using a JSON Web Token (JWT)
// in the Authorization header. It parses the token and adds the claims to the request context.
// If the token is invalid or missing, it returns an error.
func Authenticate(authenticator *auth.Authenticator) web.Middleware {

	// Middleware function that wraps the provided handler and authenticates the request.
	f := func(after web.Handler) web.Handler {

		// Handler function that authenticates the request.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.Auth")
			defer span.End()

			// Extract the token from the Authorization header and parse it.
			parts := strings.Split(r.Header.Get("Authorization"), " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				// If the token is missing or has an invalid format, return an error.
				err := errors.New("expected authorization header format: bearer <token>")
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			_, span = trace.StartSpan(ctx, "internal.ParseClaims")

			claims, err := authenticator.ParseClaims(parts[1])
			if err != nil {
				// If the token is invalid, return an error.
				return web.NewRequestError(err, http.StatusUnauthorized)
			}
			span.End()
			// Add the claims to the request context.
			ctx = context.WithValue(ctx, auth.Key, claims)

			// Call the next handler in the chain.
			return after(ctx, w, r)
		}
		return h
	}
	return f

}

func HasRole(roles ...string) web.Middleware {

	f := func(after web.Handler) web.Handler {

		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.HasRole")
			defer span.End()

			claims, ok := ctx.Value(auth.Key).(auth.Claims)
			if !ok {
				return errors.New("claims missing from request context")
			}

			if !claims.HasRole(roles...) {
				return ErrForbidden
			}
			return after(ctx, w, r)
		}
		return h
	}
	return f
}
