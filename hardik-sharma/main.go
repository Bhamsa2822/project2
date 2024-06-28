package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func main() {
	//postgres configuration
	connStr := "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(connStr)))
	db := bun.NewDB(sqldb, pgdialect.New())

	repo := NewPostgresRepo(db)
	service := NewService(repo)
	handler := NewCustomerHandler(service)
	r := registerRoutes(handler)

	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Println("server exited:", err)
	}
}
