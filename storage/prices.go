package storage

import (
	"time"
)

type Price struct {
	Dex string
	Price float64
	Pair string
	Datetime time.Time
	LastUpdate time.Time
}




