package storage

import (
	"context"
	"github.com/rogalev/pushlog/internal/config"
	"github.com/rogalev/pushlog/internal/logging"
	"github.com/rogalev/pushlog/internal/message"
	"go.uber.org/zap"
	"strconv"
	"sync"
	"time"
)

type MemoryStorage struct {
	sync.Mutex
	ExpirationGCPeriod int64
	keys               map[string]int64
	messages           []message.Message
}

func NewMemoryStorage(cfg config.MemoryStorageConfig) *MemoryStorage {
	return &MemoryStorage{
		ExpirationGCPeriod: cfg.ExpirationGCPeriod,
		keys:               make(map[string]int64),
	}
}

func (s *MemoryStorage) Run(ctx context.Context) error {

	logger := logging.GetInstance()
	ticker := time.NewTicker(time.Duration(s.ExpirationGCPeriod) * time.Second)

L:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Memory storage is stopping")
			break L
		case <-ticker.C:
			logger.Info("Memory storage expired key garbage collector started")
			s.deleteExpired()
		}
	}

	return nil
}

func (s *MemoryStorage) Shutdown(ctx context.Context) error {
	logger := logging.GetInstance()
	logger.Info("Memory storage shutdown")
	return nil
}

func (s *MemoryStorage) Push(m message.Message) bool {
	s.Lock()
	defer s.Unlock()

	logger := logging.GetInstance()

	if v, ok := s.keys[m.Key]; ok && v > time.Now().Unix() {
		logger.Debug("Key already in cache", zap.String("key", m.Key))
		return false
	}

	expiration, err := strconv.ParseInt(m.Expiration, 10, 64)

	if err != nil {
		logger.Error("Expiration value parsing error", zap.String("key", m.Key), zap.String("originalExpiration", m.Expiration), zap.Int64("parsedExpiration", expiration))

	}

	if expiration > 0 {
		s.keys[m.Key] = time.Now().Unix() + expiration
	}

	logger.Info("Message successfully added to memory storage", zap.String("key", m.Key), zap.Int64("expiration", expiration))

	s.messages = append(s.messages, m)

	return true
}

func (s *MemoryStorage) Shift() (message.Message, bool) {
	s.Lock()
	defer s.Unlock()

	var m message.Message

	if len(s.messages) < 1 {
		return m, false
	}

	m = s.messages[0]
	s.messages = s.messages[1:]

	return m, true
}

func (s *MemoryStorage) ShiftAll() []message.Message {
	s.Lock()
	defer s.Unlock()

	var m []message.Message

	if len(s.messages) < 1 {
		return m
	}

	m = s.messages
	s.messages = nil

	return m
}

func (s *MemoryStorage) deleteExpired() {
	s.Lock()
	defer s.Unlock()

	logger := logging.GetInstance()
	now := time.Now().Unix()

	for k, v := range s.keys {
		if v < now {
			logger.Debug("Remove expired key from storage", zap.String("key", k))
			delete(s.keys, k)
		}
	}
}

func (s *MemoryStorage) HasMessage() bool {
	s.Lock()
	defer s.Unlock()

	return len(s.messages) > 0
}
