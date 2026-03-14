package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/fibrchat/server/pkg/server"
	"github.com/fibrchat/worker/pkg/worker"
)

func main() {
	port := envInt("PORT", 4222)
	domain := envOr("DOMAIN", "localhost")
	workerPassword := randomPassword()

	srv, err := server.Start(server.Options{
		Domain:         domain,
		Port:           port,
		WorkerPassword: workerPassword,
	})
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	wrk, err := worker.Start(worker.Options{
		Domain:          domain,
		WorkerPassword:  workerPassword,
		InProcessServer: srv,
	})
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	fmt.Printf("Server is running on ws://%s:%d\n", domain, port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nShutting down...")

	wrk.Stop()
	srv.Stop()
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func envInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid value for %s: %v", key, err)
	}
	return n
}

func randomPassword() string {
	b := make([]byte, 256)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("failed to generate random password: %v", err)
	}

	return base64.RawURLEncoding.EncodeToString(b)
}
