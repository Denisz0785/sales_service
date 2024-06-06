package product

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

func AddSale(ctx context.Context, db *sqlx.DB, ns NewSale, ProductID string, now time.Time) (*Sale, error) {

	s := Sale{
		ID:          uuid.New().String(),
		ProductID:   ProductID,
		Quantity:    ns.Quantity,
		Paid:        ns.Paid,
		DateCreated: now,
	}

	const q = `
	INSERT INTO sales (sale_id, product_id, quantity, paid, date_created)
	VALUES ($1, $2, $3, $4, $5)`
	_, err := db.ExecContext(ctx, q, s.ID, s.ProductID, s.Quantity, s.Paid, s.DateCreated)

	if err != nil {
		return nil, err
	}

	return &s, nil
}

func ListSales(ctx context.Context, db *sqlx.DB, productID string) ([]Sale, error) {
	list := []Sale{}

	const q = `SELECT * FROM sales WHERE product_id = $1`
	if err := db.SelectContext(ctx, &list, q, productID); err != nil {
		return nil, errors.Wrap(err, "selecting sales")
	}
	return list, nil
}
