package main

import (
	"context"
	"log"
	"net/http"
	_ "net/http/pprof" //register pprof handlers
	"os"
	"os/signal"
	"sales_service/cmd/sales-api/internal/handlers"
	"sales_service/internal/platform/database"
	"syscall"
	"time"

	"github.com/go-faster/errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
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

	cfg.Web.ReadTimeout = viper.GetDuration("web.readtimeout")
	cfg.Web.WriteTimeout = viper.GetDuration("web.writetimeout")
	cfg.Web.ShutdownTimeout = viper.GetDuration("web.shutdowntimeout")
	cfg.Web.Address = viper.GetString("web.address")
	cfg.Web.Debug = viper.GetString("web.debug")

	// start Debug Service
	go func() {
		log.Printf("Debug service started on %s", cfg.Web.Debug)
		err := http.ListenAndServe(cfg.Web.Debug, http.DefaultServeMux)
		if err != nil {
			log.Printf("Debug service error: %v", err)
		}
	}()

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      handlers.API(log, db),
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
	}

	serverErrors := make(chan error, 1)

	go func() {

		log.Printf("listening on %s", api.Addr)
		serverErrors <- api.ListenAndServe()
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		errors.Wrap(err, "error of run and listennig server")
	case <-shutdown:
		log.Println("starting shutdown")

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
	}
	return nil
}

func initConfig() error {
	viper.AddConfigPath("internal/platform/conf")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
