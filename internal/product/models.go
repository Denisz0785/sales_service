package product

import "time"

type Product struct {
	ID          string    `db:"id" json:"id"`
	Name        string    `db:"name" json:"name"`
	Cost        int       `db:"cost" json:"cost"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Sold        int       `db:"sold" json:"sold"`
	Revenue     int       `db:"revenue" json:"revenue"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
	DateUpdated time.Time `db:"date_updated" json:"date_updated"`
}

type NewProduct struct {
	Name     string `json:"name"`
	Cost     int    `json:"cost"`
	Quantity int    `json:"quantity"`
}

type Sale struct {
	ID          string    `db:"sale_id" json:"id"`
	ProductID   string    `db:"product_id" json:"product_id"`
	Quantity    int       `db:"quantity" json:"quantity"`
	Paid        int       `db:"paid" json:"paid"`
	DateCreated time.Time `db:"date_created" json:"date_created"`
}

type NewSale struct {
	Quantity int `json:"quantity"`
	Paid     int `json:"paid"`
}
