package storage

import (
	"context"
	"errors"
	"github.com/rogalev/pushlog/internal/config"
	"github.com/rogalev/pushlog/internal/message"
)

type Storage interface {
	Push(message message.Message) bool
	Shift() (message.Message, bool)
	ShiftAll() []message.Message
	HasMessage() bool
	Run(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

func NewInstance(cfg config.Config) (Storage, error) {
	switch cfg.StorageEngine {
	case "memory":
		return NewMemoryStorage(cfg.MemoryStorage), nil
	default:
		return nil, errors.New("receiver engine not found")
	}
}
