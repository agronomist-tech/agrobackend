package main

import (
	"fmt"
	"github.com/agronomist-tech/agrobackend/handlers"
	"github.com/agronomist-tech/agrobackend/storage"
	"log"
	"net/http"
)


func main() {
	cfg := LoadConfig()

	db := storage.CreateDb(fmt.Sprintf("%s:%s", cfg.Clickhouse.Host, cfg.Clickhouse.Port))

	views := &handlers.Env{CH: db}

	log.Println(fmt.Sprintf("Start webserver on port: %d", cfg.Port))

	http.HandleFunc("/allPairs", views.AllPairs)
	http.HandleFunc("/getPrices", views.PairPrices)
	http.HandleFunc("/searchPairs", views.SearchPairs)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), nil))
}