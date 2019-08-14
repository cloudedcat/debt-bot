package handle

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/cloudedcat/finance-bot/calculator"
	"github.com/cloudedcat/finance-bot/log"
	"github.com/cloudedcat/finance-bot/model"
	tb "gopkg.in/tucnak/telebot.v2"
)

type Error struct {
	internal  error
	userError string
}

func newError(internal error, userError string) *Error {
	return &Error{internal: internal, userError: userError}
}

const (
	uErrAmountFormat    = "wrong amount format, example: 25.17"
	uErrCommandFormat   = "wrong command, example:\n/share 42.0 in Milliways with @ArthurDent @FordPrefect"
	uErrParticCollision = "lender can't be borrower in the same time"
)

// ShareDebt share debt between pointed participants e.g.
// A send: '/share 42.0 in restaurant with @B @C'
// that means A paid 42.0 for B and C. Share command spread this amount between
// A, B and C. So, B will owe A 14.0, C will owe A 14.0.
func ShareDebt(bot *tb.Bot, calc calculator.Service, logger log.Logger) {
	handler := &shareHandler{
		botHelper: &botLogHelper{Bot: bot, logger: logger},
		calc:      calc,
		logger:    logger,
	}

	bot.Handle("/share", handler.handle)
}

type shareHandler struct {
	botHelper *botLogHelper
	calc      calculator.Service
	logger    log.Logger
}

type debtCommand struct {
	Tag       string
	Amount    float64
	Lender    model.Alias
	Borrowers []model.Alias
}

func (cmd *debtCommand) generateDebts() []calculator.DebtWithAliases {
	share := cmd.Amount / float64(len(cmd.Borrowers)+1)
	var debts []calculator.DebtWithAliases
	for _, b := range cmd.Borrowers {
		d := calculator.DebtWithAliases{
			Amount:   share,
			Tag:      cmd.Tag,
			Borrower: b,
			Lender:   cmd.Lender,
		}
		debts = append(debts, d)
	}
	return debts
}

func (sh *shareHandler) handle(m *tb.Message) {
	loglInfo := formLogInfo(m, "ShareDebt")
	if m.Private() {
		sh.botHelper.Send(m.Chat, stubMessage, loglInfo)
		return
	}
	cmd, customErr := sh.parseCommand(m.Sender.Username, m.Text)
	if customErr != nil {
		sh.botHelper.Send(m.Chat, customErr.userError, loglInfo)
		return
	}
	err := sh.calc.AddDebtsByAliases(model.GroupID(m.Chat.ID), cmd.generateDebts()...)
	if err != nil {
		sh.botHelper.SendInternalError(m.Chat, loglInfo)
		return
	}
	sh.botHelper.Send(m.Chat, "debt shared", loglInfo)
}

func (sh *shareHandler) parseCommand(invoker string, text string) (*debtCommand, *Error) {
	text = sh.prepareText(text)
	re := regexp.MustCompile(`^/share ((\d+)(\.\d+)?) (in (.+) )?with (.*)$`)
	result := re.FindStringSubmatch(text)
	submatchNumber := 7
	if len(result) != submatchNumber {
		return nil, newError(errors.New("failed to match string"), uErrCommandFormat)
	}
	rawAmount, tag, rawBorrowers := result[1], result[5], result[6]
	amount, err := strconv.ParseFloat(rawAmount, 64)
	if err != nil {
		return nil, newError(err, uErrAmountFormat)
	}
	lender := model.MustBuildAlias(invoker)
	var borrowers []model.Alias
	for _, username := range strings.Split(rawBorrowers, " ") {
		b, customErr := sh.processBorrower(lender, username)
		if customErr != nil {
			return nil, customErr
		}
		borrowers = append(borrowers, b)
	}
	cmd := &debtCommand{
		Tag:       tag,
		Amount:    amount,
		Lender:    lender,
		Borrowers: borrowers,
	}
	return cmd, nil
}

func (sh *shareHandler) prepareText(text string) string {
	space := regexp.MustCompile(`\s+`)
	text = space.ReplaceAllString(text, " ")
	return strings.TrimSpace(text)
}

func (sh *shareHandler) processBorrower(lender model.Alias, rawBorrower string) (model.Alias, *Error) {
	b, err := model.BuildAlias(rawBorrower)
	if err != nil {
		return "", &Error{
			internal:  err,
			userError: fmt.Sprintf("wrong username '%s'", rawBorrower),
		}
	}
	if b == lender {
		return "", &Error{
			internal:  errors.New(uErrParticCollision),
			userError: uErrParticCollision,
		}
	}
	return b, nil
}
