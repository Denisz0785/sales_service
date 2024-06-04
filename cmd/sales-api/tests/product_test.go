package tests

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"sales_service/cmd/sales-api/internal/handlers"
	"sales_service/internal/platform/database/databasetest"
	"sales_service/internal/schema"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestProducts(t *testing.T) {
	db, teardown := databasetest.Setup(t)
	defer teardown()

	if err := schema.Seed(db); err != nil {
		t.Fatal(err)
	}
	log := log.New(os.Stderr, "TEST : ", log.LstdFlags|log.Lshortfile)

	tests := ProductTests{app: handlers.API(log, db)}

	t.Run("List", tests.List)
	t.Run("ProductCRUD", tests.ProductCRUD)
}

type ProductTests struct {
	app http.Handler
}

func (p ProductTests) List(t *testing.T) {

	req := httptest.NewRequest("GET", "/v1/products", nil)
	resp := httptest.NewRecorder()

	p.app.ServeHTTP(resp, req)

	if resp.Code != http.StatusOK {
		t.Fatalf("expected %d, actual %d", http.StatusOK, resp.Code)
	}

	var list []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&list); err != nil {
		t.Fatalf("decoding: %s", err)
	}

	want := []map[string]interface{}{
		{
			"id":           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd390a21",
			"name":         "Lego City",
			"cost":         float64(3000),
			"quantity":     float64(56),
			"date_created": "2024-05-05T12:12:12Z",
			"date_updated": "2024-05-06T14:15:12Z",
		},
		{
			"id":           "a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11",
			"name":         "Lego Chima",
			"cost":         float64(2000),
			"quantity":     float64(50),
			"date_created": "2024-05-05T12:12:12Z",
			"date_updated": "2024-05-06T14:15:12Z",
		},
	}
	if diff := cmp.Diff(want, list); diff != "" {
		t.Fatalf("mismatch (-want +got):\n%s", diff)
	}
}

func (p ProductTests) ProductCRUD(t *testing.T) {
	var created map[string]interface{}

	{
		body := strings.NewReader(`{"name": "test product3", "cost": 55, "quantity": 20}`)

		req := httptest.NewRequest("POST", "/v1/products", body)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Fatalf("expected %d, actual %d", http.StatusCreated, resp.Code)
		}

		if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if created["id"] == "" || created["id"] == nil {
			t.Fatal("missing id")
		}

		if created["date_created"] == "" || created["date_updated"] == nil {
			t.Fatal("missing date")
		}

		if created["name"] != "test product3" {
			t.Fatalf("expected %q, actual %q", "test product3", created["name"])
		}

		want := map[string]interface{}{
			"id":           created["id"],
			"name":         "test product3",
			"cost":         float64(55),
			"quantity":     float64(20),
			"date_created": created["date_created"],
			"date_updated": created["date_updated"],
		}

		if diff := cmp.Diff(want, created); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}

	{
		url := fmt.Sprintf("/v1/products/%s", created["id"])
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		p.app.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Fatalf("expected %d, actual %d", http.StatusOK, resp.Code)
		}

		var got map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
			t.Fatalf("decoding: %s", err)
		}

		if diff := cmp.Diff(created, got); diff != "" {
			t.Fatalf("mismatch (-want +got):\n%s", diff)
		}
	}
}
