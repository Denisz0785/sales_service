package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sales_service/internal/product"

	"github.com/jmoiron/sqlx"
)

// Product has methods for dealing with Products
type Product struct {
	DB *sqlx.DB
}

// List send all products as list
func (p *Product) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(p.DB)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error query to db", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("error marshalling", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		log.Println("error writing", err)
	}
}
