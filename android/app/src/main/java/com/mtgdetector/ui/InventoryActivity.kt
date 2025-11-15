package com.mtgdetector.ui

import android.os.Bundle
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.lifecycleScope
import androidx.recyclerview.widget.LinearLayoutManager
import com.mtgdetector.R
import com.mtgdetector.databinding.ActivityInventoryBinding
import com.mtgdetector.network.RetrofitClient
import kotlinx.coroutines.launch

class InventoryActivity : AppCompatActivity() {
    private lateinit var binding: ActivityInventoryBinding
    private lateinit var adapter: InventoryAdapter

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityInventoryBinding.inflate(layoutInflater)
        setContentView(binding.root)

        supportActionBar?.setDisplayHomeAsUpEnabled(true)
        supportActionBar?.title = getString(R.string.inventory)

        adapter = InventoryAdapter()
        binding.recyclerView.layoutManager = LinearLayoutManager(this)
        binding.recyclerView.adapter = adapter

        loadInventory()
    }

    private fun loadInventory() {
        binding.progressBar.visibility = View.VISIBLE
        binding.tvEmpty.visibility = View.GONE

        lifecycleScope.launch {
            try {
                val response = RetrofitClient.apiService.getInventory()

                if (response.isSuccessful && response.body() != null) {
                    val inventory = response.body()!!
                    binding.progressBar.visibility = View.GONE

                    if (inventory.inventory.isEmpty()) {
                        binding.tvEmpty.visibility = View.VISIBLE
                    } else {
                        binding.tvTotalCards.text = getString(R.string.total_cards, inventory.count)
                        adapter.submitList(inventory.inventory)
                    }
                } else {
                    binding.progressBar.visibility = View.GONE
                    Toast.makeText(
                        this@InventoryActivity,
                        "Failed to load inventory",
                        Toast.LENGTH_SHORT
                    ).show()
                }
            } catch (e: Exception) {
                binding.progressBar.visibility = View.GONE
                Toast.makeText(
                    this@InventoryActivity,
                    "Error: ${e.message}",
                    Toast.LENGTH_SHORT
                ).show()
            }
        }
    }

    override fun onSupportNavigateUp(): Boolean {
        finish()
        return true
    }
}
