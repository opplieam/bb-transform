// Package main provides a command-line application for generating a machine learning dataset from category data.
// This is intended for local use and debugging purposes.
// It connects to a database, fetches category information, transforms the data according to a predefined configuration,
// and stores the resulting dataset back into the database.
// The application uses the `store` package for database interactions and the `transform` package for dataset generation logic.
// It logs various stages of the process using the `slog` package for structured logging.
// The application takes configuration parameters for the dataset generation, such as version, shuffle flag, and train/validate/test ratios.
package main

import (
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/opplieam/bb-transform/internal/transform"

	"github.com/opplieam/bb-transform/internal/store"
)

func main() {
	var logger *slog.Logger
	logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		logger.Error("failed to load environment variables", "error", err)
		os.Exit(1)
	}

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
