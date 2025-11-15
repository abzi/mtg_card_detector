package com.mtgdetector.ui

import android.os.Bundle
import android.util.Log
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.camera.core.*
import androidx.camera.lifecycle.ProcessCameraProvider
import androidx.core.content.ContextCompat
import androidx.lifecycle.lifecycleScope
import com.google.mlkit.vision.barcode.BarcodeScannerOptions
import com.google.mlkit.vision.barcode.BarcodeScanning
import com.google.mlkit.vision.barcode.common.Barcode
import com.google.mlkit.vision.common.InputImage
import com.google.mlkit.vision.text.TextRecognition
import com.google.mlkit.vision.text.latin.TextRecognizerOptions
import com.mtgdetector.databinding.ActivityScanBinding
import com.mtgdetector.models.BulkScanRequest
import com.mtgdetector.models.ScanRequest
import com.mtgdetector.network.RetrofitClient
import kotlinx.coroutines.launch
import java.util.concurrent.ExecutorService
import java.util.concurrent.Executors

class ScanActivity : AppCompatActivity() {
    private lateinit var binding: ActivityScanBinding
    private lateinit var cameraExecutor: ExecutorService
    private var bulkMode = false
    private val scannedCards = mutableListOf<ScanRequest>()
    private var isProcessing = false
    private var camera: Camera? = null
    private var isFlashlightOn = false

    private val barcodeScanner = BarcodeScanning.getClient(
        BarcodeScannerOptions.Builder()
            .setBarcodeFormats(Barcode.FORMAT_ALL_FORMATS)
            .build()
    )

    private val textRecognizer = TextRecognition.getClient(TextRecognizerOptions.DEFAULT_OPTIONS)

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityScanBinding.inflate(layoutInflater)
        setContentView(binding.root)

        bulkMode = intent.getBooleanExtra("BULK_MODE", false)
        cameraExecutor = Executors.newSingleThreadExecutor()

        startCamera()

        binding.btnFlashlight.setOnClickListener {
            toggleFlashlight()
        }

        binding.btnCapture.setOnClickListener {
            val manualInput = binding.cardNameInput.text?.toString()?.trim()
            if (!manualInput.isNullOrEmpty()) {
                // Use manual input instead of camera scan
                scanCard(ScanRequest(cardName = manualInput))
            } else {
                // Use camera to scan
                captureAndScan()
            }
        }

