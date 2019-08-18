package handle

import (
	"fmt"

	"github.com/cloudedcat/debt-bot/bot"
	"github.com/cloudedcat/debt-bot/calculator"
	"github.com/cloudedcat/debt-bot/log"
	"github.com/cloudedcat/debt-bot/manager"
	"github.com/cloudedcat/debt-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

func notPrivateOnlyMiddleware(hdl bot.Handler) bot.Handler {
	const stubMessage = "that command doesn't work in private chat"

	return func(bot bot.Bot, m *tb.Message) {
		if !m.Private() {
			hdl(bot, m)
			return
		}
		logInfo := formLogInfo(m, "notPrivateOnlyMiddleware")
		bot.Send(m.Chat, stubMessage, logInfo)
	}
}

func formLogInfo(m *tb.Message, handlerName string) []interface{} {
	return []interface{}{
		"handler", handlerName,
		"chatID", m.Chat.ID,
		"invoker", m.Sender.Username,
		"invokerID", m.Sender.ID,
		"text", m.Text,
	}
}

type handlerAddToChat struct {
	mng    manager.Service
	logger log.Logger
}

func (h *handlerAddToChat) handle(_ bot.Bot, m *tb.Message) {
	logInfo := formLogInfo(m, "AddToChat")

	id := model.GroupID(m.Chat.ID)
	if err := h.mng.RegisterGroup(model.Group{ID: id}); err != nil {
		h.logger.IfErrorw(err, "failed to register group", logInfo)
		return
	}
	h.logger.Infow("register new group", logInfo...)
}

// AddToChat handles adding bot to new chat
func AddToChat(bot bot.Bot, mng manager.Service, logger log.Logger) {
	hdl := &handlerAddToChat{mng: mng, logger: logger}
	bot.Handle(tb.OnAddedToGroup, hdl.handle)
}

type handlerRegisterParticipant struct {
	bot    bot.Bot
	mng    manager.Service
	logger log.Logger
}

func (h *handlerRegisterParticipant) handle(bot bot.Bot, m *tb.Message) {
	logInfo := formLogInfo(m, "RegisterParticipant")
	groupID := model.GroupID(m.Chat.ID)
	if m.Sender.Username == "" {
		bot.Send(m.Chat, "please, set username in Telegram", logInfo)
		return
	}
	partic := model.Participant{
		ID:        model.ParticipantID(m.Sender.ID),
		Alias:     model.MustBuildAlias(m.Sender.Username),
		FirstName: m.Sender.FirstName,
		LastName:  m.Sender.LastName,
	}
	if err := h.mng.RegisterParticipant(groupID, partic); err != nil {
		bot.SendInternalError(m.Chat, logInfo)
		h.logger.IfErrorw(err, "failed to register participant", logInfo...)
	}
	h.logger.Infow("register new participant", logInfo...)
	bot.Send(m.Chat, fmt.Sprintf("@%v is registered", partic.Alias), logInfo)
}

// RegisterParticipant adds handler for registering new paticipant in group
func RegisterParticipant(bot bot.Bot, mng manager.Service, logger log.Logger) {
	hdl := &handlerRegisterParticipant{mng: mng, logger: logger}
	bot.Handle("/reg", notPrivateOnlyMiddleware(hdl.handle))
}

type handlerListParticipant struct {
	mng    manager.Service
	logger log.Logger
}

func (h *handlerListParticipant) handle(bot bot.Bot, m *tb.Message) {
	logInfo := formLogInfo(m, "ParticipantList")
	groupID := model.GroupID(m.Chat.ID)
	partics, err := h.mng.ListParticipant(groupID)
	if err != nil {
		bot.SendInternalError(m.Chat, logInfo)
		h.logger.IfErrorw(err, "failed to list participant", logInfo...)
		return
	}
	text := partics.AsString()
	if text == "" {
		text = "list of participants is empty"
	} else {
		text = "list of participants:\n" + text
	}
	bot.Send(m.Chat, text, logInfo)
}

// ListParticipants shows list of partisipants
func ListParticipants(bot bot.Bot, mng manager.Service, logger log.Logger) {
	hdl := &handlerListParticipant{mng: mng, logger: logger}
	bot.Handle("/list", notPrivateOnlyMiddleware(hdl.handle))
}

type handlerShowDebtHistory struct {
	calc   calculator.Service
	logger log.Logger
}

func (h *handlerShowDebtHistory) handle(bot bot.Bot, m *tb.Message) {
	logInfo := formLogInfo(m, "ParticipantList")
	groupID := model.GroupID(m.Chat.ID)
	particID := model.ParticipantID(m.Sender.ID)
	debts, err := h.calc.FindDebts(groupID, particID)
	if err != nil {
		bot.SendInternalError(m.Chat, logInfo)
		h.logger.IfErrorw(err, "failed to find debts", logInfo...)
		return
	}
	text := h.formText(particID, debts)

	if text == "" {
		text = "debt history is empty"
	} else {
		text = "debt history:\n" + text
	}
	bot.Send(m.Chat, text, logInfo)
}

func (h *handlerShowDebtHistory) formText(particID model.ParticipantID, debts []calculator.DetailedDebt) string {
	text := ""
	for _, d := range debts {
		action, whom := "", ""
		if d.BorrowerID == particID {
			action, whom = "owe", string(d.Lender.Alias)
		} else {
			action, whom = "lend", string(d.Borrower.Alias)
		}

		where := ""
		if d.Tag != "" {
			where = fmt.Sprintf("in %s", d.Tag)
		}

		text += fmt.Sprintf("%s: %s to @%s %.2f %s\n",
			d.Date.Format("02.01.2006 15:04:05"), action, whom, d.Amount, where)
	}
	return text
}

// ShowDebtHistory shows personal history of debts
func ShowDebtHistory(bot bot.Bot, calc calculator.Service, logger log.Logger) {
	hdl := &handlerShowDebtHistory{calc: calc, logger: logger}
	bot.Handle("/history", notPrivateOnlyMiddleware(hdl.handle))
}

type handlerAmnesty struct {
	calc   calculator.Service
	logger log.Logger
}

func (h *handlerAmnesty) handle(bot bot.Bot, m *tb.Message) {
	logInfo := formLogInfo(m, "Amnesty")
	groupID := model.GroupID(m.Chat.ID)
	err := h.calc.ClearDebts(groupID)
	if err != nil {
		bot.SendInternalError(m.Chat, logInfo)
		h.logger.IfErrorw(err, "failed to clear debts", logInfo...)
		return
	}
	bot.Send(m.Chat, "debts have been wiped", logInfo)
}

// Amnesty removes all debts
func Amnesty(bot bot.Bot, calc calculator.Service, logger log.Logger) {
	hdl := &handlerAmnesty{calc: calc, logger: logger}
	bot.Handle("/amnesty", notPrivateOnlyMiddleware(hdl.handle))
}
