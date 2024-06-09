package mid

import (
	"expvar"
	"net/http"
	"runtime"
	"sales_service/internal/platform/web"
)

var m = struct {
	gr  *expvar.Int
	req *expvar.Int
	err *expvar.Int
}{
	gr:  expvar.NewInt("goroutines"),
	req: expvar.NewInt("requests"),
	err: expvar.NewInt("errors"),
}

// Metrics is a middleware function that collects metrics about the application.
func Metrics() web.Middleware {
	// The middleware function takes a web.Handler and returns a web.Handler.
	f := func(before web.Handler) web.Handler {

		// The new handler function that wraps the input web.Handler and
		// handles the metrics.
		h := func(w http.ResponseWriter, r *http.Request) error {

			// Call the next handler in the chain.
			err := before(w, r)

			// Increment the number of requests.
			m.req.Add(1)

			// Update the number of goroutines every 100 requests.
			if m.req.Value()%100 == 0 {
				m.gr.Set(int64(runtime.NumGoroutine()))
			}

			// Increment the number of errors if there was an error.
			if err != nil {
				m.err.Add(1)
			}

			// Return the error to the caller.
			return err
		}

		// Return the new handler function.
		return h
	}

	// Return the middleware function.
	return f
}
