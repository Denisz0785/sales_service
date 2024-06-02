package main

import (
	"flag"
	"log"
	"os"

	"sales_service/internal/platform/database"
	"sales_service/internal/schema"

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
	}

	log.Println("started")
	defer log.Println("finished")

	cfg.DB.Name = viper.GetString("db.name")
	cfg.DB.User = viper.GetString("db.user")
	cfg.DB.Host = viper.GetString("db.host")
	cfg.DB.DisableTLS = viper.GetString("db.disableTLS")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load("./cmd/sales-api/.env"); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

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

	flag.Parse()
	switch flag.Arg(0) {
	case "migrate":
		if err := schema.Migrate(db); err != nil {
			log.Fatal("doing migrate", err)
		}
		log.Println("Migrations completed")
		return
	case "seed":
		if err := schema.Seed(db); err != nil {
			log.Fatal("doing insert data", err)
		}
		log.Println("Insert data completed")
		return
	}

}
func initConfig() error {
	viper.AddConfigPath("internal/platform/conf")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
