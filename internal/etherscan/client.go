package etherscan

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"Ethereum-fund-flow-analysis/internal/models"
)

// Client handles interactions with the Etherscan API
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Etherscan API client
func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 15 * time.Second},
	}
}

// EtherscanRequestParams contains all possible parameters for Etherscan API requests
type EtherscanRequestParams struct {
	Address         string
	ContractAddress string
	Page            int
	Offset          int
	StartBlock      int64
	EndBlock        int64
	Sort            string // "asc" or "desc"
}

// GetNormalTransactions fetches normal transactions for the given address
func (c *Client) GetNormalTransactions(params EtherscanRequestParams) ([]models.NormalTx, error) {
	endpoint := c.buildEndpoint("txlist", params)

	var response models.EtherscanResponse
	response.Result = &[]models.NormalTx{}

	if err := c.makeRequest(endpoint, &response); err != nil {
		return nil, err
	}

	result, ok := response.Result.(*[]models.NormalTx)
	if !ok {
		return nil, fmt.Errorf("failed to parse normal transactions response")
	}

	return *result, nil
}

// GetInternalTransactions fetches internal transactions for the given address
func (c *Client) GetInternalTransactions(params EtherscanRequestParams) ([]models.InternalTx, error) {
	endpoint := c.buildEndpoint("txlistinternal", params)

	var response models.EtherscanResponse
	response.Result = &[]models.InternalTx{}

	if err := c.makeRequest(endpoint, &response); err != nil {
		return nil, err
	}

	result, ok := response.Result.(*[]models.InternalTx)
	if !ok {
		return nil, fmt.Errorf("failed to parse internal transactions response")
	}

	return *result, nil
}

// GetERC20Transfers fetches ERC-20 token transfers for the given address
func (c *Client) GetERC20Transfers(params EtherscanRequestParams) ([]models.ERC20Transfer, error) {
	endpoint := c.buildEndpoint("tokentx", params)

	var response models.EtherscanResponse
	response.Result = &[]models.ERC20Transfer{}

	if err := c.makeRequest(endpoint, &response); err != nil {
		return nil, err
	}

	result, ok := response.Result.(*[]models.ERC20Transfer)
	if !ok {
		return nil, fmt.Errorf("failed to parse token transfers response")
	}

	return *result, nil
}

// GetERC721Transfers fetches ERC-721 token transfers for the given address
func (c *Client) GetERC721Transfers(params EtherscanRequestParams) ([]models.ERC721Transfer, error) {
	endpoint := c.buildEndpoint("tokennfttx", params)

	var response models.EtherscanResponse
	response.Result = &[]models.ERC721Transfer{}

	if err := c.makeRequest(endpoint, &response); err != nil {
		return nil, err
	}

	result, ok := response.Result.(*[]models.ERC721Transfer)
	if !ok {
		return nil, fmt.Errorf("failed to parse token transfers response")
	}

	return *result, nil
}

// GetERC1155Transfers fetches ERC-1155 token transfers for the given address
func (c *Client) GetERC1155Transfers(params EtherscanRequestParams) ([]models.ERC1155Transfer, error) {
	endpoint := c.buildEndpoint("token1155tx", params)

	var response models.EtherscanResponse
	response.Result = &[]models.ERC1155Transfer{}

	if err := c.makeRequest(endpoint, &response); err != nil {
		return nil, err
	}

	result, ok := response.Result.(*[]models.ERC1155Transfer)
	if !ok {
		return nil, fmt.Errorf("failed to parse token transfers response")
	}

	return *result, nil
}

// buildEndpoint constructs an Etherscan API endpoint with the provided parameters
func (c *Client) buildEndpoint(action string, params EtherscanRequestParams) string {
	url := fmt.Sprintf("%s?chainid=1&module=account&action=%s&address=%s",
		c.baseURL, action, params.Address)

	// Add optional contract address if specified (for token transfers)
	if params.ContractAddress != "" {
		url += fmt.Sprintf("&contractaddress=%s", params.ContractAddress)
	}

	// Add block range
	startBlock := params.StartBlock
	if startBlock < 0 {
		startBlock = 0
	}
	url += fmt.Sprintf("&startblock=%d", startBlock)

	endBlock := params.EndBlock
	if endBlock < 0 {
		endBlock = 99999999 // Default high value
	}
	url += fmt.Sprintf("&endblock=%d", endBlock)

	// Add pagination parameters if specified
	if params.Page > 0 {
		url += fmt.Sprintf("&page=%d", params.Page)
	}

	if params.Offset > 0 {
		url += fmt.Sprintf("&offset=%d", params.Offset)
	}

	// Add sorting parameter
	sort := params.Sort
	if sort == "" {
		sort = "asc" // Default sort order
	}
	url += fmt.Sprintf("&sort=%s", sort)

	// Add API key
	url += fmt.Sprintf("&apikey=%s", c.apiKey)

	return url
}

// makeRequest makes an HTTP request to the Etherscan API
func (c *Client) makeRequest(endpoint string, v interface{}) error {
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("error unmarshaling response: %w", err)
	}

	return nil
}
