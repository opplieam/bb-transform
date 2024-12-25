// Package lambdahandler provides handlers for AWS Lambda functions, specifically designed to process SQS events.
// It integrates with the `store` and `transform` packages to generate datasets based on configurations received through SQS messages.
// The package handles unmarshalling of SQS messages into configuration objects, triggers dataset generation, and manages errors during these processes.
// It's designed to be used in a serverless architecture where an AWS Lambda function is triggered by SQS events to perform data transformation tasks.
package lambdahandler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/opplieam/bb-transform/internal/store"
	"github.com/opplieam/bb-transform/internal/transform"
)

var (
	ErrUnmarshalConfig = errors.New("failed to unmarshal lambda config")
)

// Handler provides a struct to encapsulate the dependencies and methods required to handle SQS events.
// It includes a logger for logging and a CategoryStore for database interactions.
type Handler struct {
	log *slog.Logger
	cs  *store.CategoryStore
}

// NewHandler creates a new instance of Handler.
// It takes a CategoryStore instance as a dependency and initializes the logger with a component tag.
// Returns a pointer to the created Handler.
func NewHandler(cs *store.CategoryStore) *Handler {
	return &Handler{
		log: slog.With("component", "lambda"),
		cs:  cs,
	}
}

// HandleSQSEvent processes an SQS event containing configuration data for generating a dataset.
// It iterates through each SQS message, unmarshal the message body into a transform.Config,
// creates a new Transform instance with the unmarshalled configuration, and triggers the dataset generation process.
// Logs messages for tracking the start and completion of processing each message.
// Returns an error if unmarshalling the configuration or generating the dataset fails.
func (h *Handler) HandleSQSEvent(ctx context.Context, sqsEvent events.SQSEvent) error {
	for _, record := range sqsEvent.Records {
		h.log.Info("processing message", "message_id", record.MessageId)
		var cfg transform.Config
		if err := json.Unmarshal([]byte(record.Body), &cfg); err != nil {
			h.log.Error("failed to unmarshal config", "error", err)
			return ErrUnmarshalConfig
		}

		t := transform.NewTransform(h.cs, cfg)
		if err := t.GenerateDataset(); err != nil {
			h.log.Error("failed to generate dataset", "error", err)
			return err
		}
		h.log.Info("dataset generated successfully", "version", cfg.Version)
	}
	return nil
}
