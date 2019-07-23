package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/cloudedcat/finance-bot/handler"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	bot, err := tb.NewBot(tb.Settings{
		Token: "837157414:AAFxMc8exxbi69prpVOAMRU_EHns7pnQFa4",
		// You can also set custom API URL. If field is empty it equals to "https://api.telegram.org"
		// URL: "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panic(err)
	}

	log.Println("Bot authorized")

	bot.Handle("/hello", func(m *tb.Message) {
		bot.Send(m.Sender, "hello world")
	})
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
	inlineBtn := tb.InlineButton{
		Unique: "sad_moon",
		Text:   "🌚 Button #2",
	}
	// inlineKeys := [][]tb.InlineButton{
	// 	[]tb.InlineButton{inlineBtn},
	// 	// ...
	// }

	// b.Handle(&replyBtn, func(m *tb.Message) {
	// 	// on reply button pressed
	// })

	bot.Handle(&inlineBtn, func(c *tb.Callback) {
		// on inline button pressed (callback!)

		// always respond!
		bot.Respond(c, &tb.CallbackResponse{
			Text: "Inline text here",
		})
	})

	bot.Handle(handler.Register(bot))

	bot.Handle(tb.OnQuery, func(q *tb.Query) {
		log.Printf("Qurey: Text '%s', From: '%s'", q.Text, q.From.Username)
		names := []string{
			"Maksim Fedorov",
			"Tanya",
			"Ilia",
		}

		results := make(tb.Results, len(names)) // []tb.Result
		for i, name := range names {
			result := &tb.ArticleResult{
				Title: "T" + name,
				Text:  name,
			}

			results[i] = result
			results[i].SetResultID(strconv.Itoa(i)) // It's needed to set a unique string ID for each result
		}

		err := bot.Answer(q, &tb.QueryResponse{
			Results:   results,
			CacheTime: 60, // a minute
		})

		if err != nil {
			fmt.Println(err)
		}
	})

	bot.Start()
}
