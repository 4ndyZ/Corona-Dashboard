package main

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var Log Logger = Logger{}

func main() {
	// Create configuration object
	configuration := Configuration{}
	// Initialize logging
	logfolder := "/var/log/corona-dashboard"
	_, err := os.Stat(logfolder)
	if os.IsNotExist(err) || os.IsPermission(err) {
		logfolder = "log"
	}
	Log.Initialize(strings.Join([]string{logfolder, "/log.txt"}, ""))

	config := ""
	config1 := "/etc/corona-dashboard/config.conf"
	config2 := "config/config.conf"
	// Error checking for config files
	_, err1 := os.Stat(config1)
	_, err2 := os.Stat(config2)
	if err1 == nil {
		config = config1
	} else if err2 == nil {
		config = config2
	} else if os.IsNotExist(err1) && os.IsNotExist(err2) {
		Log.Logger.Info().Msg("No configuration file found. Using Commandline parameter.")
	} else if !os.IsNotExist(err1) && os.IsPermission(err1) {
		Log.Logger.Warn().Str("path", config1).Msg("Unable to use configuration file. No permission to access the configuration file.")
	} else if !os.IsNotExist(err2) && os.IsPermission(err2) {
		Log.Logger.Warn().Str("path", config2).Msg("Unable to use configuration file. No permission to access the configuration file.")
	} else if err1 != nil {
		Log.Logger.Warn().Str("error", err1.Error()).Msg("Error while accessing the configuration file.")
	} else if err2 != nil {
		Log.Logger.Warn().Str("error", err2.Error()).Msg("Error while accessing the configuration file.")
	}
	// Try to parse the configuration file if exists
	if config == "" {
		body, err := ioutil.ReadFile(config)
		if err != nil {
			Log.Logger.Warn().Str("error", err.Error()).Msg("Error while reading the configuration file.")
		}
		err = yaml.Unmarshal([]byte(body), &configuration)
		if err != nil {
			Log.Logger.Warn().Str("error", err.Error()).Msg("Error while parsing the configuration file.")
		}
		// Set default configuration parameter
	} else {
		configuration.InfluxDB.Name = "corona"
		configuration.InfluxDB.User = "corona"
		configuration.TimeInterval = 86400
		configuration.SingleRun = false
		configuration.FederalState = "all"
		configuration.Logging.Debug = false
	}
	// Commandline flags
	flag.StringVar(&configuration.InfluxDB.URL, "dburl", configuration.InfluxDB.URL, "Database connection URL including the port (e.g. https://myinfluxdb-server.tld:8086)")
	flag.StringVar(&configuration.InfluxDB.Name, "dbname", configuration.InfluxDB.Name, "Name of the InfluxDB database")
	flag.StringVar(&configuration.InfluxDB.User, "dbuser", configuration.InfluxDB.User, "Username of the InfluxDB database user")
	flag.StringVar(&configuration.InfluxDB.Password, "dbpassword", configuration.InfluxDB.Password, "Password for the InfluxDB database user")
	flag.IntVar(&configuration.TimeInterval, "timeinterval", configuration.TimeInterval, "Time interval when the data should be pulled of the RKI API in seconds (default: 86400)")
	flag.BoolVar(&configuration.SingleRun, "singlerun", configuration.SingleRun, "Option to run the microservice only one time and then stop afterwards. Option timeinterval will be ignored!")
	flag.StringVar(&configuration.FederalState, "state", configuration.FederalState, "Option to pull the data only from one German Federal State (e.g. Bayern) (default: all)")
	flag.BoolVar(&configuration.Logging.Debug, "debug", configuration.Logging.Debug, "Option to run the microservice in debugging mode")
	flag.Parse()
	// Check if debug log should be enabeled
	if configuration.Logging.Debug {
		Log.EnableDebug(true)
	}

	Log.Logger.Info().Msg("Starting ...")
	// Create app worker
	a := App{}
	a.Initialize(configuration.InfluxDB.URL, configuration.InfluxDB.Name, configuration.InfluxDB.User, configuration.InfluxDB.Password)
	// Setup signal catching
	sigs := make(chan os.Signal, 1)
	// Catch all signals since not explicitly listing
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGHUP)
	// Method invoked upon seeing signal
	go func() {
		s := <-sigs
		Log.Logger.Info().Str("reason", s.String()).Msg("Stopping the service.")
		Log.Rotate()
		os.Exit(1)
	}()
	// Infinite loop
	for {
		// Run the microservice
		Log.Logger.Info().Msg("Starting data refresh ... ")
		a.Run(configuration.FederalState)
		Log.Logger.Info().Msg("Finshed data refresh.")
		// Wait the provided time to befor running againg
		d := time.Second * time.Duration(configuration.TimeInterval)
		Log.Logger.Info().Interface("duration", d).Msg("Waiting for the next run.")
		time.Sleep(d)
	}

}
