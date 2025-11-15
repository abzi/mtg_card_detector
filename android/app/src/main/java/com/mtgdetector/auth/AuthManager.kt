package com.mtgdetector.auth

import android.content.Context
import android.content.SharedPreferences
import androidx.security.crypto.EncryptedSharedPreferences
import androidx.security.crypto.MasterKey
import com.mtgdetector.models.AuthRequest
import com.mtgdetector.network.RetrofitClient
import java.util.UUID

class AuthManager(context: Context) {
    private val masterKey = MasterKey.Builder(context)
        .setKeyScheme(MasterKey.KeyScheme.AES256_GCM)
        .build()

    private val sharedPreferences: SharedPreferences = EncryptedSharedPreferences.create(
        context,
        "mtg_detector_secure_prefs",
        masterKey,
        EncryptedSharedPreferences.PrefKeyEncryptionScheme.AES256_SIV,
        EncryptedSharedPreferences.PrefValueEncryptionScheme.AES256_GCM
    )

    companion object {
        private const val KEY_USER_ID = "user_id"
        private const val KEY_AUTH_TOKEN = "auth_token"
        private const val KEY_DEVICE_ID = "device_id"
    }

    fun getDeviceId(): String {
        var deviceId = sharedPreferences.getString(KEY_DEVICE_ID, null)
        if (deviceId == null) {
            deviceId = UUID.randomUUID().toString()
            sharedPreferences.edit().putString(KEY_DEVICE_ID, deviceId).apply()
        }
        return deviceId
    }

    fun isAuthenticated(): Boolean {
        return getAuthToken() != null && getUserId() != null
    }

    fun getUserId(): String? {
        return sharedPreferences.getString(KEY_USER_ID, null)
    }

    fun getAuthToken(): String? {
        return sharedPreferences.getString(KEY_AUTH_TOKEN, null)
    }

    fun saveAuthData(userId: String, token: String) {
        sharedPreferences.edit().apply {
            putString(KEY_USER_ID, userId)
            putString(KEY_AUTH_TOKEN, token)
            apply()
        }
        RetrofitClient.setAuthToken(token)
    }

    fun clearAuthData() {
        sharedPreferences.edit().apply {
            remove(KEY_USER_ID)
            remove(KEY_AUTH_TOKEN)
            apply()
        }
        RetrofitClient.setAuthToken(null)
    }

    suspend fun authenticate(): Result<Boolean> {
        return try {
            val deviceId = getDeviceId()
            val response = RetrofitClient.apiService.authenticateAnonymous(
                AuthRequest(deviceId)
            )

            if (response.isSuccessful && response.body() != null) {
                val authResponse = response.body()!!
                saveAuthData(authResponse.userId, authResponse.token)
                Result.success(true)
            } else {
                Result.failure(Exception("Authentication failed: ${response.code()}"))
            }
        } catch (e: Exception) {
            Result.failure(e)
        }
    }
}
