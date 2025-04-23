package service

import (
	"fmt"
	"math"
	"math/big"
	"strings"
	"time"

	"Ethereum-fund-flow-analysis/internal/client"
	"Ethereum-fund-flow-analysis/internal/models"
)

const weiToEther = 1e18

// AnalysisService handles the transaction analysis logic
type AnalysisService struct {
	etherscanClient *client.Client
}

// AnalysisParams contains parameters for the beneficiary analysis
type AnalysisParams struct {
	Address    string
	StartBlock int64
	EndBlock   int64
	Page       int
	Offset     int
	Sort       string
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(etherscanClient *client.Client) *AnalysisService {
	return &AnalysisService{
		etherscanClient: etherscanClient,
	}
}

// AnalyzeBeneficiaries analyzes transactions to identify beneficiaries
func (s *AnalysisService) AnalyzeBeneficiaries(params AnalysisParams) ([]models.Beneficiary, error) {
	// Prepare Etherscan request parameters
	requestParams := client.EtherscanRequestParams{
		Address:    params.Address,
		StartBlock: params.StartBlock,
		EndBlock:   params.EndBlock,
		Page:       params.Page,
		Offset:     params.Offset,
		Sort:       params.Sort,
	}

	// Get all transaction types
	normalTxs, err := s.etherscanClient.GetNormalTransactions(requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get normal transactions: %w", err)
	}

	internalTxs, err := s.etherscanClient.GetInternalTransactions(requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal transactions: %w", err)
	}

	erc20Txs, err := s.etherscanClient.GetERC20Transfers(requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}

	erc721Txs, err := s.etherscanClient.GetERC721Transfers(requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}

	erc1155Txs, err := s.etherscanClient.GetERC1155Transfers(requestParams)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}

	// Create a map to collect beneficiaries
	beneficiaryMap := make(map[string]*models.Beneficiary)
	// Process normal transactions
	for _, tx := range normalTxs {
		// Only consider outgoing transactions from our address
		if !strings.EqualFold(tx.From, params.Address) || tx.IsError == 1 {
			continue
		}

		beneficiaryAddress := tx.To

		// Parse value in wei to ether
		valueWei := new(big.Int)
		valueWei.SetString(tx.Value.String(), 10)
		valueEther := new(big.Float).Quo(new(big.Float).SetInt(valueWei), big.NewFloat(weiToEther))
		amount, _ := valueEther.Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		// Add to beneficiary map
		if _, exists := beneficiaryMap[beneficiaryAddress]; !exists {
			beneficiaryMap[beneficiaryAddress] = &models.Beneficiary{
				Address:      beneficiaryAddress,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}

		beneficiaryMap[beneficiaryAddress].Amount += amount
		beneficiaryMap[beneficiaryAddress].Transactions = append(
			beneficiaryMap[beneficiaryAddress].Transactions,
			transaction,
		)
	}

	// Process internal transactions (same logic as above)
	for _, tx := range internalTxs {
		if !strings.EqualFold(tx.From, params.Address) || tx.IsError == 1 {
			continue
		}

		beneficiaryAddress := tx.To

		valueWei := new(big.Int)
		valueWei.SetString(tx.Value.String(), 10)
		valueEther := new(big.Float).Quo(new(big.Float).SetInt(valueWei), big.NewFloat(weiToEther))
		amount, _ := valueEther.Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, exists := beneficiaryMap[beneficiaryAddress]; !exists {
			beneficiaryMap[beneficiaryAddress] = &models.Beneficiary{
				Address:      beneficiaryAddress,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}

		beneficiaryMap[beneficiaryAddress].Amount += amount
		beneficiaryMap[beneficiaryAddress].Transactions = append(
			beneficiaryMap[beneficiaryAddress].Transactions,
			transaction,
		)
	}

	// Process ERC20 transfers
	for _, tx := range erc20Txs {
		if !strings.EqualFold(tx.From, params.Address) {
			continue
		}

		// Convert raw big-integer string to human amount
		divisor := new(big.Float).SetFloat64(math.Pow10(int(tx.TokenDecimal)))
		valueToken := new(big.Float)
		valueToken.SetString(tx.Value.String())
		amount, _ := new(big.Float).Quo(valueToken, divisor).Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, ok := beneficiaryMap[tx.To]; !ok {
			beneficiaryMap[tx.To] = &models.Beneficiary{
				Address:      tx.To,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}

		// For token transfers, we don't add to the ETH amount
		beneficiaryMap[tx.To].Transactions = append(
			beneficiaryMap[tx.To].Transactions,
			transaction,
		)
	}

	// Process ERC721 transfers (NFTs)
	for _, tx := range erc721Txs {
		if !strings.EqualFold(tx.From, params.Address) {
			continue
		}

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		transaction := models.Transaction{
			TxAmount:      1, // NFTs always transfer one
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, ok := beneficiaryMap[tx.To]; !ok {
			beneficiaryMap[tx.To] = &models.Beneficiary{
				Address:      tx.To,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}
		beneficiaryMap[tx.To].Transactions = append(
			beneficiaryMap[tx.To].Transactions,
			transaction,
		)
	}

	// Process ERC1155 transfers
	for _, tx := range erc1155Txs {
		if !strings.EqualFold(tx.From, params.Address) {
			continue
		}

		// Handle each ID/value pair
		divisor := new(big.Float).SetFloat64(math.Pow10(int(tx.TokenDecimal)))
		valueToken := new(big.Float)
		valueToken.SetString(tx.TokenValue.String())
		amount, _ := new(big.Float).Quo(valueToken, divisor).Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, ok := beneficiaryMap[tx.To]; !ok {
			beneficiaryMap[tx.To] = &models.Beneficiary{
				Address:      tx.To,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}
		beneficiaryMap[tx.To].Transactions = append(
			beneficiaryMap[tx.To].Transactions,
			transaction,
		)
	}

	// Convert map to slice
	beneficiaries := make([]models.Beneficiary, 0, len(beneficiaryMap))
	for _, ben := range beneficiaryMap {
		beneficiaries = append(beneficiaries, *ben)
	}

	return beneficiaries, nil

}

// AnalyzePayers analyzes incoming transactions to identify payers.
func (s *AnalysisService) AnalyzePayers(params AnalysisParams) ([]models.Payer, error) {
	// Prepare Etherscan request parameters
	req := client.EtherscanRequestParams{
		Address:    params.Address,
		StartBlock: params.StartBlock,
		EndBlock:   params.EndBlock,
		Page:       params.Page,
		Offset:     params.Offset,
		Sort:       params.Sort,
	}

	// Fetch all transaction types
	normalTxs, err := s.etherscanClient.GetNormalTransactions(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get normal transactions: %w", err)
	}
	internalTxs, err := s.etherscanClient.GetInternalTransactions(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal transactions: %w", err)
	}
	erc20Txs, err := s.etherscanClient.GetERC20Transfers(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}
	erc721Txs, err := s.etherscanClient.GetERC721Transfers(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}
	erc1155Txs, err := s.etherscanClient.GetERC1155Transfers(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}

	// Map to collect payers
	payerMap := make(map[string]*models.Payer)

	// Process normal incoming transactions
	for _, tx := range normalTxs {
		if !strings.EqualFold(tx.To, params.Address) || tx.IsError == 1 {
			continue
		}
		// Convert wei to ether
		valueWei := new(big.Int)
		valueWei.SetString(tx.Value.String(), 10)
		valueEther := new(big.Float).Quo(new(big.Float).SetInt(valueWei), big.NewFloat(weiToEther))
		amount, _ := valueEther.Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		tran := models.Transaction{TxAmount: amount, DateTime: dateTime, TransactionID: tx.Hash}

		if _, exists := payerMap[tx.From]; !exists {
			payerMap[tx.From] = &models.Payer{Address: tx.From, Amount: 0, Transactions: []models.Transaction{}}
		}
		payerMap[tx.From].Amount += amount
		payerMap[tx.From].Transactions = append(payerMap[tx.From].Transactions, tran)
	}

	// Process internal incoming transactions
	for _, tx := range internalTxs {
		if !strings.EqualFold(tx.To, params.Address) || tx.IsError == 1 {
			continue
		}
		valueWei := new(big.Int)
		valueWei.SetString(tx.Value.String(), 10)
		valueEther := new(big.Float).Quo(new(big.Float).SetInt(valueWei), big.NewFloat(weiToEther))
		amount, _ := valueEther.Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		tran := models.Transaction{TxAmount: amount, DateTime: dateTime, TransactionID: tx.Hash}

		if _, exists := payerMap[tx.From]; !exists {
			payerMap[tx.From] = &models.Payer{Address: tx.From, Amount: 0, Transactions: []models.Transaction{}}
		}
		payerMap[tx.From].Amount += amount
		payerMap[tx.From].Transactions = append(payerMap[tx.From].Transactions, tran)
	}

	// Process ERC20 incoming transfers
	for _, tx := range erc20Txs {
		if !strings.EqualFold(tx.To, params.Address) {
			continue
		}
		// Convert raw big-integer string to human amount
		divisor := new(big.Float).SetFloat64(math.Pow10(int(tx.TokenDecimal)))
		valueToken := new(big.Float)
		valueToken.SetString(tx.Value.String())
		amount, _ := new(big.Float).Quo(valueToken, divisor).Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")

		tran := models.Transaction{TxAmount: amount, DateTime: dateTime, TransactionID: tx.Hash}

		if _, exists := payerMap[tx.From]; !exists {
			payerMap[tx.From] = &models.Payer{Address: tx.From, Amount: 0, Transactions: []models.Transaction{}}
		}
		payerMap[tx.From].Transactions = append(payerMap[tx.From].Transactions, tran)
	}

	// Process ERC721 incoming transfers (NFTs)
	for _, tx := range erc721Txs {
		if !strings.EqualFold(tx.To, params.Address) {
			continue
		}
		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		tran := models.Transaction{TxAmount: 1, DateTime: dateTime, TransactionID: tx.Hash}

		if _, exists := payerMap[tx.From]; !exists {
			payerMap[tx.From] = &models.Payer{Address: tx.From, Amount: 0, Transactions: []models.Transaction{}}
		}
		payerMap[tx.From].Transactions = append(payerMap[tx.From].Transactions, tran)
	}

	// Process ERC1155 incoming transfers
	for _, tx := range erc1155Txs {
		if !strings.EqualFold(tx.To, params.Address) {
			continue
		}

		// Handle each ID/value pair
		divisor := new(big.Float).SetFloat64(math.Pow10(int(tx.TokenDecimal)))
		valueToken := new(big.Float)
		valueToken.SetString(tx.TokenValue.String())
		amount, _ := new(big.Float).Quo(valueToken, divisor).Float64()

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
		tran := models.Transaction{TxAmount: amount, DateTime: dateTime, TransactionID: tx.Hash}

		if _, exists := payerMap[tx.From]; !exists {
			payerMap[tx.From] = &models.Payer{Address: tx.From, Amount: 0, Transactions: []models.Transaction{}}
		}
		payerMap[tx.From].Transactions = append(payerMap[tx.From].Transactions, tran)
	}

	// Convert map to slice
	payers := make([]models.Payer, 0, len(payerMap))
	for _, p := range payerMap {
		payers = append(payers, *p)
	}

	return payers, nil
}
