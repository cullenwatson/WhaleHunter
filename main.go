package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"tw-scanner/tradingview"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	username := os.Getenv("TRADINGVIEW_USERNAME")
	password := os.Getenv("TRADINGVIEW_PASSWORD")

	if username == "" || password == "" {
		log.Fatal("Missing required environment variables TRADINGVIEW_USERNAME and/or TRADINGVIEW_PASSWORD")
	}

	creds := tradingview.Credentials{
		Username: username,
		Password: password,
	}

	authToken, err := tradingview.SignIn(creds)
	if err != nil {
		log.Fatalf("Error signing in: %v\n", err)
	}

	fmt.Println("Successfully signed in to TradingView!")
	fmt.Printf("Auth Token: %s\n", authToken)
}
