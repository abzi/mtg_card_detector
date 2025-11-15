package com.mtgdetector.models

import com.google.gson.annotations.SerializedName

data class AuthRequest(
    @SerializedName("device_id") val deviceId: String
)

data class AuthResponse(
    @SerializedName("user_id") val userId: String,
    @SerializedName("token") val token: String
)

data class Card(
    @SerializedName("id") val id: String,
    @SerializedName("scryfall_id") val scryfallId: String?,
    @SerializedName("name") val name: String,
    @SerializedName("set_code") val setCode: String,
    @SerializedName("collector_number") val collectorNumber: String,
    @SerializedName("image_uri") val imageUri: String?,
    @SerializedName("oracle_text") val oracleText: String?,
    @SerializedName("type_line") val typeLine: String?,
    @SerializedName("mana_cost") val manaCost: String?,
    @SerializedName("rarity") val rarity: String?,
    @SerializedName("created_at") val createdAt: String
)

data class ScanRequest(
    @SerializedName("card_name") val cardName: String? = null,
    @SerializedName("set_code") val setCode: String? = null,
    @SerializedName("collector_number") val collectorNumber: String? = null,
    @SerializedName("barcode") val barcode: String? = null
)

data class BulkScanRequest(
    @SerializedName("scans") val scans: List<ScanRequest>
)

data class ScanResponse(
    @SerializedName("success") val success: Boolean,
    @SerializedName("card") val card: Card?,
    @SerializedName("error") val error: String?
)

data class BulkScanResponse(
    @SerializedName("session_id") val sessionId: Int,
    @SerializedName("total_scanned") val totalScanned: Int,
    @SerializedName("successful_scans") val successfulScans: Int,
    @SerializedName("failed_scans") val failedScans: Int,
    @SerializedName("results") val results: List<ScanResponse>
)

data class InventoryItem(
    @SerializedName("id") val id: Int,
    @SerializedName("user_id") val userId: String,
    @SerializedName("card_id") val cardId: String,
    @SerializedName("quantity") val quantity: Int,
    @SerializedName("added_at") val addedAt: String,
    @SerializedName("card") val card: Card?
)

data class InventoryResponse(
    @SerializedName("inventory") val inventory: List<InventoryItem>,
    @SerializedName("count") val count: Int
)

data class ErrorResponse(
    @SerializedName("error") val error: String,
    @SerializedName("message") val message: String?
)
