// Package main คือ entrypoint ของ backend service
// รันด้วย: go run ./cmd/server
package main

import (
	"log"

	"github.com/finance-boardgame/backend/internal/api"
	"github.com/finance-boardgame/backend/internal/config"
	"github.com/finance-boardgame/engine"
)

func main() {
	cfg := config.Load()
	log.Printf("finance-boardgame backend starting | engine=v%s | addr=%s", engine.Version, cfg.HTTPAddr)

	srv := api.NewServer(cfg)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
