package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/shubham-yadavv/log-ingestor/amqp"
	"github.com/shubham-yadavv/log-ingestor/database"
)

func HandleLogIngestion(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusBadRequest)
		return
	}

	var logData database.Log
	err = json.Unmarshal(body, &logData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	err = amqp.EnqueueLog(logData)
	if err != nil {
		log.Printf("Error enqueuing log: %v", err)
		http.Error(w, "Error enqueuing log", http.StatusInternalServerError)
		return
	}

	fmt.Println("Log enqueued successfully.")
	w.Write([]byte("Log enqueued successfully."))
}
