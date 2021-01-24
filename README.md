# RKI Corona Dashboard for Grafana
This project started out as a joke of a colleague of mine. He asked another colleague why we do not have a Corona Dashboard for our Grafana.

This is my implementation of a Corona Dashboard for [Grafana](https://grafana.com/) and a microservice based on the offical numbers of the German Federal Government.
The numbers are provided by the [Robert Koch Institute (RKI)](https://en.wikipedia.org/wiki/Robert_Koch_Institute). The RKI is a German Federal Government agency and research institute responsible for disease control and prevention.

## Function
The microservice to pull the data from the RKI API is written in [Go](https://golang.org/) and parses the data to get the relevant numbers.
After parsing the data the microservice stores them in an [InfluxDB](https://www.influxdata.com/products/influxdb/) database. InfluxDB is a time series database.


## Example


## Prerequisites
You need to have a running instance of Grafana and InfluxDB. Also it is recommend to have a server where you can deploy the microservice, but it is also possible to start the microservice manually on a local machine.

## Installation
Download the prebuilt binary packages from the [release page](https://github.com/4ndyZ/Corona-Dashboard/releases) and install them on your server.

## Configuration
After installing the package configure the microservice. The configuration file is located under `/etc/corona-dashboard/config.yml`.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[GPL-3.0](https://github.com/4ndyZ/Corona-Dashboard/blob/master/LICENSE)
