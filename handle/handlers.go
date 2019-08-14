package handle

import (
	"fmt"

	"github.com/cloudedcat/finance-bot/log"
	"github.com/cloudedcat/finance-bot/manager"
	"github.com/cloudedcat/finance-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

const stubMessage = "that command doesn't work in private chat"

type botLogHelper struct {
	*tb.Bot
	logger log.Logger
}

func (b *botLogHelper) Send(to tb.Recipient, what interface{}, logInfo []interface{}, options ...interface{}) (*tb.Message, error) {
	msg, err := b.Bot.Send(to, what, options...)
	b.logger.IfErrorw(err, "failed to send message", logInfo...)
	return msg, err
}

func (b *botLogHelper) SendInternalError(to tb.Recipient, logInfo []interface{}) (*tb.Message, error) {
	msg, err := b.Bot.Send(to, "internal error")
	b.logger.IfErrorw(err, "failed to send message", logInfo...)
	return msg, err
}

// AddToChat handles adding bot to new chat
func AddToChat(bot *tb.Bot, mng manager.Service, logger log.Logger) {
	const handlerName = "AddToChat"

	bot.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		additionalInfo := formLogInfo(m, handlerName)

		id := model.GroupID(m.Chat.ID)
		if err := mng.RegisterGroup(model.Group{ID: id}); err != nil {
			logger.IfErrorw(err, "failed to register group", additionalInfo)
			return
		}
		logger.Infow("register new group", additionalInfo...)
	})
}

// RegisterParticipant adds handler for registering new paticipant in group
func RegisterParticipant(bot *tb.Bot, mng manager.Service, logger log.Logger) {
	const handlerName = "RegisterParticipant"
	botHelper := &botLogHelper{Bot: bot, logger: logger}

	bot.Handle("/reg", func(m *tb.Message) {
		logInfo := formLogInfo(m, handlerName)
		if m.Private() {
			botHelper.Send(m.Chat, stubMessage, logInfo)
			return
		}

		groupID := model.GroupID(m.Chat.ID)
		partic := model.Participant{
			ID:        model.ParticipantID(m.Sender.ID),
			Alias:     model.MustBuildAlias(m.Sender.Username),
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
		}
		if err := mng.RegisterParticipant(groupID, partic); err != nil {
			botHelper.SendInternalError(m.Chat, logInfo)
			logger.IfErrorw(err, "failed to register participant", logInfo...)
		}
		logger.Infow("register new participant", logInfo...)
		botHelper.Send(m.Chat, fmt.Sprintf("@%v is registered", partic.Alias), logInfo)
	})
}

// ParticipantList shows list of partisipants
func ParticipantList(bot *tb.Bot, mng manager.Service, logger log.Logger) {
	const handlerName = "ParticipantList"
	botHelper := &botLogHelper{Bot: bot, logger: logger}

	bot.Handle("/list", func(m *tb.Message) {
		logInfo := formLogInfo(m, handlerName)

		if m.Private() {
			botHelper.Send(m.Chat, stubMessage, logInfo)
			return
		}
		groupID := model.GroupID(m.Chat.ID)
		partics, err := mng.ListParticipant(groupID)
		if err != nil {
			botHelper.SendInternalError(m.Chat, logInfo)
			logger.IfErrorw(err, "failed to list participant", logInfo...)
			return
		}
		text := partics.AsString()
		if text == "" {
			text = "list of participants is empty"
		} else {
			text = "list of participants:\n" + text
		}
		botHelper.Send(m.Chat, text, logInfo)
	})
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
