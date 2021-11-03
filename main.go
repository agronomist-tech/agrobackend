package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/agronomist-tech/agrobackend/storage"
	"github.com/agronomist-tech/agrobackend/handlers"
)


func main() {
	cfg := LoadConfig()

	db := storage.CreateDb(fmt.Sprintf("%s:%s", cfg.Clickhouse.Host, cfg.Clickhouse.Port))

	views := &handlers.Env{CH: db}

	http.HandleFunc("/allPairs", views.AllPairs)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}