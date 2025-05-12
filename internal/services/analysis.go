package service

import (
	"fmt"

	"Ethereum-fund-flow-analysis/internal/client"
	"Ethereum-fund-flow-analysis/internal/models"
)

// AnalysisService handles the transaction analysis logic
type AnalysisService struct {
	etherscanClient *client.Client
}

// AnalysisParams contains parameters for the analysis
type AnalysisParams struct {
	Address    string
  ChainId    int
	StartBlock int64
	EndBlock   int64
	Page       int
	Offset     int
	Sort       string
  ApiKey     string
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(etherscanClient *client.Client) *AnalysisService {
	return &AnalysisService{
		etherscanClient: etherscanClient,
	}
}

// AnalyzeBeneficiaries analyzes transactions to identify beneficiaries
func (s *AnalysisService) AnalyzeBeneficiaries(params AnalysisParams) ([]models.Beneficiary, error) {
	// Convert to Etherscan request params
	requestParams := client.EtherscanRequestParams{
		Address:    params.Address,
    ChainId:    params.ChainId,
		StartBlock: params.StartBlock,
		EndBlock:   params.EndBlock,
		Page:       params.Page,
		Offset:     params.Offset,
		Sort:       params.Sort,
    ApiKey:     params.ApiKey,
	}

	// Fetch all transactions concurrently
	txCollection, err := FetchAllTransactions(s.etherscanClient, requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Process transactions to find beneficiaries (outgoing = true)
	beneficiaryMap := ProcessTransactions(params.Address, txCollection, true)

	// Convert map to slice
	beneficiaries := make([]models.Beneficiary, 0, len(beneficiaryMap))
	for _, ben := range beneficiaryMap {
		beneficiaries = append(beneficiaries, models.Beneficiary{
			Address:      ben.Address,
			Amount:       ben.Amount,
			Transactions: ben.Transactions,
		})
	}

	return beneficiaries, nil
}

// AnalyzePayers analyzes incoming transactions to identify payers
func (s *AnalysisService) AnalyzePayers(params AnalysisParams) ([]models.Payer, error) {
	// Convert to Etherscan request params
	requestParams := client.EtherscanRequestParams{
		Address:    params.Address,
		StartBlock: params.StartBlock,
		EndBlock:   params.EndBlock,
		Page:       params.Page,
		Offset:     params.Offset,
		Sort:       params.Sort,
    ApiKey:     params.ApiKey,
	}

	// Fetch all transactions concurrently
	txCollection, err := FetchAllTransactions(s.etherscanClient, requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch transactions: %w", err)
	}

	// Process transactions to find payers (outgoing = false)
	payerMap := ProcessTransactions(params.Address, txCollection, false)

	// Convert map to slice
	payers := make([]models.Payer, 0, len(payerMap))
	for _, p := range payerMap {
		payers = append(payers, models.Payer{
			Address:      p.Address,
			Amount:       p.Amount,
			Transactions: p.Transactions,
		})
	}

	return payers, nil
}
