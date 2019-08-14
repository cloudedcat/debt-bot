package handle

import (
	"fmt"

	"github.com/cloudedcat/finance-bot/calculator"
	"github.com/cloudedcat/finance-bot/log"
	"github.com/cloudedcat/finance-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

type calculateHandler struct {
	botHelper *botLogHelper
	calc      calculator.Service
	logger    log.Logger
}

func (hdl *calculateHandler) handle(m *tb.Message) {
	logInfo := formLogInfo(m, "Calculate")

	if m.Private() {
		hdl.botHelper.Send(m.Chat, stubMessage, logInfo)
		return
	}
	groupID := model.GroupID(m.Chat.ID)
	finalDebts, err := hdl.calc.CalculateDebts(groupID)
	if err != nil {
		hdl.botHelper.SendInternalError(m.Chat, logInfo)
		hdl.logger.IfErrorw(err, "failed to calculate debts", logInfo...)
		return
	}
	hdl.botHelper.Send(m.Chat, hdl.formMessage(finalDebts), logInfo)
}

func (hdl *calculateHandler) formMessage(debts []calculator.FinalDebt) (resp string) {
	if len(debts) == 0 {
		resp = "there ain't debts"
	}
	resp = "list of debts:\n"
	for _, debt := range debts {
		resp += fmt.Sprintf("	@%s -> @%s - %.2f",
			debt.Borrower.Alias, debt.Lender.Alias, debt.Amount)
	}
	return resp
}

// Calculate shows debt for each borrower
func Calculate(bot *tb.Bot, calc calculator.Service, logger log.Logger) {
	handler := &calculateHandler{
		botHelper: &botLogHelper{Bot: bot, logger: logger},
		calc:      calc,
		logger:    logger,
	}
	bot.Handle("/calc", handler.handle)
}
