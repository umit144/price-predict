# Cryptocurrency Price Predictor

This is a simple command-line tool that fetches the current price of a cryptocurrency and makes a basic price prediction based on recent trends.

## Features

- Fetches real-time cryptocurrency data from CoinMarketCap API
- Predicts future price based on 24-hour and 7-day trends
- Displays results in a neat table format

## Prerequisites

- Go 1.16 or higher
- CoinMarketCap API key

## Installation

1. Clone this repository:
   ```
   git clone https://github.com/umit144/price-predict
   cd price-predict
   ```

2. Set your CoinMarketCap API key:
   Open `main.go` and replace `API_KEY` with your actual API key:
   ```go
   const CMC_API_KEY = "API_KEY"
   ```

3. Build the project:
   ```
   go build
   ```

## Usage

Run the compiled binary:

```
./price-predict
```

When prompted, enter the symbol of the cryptocurrency you want to predict (e.g., BTC, ETH).

## Running Tests

To run the tests, use the following command:

```
go test
```

## Disclaimer

This tool is for educational purposes only. The predictions are based on a very simple model and should not be used for actual trading decisions. Cryptocurrency markets are highly volatile and unpredictable.

## License

This project is open source and available under the [MIT License](LICENSE).