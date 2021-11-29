package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type pair struct {
	Price      float64   `json:"price"`
	Pair       string    `json:"pair"`
	LastUpdate time.Time `json:"last_update"`
}

var Periods = [4]string{"24H", "7D", "1M", "3M"}

func (env *Env) AllPairs(w http.ResponseWriter, r *http.Request) {
	rows, err := env.CH.Query("SELECT toFloat64(price), pair, changeDate FROM prices WHERE dex='dexlab' AND changeDate > now() - INTERVAL 10 DAY ORDER BY pair, changeDate DESC LIMIT 1 BY pair")

	if err != nil {
		fmt.Println(err)
	}

	var pairs []pair
	for rows.Next() {
		var p pair

		err := rows.Scan(&p.Price, &p.Pair, &p.LastUpdate)
		if err != nil {
			fmt.Println(err)
		}
		pairs = append(pairs, p)
	}

	body, err := json.Marshal(pairs)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

type pairMarkers struct {
	Dates  []string  `json:"dates"`
	Prices []float64 `json:"prices"`
}

func (env *Env) PairPrices(w http.ResponseWriter, r *http.Request) {
	pair := r.URL.Query().Get("pair")
	period := r.URL.Query().Get("period")

	if pair == "" {
		pair = "SOL/USDC"
	}

	if period == "" {
		period = "24H"
	} else {
		exist := false

		for _, el := range Periods {
			if el == period {
				exist = true
			}
		}

		if exist == false {
			period = "24H"
		}
	}

	var rows *sql.Rows
	var err error

	if period == "24H" {
		rows, err = env.CH.Query("SELECT avg(price) as price, toStartOfTenMinutes(changeDate) AS change FROM prices WHERE pair = ? AND changeDate > now() - INTERVAL 24 HOUR GROUP BY toStartOfTenMinutes(changeDate) ORDER BY toStartOfTenMinutes(changeDate) ASC", pair)
	} else if period == "7D" {
		rows, err = env.CH.Query("select avg(price) as price, toStartOfHour(changeDate) as change from prices where pair = ? and changeDate > now() - INTERVAL 7 DAY group By toStartOfHour(changeDate) order by toStartOfHour(changeDate) asc", pair)
	} else if period == "1M" {
		rows, err = env.CH.Query("select avg(price) as price, toStartOfDay(changeDate) as change from prices where pair = ? and changeDate > now() - INTERVAL 1 MONTH group By toStartOfDay(changeDate) order by toStartOfDay(changeDate) asc", pair)
	} else if period == "3M" {
		rows, err = env.CH.Query("select avg(price) as price, toStartOfDay(changeDate) as change from prices where pair = ? and changeDate > now() - INTERVAL 3 MONTH group By toStartOfDay(changeDate) order by toStartOfDay(changeDate) asc", pair)
	}

	if err != nil {
		log.Println("Problem in SQL: ", err)
	}

	markers := pairMarkers{}
	for rows.Next() {
		var price float64
		var date string

		err := rows.Scan(&price, &date)
		if err != nil {
			fmt.Println(err)
		}

		markers.Prices = append(markers.Prices, price)
		markers.Dates = append(markers.Dates, date)
	}

	body, err := json.Marshal(markers)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}


func (env *Env) SearchPairs(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("query")

	var rows *sql.Rows
	var err error

	rows, err = env.CH.Query("SELECT DISTINCT pair FROM prices WHERE lowerUTF8(pair) LIKE ?", "%" + query + "%")
	if err != nil {
		fmt.Println(err)
	}

	var result []string

	for rows.Next() {
		var p string

		err := rows.Scan(&p)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, p)
	}

	body, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}