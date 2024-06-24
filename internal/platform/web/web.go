package web

import (
	"context"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/tracecontext"
	"go.opencensus.io/trace"
)

type ctxKey int

const KeyValues ctxKey = 1

type Values struct {
	StatusCode int
	Start      time.Time
	TraceID    string
}

// Handler is a function type that handles HTTP requests.
type Handler func(context.Context, http.ResponseWriter, *http.Request) error

// App represents the web application.
type App struct {
	// mux is the router for handling HTTP requests.
	mux *chi.Mux
	// Log is the logger for logging information.
	log      *log.Logger
	mw       []Middleware
	och      *ochttp.Handler
	shutdown chan os.Signal
}

// NewApp creates a new web application.
func NewApp(shutdown chan os.Signal, logger *log.Logger, mw ...Middleware) *App {
	app := &App{
		mux:      chi.NewRouter(), // Initialize a new router.
		log:      logger,
		mw:       mw,
		shutdown: shutdown,
	}

	app.och = &ochttp.Handler{
		Propagation: &tracecontext.HTTPFormat{},
		Handler:     app.mux,
	}

	return app
}

// Handle registers a new route with a matcher for the HTTP method
// and the pattern.
func (a *App) Handle(method, pattern string, h Handler, mw ...Middleware) {

	h = wrapMiddleware(mw, h)

	h = wrapMiddleware(a.mw, h)

	// fn is the actual handler function that will be registered with the router.
	// It calls the handler function h and handles any errors.
	fn := func(w http.ResponseWriter, r *http.Request) {

		ctx, span := trace.StartSpan(r.Context(), "internal.platform.web")
		defer span.End()

		v := Values{
			Start:   time.Now(),
			TraceID: span.SpanContext().TraceID.String(),
		}
		ctx = context.WithValue(ctx, KeyValues, &v)

		// Call the handler function h with the request and response objects.
		if err := h(ctx, w, r); err != nil {
			// If there is an error, create an ErrorResponse object with the error message.
			a.log.Printf("%s: Unhandled error %+v", v.TraceID, err)
			if IsShutdown(err) {
				a.SignalShutdown()
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
	a.och.ServeHTTP(w, r)
}

// SignalShutdown is used to gracefully shut down the server.
func (a *App) SignalShutdown() {
	a.log.Println("initiating shutdown")
	a.shutdown <- syscall.SIGSTOP
}