        binding.btnDone.setOnClickListener {
            if (bulkMode && scannedCards.isNotEmpty()) {
                submitBulkScan()
            } else {
                finish()
            }
        }
    }

    private fun toggleFlashlight() {
        camera?.let {
            isFlashlightOn = !isFlashlightOn
            it.cameraControl.enableTorch(isFlashlightOn)
            binding.btnFlashlight.text = if (isFlashlightOn) "Light ON" else "Light"
        }
    }

    private lateinit var imageCapture: ImageCapture

    private fun startCamera() {
        val cameraProviderFuture = ProcessCameraProvider.getInstance(this)

        cameraProviderFuture.addListener({
            val cameraProvider = cameraProviderFuture.get()

            val preview = Preview.Builder()
                .build()
                .also {
                    it.setSurfaceProvider(binding.previewView.surfaceProvider)
                }

            imageCapture = ImageCapture.Builder()
                .setCaptureMode(ImageCapture.CAPTURE_MODE_MINIMIZE_LATENCY)
                .build()

            val cameraSelector = CameraSelector.DEFAULT_BACK_CAMERA

            try {
                cameraProvider.unbindAll()
                camera = cameraProvider.bindToLifecycle(
                    this, cameraSelector, preview, imageCapture
                )
            } catch (exc: Exception) {
                Log.e(TAG, "Camera binding failed", exc)
            }
        }, ContextCompat.getMainExecutor(this))
    }

    @androidx.camera.core.ExperimentalGetImage
    private fun captureAndScan() {
        if (isProcessing) return
        if (!::imageCapture.isInitialized) {
            Toast.makeText(this, "Camera not ready", Toast.LENGTH_SHORT).show()
            return
        }

        isProcessing = true
        binding.progressBar.visibility = View.VISIBLE
        binding.statusText.text = "Processing..."

        imageCapture.takePicture(
            cameraExecutor,
            object : ImageCapture.OnImageCapturedCallback() {
                override fun onCaptureSuccess(imageProxy: ImageProxy) {
                    processImage(imageProxy)
                }

                override fun onError(exception: ImageCaptureException) {
                    runOnUiThread {
                        binding.progressBar.visibility = View.GONE
                        binding.statusText.text = "Capture failed"
                        isProcessing = false
                    }
                }
            }
        )
    }

    @androidx.annotation.OptIn(androidx.camera.core.ExperimentalGetImage::class)
    private fun processImage(imageProxy: ImageProxy) {
        val mediaImage = imageProxy.image
        if (mediaImage != null) {
            val image = InputImage.fromMediaImage(mediaImage, imageProxy.imageInfo.rotationDegrees)

            // Try barcode first
            barcodeScanner.process(image)
                .addOnSuccessListener { barcodes ->
                    if (barcodes.isNotEmpty()) {
                        handleBarcodeDetected(barcodes.first())
                    } else {
                        // Fall back to text recognition
                        recognizeText(image)
                    }
                    imageProxy.close()
                }
                .addOnFailureListener {
                    // Fall back to text recognition
                    recognizeText(image)
                    imageProxy.close()
                }
        } else {
            imageProxy.close()
        }
    }

    private fun recognizeText(image: InputImage) {
        textRecognizer.process(image)
            .addOnSuccessListener { visionText ->
                // Extract card name from recognized text
                val text = visionText.text
                val cardName = extractCardName(text)

                if (cardName.isNotEmpty()) {
                    scanCard(ScanRequest(cardName = cardName))
                } else {
                    runOnUiThread {
                        binding.progressBar.visibility = View.GONE
                        binding.statusText.text = "No card detected"
                        isProcessing = false
                    }
                }
            }
            .addOnFailureListener { e ->
                runOnUiThread {
                    binding.progressBar.visibility = View.GONE
                    binding.statusText.text = "Scan failed"
                    isProcessing = false
                }
            }
    }

    private fun handleBarcodeDetected(barcode: Barcode) {
        val barcodeValue = barcode.rawValue ?: ""
        scanCard(ScanRequest(barcode = barcodeValue))
    }

    private fun extractCardName(text: String): String {
        // Simple heuristic: take the longest line as card name
        val lines = text.split("\n").map { it.trim() }
        return lines.maxByOrNull { it.length } ?: ""
    }

    private fun scanCard(request: ScanRequest) {
        lifecycleScope.launch {
            try {
                val response = RetrofitClient.apiService.scanCard(request)

                runOnUiThread {
                    binding.progressBar.visibility = View.GONE
                    isProcessing = false

                    if (response.isSuccessful && response.body()?.success == true) {
                        val card = response.body()?.card
                        binding.statusText.text = "Added: ${card?.name}"
                        Toast.makeText(this@ScanActivity, "Card added!", Toast.LENGTH_SHORT).show()

                        // Clear the manual input field after successful scan
                        binding.cardNameInput.text?.clear()

                        if (bulkMode) {
                            scannedCards.add(request)
                            binding.btnDone.visibility = View.VISIBLE
                        } else {
                            finish()
                        }
                    } else {
                        binding.statusText.text = "Card not found"
                    }
                }
            } catch (e: Exception) {
                runOnUiThread {
                    binding.progressBar.visibility = View.GONE
                    binding.statusText.text = "Error: ${e.message}"
                    isProcessing = false
                }
            }
        }
    }

    private fun submitBulkScan() {
        binding.progressBar.visibility = View.VISIBLE
        lifecycleScope.launch {
            try {
                val response = RetrofitClient.apiService.scanBulk(
                    BulkScanRequest(scannedCards)
                )

                runOnUiThread {
                    binding.progressBar.visibility = View.GONE
                    if (response.isSuccessful) {
                        val result = response.body()
                        Toast.makeText(
                            this@ScanActivity,
                            "Scanned ${result?.successfulScans}/${result?.totalScanned} cards",
                            Toast.LENGTH_LONG
                        ).show()
                        finish()
                    } else {
                        Toast.makeText(this@ScanActivity, "Bulk scan failed", Toast.LENGTH_SHORT).show()
                    }
                }
            } catch (e: Exception) {
                runOnUiThread {
                    binding.progressBar.visibility = View.GONE
                    Toast.makeText(this@ScanActivity, "Error: ${e.message}", Toast.LENGTH_SHORT).show()
                }
            }
        }
    }

    override fun onDestroy() {
        super.onDestroy()
        cameraExecutor.shutdown()
        barcodeScanner.close()
        textRecognizer.close()
    }

    companion object {
        private const val TAG = "ScanActivity"
    }
}
