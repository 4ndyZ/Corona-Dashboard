package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func (a *App) GetData(federalState string) (string, error) {
	var httpAPI string
	// API string for all German Federal States
	if strings.ToLower(federalState) == "all" {
		httpAPI = "https://services7.arcgis.com/mOBPykOjAyBO2ZKk/arcgis/rest/services/RKI_Landkreisdaten/FeatureServer/0/query?where=1%3D1&outFields=county,BL,cases,cases_per_100k,cases7_per_100k,deaths,last_update&returnGeometry=false&outSR=4326&f=json"
		// API string for a special German Federal State
	} else {
		httpAPI = strings.Join([]string{"https://services7.arcgis.com/mOBPykOjAyBO2ZKk/arcgis/rest/services/RKI_Landkreisdaten/FeatureServer/0/query?where=BL%20%3D%20%27", strings.ToLower(federalState), "%27&outFields=county,BL,cases,cases_per_100k,cases7_per_100k,deaths,last_update&returnGeometry=false&outSR=4326&f=json"}, "")
	}
	// Perform API HTTP call
	resp, err := http.Get(httpAPI)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	// Read response from HTTP call
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func (a *App) ParseData(body *string) ([]Landkreis, error) {
	// Regex to find important data
	re := regexp.MustCompile(`\{"attributes".*?\}`)
	// Slice with parsed data
	var landkreise []Landkreis
	// Search for data using Regex
	for _, value := range re.FindAllString(*body, -1) {
		// Cut off suffix
		value = strings.TrimPrefix(value, `{"attributes":`)
		// Structs
		var landkreisjson LandkreisJSON
		var landkreis Landkreis
		// Unmarshal the JSON data to struct
		err := json.Unmarshal([]byte(value), &landkreisjson)
		if err != nil {
			return landkreise, err
		}
		// Convert JSON struct to my struct
		landkreis.Name = landkreisjson.Name
		landkreis.Bundesland = landkreisjson.Bundesland
		landkreis.Faelle = landkreisjson.Faelle
		landkreis.FaellePer100k = landkreisjson.FaellePer100k
		landkreis.FaellePer100k7d = landkreisjson.FaellePer100k7d
		landkreis.Tode = landkreisjson.Tode
		// Parse time
		layout := "02.01.2006, 15:04 Uhr"
		t, err := time.Parse(layout, landkreisjson.LastUpdate)
		if err != nil {
			return landkreise, err
		}
		landkreis.LastUpdate = t
		// Append to parsed data the slice
		landkreise = append(landkreise, landkreis)
	}
	return landkreise, nil
}
