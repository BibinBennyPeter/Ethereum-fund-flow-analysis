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

// FetchTask defines a generic transaction fetch operation
type FetchTask[T any] struct {
	Name     string                                                                  // Name of the transaction type for error reporting
	Fetcher  func(params client.EtherscanRequestParams) ([]T, error)                 // Function to fetch transactions
	Assigner func(collection *TransactionCollection, result []T)                     // Function to assign results to the collection
}

// FetchAllTransactions concurrently fetches all transaction types for an address
func FetchAllTransactions(ethClient *client.Client, params client.EtherscanRequestParams) (TransactionCollection, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	result := TransactionCollection{
		Errors: []error{},
	}

	// Define all fetch tasks using generics
	fetchTasks := []any{
		// Normal transactions task
		FetchTask[models.NormalTx]{
			Name:    "normal transactions",
			Fetcher: ethClient.GetNormalTransactions,
			Assigner: func(collection *TransactionCollection, txs []models.NormalTx) {
				collection.NormalTxs = txs
			},
		},
		// Internal transactions task
		FetchTask[models.InternalTx]{
			Name:    "internal transactions",
			Fetcher: ethClient.GetInternalTransactions,
			Assigner: func(collection *TransactionCollection, txs []models.InternalTx) {
				collection.InternalTxs = txs
			},
		},
		// ERC20 transfers task
		FetchTask[models.ERC20Transfer]{
			Name:    "ERC20 transfers",
			Fetcher: ethClient.GetERC20Transfers,
			Assigner: func(collection *TransactionCollection, txs []models.ERC20Transfer) {
				collection.ERC20Txs = txs
			},
		},
		// ERC721 transfers task
		FetchTask[models.ERC721Transfer]{
			Name:    "ERC721 transfers",
			Fetcher: ethClient.GetERC721Transfers,
			Assigner: func(collection *TransactionCollection, txs []models.ERC721Transfer) {
				collection.ERC721Txs = txs
			},
		},
		// ERC1155 transfers task
		FetchTask[models.ERC1155Transfer]{
			Name:    "ERC1155 transfers",
			Fetcher: ethClient.GetERC1155Transfers,
			Assigner: func(collection *TransactionCollection, txs []models.ERC1155Transfer) {
				collection.ERC1155Txs = txs
			},
		},
	}

	// Execute all fetch tasks concurrently
	for _, task := range fetchTasks {
		wg.Add(1)
		
		// Use type assertions to handle different transaction types
		switch t := task.(type) {
		case FetchTask[models.NormalTx]:
			go executeTask(&wg, &mu, &result, params, t)
		case FetchTask[models.InternalTx]:
			go executeTask(&wg, &mu, &result, params, t)
		case FetchTask[models.ERC20Transfer]:
			go executeTask(&wg, &mu, &result, params, t)
		case FetchTask[models.ERC721Transfer]:
			go executeTask(&wg, &mu, &result, params, t)
		case FetchTask[models.ERC1155Transfer]:
			go executeTask(&wg, &mu, &result, params, t)
		}
	}

	// Wait for all goroutines to complete
	wg.Wait()

	// Check if there were any errors
	if len(result.Errors) > 0 {
		return result, fmt.Errorf("failed to fetch some transactions: %v", result.Errors)
	}

	return result, nil
}

// executeTask runs a fetch task and safely updates the result collection
func executeTask[T any](
	wg *sync.WaitGroup,
	mu *sync.Mutex,
	result *TransactionCollection,
	params client.EtherscanRequestParams,
	task FetchTask[T],
) {
	defer wg.Done()
	
	// Execute the fetch operation
	txs, err := task.Fetcher(params)
	
	// Safely update the result collection
	mu.Lock()
	defer mu.Unlock()
	
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("%s: %w", task.Name, err))
		return
	}
	
	// Assign results to the appropriate field in the collection
	task.Assigner(result, txs)
}
