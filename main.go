package main

import (
	"time"

	"github.com/cloudedcat/finance-bot/bunt"
	"github.com/cloudedcat/finance-bot/calculator"
	"github.com/cloudedcat/finance-bot/handle"
	"github.com/cloudedcat/finance-bot/log"
	"github.com/cloudedcat/finance-bot/manager"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	logger := log.NewZapLogger()
	logger.Infow("Bot initializing")
	db, err := bunt.Open(":memory:")
	if err != nil {
		logger.Fatalw(err.Error())
	}
	groups := bunt.NewGroupRepository(db)
	debts := bunt.NewDebtRepository(db)
	partics := bunt.NewParticipantRepository(db)

	managerService := manager.NewService(groups, partics)
	_ = calculator.NewService(debts, partics) // NYI

	bot, err := tb.NewBot(tb.Settings{
		Token:  "837157414:AAFxMc8exxbi69prpVOAMRU_EHns7pnQFa4",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		logger.Fatalw(err.Error())
	}

	logger.Infow("Bot authorized")

	handle.AddToChat(bot, managerService, logger)
	handle.RegisterParticipant(bot, managerService, logger)

	bot.Start()
	// bot.Handle("/hello", func(m *tb.Message) {
	// 	bot.Send(m.Sender, "hello world")
	// })
	// replyBtn := tb.ReplyButton{Text: "🌕 Button #1"}
	// replyKeys := [][]tb.ReplyButton{
	// 	[]tb.ReplyButton{replyBtn},
	// 	// ...
	// }

	// And this one — just under the message itself.
	// Pressing it will cause the client to send
	// the bot a callback.
	//
	// Make sure Unique stays unique as it has to be
	// for callback routing to work.
	// inlineBtn := tb.InlineButton{
	// 	Unique: "sad_moon",
	// 	Text:   "🌚 Button #2",
	// }
	// inlineKeys := [][]tb.InlineButton{
	// 	[]tb.InlineButton{inlineBtn},
	// 	// ...
	// }

	// b.Handle(&replyBtn, func(m *tb.Message) {
	// 	// on reply button pressed
	// })

	// bot.Handle(&inlineBtn, func(c *tb.Callback) {
	// 	// on inline button pressed (callback!)

	// 	// always respond!
	// 	bot.Respond(c, &tb.CallbackResponse{
	// 		Text: "Inline text here",
	// 	})
	// })

	// bot.Handle(handler.Register(bot))

	// bot.Handle(tb.OnQuery, func(q *tb.Query) {
	// 	log.Printf("Qurey: Text '%s', From: '%s'", q.Text, q.From.Username)
	// 	names := []string{
	// 		"Maksim Fedorov",
	// 		"Tanya",
	// 		"Ilia",
	// 	}

	// 	results := make(tb.Results, len(names)) // []tb.Result
	// 	for i, name := range names {
	// 		result := &tb.ArticleResult{
	// 			Title: "T" + name,
	// 			Text:  name,
	// 		}

	// 		results[i] = result
	// 		results[i].SetResultID(strconv.Itoa(i)) // It's needed to set a unique string ID for each result
	// 	}

	// 	err := bot.Answer(q, &tb.QueryResponse{
	// 		Results:   results,
	// 		CacheTime: 60, // a minute
	// 	})

	// 	if err != nil {
	// 		fmt.Println(err)
	// 	}
	// })
}
