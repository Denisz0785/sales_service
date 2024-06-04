package product

import (
	"context"
	"database/sql"
	"time"

	"github.com/go-faster/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

var (
	ErrNotFound  = errors.New("product not found")
	ErrInvalidID = errors.New("invalid product ID format")
)

// List retrieves all products from the database.
func List(ctx context.Context, db *sqlx.DB) ([]Product, error) {
	// Create a slice to store the retrieved data.
	var list []Product

	// Define the SQL query to retrieve all products.
	const query = `select * from products`

	// Use the Select method of the sqlx.DB connection to execute the query
	// and store the result in the list variable.
	if err := db.SelectContext(ctx, &list, query); err != nil {
		return nil, err
	}

	// Return the retrieved list of products and nil for the error.
	return list, nil
}

// Retrieve retrieves a single product from the database
func Retrieve(ctx context.Context, db *sqlx.DB, id string) (*Product, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, ErrInvalidID
	}

	// Create a new Product variable to store the retrieved data.
	var p Product

	// Define the SQL query to retrieve a single product by ID.
	const q = `select * from products where id = $1`

	// Execute the query to retrieve a single product by ID.
	if err := db.GetContext(ctx, &p, q, id); err != nil {

		// If it is, return the ErrNotFound error.
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		// If there is an error, return the empty Product and the error.
		return nil, err
	}

	// Return the retrieved Product and nil for the error.
	return &p, nil
}

// Create inserts a new product into the database
func Create(ctx context.Context, db *sqlx.DB, newProduct NewProduct, currentTime time.Time) (*Product, error) {
	product := &Product{
		ID:          uuid.New().String(),
		Name:        newProduct.Name,
		Cost:        newProduct.Cost,
		Quantity:    newProduct.Quantity,
		DateCreated: currentTime,
		DateUpdated: currentTime,
	}

	const query = `INSERT INTO products(id, name, cost, quantity, date_created, date_updated) VALUES($1, $2, $3, $4, $5, $6) RETURNING *`
	productFromDB := make([]Product, 1)

	err := db.SelectContext(ctx, &productFromDB, query, product.ID, product.Name, product.Cost, product.Quantity, product.DateCreated, product.DateUpdated)
	if err != nil {
		return nil, errors.Wrapf(err, "inserting product: %v", product)
	}

	return &productFromDB[0], nil
}
