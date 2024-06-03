package main

import (
	"flag"
	"log"
	"os"

	"sales_service/internal/platform/database"
	"sales_service/internal/schema"

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

	var cfg struct {
		DB struct {
			User       string
			Password   string
			Host       string
			Name       string
			DisableTLS string
		}
	}

	log.Println("started")
	defer log.Println("finished")

	cfg.DB.Name = viper.GetString("db.name")
	cfg.DB.User = viper.GetString("db.user")
	cfg.DB.Host = viper.GetString("db.host")
	cfg.DB.DisableTLS = viper.GetString("db.disableTLS")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	if err := initConfig(); err != nil {
		return errors.Wrap(err, "error initializing configs")
	}

	if err := godotenv.Load("./cmd/sales-api/.env"); err != nil {
		return errors.Wrap(err, "error loading env variables")
	}

	db, err := database.OpenDB(database.Config{
		Host:       cfg.DB.Host,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		Name:       cfg.DB.Name,
		DisableTLS: cfg.DB.DisableTLS,
	})
	if err != nil {
		return errors.Wrap(err, "error connecting to db")
	}
	defer db.Close()

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			return errors.Wrap(err, "doing migrate")
		}
		log.Println("Migrations completed")
		return nil
	case "seed":
		if err := schema.Seed(db); err != nil {
			return errors.Wrap(err, "doing insert data")
		}
		log.Println("Insert data completed")
		return nil
	}
	return nil

}
func initConfig() error {
	viper.AddConfigPath("internal/platform/conf")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
