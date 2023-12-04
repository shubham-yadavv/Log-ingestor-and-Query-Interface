package amqp

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/shubham-yadavv/log-ingestor/database"
	"github.com/streadway/amqp"
)

const (
	queueName = "logs"
)

var amqpChannel *amqp.Channel

func init() {
	err := InitializeAMQP()
	if err != nil {
		log.Fatalf("Failed to initialize AMQP: %v", err)
	}
}

func InitializeAMQP() error {

	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }

	amqpURL := os.Getenv("AMQP_URL")
	if amqpURL == "" {
		amqpURL = "amqp://localhost:5672"
	}
	amqpConn, err := amqp.Dial(amqpURL)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	channel, err := amqpConn.Channel()
	if err != nil {
		return fmt.Errorf("failed to open a channel: %v", err)
	}

	// Declare the queue
	_, err = channel.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare a queue: %v", err)
	}

	// // Defer the closure of the connection after obtaining the channel
	// defer func() {
	// 	if err := amqpConn.Close(); err != nil {
	// 		log.Printf("Error closing AMQP connection: %v", err)
	// 	}
	// }()

	amqpChannel = channel
	return nil
}

func GetChannel() *amqp.Channel {
	return amqpChannel
}

func EnqueueLog(log database.Log) error {
	if amqpChannel == nil {
		return fmt.Errorf("AMQP channel is nil")
	}

	logJSON, err := json.Marshal(log)
	if err != nil {
		return fmt.Errorf("error marshaling log: %v", err)
	}

	err = amqpChannel.Publish(
		"",        // exchange
		queueName, // routing key (queue name)
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        logJSON,
		},
	)

	return err
}
