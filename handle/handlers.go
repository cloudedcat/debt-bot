package handle

import (
	"fmt"

	"github.com/cloudedcat/debt-bot/bot"
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

type addToChatHandler struct {
	mng    manager.Service
	logger log.Logger
}

func (h *addToChatHandler) handle(_ bot.Bot, m *tb.Message) {
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
	hdl := &addToChatHandler{mng: mng, logger: logger}
	bot.Handle(tb.OnAddedToGroup, hdl.handle)
}

type registerParticipantHandler struct {
	bot    bot.Bot
	mng    manager.Service
	logger log.Logger
}

func (h *registerParticipantHandler) handle(bot bot.Bot, m *tb.Message) {
	logInfo := formLogInfo(m, "RegisterParticipant")
	groupID := model.GroupID(m.Chat.ID)
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
	hdl := registerParticipantHandler{mng: mng, logger: logger}
	bot.Handle("/reg", notPrivateOnlyMiddleware(hdl.handle))
}

type participantListHandler struct {
	mng    manager.Service
	logger log.Logger
}

func (h *participantListHandler) handle(bot bot.Bot, m *tb.Message) {
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

// ParticipantList shows list of partisipants
func ParticipantList(bot bot.Bot, mng manager.Service, logger log.Logger) {
	hdl := &participantListHandler{mng: mng, logger: logger}

	bot.Handle("/list", notPrivateOnlyMiddleware(hdl.handle))
}
