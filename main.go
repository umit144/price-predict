package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	CMC_API_KEY_ENV = "CMC_API_KEY"
)

type CMCResponse struct {
	Data map[string]CryptoData `json:"data"`
}

type CryptoData struct {
	Quote map[string]QuoteData `json:"quote"`
}

type QuoteData struct {
	Price            float64 `json:"price"`
	PercentChange24h float64 `json:"percent_change_24h"`
	PercentChange7d  float64 `json:"percent_change_7d"`
	LastUpdated      string  `json:"last_updated"`
}

type CryptoPredictor struct {
	Symbol           string
	CurrentPrice     float64
	PercentChange24h float64
	PercentChange7d  float64
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter the cryptocurrency symbol (e.g., BTC, ETH): ")
	symbol, _ := reader.ReadString('\n')
	symbol = strings.TrimSpace(strings.ToUpper(symbol))

	cryptoPredictor := &CryptoPredictor{Symbol: symbol}

	fmt.Printf("Processing data for %s...\n", symbol)

	done := make(chan bool)
	go showLoadingAnimation(done)

	err := fetchCurrentData(cryptoPredictor)
	if err != nil {
		done <- true
		<-done
		log.Fatalf("Error fetching current data: %v", err)
	}

	prediction, err := predictPrice(cryptoPredictor)
	done <- true
	<-done

	if err != nil {
		log.Fatalf("Error predicting price: %v", err)
	}

	displayResults(cryptoPredictor, prediction)
}

func fetchCurrentData(cp *CryptoPredictor) error {
	url := fmt.Sprintf("https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=%s", cp.Symbol)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	apiKey := os.Getenv(CMC_API_KEY_ENV)
	if apiKey == "" {
		return fmt.Errorf("CMC_API_KEY environment variable is not set")
	}

	req.Header.Set("X-CMC_PRO_API_KEY", apiKey)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error accessing API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned error code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %w", err)
	}

	var cmcResp CMCResponse
	if err := json.Unmarshal(body, &cmcResp); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	cryptoData, ok := cmcResp.Data[cp.Symbol]
	if !ok {
		return fmt.Errorf("no data found for '%s'", cp.Symbol)
	}

	usdQuote, ok := cryptoData.Quote["USD"]
	if !ok {
		return fmt.Errorf("no USD price found for '%s'", cp.Symbol)
	}

	cp.CurrentPrice = usdQuote.Price
	cp.PercentChange24h = usdQuote.PercentChange24h
	cp.PercentChange7d = usdQuote.PercentChange7d

	return nil
}

func predictPrice(cp *CryptoPredictor) (float64, error) {
	shortTermTrend := cp.PercentChange24h / 100
	longTermTrend := cp.PercentChange7d / 100

	weightedTrend := (shortTermTrend * 0.7) + (longTermTrend * 0.3)

	randomFactor := (rand.Float64() - 0.5) * 0.02

	prediction := cp.CurrentPrice * (1 + weightedTrend + randomFactor)

	return prediction, nil
}

func displayResults(cp *CryptoPredictor, prediction float64) {
	change := (prediction - cp.CurrentPrice) / cp.CurrentPrice * 100

	headers := []string{"Metric", "Value"}
	data := [][]string{
		{"Symbol", cp.Symbol},
		{"Current Price", fmt.Sprintf("$%.2f", cp.CurrentPrice)},
		{"Predicted Price", fmt.Sprintf("$%.2f", prediction)},
		{"Predicted Change", fmt.Sprintf("%.2f%%", change)},
		{"24h Change", fmt.Sprintf("%.2f%%", cp.PercentChange24h)},
		{"7d Change", fmt.Sprintf("%.2f%%", cp.PercentChange7d)},
	}

	printTable(headers, data)
}

func printTable(headers []string, data [][]string) {
	colWidths := make([]int, len(headers))
	for i, header := range headers {
		colWidths[i] = len(header)
	}
	for _, row := range data {
		for i, cell := range row {
			if len(cell) > colWidths[i] {
				colWidths[i] = len(cell)
			}
		}
	}

	printTableRow(headers, colWidths)
	printTableSeparator(colWidths)

	for _, row := range data {
		printTableRow(row, colWidths)
	}
}

func printTableRow(row []string, colWidths []int) {
	fmt.Print("|")
	for i, cell := range row {
		fmt.Printf(" %-*s |", colWidths[i], cell)
	}
	fmt.Println()
}

func printTableSeparator(colWidths []int) {
	fmt.Print("+")
	for _, width := range colWidths {
		fmt.Print(strings.Repeat("-", width+2))
		fmt.Print("+")
	}
	fmt.Println()
}

func showLoadingAnimation(done chan bool) {
	frames := []string{"|", "/", "-", "\\"}
	for {
		select {
		case <-done:
			fmt.Print("\r")
			done <- true
			return
		default:
			for _, frame := range frames {
				fmt.Printf("\rProcessing %s", frame)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
