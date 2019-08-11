package main

import (
	"time"

	"github.com/cloudedcat/finance-bot/bunt"
	"github.com/cloudedcat/finance-bot/calculator"
	"github.com/cloudedcat/finance-bot/handle"
	"github.com/cloudedcat/finance-bot/log"
	"github.com/cloudedcat/finance-bot/manager"
	"github.com/cloudedcat/finance-bot/model"
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
	_ = calculator.NewService(debts, partics) // NYI

	bot, err := tb.NewBot(tb.Settings{
		Token:  config.BotToken,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		logger.Fatalw(err.Error())
	}

	logger.Infow("Bot authorized")

	handle.AddToChat(bot, managerService, logger)
	handle.RegisterParticipant(bot, managerService, logger)

	bot.Start()
}
