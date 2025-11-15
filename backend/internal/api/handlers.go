package api

import (
	"encoding/json"
	"net/http"

	"github.com/abzi/mtg_card_detector/internal/auth"
	"github.com/abzi/mtg_card_detector/internal/database"
	"github.com/abzi/mtg_card_detector/internal/inventory"
	"github.com/abzi/mtg_card_detector/internal/middleware"
	"github.com/abzi/mtg_card_detector/internal/models"
)

type Handler struct {
	authService      *auth.Service
	inventoryService *inventory.Service
	db               *database.DB
}

func NewHandler(authService *auth.Service, inventoryService *inventory.Service, db *database.DB) *Handler {
	return &Handler{
		authService:      authService,
		inventoryService: inventoryService,
		db:               db,
	}
}

// HandleAnonymousAuth creates or retrieves anonymous user
func (h *Handler) HandleAnonymousAuth(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DeviceID string `json:"device_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.DeviceID == "" {
		respondError(w, http.StatusBadRequest, "device_id is required")
		return
	}

	authResp, err := h.authService.GenerateAnonymousUser(req.DeviceID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to authenticate")
		return
	}

	respondJSON(w, http.StatusOK, authResp)
}

// HandleSingleScan processes a single card scan
func (h *Handler) HandleSingleScan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req models.ScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	result, err := h.inventoryService.ProcessSingleScan(userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// HandleBulkScan processes multiple card scans
func (h *Handler) HandleBulkScan(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	var req models.BulkScanRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if len(req.Scans) == 0 {
		respondError(w, http.StatusBadRequest, "scans array cannot be empty")
		return
	}

	result, err := h.inventoryService.ProcessBulkScan(userID, &req)
	if err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusOK, result)
}

// HandleGetInventory retrieves user's inventory
func (h *Handler) HandleGetInventory(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		respondError(w, http.StatusUnauthorized, "user not authenticated")
		return
	}

	inventory, err := h.inventoryService.GetInventory(userID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to retrieve inventory")
		return
	}

	respondJSON(w, http.StatusOK, map[string]interface{}{
		"inventory": inventory,
		"count":     len(inventory),
	})
}

// HandleGetCard retrieves card details by ID
func (h *Handler) HandleGetCard(w http.ResponseWriter, r *http.Request) {
	cardID := r.URL.Query().Get("id")
	if cardID == "" {
		respondError(w, http.StatusBadRequest, "card id is required")
		return
	}

	card, err := h.db.GetCardByID(cardID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to retrieve card")
		return
	}

	if card == nil {
		respondError(w, http.StatusNotFound, "card not found")
		return
	}

	respondJSON(w, http.StatusOK, card)
}

// HandleHealthCheck returns API health status
func (h *Handler) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	respondJSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func respondError(w http.ResponseWriter, status int, message string) {
	respondJSON(w, status, models.ErrorResponse{
		Error:   message,
		Message: message,
	})
}
