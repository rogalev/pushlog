package service

import (
	"context"
	"fmt"
	"github.com/rogalev/pushlog/internal/publisher"
	"github.com/rogalev/pushlog/internal/receiver"
	"github.com/rogalev/pushlog/internal/storage"
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	Receiver  receiver.Receiver
	Publisher publisher.Publisher
	Storage   storage.Storage
}

func NewInstance(receiver receiver.Receiver, publisher publisher.Publisher, storage storage.Storage) *Service {
	return &Service{
		Receiver:  receiver,
		Publisher: publisher,
		Storage:   storage,
	}
}

func (s *Service) Run() {

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		cancel()
	}()

	g, gCtx := errgroup.WithContext(ctx)

	g.Go(func() error {
		return s.Receiver.Run(ctx, s.Storage)
	})
	g.Go(func() error {
		<-gCtx.Done()
		return s.Receiver.Shutdown(context.Background())
	})

	g.Go(func() error {
		return s.Publisher.Run(ctx, s.Storage)
	})
	g.Go(func() error {
		<-gCtx.Done()
		return s.Publisher.Shutdown(context.Background())
	})

	g.Go(func() error {
		return s.Storage.Run(ctx)
	})
	g.Go(func() error {
		<-gCtx.Done()
		return s.Storage.Shutdown(context.Background())
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}
}
