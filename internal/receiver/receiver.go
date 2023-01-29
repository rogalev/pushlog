package receiver

import (
	"context"
	"errors"
	"github.com/rogalev/pushlog/internal/config"
	"github.com/rogalev/pushlog/internal/storage"
)

type Receiver interface {
	Run(ctx context.Context, s storage.Storage) error
	Shutdown(ctx context.Context) error
}

func NewInstance(cfg config.Config) (Receiver, error) {
	switch cfg.ReceiverEngine {
	case "http":
		return NewHttpReceiver(cfg.HttpReceiver), nil
	default:
		return nil, errors.New("receiver engine not found")
	}
}
