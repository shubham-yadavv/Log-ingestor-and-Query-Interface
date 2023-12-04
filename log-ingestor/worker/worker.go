package worker

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/shubham-yadavv/log-ingestor/amqp"
	"github.com/shubham-yadavv/log-ingestor/database"
)

func WorkerMain() {
	channel := amqp.GetChannel()
	if channel == nil {
		log.Fatal("AMQP channel is nil")
		return
	}

	msgs, err := channel.Consume("logs", "", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var logData database.Log

			if err := json.Unmarshal(d.Body, &logData); err != nil {
				log.Println("Error decoding JSON:", err)
				continue
			}

			if err := database.DB.Create(&logData).Error; err != nil {
				log.Println("Error creating log record:", err)
				continue
			}

			fmt.Printf("Processing log: %+v\n", logData)
		}
	}()

	fmt.Println("Worker is waiting for logs. To exit press CTRL+C")
	<-forever
}
