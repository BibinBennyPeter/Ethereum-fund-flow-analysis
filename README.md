# Ethereum Fund Flow Analysis

_Ethereum Fund Flow Analysis_ is a Go-based project that leverages the Etherscan API to analyze how funds move through the Ethereum blockchain. It provides both **inflow (payer)** and **outflow (beneficiary)** analyses for any given address, tracing direct transfers, smart-contract interactions, and token movements.

## Key Features

- **Beneficiary (/beneficiary)**: Identify final recipients (outflow) of funds sent by a target address.
- **Payer (/payer)**: Determine sources (inflow) of funds received by a target address.
- **Flexible Fetch Filtering**: Use Etherscan query parameters (block range, pagination, sort order) to limit which transactions are fetched:
  - `address`, `sblock`, `eblock`, `page`, `offset`, `sort` (asc|desc)
  - Example: `?address=0x123...&start_block=12000000&end_block=12001000&page=1&offset=50&sort=asc`
- **Custom Post-Fetch Filters**: Further refine results in‑memory by transaction amount, zero‑value inclusion, sorting, and limit:
  - `min` (min amount), `max` (max amount), `limit` (max results), `with_zero_txs` (true|false)
  - Example: `?address=0x123...&min=0.05&max=1.0&limit=20&with_zero_txs=false`
- **Concurrent Fetching**: Parallel calls to Etherscan for normal, internal, ERC‑20, ERC‑721, and ERC‑1155 transactions maximize throughput.
- **Arkham Intel Alignment**: Outflow &gt; Beneficiary, Inflow &gt; Payer (following Arkham Intel Tracer terminology).

## API Endpoints

| Method | Path               | Description                                   |
|--------|--------------------|-----------------------------------------------|
| GET    | `/beneficiary`     | Returns outflow analysis (beneficiaries).     |
| GET    | `/payer`           | Returns inflow analysis (payers).             |

**Common Query Parameters**:
```
address        (string, required)    // target Ethereum address
sblock         (int64, optional)     // block number to start searching for transactions, default 0
eblock         (int64, optional)     // block number to stop searching for transactions, default -1 (no limit)
page           (int,   optional)     // default 1
offset         (int,   optional)     // default 100
sort           (string,optional)     // "asc" or "desc", default "desc"
min            (float, optional)     // minimum tx amount, default 0
max            (float, optional)     // maximum tx amount, default -1 (no limit)
limit          (int,   optional)     // max number of final results, default 100
with_zero_txs  (bool,  optional)     // include zero-amount entries, default true
```

**Example Request**:
```
GET /beneficiary?address=0x8C8D7C46219D9205f056f28fee5950aD564d7465&sblock=21100000&eblock=22100000&min=0.1
```

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/BibinBennyPeter/Ethereum-fund-flow-analysis.git
   cd Ethereum-fund-flow-analysis
   ```

2. **Set your Etherscan API key**:
   ```bash
   export ETHERSCAN_API_KEY=YourApiKeyHere
   ```

3. **Install dependencies & build**:
   ```bash
   go mod tidy
   go build -o ethereum-fund-analysis ./cmd/api/main.go
   ```

4. **Run the server** (default listens on `:8080`):
   ```bash
   ./ethereum-fund-analysis
   ```

## Arkham Intel Context

This project follows Arkham Intel Tracer’s definition:
- **Outflow (Beneficiary)**: Addresses receiving funds from the target (maps to `/beneficiary`).
- **Inflow (Payer)**: Addresses sending funds to the target (maps to `/payer`).

You can visualize flows using Arkham’s tracer:
```
https://intel.arkm.com/tracer?address=<target_address>
```

*Happy fund‐flow analyzing!*

