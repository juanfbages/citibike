package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/juanfbages/utilities"
)

type Stations struct {
	Station []Station `json:"features"`
}

type Station struct {
	StationProperties StationProperties `json:"properties"`
	Geometry          Geometry          `json:"geometry"`
}

type Geometry struct {
	Type        string    `json:"type"`
	Coordinates []float64 `json:"coordinates"`
}

type StationProperties struct {
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

func OSMURLBuilder(lat float64, lon float64) string {
	baseUrl := "https://nominatim.openstreetmap.org/search?format=json&addressdetails=1&q="
	params := fmt.Sprintf("%f%%2C%f", lat, lon)

	return baseUrl + params
}

func main() {

	statusBody := utilities.CurlURL("https://layer.bicyclesharing.net/map/v1/nyc/stations")

	var stationStatus Stations
	jsonStatusErr := json.Unmarshal(statusBody, &stationStatus)
	if jsonStatusErr != nil {
		log.Fatal(jsonStatusErr)
	}

	pointsByLocation := make(map[string][]StationProperties)

	for _, station := range stationStatus.Station {
		lon := station.Geometry.Coordinates[0]
		lat := station.Geometry.Coordinates[1]

		osmURL := OSMURLBuilder(lat, lon)
		osmBody := utilities.CurlURL(osmURL)

		var osmPlaces []Place
		jsonOSMErr := json.Unmarshal(osmBody, &osmPlaces)
		if jsonOSMErr != nil {
			log.Fatal(jsonOSMErr)
		}

		neighbourhood := osmPlaces[0].Address.Neighbourhood
		pointsByLocation[neighbourhood] = append(pointsByLocation[neighbourhood], station.StationProperties)

		if neighbourhood == "Upper West Side" {
			fmt.Printf(
				"\nStation: %v\n\tNeighbourhood: %v\n\tBikeAction: %v\n\tBikePoints: %v\n",
				station.StationProperties.Name,
				osmPlaces[0].Address.Neighbourhood,
				station.StationProperties.BikeAngelsAction,
				station.StationProperties.BikeAngelsPoints)
		}
	}
}
