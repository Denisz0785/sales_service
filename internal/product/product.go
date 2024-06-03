package product

import "github.com/jmoiron/sqlx"

// List retrieves all products from the database.
func List(db *sqlx.DB) ([]Product, error) {
	// Create a slice to store the retrieved data.
	var list []Product

	// Define the SQL query to retrieve all products.
	const query = `select * from products`

	// Use the Select method of the sqlx.DB connection to execute the query
	// and store the result in the list variable.
	if err := db.Select(&list, query); err != nil {
		return nil, err
	}

	// Return the retrieved list of products and nil for the error.
	return list, nil
}

// Retrieve retrieves a single product from the database
func Retrieve(db *sqlx.DB, id string) (*Product, error) {
	// Create a new Product variable to store the retrieved data.
	var p Product

	// Define the SQL query to retrieve a single product by ID.
	const q = `select * from products where id = $1`

	// Call the Select method of the *sqlx.DB connection to execute the query
	// and store the result in the Product variable.
	// The $1 placeholder is replaced with the value of the id parameter.
	if err := db.Get(&p, q, id); err != nil {
		// If there is an error, return the empty Product and the error.
		return nil, err
	}

	// Return the retrieved Product and nil for the error.
	return &p, nil
}
