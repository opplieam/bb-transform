package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"
	"github.com/opplieam/bb-transform/internal/transform"

	"github.com/opplieam/bb-transform/internal/store"
)

func main() {

	var logger *slog.Logger
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	logger.Info("connecting to database")
	db, err := store.NewDB()
	if err != nil {
		logger.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	logger.Info("connected to database")

	cs := store.NewCategoryStore(db)
	tCfg := transform.Config{
		Version:       "v1",
		Shuffle:       true,
		TrainRatio:    60,
		ValidateRatio: 20,
		TestRatio:     20,
	}
	t := transform.NewTransform(cs, tCfg)
	if err = t.GenerateDataset(); err != nil {
		logger.Error("failed to generate dataset", "error", err)
		os.Exit(1)
	}

}
