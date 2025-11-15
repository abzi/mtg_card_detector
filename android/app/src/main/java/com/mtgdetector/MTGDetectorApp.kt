package com.mtgdetector

import android.app.Application
import com.mtgdetector.auth.AuthManager
import com.mtgdetector.network.RetrofitClient

class MTGDetectorApp : Application() {
    lateinit var authManager: AuthManager
        private set

    override fun onCreate() {
        super.onCreate()
        authManager = AuthManager(this)

        // Set auth token if available
        authManager.getAuthToken()?.let {
            RetrofitClient.setAuthToken(it)
        }
    }
}
