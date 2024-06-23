package mid

import (
	"context"
	"errors"
	"log"
	"net/http"
	"sales_service/internal/platform/web"
	"time"

	"go.opencensus.io/trace"
)

// Logger is a middleware that logs the details of each HTTP request.
// It logs the status code, the time it took to process the request, the HTTP method, and the URL.
// It is typically used for debugging and logging purposes.
func Logger(log *log.Logger) web.Middleware {

	// Middleware function that logs the details of each HTTP request
	f := func(before web.Handler) web.Handler {

		// Handler function that logs the details of each HTTP request
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			ctx, span := trace.StartSpan(ctx, "internal.mid.Log")
			defer span.End()

			// Get the web values from the request context
			v, ok := ctx.Value(web.KeyValues).(*web.Values)

			// If the web values are missing from the context, return an error
			if !ok {
				return errors.New("web value missing from context")
			}

			// Execute the next handler in the chain
			err := before(ctx, w, r)

			// Log the details of the HTTP request
			log.Printf("%s : (%d) : %s %s -> %s (%s)",
				v.TraceID, v.StatusCode, r.Method, r.URL.Path,
				r.RemoteAddr, time.Since(v.Start),
			)

			// Return the error, if any, from the next handler in the chain
			return err
		}

		// Return the handler function
		return h
	}

	// Return the middleware function
	return f
}
