package service

import (
	"strings"

	"Ethereum-fund-flow-analysis/internal/models"
	"Ethereum-fund-flow-analysis/internal/utils"
)

// ProcessTransactions processes all transaction types into either beneficiaries or payers
func ProcessTransactions(address string, txCollection TransactionCollection, isOutgoing bool) map[string]*models.EntityWithTransactions {
	// This map will hold either beneficiaries or payers
	entityMap := make(map[string]*models.EntityWithTransactions)

	// Process normal transactions
	for _, tx := range txCollection.NormalTxs {
		// Skip failed transactions
		if tx.IsError == 1 {
			continue
		}

		var counterpartyAddress string
		if isOutgoing {
			// Only consider outgoing transactions from address
			if !strings.EqualFold(tx.From, address) {
				continue
			}
			counterpartyAddress = tx.To
		} else {
			// Only consider incoming transactions to address
			if !strings.EqualFold(tx.To, address) {
				continue
			}
			counterpartyAddress = tx.From
		}

		// Convert value
		amount := utils.ConvertWeiToEther(tx.Value.String())

		// Create transaction record
		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := utils.FormatTimestamp(timestamp)

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		// Add to entity map
		if _, exists := entityMap[counterpartyAddress]; !exists {
			entityMap[counterpartyAddress] = &models.EntityWithTransactions{
				Address:      counterpartyAddress,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}

		entityMap[counterpartyAddress].Amount += amount
		entityMap[counterpartyAddress].Transactions = append(
			entityMap[counterpartyAddress].Transactions,
			transaction,
		)
	}

	// Process internal transactions (similar logic)
	for _, tx := range txCollection.InternalTxs {
		if tx.IsError == 1 {
			continue
		}

		var counterpartyAddress string
		if isOutgoing {
			if !strings.EqualFold(tx.From, address) {
				continue
			}
			counterpartyAddress = tx.To
		} else {
			if !strings.EqualFold(tx.To, address) {
				continue
			}
			counterpartyAddress = tx.From
		}

		amount := utils.ConvertWeiToEther(tx.Value.String())
		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := utils.FormatTimestamp(timestamp)

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, exists := entityMap[counterpartyAddress]; !exists {
			entityMap[counterpartyAddress] = &models.EntityWithTransactions{
				Address:      counterpartyAddress,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}

		entityMap[counterpartyAddress].Amount += amount
		entityMap[counterpartyAddress].Transactions = append(
			entityMap[counterpartyAddress].Transactions,
			transaction,
		)
	}

	// Process ERC20 transfers
	for _, tx := range txCollection.ERC20Txs {
		var counterpartyAddress string
		if isOutgoing {
			if !strings.EqualFold(tx.From, address) {
				continue
			}
			counterpartyAddress = tx.To
		} else {
			if !strings.EqualFold(tx.To, address) {
				continue
			}
			counterpartyAddress = tx.From
		}

		amount := utils.ConvertTokenValueWithDecimals(tx.Value.String(), tx.TokenDecimal)
		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := utils.FormatTimestamp(timestamp)

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, ok := entityMap[counterpartyAddress]; !ok {
			entityMap[counterpartyAddress] = &models.EntityWithTransactions{
				Address:      counterpartyAddress,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}

		// For token transfers, we don't add to the ETH amount
		entityMap[counterpartyAddress].Transactions = append(
			entityMap[counterpartyAddress].Transactions,
			transaction,
		)
	}

	// Process ERC721 transfers (NFTs)
	for _, tx := range txCollection.ERC721Txs {
		var counterpartyAddress string
		if isOutgoing {
			if !strings.EqualFold(tx.From, address) {
				continue
			}
			counterpartyAddress = tx.To
		} else {
			if !strings.EqualFold(tx.To, address) {
				continue
			}
			counterpartyAddress = tx.From
		}

		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := utils.FormatTimestamp(timestamp)

		transaction := models.Transaction{
			TxAmount:      1, // NFTs always transfer one
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, ok := entityMap[counterpartyAddress]; !ok {
			entityMap[counterpartyAddress] = &models.EntityWithTransactions{
				Address:      counterpartyAddress,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}
		entityMap[counterpartyAddress].Transactions = append(
			entityMap[counterpartyAddress].Transactions,
			transaction,
		)
	}

	// Process ERC1155 transfers
	for _, tx := range txCollection.ERC1155Txs {
		var counterpartyAddress string
		if isOutgoing {
			if !strings.EqualFold(tx.From, address) {
				continue
			}
			counterpartyAddress = tx.To
		} else {
			if !strings.EqualFold(tx.To, address) {
				continue
			}
			counterpartyAddress = tx.From
		}

		amount := utils.ConvertTokenValueWithDecimals(tx.TokenValue.String(), tx.TokenDecimal)
		timestamp := tx.TimeStamp.Time().Unix()
		dateTime := utils.FormatTimestamp(timestamp)

		transaction := models.Transaction{
			TxAmount:      amount,
			DateTime:      dateTime,
			TransactionID: tx.Hash,
		}

		if _, ok := entityMap[counterpartyAddress]; !ok {
			entityMap[counterpartyAddress] = &models.EntityWithTransactions{
				Address:      counterpartyAddress,
				Amount:       0,
				Transactions: []models.Transaction{},
			}
		}
		entityMap[counterpartyAddress].Transactions = append(
			entityMap[counterpartyAddress].Transactions,
			transaction,
		)
	}

	return entityMap
}
