package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler is a function type that handles HTTP requests.

type Handler func(http.ResponseWriter, *http.Request) error

// App represents the web application.
type App struct {
	// mux is the router for handling HTTP requests.
	mux *chi.Mux
	// Log is the logger for logging information.
	log *log.Logger
}

// NewApp creates a new web application.
func NewApp(log *log.Logger) *App {
	return &App{
		mux: chi.NewRouter(), // Initialize a new router.
		log: log,             // Set the logger.
	}
}

// Handle registers a new route with a matcher for the HTTP method
// and the pattern.
func (a *App) Handle(method, pattern string, h Handler) {

	// fn is the actual handler function that will be registered with the router.
	// It calls the handler function h and handles any errors.
	fn := func(w http.ResponseWriter, r *http.Request) {

		// Call the handler function h with the request and response objects.
		if err := h(w, r); err != nil {
			// If there is an error, create an ErrorResponse object with the error message.
			a.log.Printf("Error handling request: %s", err)

			if err := RespondError(w, err); err != nil {
				a.log.Printf("Error writing response: %s", err)
			}
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
