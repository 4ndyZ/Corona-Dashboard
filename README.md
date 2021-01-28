# RKI Corona Dashboard for Grafana
This project started out as a joke of a colleague of mine. He asked another colleague why we do not have a Corona Dashboard for our Grafana.

This is my implementation of a Corona Dashboard for [Grafana](https://grafana.com/) and a microservice based on the offical numbers of the German Federal Government.
The numbers are provided by the [Robert Koch Institute (RKI)](https://en.wikipedia.org/wiki/Robert_Koch_Institute). The RKI is a German Federal Government agency and research institute responsible for disease control and prevention.

## Function
The microservice to pull the data from the RKI API is written in [Go](https://golang.org/) and parses the data to get the relevant numbers.
After parsing the data the microservice stores them in an [InfluxDB](https://www.influxdata.com/products/influxdb/) database. InfluxDB is a time series database.

## Example
![Example Grafana Dashboard](https://raw.githubusercontent.com/4ndyZ/Corona-Dashboard/main/.github/example.png)

## Prerequisites
You need to have a running instance of Grafana and InfluxDB. Also it is recommend to have a server where you can deploy the microservice, but it is also possible to start the microservice manually on a local machine.

## Installation and configuration
Download the prebuilt binary packages from the [release page](https://github.com/4ndyZ/Corona-Dashboard/releases) and install them on your server.

### Installation
#### Linux
###### DEB Package
If you are running a Debian-based Linux Distribution choose the `.deb` Package for your operating system architecture and download it. You are able to use curl to download the package.

Now you are able to install the package using APT.
`sudo apt install -f Corona-Dashboard-vX.X-.linux.XXXX.deb`

After installing the package configure the microservice. The configuration file is located under `/etc/corona-dashboard/config.yml`.

No you you are able to enable the Systemd service using `systemctl`.
`sudo systemctl enable corona-dashboard`

And start the service also using `systemctl`.
`sudo systemctl start corona-dashboard`

###### RPM Package
When running a RHEL-based Linux Distribution choose the `.rpm` package for your operating system architecture and download it.

Now you are able to install the package.
`sudo rpm -i Corona-Dashboard-vX.X-.linux.XXXX.rpm`

After installing the package configure the microservice. The configuration file is located under `/etc/corona-dashboard/config.yml`.

No you you are able to enable the Systemd service using `systemctl`.
`sudo systemctl enable corona-dashboard`

And start the service also using `systemctl`.
`sudo systemctl start corona-dashboard`

#### Windows/Other
If you plan to run the microservice on Windows or another OS the whole process is a bit more complicated because there is no installation package avaible only prebuilt binaries.

Download the prebuilt binary for your operating system.

Exctract the prebuilt binary and change the configuration file located under `config/config.conf`.

After successful changing the configuration file you are able to run the prebuilt binary.

### Configuration
The microservice tries to access the configuration file located under `/etc/corona-dashboard/config.conf`. It the configuration file is not accessable or found the service will fallback to the local file located unter `config/config.conf`.

### Logging
The microservice while try to put the log file in the `/var/log/corona-dashboard` folder. If the service is not able to access or find that folder, the logging file gets created in the local folder `logs`.

If you want to enable debug messages please change the configuration file  or run the microservice with the commandline parameter `-debug`.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

## License
[GPL-3.0](https://github.com/4ndyZ/Corona-Dashboard/blob/master/LICENSE)
