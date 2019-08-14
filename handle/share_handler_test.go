package handle

import (
	"testing"

	"github.com/cloudedcat/finance-bot/model"
	"github.com/google/go-cmp/cmp"
)

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
	handler := &shareHandler{}
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
