package main
import(
  "fmt"
  "Ethereum-fund-flow-analysis/internal/config"
)

func main()  {
  cfg, err := config.Load()
  if err != nil {
    fmt.Printf("Error:%v\n", err)
    return
  }

  fmt.Printf("Config:%v\n", cfg)
}
