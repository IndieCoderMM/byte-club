package main

import (
	"log"
	"time"

	"github.com/indiecodermm/go-social/internal/db"
	"github.com/indiecodermm/go-social/internal/env"
	"github.com/indiecodermm/go-social/internal/store"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://user:password@localhost/go-social?sslmode=disable")

	conn, err := db.New(addr, 1, 1, time.Minute)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	store := store.NewStore(conn)

	db.Seed(store)
}
