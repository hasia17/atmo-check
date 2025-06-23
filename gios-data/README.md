# gios-data

**gios-data** is a microservice in the _atmo-check_ ecosystem, responsible for fetching and providing air quality data from the public GIOS API. The service is built with Spring Boot, uses MongoDB for storage, and is container-ready with Docker.

## Features

- Fetches air quality data (stations, measurements) from the GIOS API.
- Exposes REST endpoints for data access.
- Persists data in MongoDB.
- Easily integrable with other _atmo-check_ services.

## Requirements

- Java 17+
- Maven 3.8+
- MongoDB (local or cloud instance)
- Docker (optional)