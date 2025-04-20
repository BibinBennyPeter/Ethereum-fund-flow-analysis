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
