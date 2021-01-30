package main

import (
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/log"
	"strings"
	"sync"
)

// App struct to hold refs and database info
type App struct {
	Url      *string
	Name     *string
	User     *string
	Password *string
}

// Initialize app struct with database info
func (a *App) Initialize(url string, name string, user string, password string) {
	a.Url = &url
	a.Name = &name
	a.User = &user
	a.Password = &password
}

func (a *App) Run(federalState string) {
	// Get Corona data from the RKI API
	countydata, err := a.GetDataCounty(federalState)
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Error while getting the data from the RKI county API.")
	}
	// Parse Corona data
	counties, err := a.ParseDataCounty(&countydata)
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Error while parsing the data from the RKI county API.")
	}
	// Get Vaccination data from the Vaccination dashboard
	vaccinationdata, err := a.GetDataVaccination()
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Error while getting the data from the vaccination dashboard.")
	}
	// Parse vaccination data
	vaccinations, err := a.ParseDataVaccination(&vaccinationdata)
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Error while parsing the data from the vaccination dashboard.")
	}
	// Parse
	// Create InfluxDB client
	log.Log = nil // Disable log output of the InfluxDB Client
	client := influxdb2.NewClientWithOptions(*a.Url, strings.Join([]string{*a.User, ":", *a.Password}, ""), influxdb2.DefaultOptions().SetBatchSize(50))
	writeAPI := client.WriteAPI("", strings.Join([]string{*a.Name, "/autogen"}, ""))
	// Create wait grop for error channel
	var wg sync.WaitGroup
	// Create go proc for reading and logging errors
	errChannel := writeAPI.Errors()
	go func() {
		for err := range errChannel {
			defer wg.Done()
			wg.Add(1)
			Log.Logger.Warn().Str("error", err.Error()).Msg("Error while writing the data to InfluxDB.")
		}
	}()
	// Write county data to InfluxDB
	for _, county := range counties {
		// Create point
		p := influxdb2.NewPoint(
			"county",
			map[string]string{
				"County": county.Name,
				"State":  county.State,
			},
			map[string]interface{}{
				"Cases":          county.Cases,
				"CasesPer100k":   county.CasesPer100k,
				"CasesPer100k7d": county.CasesPer100k7d,
				"Deaths":         county.Deaths,
			},
			county.LastUpdate)
		// Write asynchronously
		writeAPI.WritePoint(p)
		// Debug
		Log.Logger.Debug().
			Str("county", county.Name).
			Str("state", county.State).
			Int("cases", county.Cases).
			Float64("casesper100k", county.CasesPer100k).
			Float64("casesper100k7d", county.CasesPer100k7d).
			Int("deaths", county.Deaths).
			Time("lastupdate", county.LastUpdate).
			Msg("Store county entry to InfluxDB.")
	}
	// Write vaccination data to InfluxDB
	for _, vaccination := range vaccinations {
		// Create point
		p1 := influxdb2.NewPoint(
			"vaccination",
			map[string]string{
				"Manufacturer": "All",
			},
			map[string]interface{}{
				"Doses":    vaccination.Doses.All,
			},
			vaccination.LastUpdate)
		p2 := influxdb2.NewPoint(
			"vaccination",
			map[string]string{
				"Manufacturer": "Biontech",
			},
			map[string]interface{}{
				"Doses": vaccination.Doses.Biontech,
			},
			vaccination.LastUpdate)
		p3 := influxdb2.NewPoint(
			"vaccination",
			map[string]string{
				"Manufacturer": "Moderna",
			},
			map[string]interface{}{
				"Doses": vaccination.Doses.Moderna,
			},
			vaccination.LastUpdate)
		p4 := influxdb2.NewPoint(
			"vaccination",
			map[string]string{
				"Manufacturer": "AstraZeneca",
			},
			map[string]interface{}{
				"Doses": vaccination.Doses.AstraZeneca,
			},
			vaccination.LastUpdate)
		p5 := influxdb2.NewPoint(
			"vaccination",
			map[string]string{
				"Typ": "FirstTime",
			},
			map[string]interface{}{
				"People": vaccination.People.FirstTime,
				"Rate":    vaccination.Rate.FirstTime,
			},
			vaccination.LastUpdate)
		p6 := influxdb2.NewPoint(
			"vaccination",
			map[string]string{
				"Typ": "Full",
			},
			map[string]interface{}{
				"People": vaccination.People.Full,
				"Rate":    vaccination.Rate.Full,
			},
			vaccination.LastUpdate)
		// Write asynchronously
		writeAPI.WritePoint(p1)
		writeAPI.WritePoint(p2)
		writeAPI.WritePoint(p3)
		writeAPI.WritePoint(p4)
		writeAPI.WritePoint(p5)
		writeAPI.WritePoint(p6)
		// Debug
		Log.Logger.Debug().
			Int("doses-all", vaccination.Doses.All).
			Int("doses-biontech", vaccination.Doses.Biontech).
			Int("doses-moderna", vaccination.Doses.Moderna).
			Int("doses-astrazeneca", vaccination.Doses.AstraZeneca).
			Int("people-firsttime", vaccination.People.FirstTime).
			Int("people-full", vaccination.People.Full).
			Float64("rate-firsttime", vaccination.Rate.FirstTime).
			Float64("rate-full", vaccination.Rate.Full).
			Msg("Store vaccination entry to InfluxDB.")
	}
	// Force all unwritten data to be sent
	writeAPI.Flush()
	// Ensures background processes finishes
	client.Close()
	// Wait for error go proc
	wg.Wait()
}
