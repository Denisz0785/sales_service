package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sales_service/internal/platform/web"
	"sales_service/internal/product"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
)

// Product has methods for dealing with Products
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List send all products as list
func (p *Product) List(w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	list, err := product.List(ctx, p.DB)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error query to db", err)
		return err
	}

	return web.Respond(w, list, http.StatusOK)
}

func (p *Product) Retrieve(w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")
	ctx := r.Context()

	prod, err := product.Retrieve(ctx, p.DB, id)

	if err != nil {
		switch {
		case errors.Is(err, product.ErrNotFound):
			return web.NewRequestError(err, http.StatusNotFound)
		case errors.Is(err, product.ErrInvalidID):
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrapf(err, "looking up product %q", id)
		}

	}

	// Return the product
	return web.Respond(w, prod, http.StatusOK)
}

func (p *Product) Create(w http.ResponseWriter, r *http.Request) error {

	var newProduct product.NewProduct
	ctx := r.Context()

	if err := web.Decode(r, &newProduct); err != nil {
		return err
	}

	prod, err := product.Create(ctx, p.DB, newProduct, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(w, prod, http.StatusCreated)

}

func (p *Product) AddSale(w http.ResponseWriter, r *http.Request) error {
	var newSale product.NewSale
	if err := web.Decode(r, &newSale); err != nil {
		return errors.Wrap(err, "decode new sale")
	}

	productID := chi.URLParam(r, "id")
	fmt.Println("productID", productID)

	sale, err := product.AddSale(r.Context(), p.DB, newSale, productID, time.Now())

	if err != nil {
		return errors.Wrap(err, "add sale")
	}

	return web.Respond(w, sale, http.StatusCreated)
}

func (p *Product) ListSales(w http.ResponseWriter, r *http.Request) error {
	productID := chi.URLParam(r, "id")
	list, err := product.ListSales(r.Context(), p.DB, productID)
	if err != nil {
		return errors.Wrap(err, " gettinglist sales")
	}
	return web.Respond(w, list, http.StatusOK)
}
