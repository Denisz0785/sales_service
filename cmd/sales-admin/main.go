package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
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

// run is the entry point of the application.
func main() {
	// Run the application and handle any errors.
	if err := run(); err != nil {
		// Log the error and exit the program.
		log.Fatal(err)
	}
}

// run is the entry point of the application.
// It handles the command line arguments and calls the appropriate functions.
func run() error {
	// Define the configuration struct
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

	// Get the command line arguments
	cfg.Args = os.Args[1:]

	// Check if the required arguments are provided
	if len(cfg.Args) < 1 {
		return errors.New("should input email and password")
	}

	// Initialize the configuration
	if err := initConfig(); err != nil {
		return errors.Wrap(err, "error initializing configs")
	}

	// Load the environment variables from the .env file
	if err := godotenv.Load("./cmd/sales-api/.env"); err != nil {
		return errors.Wrap(err, "error loading env variables")
	}

	// Set the database configuration values from the environment variables
	cfg.DB.Name = viper.GetString("db.name")
	cfg.DB.User = viper.GetString("db.user")
	cfg.DB.Host = viper.GetString("db.host")
	cfg.DB.DisableTLS = viper.GetString("db.disableTLS")
	cfg.DB.Password = os.Getenv("DB_PASSWORD")

	// Create the database configuration
	dbConfig := database.Config{
		Host:       cfg.DB.Host,
		Name:       cfg.DB.Name,
		User:       cfg.DB.User,
		Password:   cfg.DB.Password,
		DisableTLS: cfg.DB.DisableTLS,
	}

	// Log the start of the application
	log.Println("started")
	defer log.Println("finished")

	// Switch on the command line argument and call the corresponding function
	var err error
	switch cfg.Args[0] {
	case "migrate":
		err = migrate(dbConfig)

	case "seed":
		err = seed(dbConfig)

	case "useradd":
		err = useradd(dbConfig, cfg.Args[1], cfg.Args[2])

	case "keygen":
		err = keygen(cfg.Args[1])

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

// keygen generates a new RSA private key and writes it to the specified file path.
func keygen(path string) error {
	// Check if the file path is empty.
	if path == "" {
		return errors.New("path is empty")
	}

	// Generate a new RSA private key.
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		// Return an error if the key generation fails.
		return errors.Wrap(err, "generating key")
	}

	// Create a new file to write the key to.
	file, err := os.Create(path)
	if err != nil {
		// Return an error if the file creation fails.
		return errors.Wrap(err, "creating file")
	}
	defer file.Close() // Close the file when done.

	// Create a PEM block containing the private key.
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	// Encode the key and write it to the file.
	if err := pem.Encode(file, block); err != nil {
		// Return an error if the key encoding or writing fails.
		return errors.Wrap(err, "encoding key")
	}

	// Return nil if the key generation, file creation, and key writing are successful.
	return nil
}

func initConfig() error {
	viper.AddConfigPath("internal/platform/conf")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
