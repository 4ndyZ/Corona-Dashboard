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
	data, err := a.GetData(federalState)
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Error while getting the data from the RKI API.")
	}
	// Parse Corona data
	landkreise, err := a.ParseData(&data)
	if err != nil {
		Log.Logger.Warn().Str("error", err.Error()).Msg("Error while parsing the data from the RKI API.")
	}
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
	// Write data to InfluxDB
	for _, landkreis := range landkreise {
		// Create point
		p := influxdb2.NewPointWithMeasurement(landkreis.Name).
			AddTag("Bundesland", landkreis.Bundesland).
			AddField("Faelle", landkreis.Faelle).
			AddField("FaellePer100k", landkreis.FaellePer100k).
			AddField("FaellePer100k7d", landkreis.FaellePer100k7d).
			AddField("Tode", landkreis.Tode).
			SetTime(landkreis.LastUpdate)
		// Write asynchronously
		writeAPI.WritePoint(p)
		// Debug
		Log.Logger.Debug().Str("landkreis", landkreis.Name).Str("bundesland", landkreis.Bundesland).Int("faelle", landkreis.Faelle).Float64("faelleper100k", landkreis.FaellePer100k).Float64("faelleper100k7d", landkreis.FaellePer100k7d).Int("tode", landkreis.Tode).Time("lastupdate", landkreis.LastUpdate).Msg("Store entry to InfluxDB.")
	}
	// Force all unwritten data to be sent
	writeAPI.Flush()
	// Ensures background processes finishes
	client.Close()
	// Wait for error go proc
	wg.Wait()
}
