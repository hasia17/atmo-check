package com.example.atmo_check_mobile

import android.annotation.SuppressLint
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity
import android.widget.TextView

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

class MainActivity : AppCompatActivity() {

    private val apiClient = AirQualityApiClient()

    @SuppressLint("SetTextI18n")
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)

        val textView = TextView(this).apply {
            text = "Loading data..."
            textSize = 24f
        }
        setContentView(textView)

        CoroutineScope(Dispatchers.IO).launch {
            loadAirQualityData(textView)
        }
    }

    @SuppressLint("SetTextI18n")
    private suspend fun loadAirQualityData(textView: TextView) {
        try {
            println("Trying to load air quality data...")
            val data = apiClient.getVoivodeshipData("DOLNOSLASKIE")

            withContext(Dispatchers.Main) {
                if (data != null) {
                    val firstParam = data.parameters.firstOrNull()
                    textView.text = "Voivodeship: ${data.voivodeship}\n" +
                            "Parameters count: ${data.parameters.size}\n" +
                            "First parameter: ${firstParam?.name ?: "None"}\n" +
                            "Average value: ${firstParam?.averageValue ?: "N/A"}"
                } else {
                    textView.text = "Data is null - check Logcat for details"
                }
            }
        } catch (e: Exception) {
            println("Error in MainActivity: ${e.message}")
            e.printStackTrace()

            withContext(Dispatchers.Main) {
                textView.text = "Exception: ${e.message}"
            }
        }
    }
}
