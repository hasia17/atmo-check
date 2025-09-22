package com.example.atmo_check_mobile

import android.annotation.SuppressLint
import android.os.Bundle
import androidx.appcompat.app.AppCompatActivity

import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import kotlinx.coroutines.withContext

import android.widget.Spinner
import android.widget.TextView
import android.widget.ArrayAdapter
import android.view.View
import android.widget.AdapterView

class MainActivity : AppCompatActivity() {

    private val apiClient = AirQualityApiClient()
    private var currentData: AggregatedVoivodeshipData? = null

    // UI elements
    private lateinit var voivodeshipSpinner: Spinner
    private lateinit var parametersSpinner: Spinner
    private lateinit var parametersLabel: TextView
    private lateinit var detailsText: TextView

    @SuppressLint("SetTextI18n")
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)

        initializeUI()
        setupVoivodeshipSpinner()
    }

    private fun initializeUI() {
        voivodeshipSpinner = findViewById(R.id.voivodeshipSpinner)
        parametersSpinner = findViewById(R.id.parametersSpinner)
        parametersLabel = findViewById(R.id.parametersLabel)
        detailsText = findViewById(R.id.detailsText)
    }

    private fun setupVoivodeshipSpinner() {
        val voivodeships = listOf(
            "DOLNOSLASKIE", "KUJAWSKO_POMORSKIE", "LUBELSKIE",
            "LUBUSKIE", "LODZKIE", "MALOPOLSKIE", "MAZOWIECKIE",
            "OPOLSKIE", "PODKARPACKIE", "PODLASKIE", "POMORSKIE",
            "SLASKIE", "SWIETOKRZYSKIE", "WARMINSKO_MAZURSKIE",
            "WIELKOPOLSKIE", "ZACHODNIOPOMORSKIE"
        )

        val adapter = ArrayAdapter(this, android.R.layout.simple_spinner_item, voivodeships)
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
        voivodeshipSpinner.adapter = adapter

        voivodeshipSpinner.onItemSelectedListener = object : AdapterView.OnItemSelectedListener {
            override fun onItemSelected(parent: AdapterView<*>, view: View?, position: Int, id: Long) {
                val selectedVoivodeship = voivodeships[position]
                loadDataForVoivodeship(selectedVoivodeship)
            }

            override fun onNothingSelected(parent: AdapterView<*>) {
            }
        }
    }

    private fun loadDataForVoivodeship(voivodeship: String) {
        CoroutineScope(Dispatchers.IO).launch {
            try {
                println("Loading data for: $voivodeship")
                val data = apiClient.getVoivodeshipData(voivodeship)

                withContext(Dispatchers.Main) {
                    if (data != null) {
                        currentData = data
                        setupParametersSpinner(data.parameters)
                        showParametersUI()
                    } else {
                        hideParametersUI()
                        detailsText.text = "Failed to load data for $voivodeship"
                        detailsText.visibility = View.VISIBLE
                    }
                }
            } catch (e: Exception) {
                println("Error loading data: ${e.message}")
                withContext(Dispatchers.Main) {
                    hideParametersUI()
                    detailsText.text = "Error: ${e.message}"
                    detailsText.visibility = View.VISIBLE
                }
            }
        }
    }

    private fun showParametersUI() {
        parametersLabel.visibility = View.VISIBLE
        parametersSpinner.visibility = View.VISIBLE
    }

    private fun hideParametersUI() {
        parametersLabel.visibility = View.GONE
        parametersSpinner.visibility = View.GONE
        detailsText.visibility = View.GONE
    }

    private fun setupParametersSpinner(parameters: List<Parameter>) {
        val parameterNames = parameters.map { "${it.name} (${it.id})" }

        val adapter = ArrayAdapter(this, android.R.layout.simple_spinner_item, parameterNames)
        adapter.setDropDownViewResource(android.R.layout.simple_spinner_dropdown_item)
        parametersSpinner.adapter = adapter

        parametersSpinner.onItemSelectedListener = object : AdapterView.OnItemSelectedListener {
            override fun onItemSelected(
                parent: AdapterView<*>,
                view: View?,
                position: Int,
                id: Long
            ) {
                val selectedParameter = parameters[position]
                showParameterDetails(selectedParameter)
            }

            override fun onNothingSelected(parent: AdapterView<*>) {
                detailsText.visibility = View.GONE
            }
        }
    }

        private fun showParameterDetails(parameter: Parameter) {
            val details = buildString {
                appendLine("Parameter: ${parameter.name}")
                appendLine("ID: ${parameter.id}")
                appendLine("Unit: ${parameter.unit ?: "N/A"}")
                appendLine("Source: ${parameter.source ?: "N/A"}")
                appendLine()
                appendLine("Average: ${parameter.averageValue ?: "N/A"}")
                appendLine("Min: ${parameter.minValue ?: "N/A"}")
                appendLine("Max: ${parameter.maxValue ?: "N/A"}")
                appendLine("Measurements: ${parameter.measurementCount ?: "N/A"}")
                appendLine()
                appendLine("Latest value: ${parameter.latestValue ?: "N/A"}")
                appendLine("Latest time: ${parameter.latestTimestamp ?: "N/A"}")
            }

            detailsText.text = details
            detailsText.visibility = View.VISIBLE
        }
    }


//    @SuppressLint("SetTextI18n")
//    private suspend fun loadAirQualityData(textView: TextView) {
//        try {
//            println("Trying to load air quality data...")
//            val data = apiClient.getVoivodeshipData("DOLNOSLASKIE")
//
//            withContext(Dispatchers.Main) {
//                if (data != null) {
//                    val firstParam = data.parameters.firstOrNull()
//                    textView.text = "Voivodeship: ${data.voivodeship}\n" +
//                            "Parameters count: ${data.parameters.size}\n" +
//                            "First parameter: ${firstParam?.name ?: "None"}\n" +
//                            "Average value: ${firstParam?.averageValue ?: "N/A"}"
//                } else {
//                    textView.text = "Data is null - check Logcat for details"
//                }
//            }
//        } catch (e: Exception) {
//            println("Error in MainActivity: ${e.message}")
//            e.printStackTrace()
//
//            withContext(Dispatchers.Main) {
//                textView.text = "Exception: ${e.message}"
//            }
//        }

