//go:generate mockgen -destination mock_bot/mock_bot.go github.com/cloudedcat/debt-bot/bot Bot

package bot

import (
	"github.com/cloudedcat/debt-bot/log"

	tb "gopkg.in/tucnak/telebot.v2"
)

type Bot interface {
	Send(to tb.Recipient, what interface{}, logInfo []interface{}, options ...interface{}) (*tb.Message, error)
	SendInternalError(to tb.Recipient, logInfo []interface{}) (*tb.Message, error)
	Handle(endpoint interface{}, handler Handler)
	Start()
}

func NewTelegramBot(settings tb.Settings, logger log.Logger) (*TelegramBot, error) {
	bot, err := tb.NewBot(settings)
	if err != nil {
		return nil, err
	}
	return &TelegramBot{bot: bot, logger: logger}, nil
}

type TelegramBot struct {
	bot    *tb.Bot
	logger log.Logger
}

func (b *TelegramBot) Send(to tb.Recipient, what interface{}, logInfo []interface{}, options ...interface{}) (*tb.Message, error) {
	msg, err := b.bot.Send(to, what, options...)
	b.logger.IfErrorw(err, "failed to send message", logInfo...)
	return msg, err
}

func (b *TelegramBot) SendInternalError(to tb.Recipient, logInfo []interface{}) (*tb.Message, error) {
	return b.Send(to, "internal error", logInfo)
}

type Handler func(Bot, *tb.Message)

func (b *TelegramBot) Handle(endpoint interface{}, handler Handler) {
	normalHandler := func(msg *tb.Message) { handler(b, msg) }
	b.bot.Handle(endpoint, normalHandler)
}

func (b *TelegramBot) Start() {
	b.bot.Start()
}
