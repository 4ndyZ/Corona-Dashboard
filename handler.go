package main

import (
	"encoding/csv"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
	_ "time/tzdata"
)

func (a *App) GetDataCounty(state string) (string, error) {
	var httpAPI string
	// API string for all German Federal states
	if strings.ToLower(state) == "all" {
		httpAPI = "https://services7.arcgis.com/mOBPykOjAyBO2ZKk/arcgis/rest/services/RKI_Landkreisdaten/FeatureServer/0/query?where=1%3D1&outFields=county,BL,cases,cases_per_100k,cases7_per_100k,deaths,last_update&returnGeometry=false&outSR=4326&f=json"
		// API string for a special German Federal state
	} else {
		httpAPI = strings.Join([]string{"https://services7.arcgis.com/mOBPykOjAyBO2ZKk/arcgis/rest/services/RKI_Landkreisdaten/FeatureServer/0/query?where=BL%20%3D%20%27", strings.ToLower(state), "%27&outFields=county,BL,cases,cases_per_100k,cases7_per_100k,deaths,last_update&returnGeometry=false&outSR=4326&f=json"}, "")
	}
	// Perform API HTTP call
	resp, err := http.Get(httpAPI)
	if err != nil {
		return "", err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {

		}
	}(resp.Body)
	// Read response from HTTP call
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (a *App) GetDataVaccination() ([][]string, error) {
	// Data
	httpCSV := "https://impfdashboard.de/static/data/germany_vaccinations_timeseries_v2.tsv"
	// Perform HTTP call
	resp, err := http.Get(httpCSV)
	if err != nil {
		return [][]string{}, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	reader := csv.NewReader(resp.Body)
	reader.Comma = '\t'
	body, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	return body, nil
}

func (a *App) ParseDataCounty(body *string) ([]County, error) {
	// Constants
	loc, _ := time.LoadLocation("Europe/Berlin")
	layout := "02.01.2006, 15:04 Uhr"
	// Regex to find important data
	re := regexp.MustCompile(`\{"attributes".*?\}`)
	// Slice with parsed data
	var counties []County
	// Search for data using Regex
	for _, value := range re.FindAllString(*body, -1) {
		// Cut off suffix
		value = strings.TrimPrefix(value, `{"attributes":`)
		// Structs
		var countyJSON CountyJSON
		var county County
		// Unmarshal the JSON data to struct
		err := json.Unmarshal([]byte(value), &countyJSON)
		if err != nil {
			return counties, err
		}
		// Convert JSON struct to my struct
		county.Name = countyJSON.Name
		county.State = countyJSON.State
		county.Cases = countyJSON.Cases
		county.CasesPer100k = countyJSON.CasesPer100k
		county.CasesPer100k7d = countyJSON.CasesPer100k7d
		county.Deaths = countyJSON.Deaths
		// Parse date
		t, err := time.ParseInLocation(layout, countyJSON.LastUpdate, loc)
		if err != nil {
			return counties, err
		}
		county.LastUpdate = t
		// Append to parsed data the slice
		counties = append(counties, county)
	}
	return counties, nil
}

func (a *App) ParseDataVaccination(body *[][]string) ([]Vaccination, error) {
	// Constants
	loc, _ := time.LoadLocation("Europe/Berlin")
	layout := "2006-01-02"
	// Slice with parsed data
	var vaccinations []Vaccination
	//
	b := *body
	//
	for idx, row := range b {
		// Struct
		var vaccination Vaccination
		// Skip first line (CSV header)
		if idx == 0 {
			continue
		}
		for key, value := range row {
			// Go through the columns of each row
			switch b[0][key] {
			case "dosen_kumulativ":
				vaccination.Doses.All = a.MustStringToInt(value)
			case "dosen_biontech_kumulativ":
				vaccination.Doses.Biontech = a.MustStringToInt(value)
			case "dosen_moderna_kumulativ":
				vaccination.Doses.Moderna = a.MustStringToInt(value)
			case "dosen_astra_kumulativ":
				vaccination.Doses.AstraZeneca = a.MustStringToInt(value)
			case "dosen_johnson_kumulativ":
				vaccination.Doses.Johnson = a.MustStringToInt(value)
			case "dosen_novavax_kumulativ":
				vaccination.Doses.Novavax = a.MustStringToInt(value)
			case "personen_erst_kumulativ":
				vaccination.People.FirstTime = a.MustStringToInt(value)
			case "personen_voll_kumulativ":
				vaccination.People.Full = a.MustStringToInt(value)
			case "personen_auffrisch_kumulativ":
				vaccination.People.Refreshment = a.MustStringToInt(value)
			case "impf_quote_erst":
				vaccination.Rate.FirstTime = a.MustStringToFloat(value)
			case "impf_quote_voll":
				vaccination.Rate.Full = a.MustStringToFloat(value)
			case "date":
				// Parse date
				t, err := time.ParseInLocation(layout, value, loc)
				if err != nil {
					return vaccinations, err
				}
				// Add 10h to the date because data gets a refresh very day at 10:00 e
				hours, _ := time.ParseDuration("10h")
				t = t.Add(hours)
				vaccination.LastUpdate = t
			default:
			}
		}
		// Append to parsed data the slice
		vaccinations = append(vaccinations, vaccination)
	}
	return vaccinations, nil
}
