package main

import (
    "fmt"
    "io/ioutil"

    "github.com/NicoNex/echotron"
)

type bot struct {
    chatId int64
    echotron.Api
}

func newBot(engine echotron.Api, chatId int64) echotron.Bot {
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

    dsp := echotron.NewDispatcher(string(token), newBot)
    dsp.Run()
}
