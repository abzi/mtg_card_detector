package com.mtgdetector.network

import com.mtgdetector.models.*
import retrofit2.Response
import retrofit2.http.*

interface ApiService {
    @POST("auth/anonymous")
    suspend fun authenticateAnonymous(@Body request: AuthRequest): Response<AuthResponse>

    @POST("cards/scan")
    suspend fun scanCard(@Body request: ScanRequest): Response<ScanResponse>

    @POST("cards/scan/bulk")
    suspend fun scanBulk(@Body request: BulkScanRequest): Response<BulkScanResponse>

    @GET("inventory")
    suspend fun getInventory(): Response<InventoryResponse>

    @GET("cards")
    suspend fun getCard(@Query("id") cardId: String): Response<Card>
}
