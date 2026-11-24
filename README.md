# Atmo-check

System do monitorowania jakości powietrza w Polsce. Zbiera pomiary z Open-Meteo i OpenAQ, łączy je i wylicza średnie dla każdego województwa, a wyniki pokazuje na mapie.

## Architektura

```
            ┌──────────────────┐
            │  atmo-check-page │  Angular + nginx  (:80)
            └────────┬─────────┘
                     │
            ┌────────▼─────────┐
            │    aggregator    │  Go  (:8082)
            └───┬──────────┬───┘
                │          │
   ┌────────────▼───┐  ┌───▼────────────┐
   │ open-meteo-data│  │  openaq-data   │
   │ Java/Spring    │  │  Go            │
   │ (:8083)        │  │  (:3001→3000)  │
   └───────┬────────┘  └───────┬────────┘
           └───────┬───────────┘
              ┌────▼────┐
              │ MongoDB │  (:27017)
              └─────────┘
```

| Katalog | Technologia | Opis |
|---|---|---|
| `aggregator` | Go 1.25 | Agreguje dane z obu źródeł i grupuje stacje według województw; dane trzyma w cache |
| `open-meteo-data` | Java 21, Spring Boot | Pobiera i zapisuje dane z Open-Meteo |
| `openaq-data` | Go 1.24 | Pobiera i zapisuje dane z OpenAQ (spec: `openaq-data/openapi.yaml`) |
| `atmo-check-page` | Angular 21 | Interfejs webowy z mapą wartości parametrów |
| `atmo-check-mobile` | Kotlin | Aplikacja mobilna (wczesny etap, nie jest w docker-compose) |

Obsługiwane parametry: PM10, PM2.5, CO, CO2, NO2, SO2, O3, CH4.

## API agregatora

- `GET /aggregatedData`: dane dla wszystkich województw
- `GET /aggregatedData/{voivodeship}`: dane dla jednego województwa

## Uruchomienie

1. Utwórz plik `.env` w katalogu głównym:

   ```env
   OPENAQ_API_KEY=<twój klucz API OpenAQ>
   MONGO_URI=mongodb://mongodb:27017
   ```

2. Uruchom wszystkie usługi:

   ```bash
   docker compose up --build
   ```

3. Otwórz `http://localhost`.

## Rozwój lokalny

```bash
# aggregator
cd aggregator && go run . && go test ./...

# openaq-data
cd openaq-data && go run .

# open-meteo-data
cd open-meteo-data && ./mvnw spring-boot:run

# atmo-check-page
cd atmo-check-page && npm install && npm start   # http://localhost:4200
```

## CI

Plik `.github/workflows/aggregator-ci.yml` buduje i testuje `aggregator` (`go test -race`) oraz sprawdza podatności przez `govulncheck`. Uruchamia się przy zmianach w `aggregator/`.
