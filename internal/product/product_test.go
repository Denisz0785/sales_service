package product

import (
	"sales_service/internal/platform/database/databasetest"
	"sales_service/internal/schema"

	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	NewProduct := NewProduct{
		Name:     "test product",
		Cost:     10,
		Quantity: 20,
	}
	now := time.Date(2024, 5, 5, 5, 5, 5, 0, time.UTC)
	product1, err := Create(db, NewProduct, now)
	if err != nil {
		t.Fatal(err)
	}

	product2, err := Retrieve(db, product1.ID)
	if err != nil {
		t.Fatal(err)
	}

	if diff := cmp.Diff(product1, product2); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
}

func TestProductList(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}

	products, err := List(db)
	if err != nil {
		t.Fatal(err)
	}

	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
}
