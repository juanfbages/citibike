package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Stations struct {
	Station []Station `json:"features"`
}

type Station struct {
	StationProperties Properties `json:"properties"`
	Geometry          Geometry   `json:"geometry"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type Properties struct {
	StationID        string `json:"station_id"`
	Name             string `json:"name"`
	BikeAngelsAction string `json:"bike_angels_action"`
	BikeAngelsPoints int    `json:"bike_angels_points"`
	BikeAngelsDigits int    `json:"bike_angels_digits"`
}

type Place struct {
	PlaceID int     `json:"place_id"`
	Lat     string  `json:"lat"`
	Lon     string  `json:"lon"`
	Address Address `json:"address"`
}

type Address struct {
	Road          string `json:"road"`
	Neighbourhood string `json:"neighbourhood"`
	Suburb        string `json:"suburb"`
	County        string `json:"county"`
	City          string `json:"city"`
	State         string `json:"state"`
	Postcode      string `json:"postcode"`
}

func CurlURL(url string) (body []byte) {

	Client := http.Client{Timeout: time.Second * 5}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("User-Agent", "curler")

	res, getErr := Client.Do(req)
	if getErr != nil {
		log.Fatal(getErr)
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Fatal(readErr)
	}

	return body

}

func OSMURLBuilder(lat float64, lon float64) string {
	baseUrl := "https://nominatim.openstreetmap.org/search?format=json&addressdetails=1&q="
	params := fmt.Sprintf("%f%%2C%f", lat, lon)

	return baseUrl + params
}

func main() {

	statusBody := CurlURL("https://layer.bicyclesharing.net/map/v1/nyc/stations")

	var stationStatus Stations
	jsonStatusErr := json.Unmarshal(statusBody, &stationStatus)
	if jsonStatusErr != nil {
		log.Fatal(jsonStatusErr)
	}

	for _, station := range stationStatus.Station {
		lon := station.Geometry.Coordinates[0]
		lat := station.Geometry.Coordinates[1]

		osmURL := OSMURLBuilder(lat, lon)
		osmBody := CurlURL(osmURL)

		var osmPlaces []Place
		jsonOSMErr := json.Unmarshal(osmBody, &osmPlaces)
		if jsonOSMErr != nil {
			log.Fatal(jsonOSMErr)
		}

		fmt.Printf(
			"Station: %v, Neighbourhood: %v, Zipcode: %v\n",
			station.StationProperties.Name,
			osmPlaces[0].Address.Neighbourhood,
			osmPlaces[0].Address.Postcode)

	}
}
