package web

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

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
// and the pattern. The handler function will be called when the
// pattern is matched.

func (a *App) Handle(method, pattern string, fn http.HandlerFunc) {
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
