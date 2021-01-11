package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/IncSW/geoip2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Params struct {
	Address string `json:"address"`
}

type Location struct {
	Timezone  string
	Latitude  float64
	Longitude float64
}

type Details struct {
	Continent string
	Country   string
	City      string
	Location  Location
}

var args = struct {
	bindAddr                 **net.TCPAddr
	appUrlPath               *string
	metricsUrlPath           *string
	geolite2DatabasePath     *string
	geolite2DatabaseFilename *string
	geolite2Locale           *string
	metrics                  *bool
	debug                    *bool
}{
	kingpin.Flag("bind", "Bind address:port").
		Envar("LISTEN_ADDRESS_PORT").
		Default("0.0.0.0:8080").
		TCP(),
	kingpin.Flag("app-url-path", "App URL path").
		Envar("APP_URL_PATH").
		Default("/location").
		String(),
	kingpin.Flag("metrics-url-path", "Metrics URL path").
		Envar("METRICS_URL_PATH").
		Default("/metrics").
		String(),
	kingpin.Flag("geolite2-database-path", "GeoLite2 mmdb database file path.").
		Envar("GL2_DATABASE_PATH").
		Required().
		String(),
	kingpin.Flag("geolite2-database-filename", "GeoLite2 mmdb database filename.").
		Envar("GL2_DATABASE_FILENAME").
		Required().
		String(),
	kingpin.Flag("geolite2-locale", "GeoLite2 database file locale. Can be 'de', 'en'(default), 'es', 'fr', 'ja', 'pt-BR', 'ru”', and 'zh-CN").
		Envar("GL2_LOCALE").
		Default("en").
		String(),
	kingpin.Flag("enable-metrics", "Enable Metrics").
		Envar("METRICS").
		Bool(),
	kingpin.Flag("debug", "Debug Mode").
		Envar("DEBUG").
		Bool(),
}

func main() {
	// Parse and load the environment varialbe using Kingpin library
	kingpin.Parse()
	// Add the handler function for the primary app function of recieving and IP addresa and returning location data in JSON format
	http.HandleFunc(*args.appUrlPath, checkIP)
	// If enabled, add basic prometheus metrics handling on '/metrics'
	if *args.metrics {
		http.Handle(*args.metricsUrlPath, promhttp.Handler())
	}
	if *args.debug {
		fmt.Printf("Starting Server: listening on %s", *args.bindAddr)
	}
	// Start the http listener.
	log.Fatal(http.ListenAndServe((*args.bindAddr).String(), nil))
}

func checkIP(w http.ResponseWriter, r *http.Request) {
	var params Params

	// Take the HTTP POST Body and store the data in the params struct
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// Only act on the request if the method id POST
	switch r.Method {
	case "POST":
		w.WriteHeader(http.StatusCreated)
		w.Write(post(params.Address))
	default:
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"message": "Requested method not supported"}`))
	}
}

func post(address string) []byte {
	// Load the GeoLite2-City database from the mmdb file path configured
	reader, err := geoip2.NewCityReaderFromFile(*args.geolite2DatabasePath + "/" + *args.geolite2DatabaseFilename)
	if err != nil {
		panic(err)
	}

	// Call the database and recieve the appropriate record
	record, err := reader.Lookup(net.ParseIP(address))
	if err != nil {
		panic(err)
	}

	// Place the return data into the struct
	var details = Details{
		Continent: record.Continent.Names[*args.geolite2Locale],
		Country:   record.Country.Names[*args.geolite2Locale],
		City:      record.City.Names[*args.geolite2Locale],
		Location: Location{
			Timezone:  record.Location.TimeZone,
			Latitude:  record.Location.Latitude,
			Longitude: record.Location.Longitude,
		},
	}

	// Marshal the stuct into JSON to be returned the requester
	ret, err := json.Marshal(details)
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}

	return ret
}
