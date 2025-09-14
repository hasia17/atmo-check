package com.example.atmo_check_mobile

import kotlinx.serialization.Serializable

@Serializable
data class AggregatedVoivodeshipData(
    val voivodeship: String,
    val parameters: List<Parameter>
)

@Serializable
data class Parameter(
    val id: String?,
    val name: String?,
    val unit: String?,
    val description: String?,
    val source: String?,
    val averageValue: Double?,
    val minValue: Double?,
    val maxValue: Double?,
    val measurementCount: Int?,
    val latestValue: Double?,
    val latestTimestamp: String?
)