package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/IncSW/geoip2"
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

	http.HandleFunc("/", checkIP)
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

func checkIP(w http.ResponseWriter, r *http.Request) {
	var params Params

	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

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
	reader, err := geoip2.NewCityReaderFromFile("/tmp/test/GeoLite2-City_20210105/GeoLite2-City.mmdb")
	if err != nil {
		panic(err)
	}

	record, err := reader.Lookup(net.ParseIP(address))
	if err != nil {
		panic(err)
	}

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

	ret, err := json.Marshal(details)
	if err != nil {
		fmt.Println(err)
		return []byte(err.Error())
	}

	return ret
}
