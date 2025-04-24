package utils

import (
	"math"
	"math/big"
	"time"
)

const weiToEther = 1e18

// ConvertWeiToEther converts wei value (string) to ether (float64)
func ConvertWeiToEther(weiValue string) float64 {
	valueWei := new(big.Int)
	valueWei.SetString(weiValue, 10)
	valueEther := new(big.Float).Quo(new(big.Float).SetInt(valueWei), big.NewFloat(weiToEther))
	amount, _ := valueEther.Float64()
	return amount
}

// ConvertTokenValueWithDecimals converts token value considering its decimals
func ConvertTokenValueWithDecimals(value string, tokenDecimal uint8) float64 {
	divisor := new(big.Float).SetFloat64(math.Pow10(int(tokenDecimal)))
	valueToken := new(big.Float)
	valueToken.SetString(value)
	amount, _ := new(big.Float).Quo(valueToken, divisor).Float64()
	return amount
}

// FormatTimestamp converts Unix timestamp to formatted datetime string
func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}
