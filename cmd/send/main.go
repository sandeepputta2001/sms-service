package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"twilio-e2e/internal/twilio"
)

func main() {
	accountSID := mustEnv("TWILIO_ACCOUNT_SID")
	authToken := mustEnv("TWILIO_AUTH_TOKEN")
	from := mustEnv("TWILIO_FROM")
	to := mustEnv("TWILIO_TO")
	body := getEnv("TWILIO_BODY", "Hello from Go!")
	statusCallbackURL := os.Getenv("TWILIO_STATUS_CALLBACK_URL") // optional

	client := twilio.NewClient(accountSID, authToken)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	resp, err := client.SendSMS(ctx, from, to, body, statusCallbackURL)
	if err != nil {
		log.Fatalf("send failed: %v", err)
	}

	pretty, _ := json.MarshalIndent(resp, "", "  ")
	log.Printf("twilio response:\n%s", string(pretty))
}

func mustEnv(key string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	log.Fatalf("missing env var: %s", key)
	return ""
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

