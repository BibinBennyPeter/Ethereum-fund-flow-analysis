package service

import (
	"fmt"
	"math"
  "strings"
	"math/big"
	"time"

	"Ethereum-fund-flow-analysis/internal/etherscan"
	"Ethereum-fund-flow-analysis/internal/models"
)

const weiToEther = 1e18

// AnalysisService handles the transaction analysis logic
type AnalysisService struct {
	etherscanClient *etherscan.Client
}

// NewAnalysisService creates a new analysis service
func NewAnalysisService(etherscanClient *etherscan.Client) *AnalysisService {
	return &AnalysisService{
		etherscanClient: etherscanClient,
	}
}

// AnalyzeBeneficiaries analyzes transactions to identify beneficiaries
func (s *AnalysisService) AnalyzeBeneficiaries(address string) ([]models.Beneficiary, error) {
	// Get all transaction types
	normalTxs, err := s.etherscanClient.GetNormalTransactions(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get normal transactions: %w", err)
	}

	internalTxs, err := s.etherscanClient.GetInternalTransactions(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get internal transactions: %w", err)
	}

	erc20Txs, err := s.etherscanClient.GetERC20Transfers(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}

  erc721Txs, err := s.etherscanClient.GetERC721Transfers(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}

  erc1155Txs, err := s.etherscanClient.GetERC1155Transfers(address)
	if err != nil {
		return nil, fmt.Errorf("failed to get token transfers: %w", err)
	}

  // Create a map to collect beneficiaries
	beneficiaryMap := make(map[string]*models.Beneficiary)

  // Process normal transactions
	for _, tx := range normalTxs {
		// Only consider outgoing transactions from our address
		if !strings.EqualFold(tx.From,  address) || tx.IsError == 1 {
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
	// Process internal transactions
	for _, tx := range internalTxs {
		if !strings.EqualFold(tx.From, address) || tx.IsError == 1 {
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


  for _, tx := range erc20Txs {
    if !strings.EqualFold(tx.From, address) {
        continue
    }

    // Convert raw big‑integer string to human amount
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


  for _, tx := range erc721Txs {
    if !strings.EqualFold(tx.From, address) {
        continue
    }

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")


    transaction := models.Transaction{
        TxAmount:      1,           // NFTs always transfer one
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


  for _, tx := range erc1155Txs {
    if !strings.EqualFold(tx.From, address) {
        continue
    }

    // Handle each ID/value pair (here single per event)
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
