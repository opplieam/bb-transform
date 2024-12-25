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
