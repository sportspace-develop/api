package main

import (
	"context"
	"fmt"
	"log"

	"sport-space/internal/adapter/api/rest"
	"sport-space/internal/adapter/logger"
	"sport-space/internal/adapter/sender"
	"sport-space/internal/adapter/storage"
	"sport-space/internal/core/config"
	"sport-space/internal/core/sportspace"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	lgr, err := logger.New()
	if err != nil {
		return fmt.Errorf("failed initialize logger: %w", err)
	}

	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("failed initialize config: %w", err)
	}

	store, err := storage.New(ctx, cfg.Store)
	if err != nil {
		return fmt.Errorf("failed initialize storage: %w", err)
	}

	sender, err := sender.New(cfg.Sender, sender.SetLogger(lgr))
	if err != nil {
		return fmt.Errorf("failed initialize sender: %w", err)
	}

	sspace, err := sportspace.New(store, sender, sportspace.SetLogger(lgr))
	if err != nil {
		return fmt.Errorf("failed initialize sportspace service: %w", err)
	}

	server, err := rest.New(sspace, rest.SetLogger(lgr), rest.SetAddress(cfg.Address), rest.SetSecretKey(cfg.SecretKey))
	if err != nil {
		return fmt.Errorf("failed initalize rest api: %w", err)
	}

	if err := server.Run(); err != nil {
		return err
	}
	return err
}
