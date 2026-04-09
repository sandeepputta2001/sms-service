package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"twilio-e2e/internal/twiml"
)

func main() {
	port := getEnv("PORT", "8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = w.Write([]byte("ok\n"))
	})

	mux.HandleFunc("/sms/incoming", incomingSMSHandler)
	mux.HandleFunc("/status", statusHandler)

	srv := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	log.Printf("server listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func incomingSMSHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	from := r.Form.Get("From")
	body := r.Form.Get("Body")
	to := r.Form.Get("To")

	log.Printf("inbound sms: from=%q to=%q body=%q", from, to, body)

	reply := fmt.Sprintf("Thanks! You said: %s . Ela unnav macha , hw is it gng. Eppidi irikinge...... ", body)
	w.Header().Set("Content-Type", "application/xml; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(twiml.MessageResponse(reply)))
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad form", http.StatusBadRequest)
		return
	}

	// Twilio sends delivery status updates as form-encoded fields.
	messageSid := r.Form.Get("MessageSid")
	messageStatus := r.Form.Get("MessageStatus")
	to := r.Form.Get("To")
	from := r.Form.Get("From")

	log.Printf("status update: messageSid=%q messageStatus=%q from=%q to=%q", messageSid, messageStatus, from, to)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok\n"))
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
