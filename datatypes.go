package main

import (
	"time"
)

type County struct {
	Name           string
	State          string
	Cases          int
	CasesPer100k   float64
	CasesPer100k7d float64
	Deaths         int
	LastUpdate     time.Time
}

type CountyJSON struct {
	Name           string  `json:"county"`
	State          string  `json:"BL"`
	Cases          int     `json:"cases"`
	CasesPer100k   float64 `json:"cases_per_100k"`
	CasesPer100k7d float64 `json:"cases7_per_100k"`
	Deaths         int     `json:"deaths"`
	LastUpdate     string  `json:"last_update"`
}

type Vaccination struct {
	Doses struct {
		All         int
		Biontech    int
		Moderna     int
		AstraZeneca int
        Johnson     int
	}
	People struct {
		FirstTime int
		Full      int
	}
	Rate struct {
		FirstTime float64
		Full      float64
	}
	LastUpdate time.Time
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
