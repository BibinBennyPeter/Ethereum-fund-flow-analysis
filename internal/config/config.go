package config

import (
    "errors"
    "os"
)

type Config struct {
  EtherscanAPIKey string
  EtherscanBaseURL string
}

func Load()  (*Config, error){
  apiKey := os.Getenv("ETHERSCAN_API_KEY")
  if apiKey == "" {
    return nil, errors.New("ETHERSCAN_API_KEY environment variable is required")
  }

  baseURL := os.Getenv("ETHERSCAN_BASE_URL")
  if baseURL == "" {
    baseURL = "http://api.etherscan.io/api" // Default API
  }

  return &Config{
    EtherscanAPIKey : apiKey,
    EtherscanBaseURL : baseURL,
  }, nil
}
