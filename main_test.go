package main

import (
	"strings"
	"testing"
)

func TestFetchStockData(t *testing.T) {
	testCases := []struct {
		name         string
		stockCode    string
		expectError  bool
		expectedData string
	}{
		{
			name:         "Valid stock code",
			stockCode:    "AAPL.US",
			expectError:  false,
			expectedData: "AAPL.US quote is",
		},
		{
			name:        "Invalid stock code",
			stockCode:   "INVALID_CODE",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stockData, err := fetchStockData(tc.stockCode)

			if tc.expectError {
				if err == nil {
					t.Error("Expected an error, but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				}

				if !strings.Contains(stockData, tc.expectedData) {
					t.Errorf("Expected stock data '%s', but got '%s'", tc.expectedData, stockData)
				}
			}
		})
	}
}
