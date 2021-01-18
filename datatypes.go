package main

import (
	"time"
)

type Landkreis struct {
	Name string
	Bundesland string
	Faelle int
	FaellePer100k float64
	FaellePer100k7d float64
	Tode int
	LastUpdate time.Time
}

type LandkreisJSON struct {
	Name string `json:"county"`
	Bundesland string `json:"BL"`
	Faelle int `json:"cases"`
	FaellePer100k float64 `json:"cases_per_100k"`
	FaellePer100k7d float64 `json:"cases7_per_100k"`
	Tode int `json:"deaths"`
	LastUpdate string `json:"last_update"`
}
