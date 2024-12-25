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

type Handler struct {
	log *slog.Logger
	cs  *store.CategoryStore
}

func NewHandler(cs *store.CategoryStore) *Handler {
	return &Handler{
		log: slog.With("component", "lambda"),
		cs:  cs,
	}
}

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
