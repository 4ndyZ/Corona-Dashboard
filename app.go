package main

import (
	"fmt"
	"log"
	"os"
	"github.com/influxdata/influxdb-client-go"
)

// App struct to hold refs and database info
type App struct {
	Logger *log.Logger
	Url *string
	User *string
	Password *string
}

// Initialize app struct with database info and logger
func (a *App) Initialize(url string, user string, password string) {
	a.Logger = log.New(os.Stdout, "", log.LstdFlags)
	a.Url = &url
	a.User = &user
	a.Password = &password
}

func (a *App) Run() {
	// Get Corona data from the RKI API
	data, err := a.GetData()
	if err != nil {
		a.Logger.Fatal(err)
	}
	// Parse Corona data
	landkreise, err := a.ParseData(&data)
	if err != nil {
		a.Logger.Fatal(err)
	}
	// Create InfluxDB client
	client := influxdb2.NewClientWithOptions(*a.Url, fmt.Sprintf("%s:%s", a.User, a.Password), influxdb2.DefaultOptions().SetBatchSize(50))
    writeAPI := client.WriteAPI("", "corona/autogen")
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
		fmt.Println(landkreis)
	}
    // Force all unwritten data to be sent
    writeAPI.Flush()
    // Ensures background processes finishes
    client.Close()
}
