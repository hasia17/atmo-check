package main

import (
	"fmt"
	"log"
	"openaq-data/internal"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, proceeding without it.")
	}

	apiKey := ""
	if v := os.Getenv("OPENAQ_API_KEY"); v != "" {
		apiKey = v
	}

	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <location>", os.Args[0])
	}
	location := os.Args[1]

	s := internal.NewService(apiKey)
	data, err := s.FetchData(location)
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	for _, result := range data.Results {
		for _, sensor := range result.Sensors {
			fmt.Printf(
				"Locality: %s, Name: %s, Latitude: %f, Longitude: %f, Country: %s (%s), Sensor Id: %d, Sensor Name: %s, Parameter: %s (%s)\n",
				result.Locality,
				result.Name,
				result.Coordinates.Latitude,
				result.Coordinates.Longitude,
				result.Country.Name,
				result.Country.Code,
				sensor.Id,
				sensor.Name,
				sensor.Parameter.Name,
				sensor.Parameter.Units,
			)
		}
	}
}
