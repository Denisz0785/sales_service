package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"sales_service/internal/product"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

// Product has methods for dealing with Products
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List send all products as list
func (p *Product) List(w http.ResponseWriter, r *http.Request) {

	list, err := product.List(p.DB)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error query to db", err)
		return
	}

	data, err := json.Marshal(list)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshalling", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing", err)
	}
}

func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) {

	id := chi.URLParam(r, "id")

	prod, err := product.Retrieve(p.DB, id)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error query to db", err)
		return
	}

	data, err := json.Marshal(prod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshalling", err)
		return
	}

	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing", err)
	}
}

func (p *Product) Create(w http.ResponseWriter, r *http.Request) {

	var newProduct product.NewProduct
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		p.Log.Println("error decoding", err)
		return
	}

	prod, err := product.Create(p.DB, newProduct, time.Now())

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error query to db", err)
		return
	}

	data, err := json.Marshal(prod)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error marshalling", err)
		return
	}

	w.Header().Set("content-type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(data); err != nil {
		p.Log.Println("error writing", err)
	}
}
