.PHONY: psql api

psql:
	docker exec -it sales_db psql postgres://postgres:example@db:5432/postgres
api:
	go run ./cmd/sales-api/main.go
