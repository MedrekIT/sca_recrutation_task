package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/MedrekIT/sca-recrutation-task/server/internal/api"
	"github.com/MedrekIT/sca-recrutation-task/server/internal/database"
	"github.com/joho/godotenv"
	"github.com/rs/cors"

	_ "github.com/lib/pq"
)

func main() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	godotenv.Load()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}

	dbPath := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbPath)
	if err != nil {
		log.Fatalf("Error while opening database connection - %v\n", err)
	}
	defer db.Close()
	db.SetMaxOpenConns(20)
	db.SetConnMaxIdleTime(30 * time.Second)

	cfg := api.ApiConfig{
		Port: fmt.Sprintf(":%s", serverPort),
		Db:   database.New(db),
	}

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:5173"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	handler := c.Handler(api.Routes(cfg))
	server := &http.Server{
		Addr:    cfg.Port,
		Handler: handler,
	}

	go func() {
		log.Printf("Server running and listening on port %s\n", cfg.Port[1:])
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Printf("Error - %v\n", err)
			}
		}
	}()

	<-ch

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	log.Println("Shutting down...")
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Error while closing server - %v\n", err)
	}
	log.Println("Server closed")
}
