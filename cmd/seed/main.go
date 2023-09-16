package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
)

// Order Define the order struct for JSON marshaling
type Order struct {
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

	// Set up Kafka producer configuration
	producer, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": broker})
	if err != nil {
		log.Fatalf("Error creating Kafka producer: %v", err)
	}
	defer producer.Close()

	// Sample list of orders
	orders := []Order{
		{"1", "buy", "AAPL", 10, 150.0},
		{"2", "sell", "AAPL", 5, 2700.0},
		{"3", "buy", "TSLA", 15, 700.0},
		{"4", "sell", "TSLA", 20, 250.0},
		{"5", "buy", "AAPL", 10, 3200.0},
		{"6", "sell", "AAPL", 10, 300.0},
		{"7", "buy", "TSLA", 10, 500.0},
		{"8", "sell", "TSLA", 10, 500.0},
		{"9", "buy", "AAPL", 10, 200.0},
		{"10", "sell", "AAPL", 10, 50.0},
	}

	for _, order := range orders {
		orderJSON, err := json.Marshal(order)
		if err != nil {
			log.Printf("Error marshaling order: %v", err)
			continue
		}

		message := &kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(orderJSON),
		}

		err = producer.Produce(message, nil)
		if err != nil {
			log.Printf("Error sending message to Kafka: %v", err)
		} else {
			fmt.Println("Message sent to Kafka")
		}
	}

	fmt.Println("Data seeding completed. Exiting...")

}
