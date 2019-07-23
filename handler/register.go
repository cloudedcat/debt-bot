package handler

import (
	"log"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Register is handler for registering new paticipant in group
func Register(bot *tb.Bot) (string, func(m *tb.Message)) {
	return "/reg", func(m *tb.Message) {
		log.Printf("/reg %s in chat %s", m.Sender.Username, m.Chat.Username)
	}
}
