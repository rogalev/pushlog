package publisher

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/rogalev/pushlog/internal/config"
	"github.com/rogalev/pushlog/internal/logging"
	"github.com/rogalev/pushlog/internal/message"
	"github.com/rogalev/pushlog/internal/storage"
	"go.uber.org/zap"
	"io/ioutil"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

type TelegramPublisher struct {
	token      string
	chat       int64
	sendAsFile bool
	tmpFileDir string

	botApi *tgbotapi.BotAPI
}

func NewTelegramPublisher(cfg config.TelegramPublisherConfig) *TelegramPublisher {
	return &TelegramPublisher{
		token:      cfg.Token,
		chat:       cfg.Chat,
		sendAsFile: cfg.SendAsFile,
		tmpFileDir: cfg.TmpFileDir,
	}
}

func (p *TelegramPublisher) Run(ctx context.Context, s storage.Storage) error {
	logger := logging.GetInstance()

	bot, err := tgbotapi.NewBotAPI(p.token)
	if err != nil {
		return err
	}
	p.botApi = bot

	ticker := time.NewTicker(50 * time.Millisecond)

L:
	for {
		select {
		case <-ctx.Done():
			logger.Info("Telegram publisher is stopping")
			break L
		case <-ticker.C:

			if !s.HasMessage() {
				break
			}

			if p.sendAsFile {
				if err := p.sendFile(s.ShiftAll()); err != nil {
					logger.Error("Send telegram file error", zap.Error(err))
				}
			} else {
				if m, ok := s.Shift(); ok {
					if err := p.sendMessage(m); err != nil {
						logger.Error("Send telegram message error", zap.Error(err))
					}
				}
			}
		}
	}

	return nil
}

func (p *TelegramPublisher) Shutdown(ctx context.Context) error {
	logger := logging.GetInstance()
	logger.Info("Telegram publisher shutdown")
	return nil
}

func (p *TelegramPublisher) sendMessage(m message.Message) error {

	data := m.Body

	if utf8.RuneCountInString(data) > 4096 {
		data = string([]rune(data)[:4096])
	}

	msg := tgbotapi.NewMessage(p.chat, data)
	_, err := p.botApi.Send(msg)
	return err
}

func (p *TelegramPublisher) sendFile(list []message.Message) error {

	file, err := ioutil.TempFile(p.tmpFileDir, "telegram_log_file.*.log")

	if err != nil {
		return err
	}

	defer func() {
		_ = os.Remove(file.Name())
	}()

	for _, m := range list {

		data := strings.Replace(m.Body, `\n`, "\n", -1)

		if _, err = file.WriteString(data); err != nil {
			return err
		}
	}

	document := tgbotapi.NewDocument(p.chat, tgbotapi.FilePath(file.Name()))

	if _, err = p.botApi.Send(document); err != nil {
		return err
	}

	return nil
}
