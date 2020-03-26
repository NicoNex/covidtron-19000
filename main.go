/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020 Nicolò Santamaria, Michele Dimaggio.
 *
 * Covidtron-19000 is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Covidtron-19000 is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"io/ioutil"

	"github.com/NicoNex/covidtron-19000/cache"
	"github.com/NicoNex/covidtron-19000/c19"
	"github.com/NicoNex/echotron"
)

const BOT_NAME = "covidtron-19000"

type bot struct {
	chatId int64
	echotron.Api
}

var cc *cache.Cache

type botStates int

const (
	idle = iota
	regione
	provincia
)

var botState = botStates.idle

func NewBot(engine echotron.Api, chatId int64) echotron.Bot {
	go cc.SaveSession(chatId)

	return &bot{
		chatId,
		engine,
	}
}

func (b *bot) Update(update *echotron.Update) {
	switch botState {
	case botStates.idle:
		if update.Message.Text == "/start" {
			sendIntroduction()

		} else if update.Message.Text == "/andamento" {
			b.SendMessageOptions(c19.GetAndamentoMsg(), b.chatId, echotron.PARSE_MARKDOWN)

		} else if update.Message.Text == "/regione" {
			b.SendMessage("Inserisci il nome di una regione.", b.chatId)
			botState = botStates.regione

		} else if update.Message.Text == "/provincia" {
			b.SendMessage("Inserisci il nome di una provincia.", b.chatId)
			botState = botStates.provincia
		}

	case botStates.regione:
		if update.Message.Text == "/cancel" {
			b.SendMessage("Operazione annullata.", b.chatId)

		} else {
			b.SendMessageOptions(c19.GetRegioneMsg(update.Message.Text), b.chatId, echotron.PARSE_MARKDOWN)
		}

		botState = botStates.idle

	case botStates.provincia:
		if update.Message.Text == "/cancel" {
			b.SendMessage("Operazione annullata.", b.chatId)

		} else {
			b.SendMessageOptions(c19.GetProvinciaMsg(update.Message.Text), b.chatId, echotron.PARSE_MARKDOWN)
		}

		botState = botStates.idle
	}
}

func sendIntroduction() {
	msg := `
	*Benvenuto in Covidtron-19000!*

	*Comandi:*
	/start: visualizza questo messaggio
	/andamento: visualizza andamento nazionale
	/regione: visualizza andamento regione
	/provincia: visualizza andamento provincia
	/cancel: annulla l'operazione in corso

	Bot creato da @NicoNex e @Dj_Mike238.
	Basato su [echotron](https://github.com/NicoNex/echotron).

	Icona creata da [Nhor Phai](https://www.flaticon.com/authors/nhor-phai) su [Flaticon](https://www.flaticon.com)
	`

	b.SendMessageOptions(msg, b.chatId, echotron.PARSE_MARKDOWN | echotron.DISABLE_WEB_PAGE_PREVIEW)
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
