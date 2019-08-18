package handle

import (
	"testing"

	"github.com/cloudedcat/debt-bot/bot/mock_bot"
	"github.com/cloudedcat/debt-bot/log"
	"github.com/cloudedcat/debt-bot/manager/mock_manager"
	"github.com/cloudedcat/debt-bot/model"
	"github.com/cloudedcat/debt-bot/testset"
	"github.com/golang/mock/gomock"

	tb "gopkg.in/tucnak/telebot.v2"
)

func TestAddToChatHandler(t *testing.T) {
	// Arrange
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bot := mock_bot.NewMockBot(ctrl)
	mng := mock_manager.NewMockService(ctrl)
	logger := log.NewZapLogger()
	hdl := handlerAddToChat{
		mng:    mng,
		logger: logger,
	}
	m := testTextMessage(testUser(1, "alice", "Alice"), "")

	// Assert
	mng.EXPECT().RegisterGroup(*model.BuildGroup(model.GroupID(testChat().ID))).Return(nil)
	// Act
	hdl.handle(bot, m)
}

func TestRegisterParticipantHandler(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bot := mock_bot.NewMockBot(ctrl)
	mng := mock_manager.NewMockService(ctrl)
	logger := log.NewZapLogger()
	hdl := handlerRegisterParticipant{
		mng:    mng,
		logger: logger,
	}
	m := testTextMessage(testUser(1, "alice", "Alice"), "/reg")
	expectedPartic := model.Participant{1, "alice", "Alice", ""}

	mng.EXPECT().RegisterParticipant(testset.GroupID, expectedPartic).Return(nil)
	bot.EXPECT().Send(testChat(), gomock.Any(), gomock.Any())
	hdl.handle(bot, m)
}

func TestRegisterParticipantHandlerNoUsername(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bot := mock_bot.NewMockBot(ctrl)
	mng := mock_manager.NewMockService(ctrl)
	logger := log.NewZapLogger()
	hdl := handlerRegisterParticipant{
		mng:    mng,
		logger: logger,
	}
	m := testTextMessage(testUser(1, "", "Alice"), "/reg")

	bot.EXPECT().Send(testChat(), gomock.Any(), gomock.Any()).Do(
		func(to tb.Recipient, what interface{}, logInfo []interface{}, options ...interface{}) (*tb.Message, error) {
			if what.(string) != "please, set username in Telegram" {
				t.Fatal("Handler sent unexpected message:", what)
			}
			return &tb.Message{}, nil
		})
	hdl.handle(bot, m)
}

func testTextMessage(sender *tb.User, text string) *tb.Message {
	return &tb.Message{
		ID:     1,
		Sender: sender,
		Chat:   testChat(),
		Text:   text,
	}
}

func testUser(id int, username string, firstname string) *tb.User {
	return &tb.User{
		ID:        id,
		Username:  username,
		FirstName: firstname,
	}
}

func testChat() *tb.Chat {
	return &tb.Chat{
		ID:    int64(testset.GroupID),
		Title: "Boston Tea Party",
	}
}
