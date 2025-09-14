package com.example.atmo_check_mobile

import io.ktor.client.*
import io.ktor.client.call.*
import io.ktor.client.plugins.contentnegotiation.*
import io.ktor.client.request.*
import io.ktor.serialization.kotlinx.json.*
import kotlinx.serialization.json.Json

class AirQualityApiClient {
    private val client = HttpClient {
        install(ContentNegotiation) {
            json(Json {
                ignoreUnknownKeys = true
            })
        }
    }

    suspend fun getVoivodeshipData(voivodeship: String): AggregatedVoivodeshipData? {
        return try {
            println("Trying to fetch data from: http://10.0.2.2:8082/aggregator/airQuality/$voivodeship")
            val response = client.get("http://10.0.2.2:8082/aggregator/airQuality/$voivodeship").body<AggregatedVoivodeshipData>()
            println("Data fetched successfully: $response")
            response
        } catch (e: Exception) {
            println("Error in API client: ${e.message}")
            e.printStackTrace()
            null
        }
    }
}