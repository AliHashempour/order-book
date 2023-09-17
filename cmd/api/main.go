package main

import (
	"book-orders/pkg/model"
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

func main() {
	// Set up the PostgresSQL database connection
	dsn, err := setUpDatabaseConfig()
	if err != nil {
		log.Fatalf("Error setting up database configuration: %v", err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Start the Gin server
	r := setupRouter(db)
	port := 9090
	fmt.Printf("Starting server on port %d...\n", port)
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatalf("Error starting server: %v", err)
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
		var buyOrders, sellOrders []model.Order

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
func formatOrders(orders []model.Order) [][]string {
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
	r.POST("/orders", getOrdersHandler(db))
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
