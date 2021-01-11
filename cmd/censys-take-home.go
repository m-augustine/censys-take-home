package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/IncSW/geoip2"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

func main() {
	// Add the handler function for the primary app function of recieving and IP addresa and returning location data in JSON format
	http.HandleFunc("/", checkIP)
	// Add basic prometheus metrics handling on '/metrics'
	http.Handle("/metrics", promhttp.Handler())
	println("Starting Server: listneing on 0.0.0.0:8080")
	// Start the http listener.
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
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
		w.Write([]byte(`{"message": "Can't find method requested"}`))
	}
}

func post(address string) []byte {
	// Load the GeoLite2-City database from the mmdb file path configured 
	reader, err := geoip2.NewCityReaderFromFile("/tmp/test/GeoLite2-City_20210105/GeoLite2-City.mmdb")
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
		Continent: record.Continent.Names["en"],
		Country:   record.Country.Names["en"],
		City:      record.City.Names["en"],
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
