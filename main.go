// Package main provides the entry point for an AWS Lambda function designed to generate machine learning datasets.
// It initializes a database connection, sets up a handler for processing SQS events, and starts the Lambda function.
// The application uses the `lambdahandler` package to handle SQS events, which trigger the dataset generation process.
// It utilizes the `store` package for database interactions, specifically for creating and managing category data.
// Logging is performed using the `slog` package, providing structured logs for monitoring and debugging.
// This setup is intended for deployment in an AWS environment where the Lambda function is invoked by SQS messages.
package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/opplieam/bb-transform/internal/lambdahandler"
	"github.com/opplieam/bb-transform/internal/store"
)

func init() {
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)
}

func main() {
	slog.Info("connecting to database")
	db, err := store.NewDB()
	if err != nil {
		slog.Error("failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()
	slog.Info("connected to database")

	cs := store.NewCategoryStore(db)
	lh := lambdahandler.NewHandler(cs)

	lambda.Start(lh.HandleSQSEvent)
}
