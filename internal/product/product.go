package product

import "github.com/jmoiron/sqlx"

// List returns all Products
func List(db *sqlx.DB) ([]Product, error) {

	list := []Product{}

	const q = `select * from products`

	if err := db.Select(&list, q); err != nil {
		return nil, err
	}

	return list, nil
}
