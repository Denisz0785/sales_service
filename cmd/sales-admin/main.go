package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"sales_service/internal/platform/database"
	"sales_service/internal/schema"
	"sales_service/internal/user"

	"github.com/go-faster/errors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/spf13/viper"

	"sales_service/internal/platform/auth"
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
		Args []string
	}

	cfg.Args = os.Args[1:]

	if len(cfg.Args) < 1 {
		return errors.New("should imput email and password")
	}

	if err := initConfig(); err != nil {
		return errors.Wrap(err, "error initializing configs")
	}

	if err := godotenv.Load("./cmd/sales-api/.env"); err != nil {
		return errors.Wrap(err, "error loading env variables")
	}

	cfg.DB.Name = viper.GetString("db.name")
	cfg.DB.User = viper.GetString("db.user")
	cfg.DB.Host = viper.GetString("db.host")
	cfg.DB.DisableTLS = viper.GetString("db.disableTLS")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	dbConfig := database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	}

	log.Println("started")
	defer log.Println("finished")

	var err error
	switch cfg.Args[0] {
	case "migrate":
		err = migrate(dbConfig)

	case "seed":
		err = seed(dbConfig)

	case "useradd":
		err = useradd(dbConfig, cfg.Args[1], cfg.Args[2])

	default:
		err = errors.New("invalid command")
	}
	if err != nil {
		return err
	}
	return nil
}

func migrate(cfg database.Config) error {
	db, err := database.OpenDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := schema.Migrate(db); err != nil {
		return err
	}
	log.Println("Migrations completed")
	return nil
}

func seed(cfg database.Config) error {
	db, err := database.OpenDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()
	if err := schema.Seed(db); err != nil {
		return err
	}
	log.Println("Insert data completed")
	return nil
}

func useradd(cfg database.Config, email, password string) error {
	db, err := database.OpenDB(cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	if email == "" || password == "" {
		return errors.New("email or password is empty")
	}

	fmt.Printf("Admin user will be created with email: %s and password: %s\n", email, password)
	fmt.Print("Do you want to continue? [1/0]: ")

	var confirm bool

	_, err = fmt.Scanf("%t\n", &confirm)
	if err != nil {
		return errors.Wrap(err, "read confirm")
	}
	if !confirm {
		fmt.Println("Aborted")
		return nil
	}

	ctx := context.Background()

	nu := user.NewUser{

		Email:           email,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []string{auth.RoleAdmin, auth.RoleUser},
	}

	u, err := user.Create(ctx, db, nu, time.Now())
	if err != nil {
		return err
	}
	fmt.Println("User was created with id:", u.ID)
	return nil
}

func initConfig() error {
	viper.AddConfigPath("internal/platform/conf")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
