package main

import (
	"log"
	"time"

	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	b, err := tb.NewBot(tb.Settings{
		Token: "837157414:AAFxMc8exxbi69prpVOAMRU_EHns7pnQFa4",
		// You can also set custom API URL. If field is empty it equals to "https://api.telegram.org"
		// URL: "http://195.129.111.17:8012",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		log.Panic(err)
		Fuck you nigga
	}

	b.Handle("/hello", func(m *tb.Message) {
		b.Send(m.Sender, "hello world")
	})
	replyBtn := tb.ReplyButton{Text: "🌕 Button #1"}
	replyKeys := [][]tb.ReplyButton{
		[]tb.ReplyButton{replyBtn},
		// ...
	}

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
	inlineKeys := [][]tb.InlineButton{
		[]tb.InlineButton{inlineBtn},
		// ...
	}

	b.Handle(&replyBtn, func(m *tb.Message) {
		// on reply button pressed
	})

	b.Handle(&inlineBtn, func(c *tb.Callback) {
		// on inline button pressed (callback!)

		// always respond!
		b.Respond(c, &tb.CallbackResponse{})
	})

	

	// Command: /start <PAYLOAD>
	b.Handle("/start", func(m *tb.Message) {
		if !m.Private() {
			return
		}

		b.Send(m.Sender, "Hello!", &tb.ReplyMarkup{
			ReplyKeyboard:  replyKeys,
			InlineKeyboard: inlineKeys,
		})
	})

	b.Start()
}
