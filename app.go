package main

import (
	"fmt"
	"sync"
	"github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/log"
)

// App struct to hold refs and database info
type App struct {
	Url *string
	Name *string
	User *string
	Password *string
}

// Initialize app struct with database info
func (a *App) Initialize(url string, name string, user string, password string) {
	a.Url = &url
	a.Name = &name
	a.User = &user
	a.Password = &password
}

func (a *App) Run() {
	// Get Corona data from the RKI API
	data, err := a.GetData()
	if err != nil {
		Log.Logger.Warn(err)
	}
	// Parse Corona data
	landkreise, err := a.ParseData(&data)
	if err != nil {
		Log.Logger.Warn(err)
	}
	// Create InfluxDB client
	log.Log = nil // Disable log output of the InfluxDB Client
	client := influxdb2.NewClientWithOptions(*a.Url, fmt.Sprintf("%s:%s", *a.User, *a.Password), influxdb2.DefaultOptions().SetBatchSize(50))
    writeAPI := client.WriteAPI("", fmt.Sprintf("%s/autotgen", *a.Name))
	// Create wait grop for error channel
	var wg sync.WaitGroup
	// Create go proc for reading and logging errors
	errChannel := writeAPI.Errors()
	go func() {
		for err := range errChannel {
			defer wg.Done()
			wg.Add(1)
			Log.Logger.Warnw("Error while writing the data to InfluxDB.", "message", err.Error())
		}
	}()
	// Write data to InfluxDB
	for _, landkreis := range landkreise {
		// Create point
		p := influxdb2.NewPointWithMeasurement(landkreis.Name).
			AddTag("Bundesland", landkreis.Bundesland).
			AddField("Faelle", landkreis.Faelle).
			AddField("FaellePer100k", landkreis.FaellePer100k).
			AddField("FaellePer100k7d",landkreis.FaellePer100k7d).
			AddField("Tode", landkreis.Tode).
			SetTime(landkreis.LastUpdate)
		// Write asynchronously
		writeAPI.WritePoint(p)
		// Debug
		Log.Logger.Debugw("Store entry to InfluxDB.", "landkreis", landkreis.Name, "bundesland", landkreis.Bundesland, "faelle", landkreis.Faelle, "faelleper100k", landkreis.FaellePer100k, "faelleper100k7d", landkreis.FaellePer100k7d, "tode", landkreis.Tode, "lastupdate", landkreis.LastUpdate)
	}
    // Force all unwritten data to be sent
    writeAPI.Flush()
    // Ensures background processes finishes
    client.Close()
	// Wait for error go proc
	wg.Wait()
}
