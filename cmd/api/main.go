package main

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Order represents an order structure for database mapping
type Order struct {
	gorm.Model
	OrderID string  `json:"order_id"`
	Side    string  `json:"side"`
	Symbol  string  `json:"symbol"`
	Amount  float64 `json:"amount"`
	Price   float64 `json:"price"`
}

func main() {
	// Set up the PostgresSQL database connection
	dsn, envErr := setUpDatabaseConfig()
	if envErr != nil {
		log.Fatalf("Error setting up database configuration: %v", envErr)
	}

	db, dbErr := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if dbErr != nil {
		log.Fatalf("Error connecting to the database: %v", dbErr)
	}

	// Start the Gin server
	r := setupRouter(db)
	port := 9090
	fmt.Printf("Starting server on port %d...\n", port)
	if serverErr := r.Run(fmt.Sprintf(":%d", port)); serverErr != nil {
		log.Fatalf("Error starting server: %v", serverErr)
	}
}

// getHomeMessage returns a handler function for the / endpoint
func getHomeMessage(c *gin.Context) {
	message := gin.H{
		"msg": "home page",
	}
	c.IndentedJSON(http.StatusOK, message)
}

// getOrdersHandler returns a handler function for the /get_orders endpoint
func getOrdersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		symbol := c.PostForm("symbol")
		limitStr := c.DefaultPostForm("limit", "100")
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
			return
		}
		// Enforce a maximum limit of 1000
		if limit > 1000 {
			limit = 1000
		}
		var buyOrders, sellOrders []Order

		// Query buy and sell orders based on the symbol
		db.Table("orders").Where("side = ? AND symbol = ?", "buy", symbol).
			Order("price desc").Limit(limit).Find(&buyOrders)

		db.Table("orders").Where("side = ? AND symbol = ?", "sell", symbol).
			Order("price asc").Limit(limit).Find(&sellOrders)

		// Prepare the response JSON
		response := gin.H{
			"bids": formatOrders(buyOrders),
			"asks": formatOrders(sellOrders),
		}
		c.IndentedJSON(http.StatusOK, response)
	}
}

// formatOrders formats a list of orders into the expected JSON format.
func formatOrders(orders []Order) [][]string {
	result := make([][]string, len(orders))
	for i, order := range orders {
		result[i] = []string{fmt.Sprintf("%.2f", order.Price), fmt.Sprintf("%.2f", order.Amount)}
	}
	return result
}

// LoadEnv loads the environment variables from the .env file
func LoadEnv() error {
	if err := godotenv.Load("../../.env"); err != nil {
		return err
	}
	return nil
}

// setupRouter sets up a Gin router with the routes we need
func setupRouter(db *gorm.DB) *gin.Engine {
	// Set up the Gin router
	r := gin.Default()
	// Define a GET simple endpoint
	r.GET("/", getHomeMessage)
	// Define a POST endpoint to get buy and sell orders with a limit
	r.POST("/get_orders", getOrdersHandler(db))
	return r
}

// setUpDatabaseConfig sets up the PostgresSQL database connection using environment variables
func setUpDatabaseConfig() (string, error) {
	// Set up the PostgresSQL database connection using environment variables
	err := LoadEnv()
	if err != nil {
		return "", errors.New("error loading .env file")
	}
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	return dsn, nil
}
