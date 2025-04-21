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

// GetNormalTransactions fetches normal transactions for the given address
func (c *Client) GetNormalTransactions(address string) ([]models.NormalTx, error) {
	endpoint := fmt.Sprintf("%s?chainid=1&module=account&action=txlist&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
		c.baseURL, address, c.apiKey)
	
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
func (c *Client) GetInternalTransactions(address string) ([]models.InternalTx, error) {
    endpoint := fmt.Sprintf("%s?chainid=1&module=account&action=txlistinternal&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
        c.baseURL, address, c.apiKey)
    
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
func (c *Client) GetERC20Transfers(address string) ([]models.ERC20Transfer, error) {
    endpoint := fmt.Sprintf("%s?chainid=1&module=account&action=tokentx&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
        c.baseURL, address, c.apiKey)
    
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
func (c *Client) GetERC721Transfers(address string) ([]models.ERC721Transfer, error) {
    endpoint := fmt.Sprintf("%s?chainid=1&module=account&action=tokennfttx&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
        c.baseURL, address, c.apiKey)
    
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
func (c *Client) GetERC1155Transfers(address string) ([]models.ERC1155Transfer, error) {
    endpoint := fmt.Sprintf("%s?chainid=1&module=account&action=tokentx&address=%s&startblock=0&endblock=99999999&sort=asc&apikey=%s",
        c.baseURL, address, c.apiKey)
    
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
