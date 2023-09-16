package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"os/signal"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Order Define the Order struct for database mapping
type Order struct {
	gorm.Model
	OrderID string  `json:"order_id"`
	Side    string  `json:"side"`
	Symbol  string  `json:"symbol"`
	Amount  float64 `json:"amount"`
	Price   float64 `json:"price"`
}

var (
	broker = "0.0.0.0:9092"
	topic  = "orders"
)

func main() {
	// Set up Kafka consumer configuration
	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"group.id":          "orders-consumer",
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Error creating Kafka consumer: %v", err)
	}
	defer func(consumer *kafka.Consumer) {
		err := consumer.Close()
		if err != nil {
			log.Printf("Error closing Kafka consumer: %v", err)
		}
	}(consumer)

	// Subscribe to the Kafka topic
	if err := consumer.SubscribeTopics([]string{topic}, nil); err != nil {
		log.Fatalf("Error subscribing to topic: %v", err)
	}

	// Set up the PostgresSQL database connection
	LoadEnv()
	dsn := fmt.Sprintf(
		"user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.AutoMigrate(&Order{})
	if err != nil {
		log.Fatalf("Error auto-migrating the database: %v", err)
	} // Auto-create the "orders" table

	// Create a signal channel to exit the program
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Consume and process Kafka messages
	run := true
	for run {
		select {
		case <-signals:
			run = false // Exit on Ctrl+C
		default:
			ev := consumer.Poll(1000)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				// Unmarshal the Kafka message into an Order struct
				var order Order
				if err := json.Unmarshal(e.Value, &order); err != nil {
					log.Printf("Error decoding message: %v", err)
				} else {
					// Insert the order into the database
					db.Create(&order)
					fmt.Printf("Received and inserted order: %+v\n", order)
				}
			case kafka.Error:
				log.Printf("Kafka error: %v", e)
			}
		}
	}

	fmt.Println("Consumer is shutting down...")
}
func LoadEnv() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
