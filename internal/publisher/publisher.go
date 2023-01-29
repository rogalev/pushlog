package publisher

import (
	"errors"
	"github.com/rogalev/pushlog/internal/config"
	"github.com/rogalev/pushlog/internal/storage"
)
import "context"

type Publisher interface {
	Run(ctx context.Context, s storage.Storage) error
	Shutdown(ctx context.Context) error
}

func NewInstance(cfg config.Config) (Publisher, error) {
	switch cfg.PublisherEngine {
	case "telegram":
		return NewTelegramPublisher(cfg.TelegramPublisher), nil
	default:
		return nil, errors.New("publisher engine not found")
	}
}
