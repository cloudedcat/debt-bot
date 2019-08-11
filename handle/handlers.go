package handle

import (
	"fmt"

	"github.com/cloudedcat/finance-bot/log"
	"github.com/cloudedcat/finance-bot/manager"
	"github.com/cloudedcat/finance-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

// AddToChat handles adding bot to new chat
func AddToChat(bot *tb.Bot, mng manager.Service, logger log.Logger) {
	const handlerName = "AddToChat"

	bot.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		additionalInfo := formAdditionalInfo(m, handlerName)

		id := model.GroupID(m.Chat.ID)
		if err := mng.RegisterGroup(model.Group{ID: id}); err != nil {
			additionalInfo = append(additionalInfo, "error", err.Error())
			logger.Errorw("failed to register group", additionalInfo)
			return
		}
		logger.Infow("register new group", additionalInfo...)
	})
}

// RegisterParticipant add handler for registering new paticipant in group
func RegisterParticipant(bot *tb.Bot, mng manager.Service, logger log.Logger) {
	const handlerName = "RegisterParticipant"

	bot.Handle("/reg", func(m *tb.Message) {
		additionalInfo := formAdditionalInfo(m, handlerName)

		groupID := model.GroupID(m.Chat.ID)
		partic := model.Participant{
			ID:        model.ParticipantID(m.Sender.ID),
			Alias:     model.Alias(m.Sender.Username),
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
		}
		if m.Private() {
			stubInPrivateChat(bot, m, logger, additionalInfo)
			return
		}
		if err := mng.RegisterParticipant(groupID, partic); err != nil {
			logger.IfErrorw(err, "failed to register participant", additionalInfo...)
		}
		logger.Infow("register new participant", additionalInfo...)
		_, err := bot.Send(m.Chat, fmt.Sprintf("@%v is registered", partic.Alias))
		if err != nil {
			logger.IfErrorw(err, "failed to send msg about registration", additionalInfo...)
			return
		}
	})
}

func ParticipantList(bot *tb.Bot, mng manager.Service, logger log.Logger) {
	const handlerName = "RegisterParticipant"

	bot.Handle("/list", func(m *tb.Message) {
		additionalInfo := formAdditionalInfo(m, handlerName)

		if m.Private() {
			stubInPrivateChat(bot, m, logger, additionalInfo)
			return
		}
		groupID := model.GroupID(m.Chat.ID)
		partics, err := mng.ListParticipant(groupID)
		if err != nil {
			logger.IfErrorw(err, "failed to list participant", additionalInfo...)
			return
		}
		text := partics.AsString()
		if text == "" {
			text = "list of participants is empty"
		} else {
			text = "list of participants:\n" + text
		}
		if _, err := bot.Send(m.Chat, text); err != nil {
			logger.IfErrorw(err, "failed to send message", additionalInfo...)
			return
		}
	})
}

func formAdditionalInfo(m *tb.Message, handlerName string) []interface{} {
	return []interface{}{
		"handler", handlerName,
		"chatID", m.Chat.ID,
		"invoker", m.Sender.Username,
		"invokerID", m.Sender.ID,
	}
}

func stubInPrivateChat(bot *tb.Bot, m *tb.Message, logger log.Logger, additionalInfo []interface{}) {
	_, err := bot.Send(m.Chat, fmt.Sprintf("the command doesn't work in private chat"))
	if err != nil {
		logger.IfErrorw(err, "failed to send stub message", additionalInfo...)
		return
	}
	logger.Infow("show stub message in private chat", additionalInfo...)
	return
}
