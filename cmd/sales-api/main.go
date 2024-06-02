package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sales_service/cmd/sales-api/internal/handlers"
	"sales_service/internal/platform/database"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {

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
			ReadTimeout     time.Duration
			WriteTimeout    time.Duration
			ShutdownTimeout time.Duration
		}
	}

	log.Println("started")
	defer log.Println("finished")

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load("./cmd/sales-api/.env"); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
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
		log.Fatalf("error connect to DB %v", err)
	}
	defer db.Close()

	ps := handlers.Product{DB: db}

	cfg.Web.ReadTimeout = viper.GetDuration("web.readtimeout")
	cfg.Web.WriteTimeout = viper.GetDuration("web.writetimeout")
	cfg.Web.ShutdownTimeout = viper.GetDuration("web.shutdowntimeout")
	cfg.Web.Address = viper.GetString("web.address")

	api := http.Server{
		Addr:         cfg.Web.Address,
		Handler:      http.HandlerFunc(ps.List),
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
		log.Printf("error of listenig %v", err)
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
			log.Fatalf("could not stop server gracefully %v", err)
		}
	}
}

func initConfig() error {
	viper.AddConfigPath("internal/platform/conf")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
