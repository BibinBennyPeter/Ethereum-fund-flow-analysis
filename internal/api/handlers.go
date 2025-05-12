package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"Ethereum-fund-flow-analysis/internal/client"
	"Ethereum-fund-flow-analysis/internal/config"
	"Ethereum-fund-flow-analysis/internal/models"
	"Ethereum-fund-flow-analysis/internal/services"
)

// Handler contains the dependencies needed by the API handlers
type Handler struct {
	analysisService *service.AnalysisService
}

// NewHandler creates a new API handler
func NewHandler(cfg *config.Config) *Handler {
	etherscanClient := client.NewClient(cfg.EtherscanBaseURL, cfg.EtherscanAPIKey)
	analysisService := service.NewAnalysisService(etherscanClient)

	return &Handler{
		analysisService: analysisService,
	}
}

// FilterAndSortParams defines the parameters for filtering and sorting transactions
type FilterAndSortParams struct {
	// Etherscan API params (applied at the data source)
	Address    string
  ChainId    int
	StartBlock int64
	EndBlock   int64
	Page       int
	Offset     int
	Sort       string // "asc" or "desc" for the API call
  ApiKey     string

	// Custom filtering params (applied after fetching data)
	MinAmount   float64
	MaxAmount   float64
	SortBy      string // "amount"
	Limit       int    // Maximum number of results to return
	WithZeroTxs bool   // Include entries with zero amount transactions
}


