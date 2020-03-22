package main

import (
	"fmt"
	"io/ioutil"

	"github.com/NicoNex/covidtron-19000/cache"
	"github.com/NicoNex/echotron"
)

const BOT_NAME = "covidtron-19000"

type bot struct {
	chatId int64
	echotron.Api
}

var cc *cache.Cache

func NewBot(engine echotron.Api, chatId int64) echotron.Bot {
	go cc.SaveSession(chatId)

	return &bot{
		chatId,
		engine,
	}
}

func (b *bot) Update(update *echotron.Update) {
	if update.Message.Text == "/start" {
		b.SendMessage("Hello world", b.chatId)
	}
}

func main() {
	token, err := ioutil.ReadFile("./token")
	if err != nil {
		fmt.Println("error: could not find token file")
		return
	}

	cc = cache.NewCache(BOT_NAME)
	dsp := echotron.NewDispatcher(string(token), NewBot)

	for _, id := range cc.GetSessions() {
		dsp.AddSession(id)
	}

	dsp.Run()
}
