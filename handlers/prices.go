package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


type pair struct {
	Price float64 `json:"price"`
	Pair string `json:"pair"`
	LastUpdate time.Time `json:"last_update"`
}


func (env *Env) AllPairs(w http.ResponseWriter, r *http.Request) {
	rows, err := env.CH.Query("select price, pair, changeDate from prices where dex='dexlab' order by pair desc, changeDate desc limit 1 by pair")

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