// Package main provides the entry point for an AWS Lambda function designed to generate machine learning datasets.
// It initializes a database connection, sets up a handler for processing SQS events, and starts the Lambda function.
// The application uses the `lambdahandler` package to handle SQS events, which trigger the dataset generation process.
// It utilizes the `store` package for database interactions, specifically for creating and managing category data.
// Logging is performed using the `slog` package, providing structured logs for monitoring and debugging.
// This setup is intended for deployment in an AWS environment where the Lambda function is invoked by SQS messages.
package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/opplieam/bb-transform/internal/lambdahandler"
	"github.com/opplieam/bb-transform/internal/store"
	"github.com/opplieam/bb-transform/internal/transform"

	_ "github.com/lib/pq"
)

func initLogger() *slog.Logger {
	var logger *slog.Logger
	if os.Getenv("DEBUG") == "true" {
		logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	slog.SetDefault(logger)
	return logger
}

func run() error {
	logger := initLogger()
	err := godotenv.Load()
	if err != nil {
		logger.Info("no .env file")
	}
	logger.Info("connecting to database")
	db, err := store.NewDB()
	if err != nil {
		return err
	}
	defer db.Close()
	logger.Info("connected to database")

	cs := store.NewCategoryStore(db)

	if os.Getenv("ENV") == "dev" {
		const (
			trainRatio    = 60
			validateRatio = 20
			testRatio     = 20
		)
		tCfg := transform.Config{
			Version:       "v1",
			Shuffle:       true,
			TrainRatio:    trainRatio,
			ValidateRatio: validateRatio,
			TestRatio:     testRatio,
		}
		t := transform.NewTransform(logger, cs, tCfg)
		if err = t.GenerateDataset(); err != nil {
			logger.Error("failed to generate dataset", "error", err)
			return err
		}
	} else {
		lh := lambdahandler.NewHandler(logger, cs)
		lambda.Start(lh.HandleSQSEvent)
	}
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
