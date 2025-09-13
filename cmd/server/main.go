package main

import "C"
import (
	"fmt"
	"github.com/rs/cors"
	"gopkg.in/yaml.v3"
	"log"
	"main/pkg/config"
	"main/pkg/handler"
	"main/pkg/storage" // Assuming your storage functions are here
	"net/http"         // The core HTTP package
	"os"               // To potentially read port from environment
	"regexp"
)

// --- Main Function ---

func main() {
	yamlFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Error reading YAML file: %v", err)
	}

	var cfg config.Config
	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		log.Fatalf("Error unmarshalling YAML: %v", err)
	}

	if cfg.FeaturesConfig.SaveToDB {
		err = storage.InitDB(cfg.Database)
		if err != nil {
			log.Fatalf("Failed to initialize database: %v", err)
		}
		defer storage.CloseDB()
		log.Println("Database initialized successfully.")
	} else {
		log.Println("Database not configured. Transactions API might not function as expected if DB is required.")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.HealthCheckHandler)
	mux.HandleFunc("/api/v1/transactions", handler.TransactionsHandler)

	prefilledHandler := handler.GetPrefilledExpensesHandler(cfg.FrequentExpenses)
	mux.HandleFunc("/api/v1/prefilled-expenses", prefilledHandler)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000", "http://192.168.1.11:3000"},
		AllowOriginFunc: func(origin string) bool {
			match, _ := regexp.MatchString(`^https?://(localhost|127\.0\.0\.1|192\.168\.\d{1,3}\.\d{1,3}|10\.\d{1,3}\.\d{1,3}\.\d{1,3}|172\.(1[6-9]|2[0-9]|3[0-1])\.\d{1,3}\.\d{1,3}):\d+$`, origin)
			return match
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "Authorization"},
	})

	corsMiddlewareHandler := corsMiddleware.Handler(mux)

	port := "8081"
	serverAddr := fmt.Sprintf(":%s", port)

	log.Printf("Starting HTTP server on %s", serverAddr)

	log.Printf("The available endpoints:\nHealth check endpoint: localhost:%v/health\nTransactions API endpoint: localhost:%v/api/v1/transactions",
		port, port)

	err = http.ListenAndServe(serverAddr, corsMiddlewareHandler)
	if err != nil {
		log.Fatalf("HTTP server failed to start: %v", err)
	}
}
