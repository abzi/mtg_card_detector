package com.mtgdetector.ui

import android.Manifest
import android.content.Intent
import android.content.pm.PackageManager
import android.os.Bundle
import android.view.View
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import androidx.lifecycle.lifecycleScope
import com.mtgdetector.MTGDetectorApp
import com.mtgdetector.databinding.ActivityMainBinding
import kotlinx.coroutines.launch

class MainActivity : AppCompatActivity() {
    private lateinit var binding: ActivityMainBinding
    private lateinit var authManager: com.mtgdetector.auth.AuthManager

    private val requestPermissionLauncher = registerForActivityResult(
        ActivityResultContracts.RequestPermission()
    ) { isGranted ->
        if (isGranted) {
            startScanActivity(false)
        } else {
            Toast.makeText(this, "Camera permission required", Toast.LENGTH_LONG).show()
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        authManager = (application as MTGDetectorApp).authManager

        // Authenticate on first launch
        if (!authManager.isAuthenticated()) {
            authenticateUser()
        }

        binding.btnSingleScan.setOnClickListener {
            if (checkCameraPermission()) {
                startScanActivity(false)
            } else {
                requestCameraPermission()
            }
        }

        binding.btnBulkScan.setOnClickListener {
            if (checkCameraPermission()) {
                startScanActivity(true)
            } else {
                requestCameraPermission()
            }
        }

        binding.btnViewInventory.setOnClickListener {
            startActivity(Intent(this, InventoryActivity::class.java))
        }
    }

    private fun authenticateUser() {
        binding.progressBar.visibility = View.VISIBLE
        lifecycleScope.launch {
            authManager.authenticate()
                .onSuccess {
                    binding.progressBar.visibility = View.GONE
                }
                .onFailure { error ->
                    binding.progressBar.visibility = View.GONE
                    Toast.makeText(
                        this@MainActivity,
                        "Authentication failed: ${error.message}",
                        Toast.LENGTH_LONG
                    ).show()
                }
        }
    }

    private fun checkCameraPermission(): Boolean {
        return ContextCompat.checkSelfPermission(
            this,
            Manifest.permission.CAMERA
        ) == PackageManager.PERMISSION_GRANTED
    }

    private fun requestCameraPermission() {
        requestPermissionLauncher.launch(Manifest.permission.CAMERA)
    }

    private fun startScanActivity(bulkMode: Boolean) {
        val intent = Intent(this, ScanActivity::class.java).apply {
            putExtra("BULK_MODE", bulkMode)
        }
        startActivity(intent)
    }
}
