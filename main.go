package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	bot, err := tgbotapi.NewBotAPI("finance-bot")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
}
