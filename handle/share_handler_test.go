package handle

import (
	"testing"

	"github.com/cloudedcat/debt-bot/bot/mock_bot"
	"github.com/cloudedcat/debt-bot/calculator/mock_calculator"
	"github.com/cloudedcat/debt-bot/log"
	"github.com/cloudedcat/debt-bot/model"
	"github.com/cloudedcat/debt-bot/testset"
	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
)

func TestHandlerShare(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	bot := mock_bot.NewMockBot(ctrl)
	calc := mock_calculator.NewMockService(ctrl)
	logger := log.NewZapLogger()
	hdl := handlerShare{
		calc:   calc,
		logger: logger,
	}
	m := testTextMessage(testUser(1, "trillian", "Trillian"), "/share 42.0 in Milliways with @ArthurDent @FordPrefect")

	calc.EXPECT().AddDebtsByAliases(testset.GroupID, gomock.Any(), gomock.Any()).Return(nil)
	bot.EXPECT().Send(testChat(), gomock.Any(), gomock.Any())
	hdl.handle(bot, m)
}

var testShareParseCommandCases = []struct {
	invoker string
	text    string
	result  debtCommand
	isError bool
}{
	{
		invoker: "Trillian",
		text:    "/share 42.0 in Milliways with @ArthurDent @FordPrefect",
		result: debtCommand{
			Tag:       "Milliways",
			Amount:    42.0,
			Lender:    "trillian",
			Borrowers: []model.Alias{"arthurdent", "fordprefect"},
		},
		isError: false,
	},
	{
		invoker: "Trillian",
		text:    "/share 42 with  @ArthurDent  ",
		result: debtCommand{
			Tag:       "",
			Amount:    42,
			Lender:    "trillian",
			Borrowers: []model.Alias{"arthurdent"},
		},
		isError: false,
	},
	{
		invoker: "Trillian",
		text:    "/share 42 with @Trillian",
		isError: true,
	},
	{
		invoker: "Trillian",
		text:    "/share 42 in Milliways",
		isError: true,
	},
}

func TestShareParseCommand(t *testing.T) {
	handler := &handlerShare{}
	for i, testCase := range testShareParseCommandCases {
		result, customErr := handler.parseCommand(testCase.invoker, testCase.text)
		if testCase.isError {
			if customErr == nil {
				t.Errorf("case %d expected error, but got nil", i)
			}
			continue
		}

		if customErr != nil {
			t.Errorf("case %d failed: %v", i, customErr)
			continue
		}
		if diff := cmp.Diff(testCase.result, *result); diff != "" {
			t.Errorf("Wrong result (-expected, +got):\n%s", diff)
		}
	}
}
