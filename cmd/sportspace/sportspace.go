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

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Println("Build verson: " + BuildData(buildVersion))
	fmt.Println("Build date: " + BuildData(buildDate))
	fmt.Println("Build commit: " + BuildData(buildCommit))
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.Init()
	if err != nil {
		return fmt.Errorf("failed initialize config: %w", err)
	}

	lgr, err := logger.New(cfg.LogLevel)
	if err != nil {
		return fmt.Errorf("failed initialize logger: %w", err)
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

	server, err := rest.New(sspace,
		rest.SetLogger(lgr),
		rest.SetAddress(cfg.Address),
		rest.SetSecretKey(cfg.SecretKey),
		rest.SetUploadPath(cfg.UploadPath),
		rest.SetBaseURL(cfg.BaseURL),
	)
	if err != nil {
		return fmt.Errorf("failed initalize rest api: %w", err)
	}

	if err := server.Run(); err != nil {
		return err
	}
	return err
}
func BuildData(data string) string {
	if data != "" {
		return data
	}
	return "N/A"
}
