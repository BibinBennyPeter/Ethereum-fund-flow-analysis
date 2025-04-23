package api

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strconv"
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

// FilterAndSortParams defines the parameters for filtering and sorting beneficiaries
type FilterAndSortParams struct {
	// Etherscan API params (applied at the data source)
	Address    string
	StartBlock int64
	EndBlock   int64
	Page       int
	Offset     int
	Sort  string // "asc" or "desc" for the API call

	// Custom filtering params (applied after fetching data)
	MinAmount   float64
	MaxAmount   float64
	SortBy      string // "amount"
	Limit       int    // Maximum number of beneficiaries to return
	WithZeroTxs bool   // Include beneficiaries with zero amount transactions
}

// parseQueryParams extracts and validates the filter and sort parameters from the request
func parseQueryParams(r *http.Request) (FilterAndSortParams, error) {
	query := r.URL.Query()
	params := FilterAndSortParams{
		Address:     query.Get("address"),
		MinAmount:   0,
		MaxAmount:   -1, // Negative value means no maximum limit
		StartBlock:  0,
		EndBlock:    -1, // Negative value means no end block limit
		Page:        1,  // Default to first page
		Offset:      100, // Default to 100 transactions per page
		SortBy:      "amount", // Default sort by amount
		Sort:        "desc",   // Default sort order is descending
		Limit:       100,      // Default limit is 100 beneficiaries
		WithZeroTxs: true,    // By default, exclude zero amount transactions
	}

	// Parse min amount
	if minAmtStr := query.Get("min"); minAmtStr != "" {
		minAmt, err := strconv.ParseFloat(minAmtStr, 64)
		if err != nil {
			return params, err
		}
		params.MinAmount = minAmt
	}

	// Parse max amount
	if maxAmtStr := query.Get("max"); maxAmtStr != "" {
		maxAmt, err := strconv.ParseFloat(maxAmtStr, 64)
		if err != nil {
			return params, err
		}
		params.MaxAmount = maxAmt
	}

	// Parse start block
	if startBlockStr := query.Get("sblock"); startBlockStr != "" {
		startBlock, err := strconv.ParseInt(startBlockStr, 10, 64)
		if err != nil {
			return params, err
		}
		params.StartBlock = startBlock
	}

	// Parse end block
	if endBlockStr := query.Get("eblock"); endBlockStr != "" {
		endBlock, err := strconv.ParseInt(endBlockStr, 10, 64)
		if err != nil {
			return params, err
		}
		params.EndBlock = endBlock
	}

	// Parse page
	if pageStr := query.Get("page"); pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			return params, err
		}
		if page <= 0 {
			page = 1
		}
		params.Page = page
	}

	// Parse offset (items per page)
	if offsetStr := query.Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return params, err
		}
		if offset <= 0 {
			offset = 100
		}
		params.Offset = offset
	}

	// Parse sort by
	if sortBy := query.Get("sort_by"); strings.ToLower(sortBy) != "amount" {
			params.SortBy = "amount"
	}

	// Parse sort order
	if sortOrder := query.Get("sort"); sortOrder != "" {
		// Validate sort order parameter
		switch strings.ToLower(sortOrder) {
		case "asc", "desc":
			params.Sort = strings.ToLower(sortOrder)
		default:
			params.Sort = "desc" // Default to descending if invalid
		}
	}

	// Parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return params, err
		}
		// Ensure limit is positive
		if limit <= 0 {
			limit = 100 // Default to 100 if negative or zero
		}
		params.Limit = limit
	}

	// Parse with_zero_txs
	if withZeroTxsStr := query.Get("with_zero_txs"); withZeroTxsStr != "" {
		withZeroTxs, err := strconv.ParseBool(withZeroTxsStr)
		if err != nil {
			return params, err
		}
		params.WithZeroTxs = withZeroTxs
	}

	return params, nil
}

// BeneficiaryHandler handles requests to the /beneficiary endpoint
func (h *Handler) BeneficiaryHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow GET requests
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	
	// Parse filter and sort parameters
	params, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Get the Ethereum address from query parameters
	if params.Address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return
	}
	
	// Validate the Ethereum address (basic validation)
	params.Address = strings.ToLower(params.Address)
	if !strings.HasPrefix(params.Address, "0x") || len(params.Address) != 42 {
		http.Error(w, "Invalid Ethereum address", http.StatusBadRequest)
		return
	}

	// Prepare analysis parameters
	analysisParams := service.AnalysisParams{
		Address:     params.Address,
		StartBlock:  params.StartBlock,
		EndBlock:    params.EndBlock,
		Page:        params.Page,
		Offset:      params.Offset,
		Sort:        params.Sort,
	}


	// Get beneficiaries from the service
	beneficiaries, err := h.analysisService.AnalyzeBeneficiaries(analysisParams)
	if err != nil {
		log.Printf("Error analyzing beneficiaries: %v", err)
		http.Error(w, "Failed to analyze beneficiaries", http.StatusInternalServerError)
		return
	}
	
	// Apply filtering based on amount
	filteredBeneficiaries := []models.Beneficiary{}
	for _, ben := range beneficiaries {
		// Skip beneficiaries with zero amount if not explicitly requested
		if ben.Amount == 0 && !params.WithZeroTxs {
			continue
		}

		// Apply amount filters
		if ben.Amount < params.MinAmount {
			continue
		}
		if params.MaxAmount > 0 && ben.Amount > params.MaxAmount {
			continue
		}

		filteredBeneficiaries = append(filteredBeneficiaries, ben)
	}

	// Apply custom sorting amount
	if params.SortBy == "amount" {
		applySorting(filteredBeneficiaries, params.Sort)
	}

	// Apply limit
	if len(filteredBeneficiaries) > params.Limit {
		filteredBeneficiaries = filteredBeneficiaries[:params.Limit]
	}
	
	// Create the response
	response := models.BeneficiaryResponse{
		Message: "success",
		Data:    filteredBeneficiaries,
	}
	
	// Send JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// applySorting sorts the beneficiaries based on amount
func applySorting(beneficiaries []models.Beneficiary, sortOrder string) {
	if sortOrder == "asc" {
		sort.Slice(beneficiaries, func(i, j int) bool {
			return beneficiaries[i].Amount < beneficiaries[j].Amount
		})
	} else {
		sort.Slice(beneficiaries, func(i, j int) bool {
			return beneficiaries[i].Amount > beneficiaries[j].Amount
		})
	}
}
