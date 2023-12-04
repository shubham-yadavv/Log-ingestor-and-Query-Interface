package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/shubham-yadavv/log-ingestor/amqp"
	"github.com/shubham-yadavv/log-ingestor/database"
	"github.com/shubham-yadavv/log-ingestor/handlers"
	"github.com/shubham-yadavv/log-ingestor/worker"
)

const serverAddr = ":3000"

func main() {
	database.InitializeDB()

	amqp.InitializeAMQP()


	http.HandleFunc("/", handlers.HandleLogIngestion)
	http.HandleFunc("/search", handlers.SearchLogsHandler)
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello HEE")
	})


	server := &http.Server{Addr: serverAddr}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		fmt.Println("\nShutting down gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server shutdown error:", err)
		}
	}()

	go worker.WorkerMain()

	fmt.Printf("Log Ingestor listening on %s...\n", serverAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server startup error:", err)
	}
}
