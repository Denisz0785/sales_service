package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sales_service/internal/platform/auth"
	"sales_service/internal/platform/web"
	"sales_service/internal/product"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-faster/errors"
	"github.com/jmoiron/sqlx"
	"go.opencensus.io/trace"
)

// Product has methods for dealing with Products
type Product struct {
	DB  *sqlx.DB
	Log *log.Logger
}

// List send all products as list
func (p *Product) List(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	ctx, span := trace.StartSpan(ctx, "handlers.Product.List")
	defer span.End()

	list, err := product.List(ctx, p.DB)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		p.Log.Println("error query to db", err)
		return err
	}

	return web.Respond(ctx, w, list, http.StatusOK)
}

func (p *Product) Retrieve(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	id := chi.URLParam(r, "id")

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
	return web.Respond(ctx, w, prod, http.StatusOK)
}

func (p *Product) Create(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var newProduct product.NewProduct

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from request context")
	}

	if err := web.Decode(r, &newProduct); err != nil {
		return err
	}

	fmt.Println("newProduct", newProduct)

	prod, err := product.Create(ctx, p.DB, claims, newProduct, time.Now())
	if err != nil {
		return err
	}

	return web.Respond(ctx, w, prod, http.StatusCreated)

}

func (p *Product) Update(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	var update product.UpdateProduct
	if err := web.Decode(r, &update); err != nil {
		return errors.Wrap(err, "decode update product")
	}

	claims, ok := ctx.Value(auth.Key).(auth.Claims)
	if !ok {
		return errors.New("claims missing from request context")
	}

	if err := product.Update(ctx, p.DB, claims, id, update, time.Now()); err != nil {
		switch err {
		case product.ErrNotFound:
			return web.NewRequestError(err, http.StatusNotFound)
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		case product.ErrForbidden:
			return web.NewRequestError(err, http.StatusForbidden)
		default:
			return errors.Wrap(err, "updating product")
		}
	}
	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (p *Product) Delete(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	id := chi.URLParam(r, "id")

	err := product.Delete(ctx, p.DB, id)
	if err != nil {
		switch err {
		case product.ErrInvalidID:
			return web.NewRequestError(err, http.StatusBadRequest)
		default:
			return errors.Wrap(err, "deleting product")
		}
	}

	return web.Respond(ctx, w, nil, http.StatusNoContent)
}

func (p *Product) AddSale(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	var newSale product.NewSale
	if err := web.Decode(r, &newSale); err != nil {
		return errors.Wrap(err, "decode new sale")
	}

	productID := chi.URLParam(r, "id")
	fmt.Println("productID", productID)

	sale, err := product.AddSale(ctx, p.DB, newSale, productID, time.Now())

	if err != nil {
		return errors.Wrap(err, "add sale")
	}

	return web.Respond(ctx, w, sale, http.StatusCreated)
}

func (p *Product) ListSales(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	productID := chi.URLParam(r, "id")
	list, err := product.ListSales(ctx, p.DB, productID)
	if err != nil {
		return errors.Wrap(err, " gettinglist sales")
	}
	return web.Respond(ctx, w, list, http.StatusOK)
}
