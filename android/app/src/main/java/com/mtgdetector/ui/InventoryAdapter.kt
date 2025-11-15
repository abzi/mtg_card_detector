package com.mtgdetector.ui

import android.view.LayoutInflater
import android.view.ViewGroup
import androidx.recyclerview.widget.DiffUtil
import androidx.recyclerview.widget.ListAdapter
import androidx.recyclerview.widget.RecyclerView
import com.bumptech.glide.Glide
import com.mtgdetector.R
import com.mtgdetector.databinding.ItemCardBinding
import com.mtgdetector.models.InventoryItem

class InventoryAdapter : ListAdapter<InventoryItem, InventoryAdapter.ViewHolder>(DiffCallback()) {

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int): ViewHolder {
        val binding = ItemCardBinding.inflate(
            LayoutInflater.from(parent.context),
            parent,
            false
        )
        return ViewHolder(binding)
    }

    override fun onBindViewHolder(holder: ViewHolder, position: Int) {
        holder.bind(getItem(position))
    }

    class ViewHolder(private val binding: ItemCardBinding) :
        RecyclerView.ViewHolder(binding.root) {

        fun bind(item: InventoryItem) {
            binding.tvCardName.text = item.card?.name ?: "Unknown Card"
            binding.tvSetInfo.text = "${item.card?.setCode} #${item.card?.collectorNumber}"
            binding.tvQuantity.text = binding.root.context.getString(R.string.quantity, item.quantity)

            // Load card image
            item.card?.imageUri?.let { imageUri ->
                Glide.with(binding.root.context)
                    .load(imageUri)
                    .placeholder(android.R.drawable.ic_menu_gallery)
                    .error(android.R.drawable.ic_menu_report_image)
                    .into(binding.ivCardImage)
            }
        }
    }

    class DiffCallback : DiffUtil.ItemCallback<InventoryItem>() {
        override fun areItemsTheSame(oldItem: InventoryItem, newItem: InventoryItem): Boolean {
            return oldItem.id == newItem.id
        }

        override fun areContentsTheSame(oldItem: InventoryItem, newItem: InventoryItem): Boolean {
            return oldItem == newItem
        }
    }
}
