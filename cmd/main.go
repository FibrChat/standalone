package main

import (
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
	workerPassword := envOr("WORKER_PASSWORD", "simplechat-worker")
	remotePassword := envOr("REMOTE_PASSWORD", "simplechat-remote")

	ns, err := server.Start(server.Options{
		Domain:         domain,
		Port:           port,
		WorkerPassword: workerPassword,
		RemotePassword: remotePassword,
	})
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	// Use in-process later...
	serverURL := fmt.Sprintf("nats://127.0.0.1:%d", port+1)

	wrk, err := worker.Start(worker.Options{
		Domain:         domain,
		ServerURL:      serverURL,
		WorkerPassword: workerPassword,
		RemotePassword: remotePassword,
	})
	if err != nil {
		log.Fatalf("Failed to start worker: %v", err)
	}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	fmt.Println("\nShutting down...")

	wrk.Shutdown()
	ns.Stop()
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
