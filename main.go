package main

import (
	"flag"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"time"
	"gopkg.in/yaml.v2"
)

var Log Logger = Logger{}

func main() {
	Log.Initialize(false, "/var/log/corona-dashboard/log.txt") // Still hardcoded 

	body, _ := ioutil.ReadFile("/etc/corona-dashboard/config.conf")

	configuration := Configuration{}

	err := yaml.Unmarshal([]byte(body), &configuration)
	if err != nil {
		Log.Logger.Warn("No configuration file error")
	}

	flag.StringVar(&configuration.InfluxDB.URL, "dburl", configuration.InfluxDB.URL, "Database connection URL including the port (e.g. https://myinfluxdb-server.tld:8086)")
	flag.StringVar(&configuration.InfluxDB.Name, "dbname", configuration.InfluxDB.Name, "Name of the InfluxDB database")
	flag.StringVar(&configuration.InfluxDB.User, "dbuser", configuration.InfluxDB.User, "Username of the InfluxDB database user")
	flag.StringVar(&configuration.InfluxDB.Password, "dbpassword", configuration.InfluxDB.Password, "Password for the InfluxDB database user")
	flag.IntVar(&configuration.TimeInterval, "timeinterval", configuration.TimeInterval, "Time interval when the data should be pulled of the RKI API in seconds (default: 86400)")
	flag.StringVar(&configuration.FederalState, "state", configuration.FederalState, "Option to pull the data only from one German Federal State (e.g. Bayern) (default: all)")
	flag.StringVar(&configuration.Logging.Dir, "logdir", configuration.Logging.Dir, "Directory for the log files (default: /var/logs/corona-dashboard)")
	flag.BoolVar(&configuration.Logging.Debug, "debug", configuration.Logging.Debug, "Option to run the microservice in debugging mode")

    flag.Parse()

    Log.Logger.Info("Starting ...")

	a := App{}
	a.Initialize(configuration.InfluxDB.URL, configuration.InfluxDB.Name, configuration.InfluxDB.User, configuration.InfluxDB.Password)

    // Setup signal catching
    sigs := make(chan os.Signal, 1)

    // catch all signals since not explicitly listing
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

    // Method invoked upon seeing signal
    go func() {
      s := <-sigs
      Log.Logger.Infow("Stopping the service.", "reason", s)
      os.Exit(1)
    }()

    // Infinite loop
    for {
	  // Run the microservice
	  Log.Logger.Info("Starting data refresh ... ")
      a.Run()
	  Log.Logger.Info("Finshed data refresh.")
	  // Wait the provided time to befor running againg
	  d := time.Second * time.Duration(configuration.TimeInterval)
	  Log.Logger.Infow("Waiting for the next run.", "duration:", d)
      time.Sleep(d)
    }

}
