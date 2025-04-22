package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Ethereum-fund-flow-analysis/internal/config"
	"Ethereum-fund-flow-analysis/internal/etherscan"
	"Ethereum-fund-flow-analysis/internal/models"
	"Ethereum-fund-flow-analysis/internal/services"
)

// Handler contains the dependencies needed by the API handlers
type Handler struct {
	analysisService *service.AnalysisService
}

// NewHandler creates a new API handler
func NewHandler(cfg *config.Config) *Handler {
	etherscanClient := etherscan.NewClient(cfg.EtherscanBaseURL, cfg.EtherscanAPIKey)
	analysisService := service.NewAnalysisService(etherscanClient)
	
	return &Handler{
		analysisService: analysisService,
	}
}

// BeneficiaryHandler handles requests to the /beneficiary endpoint
func (h *Handler) BeneficiaryHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Get the Ethereum address from query parameters
	address := r.URL.Query().Get("address")
	if address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}
	
	// Validate the Ethereum address (basic validation)
	address = strings.ToLower(address)
	if !strings.HasPrefix(address, "0x") || len(address) != 42 {
		http.Error(w, "Invalid Ethereum address", http.StatusBadRequest)
		return
	}
	
	// Analyze beneficiaries
	beneficiaries, err := h.analysisService.AnalyzeBeneficiaries(address)
	if err != nil {
		log.Printf("Error analyzing beneficiaries: %v", err)
		http.Error(w, "Failed to analyze beneficiaries", http.StatusInternalServerError)
		return
	}
	
	// Create the response
	response := models.BeneficiaryResponse{
		Message: "success",
		Data:    beneficiaries,
	}
	
	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
