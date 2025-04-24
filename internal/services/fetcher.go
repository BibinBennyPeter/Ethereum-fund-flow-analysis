package service

import (
	"fmt"
	"sync"

	"Ethereum-fund-flow-analysis/internal/client"
	"Ethereum-fund-flow-analysis/internal/models"
)

// TransactionCollection holds all types of transactions
type TransactionCollection struct {
	NormalTxs   []models.NormalTx
	InternalTxs []models.InternalTx
	ERC20Txs    []models.ERC20Transfer
	ERC721Txs   []models.ERC721Transfer
	ERC1155Txs  []models.ERC1155Transfer
	Errors      []error
}

// FetchAllTransactions concurrently fetches all transaction types for an address
func FetchAllTransactions(ethClient *client.Client, params client.EtherscanRequestParams) (TransactionCollection, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	result := TransactionCollection{
		Errors: []error{},
	}

	// Fetch normal transactions
	wg.Add(1)
	go func() {
		defer wg.Done()
		txs, err := ethClient.GetNormalTransactions(params)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("normal transactions: %w", err))
			return
		}
		result.NormalTxs = txs
	}()

	// Fetch internal transactions
	wg.Add(1)
	go func() {
		defer wg.Done()
		txs, err := ethClient.GetInternalTransactions(params)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("internal transactions: %w", err))
			return
		}
		result.InternalTxs = txs
	}()

	// Fetch ERC20 transfers
	wg.Add(1)
	go func() {
		defer wg.Done()
		txs, err := ethClient.GetERC20Transfers(params)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("ERC20 transfers: %w", err))
			return
		}
		result.ERC20Txs = txs
	}()

	// Fetch ERC721 transfers
	wg.Add(1)
	go func() {
		defer wg.Done()
		txs, err := ethClient.GetERC721Transfers(params)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("ERC721 transfers: %w", err))
			return
		}
		result.ERC721Txs = txs
	}()

	// Fetch ERC1155 transfers
	wg.Add(1)
	go func() {
		defer wg.Done()
		txs, err := ethClient.GetERC1155Transfers(params)
		mu.Lock()
		defer mu.Unlock()
		if err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("ERC1155 transfers: %w", err))
			return
		}
		result.ERC1155Txs = txs
	}()

	// Wait for all goroutines to complete
	wg.Wait()

	// Check if there were any errors
	if len(result.Errors) > 0 {
		return result, fmt.Errorf("failed to fetch some transactions: %v", result.Errors)
	}

	return result, nil
}
