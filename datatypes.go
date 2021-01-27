package main

import (
	"time"
)

type Landkreis struct {
	Name            string
	Bundesland      string
	Faelle          int
	FaellePer100k   float64
	FaellePer100k7d float64
	Tode            int
	LastUpdate      time.Time
}

type LandkreisJSON struct {
	Name            string  `json:"county"`
	Bundesland      string  `json:"BL"`
	Faelle          int     `json:"cases"`
	FaellePer100k   float64 `json:"cases_per_100k"`
	FaellePer100k7d float64 `json:"cases7_per_100k"`
	Tode            int     `json:"deaths"`
	LastUpdate      string  `json:"last_update"`
}

//
type Configuration struct {
	InfluxDB struct {
		URL      string `yaml:"url"`
		Name     string `yaml:"name"`
		User     string `yaml:"user"`
		Password string `yaml:"password"`
	}
	TimeInterval int    `yaml:"timeinterval-to-pull"`
	SingleRun    bool   `yaml:"single-run"`
	FederalState string `yaml:"federal-state"`
	Logging      struct {
		Dir   string `yaml:"log-dir"`
		Debug bool   `yaml:"debug"`
	}
}
