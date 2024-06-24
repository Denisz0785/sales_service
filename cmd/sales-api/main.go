package main

import (
	"context"
	"crypto/rsa"
	_ "expvar" // register the /debug/vars handler
	"log"
	"net/http"
	_ "net/http/pprof" //register pprof handlers
	"os"
	"os/signal"
	"sales_service/cmd/sales-api/internal/handlers"
	"sales_service/internal/platform/auth"
	"sales_service/internal/platform/database"
	"syscall"
	"time"

	"contrib.go.opencensus.io/exporter/zipkin"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-faster/errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	openzipkin "github.com/openzipkin/zipkin-go"
	zipkinHTTP "github.com/openzipkin/zipkin-go/reporter/http"
	"github.com/spf13/viper"
	"go.opencensus.io/trace"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
func run() error {
	log := log.New(os.Stdout, "Sales : ", log.LstdFlags|log.Lshortfile)
	var cfg struct {
		DB struct {
			User       string
			Password   string
			Host       string
			Name       string
			DisableTLS string
		}
		Web struct {
			Address         string
			Debug           string
			ReadTimeout     time.Duration
			WriteTimeout    time.Duration
			ShutdownTimeout time.Duration
		}
		Auth struct {
			PrivateKeyFile string
			KeyID          string
			Algorithm      string
		}
		Trace struct {
			URL         string
			Service     string
			Probability float64
		}
	}
	log.Println("started")
	defer log.Println("finished")

	if err := initConfig(); err != nil {
		return errors.Wrap(err, "error initializing configs")
	}

	if err := godotenv.Load("./cmd/sales-api/.env"); err != nil {
		return errors.Wrap(err, "error of loading env variables")
	}

	cfg.DB.Name = viper.GetString("db.name")
	cfg.DB.User = viper.GetString("db.user")
	cfg.DB.Host = viper.GetString("db.host")

	cfg.DB.DisableTLS = viper.GetString("db.disableTLS")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	cfg.Auth.Algorithm = viper.GetString("auth.algorithm")
	cfg.Auth.KeyID = viper.GetString("auth.keyID")
	cfg.Auth.PrivateKeyFile = viper.GetString("auth.privateKeyFile")

	authenticator, err := createAuth(
		cfg.Auth.PrivateKeyFile,
		cfg.Auth.KeyID,
		cfg.Auth.Algorithm,
	)

	if err != nil {
		return errors.Wrap(err, "error creating authenticator")
	}

	db, err := database.OpenDB(database.Config{
		Host:       cfg.DB.Host,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "error connect to DB")
	}
	defer db.Close()

	cfg.Web.Address = viper.GetString("web.address")
	cfg.Trace.Probability = viper.GetFloat64("trace.probability")
	cfg.Trace.URL = viper.GetString("trace.url")
	cfg.Trace.Service = viper.GetString("trace.service")

	// start tracing
	closer, err := registerTracer(
		cfg.Trace.Service,
		cfg.Web.Address,
		cfg.Trace.URL,
		cfg.Trace.Probability,
	)

	if err != nil {
		return errors.Wrap(err, "error registering tracer")
	}

	defer closer()

	cfg.Web.ReadTimeout = viper.GetDuration("web.readtimeout")
	cfg.Web.WriteTimeout = viper.GetDuration("web.writetimeout")
	cfg.Web.ShutdownTimeout = viper.GetDuration("web.shutdowntimeout")
	// cfg.Web.Address = viper.GetString("web.address")
	cfg.Web.Debug = viper.GetString("web.debug")

	// start Debug Service
	go func() {
		log.Printf("Debug service started on %s", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		if err != nil {
			log.Printf("Debug service error: %v", err)
		}
	}()
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(shutdown, log, db, authenticator),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {

		log.Printf("listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	select {
	case err := <-serverErrors:
		errors.Wrap(err, "error of run and listennig server")
	case sig := <-shutdown:
		log.Printf("starting shutdown %v", sig)

		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		err := api.Shutdown(ctx)
		if err != nil {
			log.Printf("error of shutdown %v", err)
			err = api.Close()
		}
		if err != nil {
			return errors.Wrap(err, "could not stop server gracefully")
		}
		if sig == syscall.SIGSTOP {
			return errors.New("integrity error detected, shutting down immediately")
		}
	}
	return nil
}

func initConfig() error {
	viper.AddConfigPath("internal/platform/conf")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

// createAuth creates an authenticator using the provided private key file, key ID, and algorithm.
// It returns the authenticator and any error encountered.
func createAuth(privateKeyFile, keyID, algorithm string) (*auth.Authenticator, error) {
	// Read the contents of the private key file.
	keyContents, err := os.ReadFile(privateKeyFile)
	if err != nil {
		return nil, errors.Wrap(err, "reading private key file")
	}

	// Parse the private key from the file contents.
	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyContents)
	if err != nil {
		return nil, errors.Wrap(err, "parsing private key file")
	}

	// Create a public key lookup function using the key ID and public key.
	public := auth.NewSimpleKeyLookupFunc(keyID, key.Public().(*rsa.PublicKey))

	// Create and return the authenticator using the key, key ID, algorithm, and public key lookup function.
	return auth.NewAuthenticator(key, keyID, algorithm, public)
}

// registerTracer registers a Zipkin tracer with the provided service name, HTTP address, trace URL, and probability of sampling.
// It returns a function to close the tracer and any error encountered.
func registerTracer(service, httpAddr, traceURL string, probability float64) (func() error, error) {
	// Create a new Zipkin endpoint with the provided service name and HTTP address.
	localEndPoint, err := openzipkin.NewEndpoint(service, httpAddr)
	if err != nil {
		return nil, errors.Wrap(err, "creating the local zipkinEndpoint")
	}

	// Create a new Zipkin HTTP reporter with the provided trace URL.
	reporter := zipkinHTTP.NewReporter(traceURL)

	// Register the Zipkin exporter with the provided endpoint.
	trace.RegisterExporter(zipkin.NewExporter(reporter, localEndPoint))

	// Apply a configuration with the provided probability of sampling.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(probability)})

	// Return a function to close the reporter and any error encountered.
	return reporter.Close, nil
}
