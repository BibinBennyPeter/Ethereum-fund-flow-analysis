package models

import (
	"time"
)

// Transaction represents a single transaction in the response
type Transaction struct {
	TxAmount      float64   `json:"tx_amount"`
	DateTime      string    `json:"date_time"`
	TransactionID string    `json:"transaction_id"`
}

// Beneficiary represents a single beneficiary with all related transactions
type Beneficiary struct {
	Address      string        `json:"beneficiary_address"`
	Amount       float64       `json:"amount"`
	Transactions []Transaction `json:"transactions"`
}

// BeneficiaryResponse is the complete response for the /beneficiary endpoint
type BeneficiaryResponse struct {
	Message string        `json:"message"`
	Data    []Beneficiary `json:"data"`
}

// Payer represents a single payer with all related transactions
type Payer struct {
	Address      string        `json:"payer_address"`
	Amount       float64       `json:"amount"`
	Transactions []Transaction `json:"transactions"`
}

// PayerResponse is the complete response for the /payer endpoint
type PayerResponse struct {
	Message string  `json:"message"`
	Data    []Payer `json:"data"`
}


// NormalTx holds info from normal tx query
type NormalTx struct {
	BlockNumber       int     `json:"blockNumber,string"`
	TimeStamp         Time    `json:"timeStamp"`
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce,string"`
	BlockHash         string  `json:"blockHash"`
	TransactionIndex  int     `json:"transactionIndex,string"`
	From              string  `json:"from"`
	To                string  `json:"to"`
	Value             *BigInt `json:"value"`
	Gas               int     `json:"gas,string"`
	GasPrice          *BigInt `json:"gasPrice"`
	IsError           int     `json:"isError,string"`
	TxReceiptStatus   string  `json:"txreceipt_status"`
	Input             string  `json:"input"`
	ContractAddress   string  `json:"contractAddress"`
	CumulativeGasUsed int     `json:"cumulativeGasUsed,string"`
	GasUsed           int     `json:"gasUsed,string"`
	Confirmations     int     `json:"confirmations,string"`
	FunctionName      string  `json:"functionName"`
	MethodId          string  `json:"methodId"`
}

// InternalTx holds info from internal tx query
type InternalTx struct {
	BlockNumber     int     `json:"blockNumber,string"`
	TimeStamp       Time    `json:"timeStamp"`
	Hash            string  `json:"hash"`
	From            string  `json:"from"`
	To              string  `json:"to"`
	Value           *BigInt `json:"value"`
	ContractAddress string  `json:"contractAddress"`
	Input           string  `json:"input"`
	Type            string  `json:"type"`
	Gas             int     `json:"gas,string"`
	GasUsed         int     `json:"gasUsed,string"`
	TraceID         string  `json:"traceId"`
	IsError         int     `json:"isError,string"`
	ErrCode         string  `json:"errCode"`
}

// ERC20Transfer holds info from ERC20 token transfer event query
type ERC20Transfer struct {
	BlockNumber       int     `json:"blockNumber,string"`
	TimeStamp         Time    `json:"timeStamp"`
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce,string"`
	BlockHash         string  `json:"blockHash"`
	From              string  `json:"from"`
	ContractAddress   string  `json:"contractAddress"`
	To                string  `json:"to"`
	Value             *BigInt `json:"value"`
	TokenName         string  `json:"tokenName"`
	TokenSymbol       string  `json:"tokenSymbol"`
	TokenDecimal      uint8   `json:"tokenDecimal,string"`
	TransactionIndex  int     `json:"transactionIndex,string"`
	Gas               int     `json:"gas,string"`
	GasPrice          *BigInt `json:"gasPrice"`
	GasUsed           int     `json:"gasUsed,string"`
	CumulativeGasUsed int     `json:"cumulativeGasUsed,string"`
	Input             string  `json:"input"`
	Confirmations     int     `json:"confirmations,string"`
}

// ERC721Transfer holds info from ERC721 token transfer event query
type ERC721Transfer struct {
	BlockNumber       int     `json:"blockNumber,string"`
	TimeStamp         Time    `json:"timeStamp"`
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce,string"`
	BlockHash         string  `json:"blockHash"`
	From              string  `json:"from"`
	ContractAddress   string  `json:"contractAddress"`
	To                string  `json:"to"`
	TokenID           *BigInt `json:"tokenID"`
	TokenName         string  `json:"tokenName"`
	TokenSymbol       string  `json:"tokenSymbol"`
	TokenDecimal      uint8   `json:"tokenDecimal,string"`
	TransactionIndex  int     `json:"transactionIndex,string"`
	Gas               int     `json:"gas,string"`
	GasPrice          *BigInt `json:"gasPrice"`
	GasUsed           int     `json:"gasUsed,string"`
	CumulativeGasUsed int     `json:"cumulativeGasUsed,string"`
	Input             string  `json:"input"`
	Confirmations     int     `json:"confirmations,string"`
}

// ERC1155Transfer holds info from ERC1155 token transfer event query
type ERC1155Transfer struct {
	BlockNumber       int     `json:"blockNumber,string"`
	TimeStamp         Time    `json:"timeStamp"`
	Hash              string  `json:"hash"`
	Nonce             int     `json:"nonce,string"`
	BlockHash         string  `json:"blockHash"`
	From              string  `json:"from"`
	ContractAddress   string  `json:"contractAddress"`
	To                string  `json:"to"`
	TokenID           *BigInt `json:"tokenID"`
	TokenName         string  `json:"tokenName"`
	TokenSymbol       string  `json:"tokenSymbol"`
	TokenDecimal      uint8   `json:"tokenDecimal,string"`
	TokenValue        uint8   `json:"tokenValue,string"`
	TransactionIndex  int     `json:"transactionIndex,string"`
	Gas               int     `json:"gas,string"`
	GasPrice          *BigInt `json:"gasPrice"`
	GasUsed           int     `json:"gasUsed,string"`
	CumulativeGasUsed int     `json:"cumulativeGasUsed,string"`
	Input             string  `json:"input"`
	Confirmations     int     `json:"confirmations,string"`
}


// EtherscanResponse is the generic response structure from Etherscan API
type EtherscanResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
