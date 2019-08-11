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
	bot.Handle(tb.OnAddedToGroup, func(m *tb.Message) {
		id := model.GroupID(m.Chat.ID)
		additionalInfo := []interface{}{"handler", "AddToChat", "groupID", id}
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
	bot.Handle("/reg", func(m *tb.Message) {
		groupID := model.GroupID(m.Chat.ID)
		partic := model.Participant{
			ID:        model.ParticipantID(m.Sender.ID),
			Alias:     model.Alias(m.Sender.Username),
			FirstName: m.Sender.FirstName,
			LastName:  m.Sender.LastName,
		}
		additionalInfo := []interface{}{
			"handler", "RegisterParticipant",
			"groupID", groupID,
			"particID", partic.ID,
			"alias", m.Sender.Username,
		}
		if m.Private() {
			stubInPrivateChat(bot, m, logger, additionalInfo)
			return
		}
		if err := mng.RegisterParticipant(groupID, partic); err != nil {
			additionalInfo = append(additionalInfo, "error", err.Error())
			logger.Errorw("failed to register participant", additionalInfo...)
		}
		logger.Infow("register new participant", additionalInfo...)
		_, err := bot.Send(m.Chat, fmt.Sprintf("registered participant @%v", partic.Alias))
		if err != nil {
			logger.Errorw("failed to send msg about registration", additionalInfo...)
		}
	})
}

func stubInPrivateChat(bot *tb.Bot, m *tb.Message, logger log.Logger, additionalInfo []interface{}) {
	_, err := bot.Send(m.Chat, fmt.Sprintf("the command doesn't work in private chat"))
	if err != nil {
		additionalInfo = append(additionalInfo, "error", err)
		logger.Errorw("failed to send stub message", additionalInfo...)
	}
	logger.Infow("show stub message in private chat", additionalInfo...)
	return
}
