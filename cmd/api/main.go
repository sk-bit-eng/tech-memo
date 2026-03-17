// cmd/api/main.go
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"tech-memo/internal/infrastructure/api"
)

func main() {
	port := getEnv("PORT", "8080")
	dbPath := getEnv("DB_PATH", "tech_memo.db")

	handler, err := api.BuildApp(dbPath)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on http://localhost%s", addr)
	log.Printf("Database: %s", dbPath)

	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
