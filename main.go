package main

import (
	"time"

	"github.com/cloudedcat/debt-bot/bot"
	"github.com/cloudedcat/debt-bot/bunt"
	"github.com/cloudedcat/debt-bot/calculator"
	"github.com/cloudedcat/debt-bot/handle"
	"github.com/cloudedcat/debt-bot/log"
	"github.com/cloudedcat/debt-bot/manager"
	"github.com/cloudedcat/debt-bot/model"
	"github.com/tidwall/buntdb"

	tb "gopkg.in/tucnak/telebot.v2"
)

func newBuntRepositories(db *buntdb.DB) (
	model.GroupRepository, model.ParticipantRepository, model.DebtRepository) {

	return bunt.NewGroupRepository(db), bunt.NewParticipantRepository(db), bunt.NewDebtRepository(db)
}

func main() {
	logger := log.NewZapLogger()
	logger.Infow("Bot initializing...")
	db, err := bunt.Open(config.DBName)
	if err != nil {
		logger.Fatalw(err.Error())
	}
	groups, partics, debts := newBuntRepositories(db)

	managerService := manager.NewService(groups, partics)
	calculatorService := calculator.NewService(debts, partics)
	_ = calculatorService

	bot, err := bot.NewTelegramBot(tb.Settings{
		Token:  config.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	}, logger)

	if err != nil {
		logger.Fatalw(err.Error())
	}

	logger.Infow("Bot authorized")

	handle.AddToChat(bot, managerService, logger)
	handle.RegisterParticipant(bot, managerService, logger)
	handle.ListParticipants(bot, managerService, logger)
	handle.ShareDebt(bot, calculatorService, logger)
	handle.Calculate(bot, calculatorService, logger)

	bot.Start()
}