// parseQueryParams extracts and validates the filter and sort parameters from the request
func parseQueryParams(r *http.Request) (FilterAndSortParams, error) {
	query := r.URL.Query()
	params := FilterAndSortParams{
		Address:     query.Get("address"),
    ChainId:     1,          // Default to 1, Ethereum Mainnet
		MinAmount:   0,
		MaxAmount:   -1, // Negative value means no maximum limit
		StartBlock:  0,
		EndBlock:    -1,       // Negative value means no end block limit
		Page:        1,        // Default to first page
		Offset:      100,      // Default to 100 transactions per page
		SortBy:      "amount", // Default sort by amount
		Sort:        "desc",   // Default sort order is descending
    ApiKey:       "",
		Limit:       100,      // Default limit is 100 results
		WithZeroTxs: true,     // By default, include zero amount transactions
	}
  
  // Parse chain in
  if chainId := query.Get("chainid"); chainId != ""{
    chainId, err := strconv.Atoi(chainId)

    if err != nil{
      return params, err
    }

    if validateChainId(chainId){
      params.ChainId = chainId
    }
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
	if sortBy := query.Get("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}

  //Parse api key 
  if apiKey := query.Get("apikey"); apiKey != "" {
    params.ApiKey = apiKey
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

// validateAddress checks if the provided Ethereum address is valid
func validateAddress(address string) error {
	address = strings.ToLower(address)
	if !strings.HasPrefix(address, "0x") || len(address) != 42 {
		return errors.New("invalid Ethereum address")
	}
	return nil
}

// validateChainId checks if the provided chain id is valid
func validateChainId(chainid int) bool {
  validChainIDs := map[int]struct{}{
    1: {}, 11155111: {}, 17000: {}, 2741: {}, 11124: {},
    33111: {}, 33139: {}, 42170: {}, 42161: {}, 421614: {},
    43114: {}, 43113: {}, 8453: {}, 84532: {}, 80094: {},
    80069: {}, 199: {}, 1028: {}, 81457: {}, 168587773: {},
    56: {}, 97: {}, 44787: {}, 42220: {}, 25: {},
    252: {}, 2522: {}, 100: {}, 59144: {}, 59141: {},
    5000: {}, 5003: {}, 4352: {}, 43521: {}, 1287: {},
    1284: {}, 1285: {}, 10: {}, 11155420: {}, 80002: {},
    137: {}, 2442: {}, 1101: {}, 534352: {}, 534351: {},
    57054: {}, 146: {}, 50104: {}, 531050104: {}, 1923: {},
    1924: {}, 167009: {}, 167000: {}, 130: {}, 1301: {},
    1111: {}, 1112: {}, 480: {}, 4801: {}, 660279: {},
    37714555429: {}, 51: {}, 50: {}, 324: {}, 300: {},
  }
  _, ok := validChainIDs[chainid]
  return ok
}

// httpHelper contains common HTTP handler operations
type httpHelper struct{}

// ensureMethod ensures the request uses the allowed HTTP method
func (h httpHelper) ensureMethod(w http.ResponseWriter, r *http.Request, method string) bool {
	if r.Method != method {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

// getValidParams extracts and validates params and address
func (h httpHelper) getValidParams(w http.ResponseWriter, r *http.Request) (FilterAndSortParams, bool) {
	// Parse filter and sort parameters
	params, err := parseQueryParams(r)
	if err != nil {
		http.Error(w, "Invalid query parameters: "+err.Error(), http.StatusBadRequest)
		return params, false
	}

	// Check address presence
	if params.Address == "" {
		http.Error(w, "Address parameter is required", http.StatusBadRequest)
		return params, false
	}

	// Validate the Ethereum address
	if err := validateAddress(params.Address); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return params, false
	}

	return params, true
}

// toAnalysisParams converts HTTP params to service params
func (h httpHelper) toAnalysisParams(params FilterAndSortParams) service.AnalysisParams {
	return service.AnalysisParams{
		Address:    params.Address,
    ChainId:    params.ChainId,
		StartBlock: params.StartBlock,
		EndBlock:   params.EndBlock,
		Page:       params.Page,
		Offset:     params.Offset,
		Sort:       params.Sort,
    ApiKey:     params.ApiKey,
	}
}

// respondWithJSON sends a JSON response
func (h httpHelper) respondWithJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// BeneficiaryHandler handles requests to the /beneficiary endpoint
func (h *Handler) BeneficiaryHandler(w http.ResponseWriter, r *http.Request) {
	helper := httpHelper{}

	// Validate HTTP method
	if !helper.ensureMethod(w, r, http.MethodGet) {
		return
	}

	// Parse and validate parameters
	params, ok := helper.getValidParams(w, r)
	if !ok {
		return
	}

	// Get beneficiaries from the service
	beneficiaries, err := h.analysisService.AnalyzeBeneficiaries(helper.toAnalysisParams(params))
	if err != nil {
		log.Printf("Error analyzing beneficiaries: %v", err)
		http.Error(w, "Failed to analyze beneficiaries", http.StatusInternalServerError)
		return
	}

	// Apply filtering and sorting
	filteredBeneficiaries := filterBeneficiaries(beneficiaries, params)

	// Create the response
	response := models.BeneficiaryResponse{
		Message: "success",
		Data:    filteredBeneficiaries,
	}

	// Send JSON response
	helper.respondWithJSON(w, response)
}

// PayerHandler handles requests to the /payer endpoint
func (h *Handler) PayerHandler(w http.ResponseWriter, r *http.Request) {
	helper := httpHelper{}

	// Validate HTTP method
	if !helper.ensureMethod(w, r, http.MethodGet) {
		return
	}

	// Parse and validate parameters
	params, ok := helper.getValidParams(w, r)
	if !ok {
		return
	}

	// Get payers from the service
	payers, err := h.analysisService.AnalyzePayers(helper.toAnalysisParams(params))
	if err != nil {
		log.Printf("Error analyzing payers: %v", err)
		http.Error(w, "Failed to analyze payers", http.StatusInternalServerError)
		return
	}

	// Apply filtering and sorting
	filteredPayers := filterPayers(payers, params)

	// Create the response
	response := models.PayerResponse{
		Message: "success",
		Data:    filteredPayers,
	}

	// Send JSON response
	helper.respondWithJSON(w, response)
}

// filterBeneficiaries applies filtering and sorting to beneficiaries based on the params
func filterBeneficiaries(beneficiaries []models.Beneficiary, params FilterAndSortParams) []models.Beneficiary {
	filtered := []models.Beneficiary{}

	// Apply filters
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

		filtered = append(filtered, ben)
	}

	// Apply sorting
	if params.SortBy == "amount" {
		if params.Sort == "asc" {
			sort.Slice(filtered, func(i, j int) bool {
				return filtered[i].Amount < filtered[j].Amount
			})
		} else {
			sort.Slice(filtered, func(i, j int) bool {
				return filtered[i].Amount > filtered[j].Amount
			})
		}
	}

	// Apply limit
	if len(filtered) > params.Limit {
		filtered = filtered[:params.Limit]
	}

	return filtered
}

// filterPayers applies filtering and sorting to payers based on the params
func filterPayers(payers []models.Payer, params FilterAndSortParams) []models.Payer {
	filtered := []models.Payer{}

	// Apply filters
	for _, payer := range payers {
		// Skip payers with zero amount if not explicitly requested
		if payer.Amount == 0 && !params.WithZeroTxs {
			continue
		}

		// Apply amount filters
		if payer.Amount < params.MinAmount {
			continue
		}
		if params.MaxAmount > 0 && payer.Amount > params.MaxAmount {
			continue
		}

		filtered = append(filtered, payer)
	}

	// Apply sorting
	if params.SortBy == "amount" {
		if params.Sort == "asc" {
			sort.Slice(filtered, func(i, j int) bool {
				return filtered[i].Amount < filtered[j].Amount
			})
		} else {
			sort.Slice(filtered, func(i, j int) bool {
				return filtered[i].Amount > filtered[j].Amount
			})
		}
	}

	// Apply limit
	if len(filtered) > params.Limit {
		filtered = filtered[:params.Limit]
	}

	return filtered
}
