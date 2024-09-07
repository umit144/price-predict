package main

import (
	"testing"
)

func TestPredictPrice(t *testing.T) {
	tests := []struct {
		name    string
		cp      *CryptoPredictor
		wantErr bool
		checkFn func(t *testing.T, got float64)
	}{
		{
			name: "Normal prediction",
			cp: &CryptoPredictor{
				Symbol:           "BTC",
				CurrentPrice:     50000,
				PercentChange24h: 5,
				PercentChange7d:  10,
			},
			wantErr: false,
			checkFn: func(t *testing.T, got float64) {
				if got <= 50000 || got >= 55000 {
					t.Errorf("Prediction %f is outside expected range", got)
				}
			},
		},
		{
			name: "Zero current price",
			cp: &CryptoPredictor{
				Symbol:           "TEST",
				CurrentPrice:     0,
				PercentChange24h: 5,
				PercentChange7d:  10,
			},
			wantErr: false,
			checkFn: func(t *testing.T, got float64) {
				if got != 0 {
					t.Errorf("Expected 0 prediction for 0 current price, got %f", got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := predictPrice(tt.cp)
			if (err != nil) != tt.wantErr {
				t.Errorf("predictPrice() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.checkFn != nil {
				tt.checkFn(t, got)
			}
		})
	}
}

func TestDisplayResults(t *testing.T) {
	cp := &CryptoPredictor{
		Symbol:           "BTC",
		CurrentPrice:     50000,
		PercentChange24h: 5,
		PercentChange7d:  10,
	}
	prediction := 52000.0

	displayResults(cp, prediction)
}
