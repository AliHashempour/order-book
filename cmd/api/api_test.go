package main

import (
	"book-orders/pkg/model"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetOrdersRoute(t *testing.T) {
	// Create a test router
	r := gin.Default()
	r.GET("/", nil)
	r.POST("/orders", nil)
	// Define test cases for the /get_orders route
	testCases := []struct {
		method       string
		path         string
		expectedCode int
	}{
		{"POST", "/orders", http.StatusOK},
		{"GET", "/", http.StatusOK},
	}
	// Iterate through the test cases
	for _, tc := range testCases {
		// Create a test request
		_, err := http.NewRequest(tc.method, tc.path, nil)
		if err != nil {
			t.Fatalf("Error creating test request: %v", err)
		}
	}
}

func TestLoadEnv(t *testing.T) {
	// Call the LoadEnv function and check for errors
	err := LoadEnv()
	assert.NoErrorf(t, err, "%s", err)
}

func TestFormatOrders(t *testing.T) {
	// Create a sample list of orders
	orders := []model.Order{
		{OrderID: "1", Side: "buy", Symbol: "AAPL", Amount: 10, Price: 150.0},
		{OrderID: "2", Side: "sell", Symbol: "AAPL", Amount: 5, Price: 2700.0},
	}
	// Call the formatOrders function to format the orders
	formatted := formatOrders(orders)
	// Check if the formatted result matches the expected format
	expected := [][]string{
		{"150.00", "10.00"},
		{"2700.00", "5.00"},
	}
	assert.Equal(t, expected, formatted, "Formatted orders should match the expected format")
}
