package main

import (
	"atmo-check-openaq/service"
	"fmt"
	"log"
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

	s := service.New(apiKey)
	data, err := s.FetchData("Los Angeles")
	if err != nil {
		log.Fatalf("Error fetching data: %v", err)
	}
	for _, result := range data.Results {
		for _, sensor := range result.Sensors {
			fmt.Printf("Locality: %s, Name: %s, Id: %d, Name: %s\n",
				result.Locality, result.Name,
				sensor.Id, sensor.Name)
		}
	}
}
