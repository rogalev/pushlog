package receiver

import (
	"context"
	"fmt"
	"github.com/rogalev/pushlog/internal/config"
	"github.com/rogalev/pushlog/internal/logging"
	"github.com/rogalev/pushlog/internal/message"
	"github.com/rogalev/pushlog/internal/storage"
	"net/http"
	"time"
)

type HttpReceiver struct {
	Host         string
	Port         int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	httpServer http.Server
}

func NewHttpReceiver(cfg config.HttpReceiverConfig) *HttpReceiver {
	return &HttpReceiver{
		Host:         cfg.Host,
		Port:         cfg.Port,
		ReadTimeout:  time.Duration(cfg.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.WriteTimeout) * time.Second,
	}
}

func (r *HttpReceiver) Run(ctx context.Context, s storage.Storage) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/push", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		m := message.Message{
			Key:        r.FormValue("key"),
			Expiration: r.FormValue("expiration"),
			Body:       r.FormValue("body"),
		}

		_ = s.Push(m)
		return
	})

	r.httpServer = http.Server{
		Addr:         fmt.Sprintf("%s:%d", r.Host, r.Port),
		ReadTimeout:  r.ReadTimeout,
		WriteTimeout: r.WriteTimeout,
		Handler:      mux,
	}

	return r.httpServer.ListenAndServe()
}

func (r *HttpReceiver) Shutdown(ctx context.Context) error {
	logger := logging.GetInstance()
	logger.Info("Http receiver shutdown")
	return r.httpServer.Shutdown(ctx)
}
