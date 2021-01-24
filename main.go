package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	 "io/ioutil"
)

func main() {


	body, _ := ioutil.ReadFile("config/config.yml")

	configuration := Configuration{}

	_ = yaml.Unmarshal([]byte(body), &configuration)
	fmt.Printf("--- t:\n%v\n\n", configuration)

	flag.StringVar(&configuration.InfluxDB.URL, "dburl", configuration.InfluxDB.URL, "Database connection URL including the port (e.g. https://myinfluxdb-server.tld:8086)")
	flag.StringVar(&configuration.InfluxDB.User, "dbuser", configuration.InfluxDB.User, "Username of the database user")
	flag.StringVar(&configuration.InfluxDB.Password, "dbpassword", configuration.InfluxDB.Password, "Password for the database user")
	flag.IntVar(&configuration.TimeInterval, "timeinterval", configuration.TimeInterval, "Time interval when the data should be pulled of the RKI API in seconds (default: 86400)")
	flag.StringVar(&configuration.FederalState, "state", configuration.FederalState, "")
	flag.StringVar(&configuration.Logging.Dir, "logdir", configuration.Logging.Dir, "")
	flag.BoolVar(&configuration.Logging.Debug, "debug", configuration.Logging.Debug, "")

	a := App{}

	a.Initialize(configuration.InfluxDB.URL, configuration.InfluxDB.User, configuration.InfluxDB.Password)

	a.Run()
}
