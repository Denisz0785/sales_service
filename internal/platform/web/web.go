package web

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
	StatusCode int
	Start      time.Time
}

// Handler is a function type that handles HTTP requests.
type Handler func(http.ResponseWriter, *http.Request) error

// App represents the web application.
type App struct {
	// mux is the router for handling HTTP requests.
	mux *chi.Mux
	// Log is the logger for logging information.
	log *log.Logger
	mw  []Middleware
}

// NewApp creates a new web application.
func NewApp(logger *log.Logger, mw ...Middleware) *App {
	return &App{
		mux: chi.NewRouter(), // Initialize a new router.
		log: logger,
		mw:  mw, // Set the logger.
	}
}

// Handle registers a new route with a matcher for the HTTP method
// and the pattern.
func (a *App) Handle(method, pattern string, h Handler) {

	h = wrapMiddleware(a.mw, h)

	// fn is the actual handler function that will be registered with the router.
	// It calls the handler function h and handles any errors.
	fn := func(w http.ResponseWriter, r *http.Request) {

		v := Values{Start: time.Now()}

		ctx := r.Context()
		ctx = context.WithValue(ctx, KeyValues, &v)
		r = r.WithContext(ctx)

		// Call the handler function h with the request and response objects.
		if err := h(w, r); err != nil {
			// If there is an error, create an ErrorResponse object with the error message.
			a.log.Printf("Error: Unhandled error %v", err)

		}
	}

	// Register the route with the router.
	a.mux.MethodFunc(method, pattern, fn)
}

// ServeHTTP implements the http.Handler interface.
//
// It serves HTTP requests by calling the ServeHTTP method of the router.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Call the ServeHTTP method of the router to handle the request.
	// The router routes the request to the appropriate handler function.
	a.mux.ServeHTTP(w, r)
}
