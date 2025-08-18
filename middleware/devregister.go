package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// StartDevMode starts the Gin server on a random port, registers with the API service, and handles graceful shutdown.
func StartDevMode(port int, appName string) {
	// Register with API service after server is up
	go func() {
		time.Sleep(200 * time.Millisecond) // Give server a moment to start
		registerDevService(appName, port)
	}()

	// Handle graceful shutdown and unregister
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-quit
		unregisterDevService(appName)
		os.Exit(0)
	}()
}

func registerDevService(appName string, port int) {
	url := fmt.Sprintf("http://localhost:%d", port)
	payload := map[string]string{
		"appName": appName,
		"url":     url,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:8080/dev/register", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to register dev service: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Dev service registration failed: %s", resp.Status)
	}
	log.Printf("Registered dev service: %s at %s", appName, url)
}

func unregisterDevService(appName string) {
	payload := map[string]string{
		"appName": appName,
	}
	body, _ := json.Marshal(payload)
	resp, err := http.Post("http://localhost:8080/dev/unregister", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Failed to unregister dev service: %v", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Dev service unregistration failed: %s", resp.Status)
	}
	log.Printf("Unregistered dev service: %s", appName)
}
