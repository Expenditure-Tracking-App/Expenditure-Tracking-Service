package main

import "C"
import (
	"context"
	"fmt"
	"github.com/rs/cors"
	"gopkg.in/yaml.v3"
	"log"
	"main/pkg/config"
	"main/pkg/handler"
	"main/pkg/storage"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"
	"time"
)

// --- Main Function ---

func main() {
	// Support CONFIG_PATH env var for custom config location
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
	}

	yamlFile, err := os.ReadFile(configPath)
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
	mux.HandleFunc("/api/v1/transactions/", handler.TransactionsHandler)

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

	// Use port from config if available, otherwise default to 8081
	port := strconv.Itoa(cfg.Server.Port)
	if port == "0" {
		port = "8081"
	}
	serverAddr := fmt.Sprintf(":%s", port)

	srv := &http.Server{
		Addr:         serverAddr,
		Handler:      corsMiddlewareHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Starting HTTP server on %s", serverAddr)

	log.Printf("The available endpoints:\nHealth check endpoint: localhost:%v/health\nTransactions API endpoint: localhost:%v/api/v1/transactions",
		port, port)

	// Start server in a goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down HTTP server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("HTTP server graceful shutdown error: %v", err)
	}

	log.Println("HTTP server exited")
}
