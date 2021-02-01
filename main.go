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
	"log"
	"os"
	"time"

	"github.com/NicoNex/covidtron-19000/c19"
	"github.com/NicoNex/covidtron-19000/cache"
	"github.com/NicoNex/echotron"
)

const BOT_NAME = "covidtron-19000"

// Recursive definition of the state-function type.
type stateFn func(*echotron.Update) stateFn

type bot struct {
	chatId int64
	state  stateFn
	echotron.Api
}

var cc *cache.Cache

func newBot(chatId int64) echotron.Bot {
	go cc.SaveSession(chatId)

	b := &bot{
		chatId: chatId,
		Api:    echotron.NewApi(readToken()),
	}
	b.state = b.handleMessage
	return b
}

func (b bot) handleRegione(update *echotron.Update) stateFn {
	b.SendMessage(
		c19.GetRegioneMsg(extractText(update)),
		b.chatId,
		echotron.PARSE_MARKDOWN,
	)
	return b.handleMessage
}

func (b bot) handleProvincia(update *echotron.Update) stateFn {
	b.SendMessage(
		c19.GetProvinciaMsg(extractText(update)),
		b.chatId,
		echotron.PARSE_MARKDOWN,
	)
	return b.handleMessage
}

func (b bot) handleMessage(update *echotron.Update) stateFn {
	switch cmd := extractText(update); cmd {
	case "/start":
		b.sendIntroduction()

	case "/andamento":
		b.SendMessage(c19.GetAndamentoMsg(), b.chatId, echotron.PARSE_MARKDOWN)

	case "/regione":
		b.SendMessage("Inserisci il nome di una regione.", b.chatId)
		return b.handleRegione

	case "/provincia":
		b.SendMessage("Inserisci il nome di una provincia.", b.chatId)
		return b.handleProvincia

	case "/users":
		b.SendMessage(fmt.Sprintf("Utenti: %d", cc.CountSessions()), b.chatId)
	}

	return b.handleMessage
}

func (b *bot) Update(update *echotron.Update) {
	if extractText(update) == "/cancel" {
		go b.SendMessage("Operazione annullata.", b.chatId)
		b.state = b.handleMessage
		return
	}

	b.state = b.state(update)
}

func (b bot) sendIntroduction() {
	b.SendMessageWithKeyboard(`*Benvenuto in Covidtron-19000!*

*Comandi:*
/start: visualizza questo messaggio
/andamento: visualizza andamento nazionale
/regione: visualizza andamento regione
/provincia: visualizza andamento provincia
/cancel: annulla l'operazione in corso

Bot creato da @NicoNex e @Dj\_Mike238.
Basato su [echotron](https://github.com/NicoNex/echotron).

Icona creata da [Nhor Phai](https://www.flaticon.com/authors/nhor-phai) su [Flaticon](https://www.flaticon.com).`,
		b.chatId,
		b.InlineKbdMarkup(
			b.InlineKbdRow(
				b.InlineKbdBtn("☕️ Offrici un caffè", "https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=HPUYKM3VJ2QMN&source=url", ""),
				b.InlineKbdBtn("👾 GitHub Repository", "https://github.com/NicoNex/covidtron-19000", ""),
			),
		),
	)
}

func ticker(tch <-chan time.Time) {
	for t := range tch {
		if t.Hour() >= 16 && t.Hour() <= 19 {
			updateData()
		}
	}
}

func updateData() {
	cc = cache.LoadCache(BOT_NAME)

	sha := cc.GetSha()
	latest := cc.GetLatestCommit()

	if latest != sha {
		c19.Update()
		cc.SaveLatestCommit(sha)
	}
}

func readToken() string {
	path := fmt.Sprintf("%s/.config/covidtron-19000/token", os.Getenv("HOME"))
	tok, err := ioutil.ReadFile(path)
	if err != nil {
		log.Println("error: could not find token file")
	}
	return string(tok)
}

func extractText(update *echotron.Update) string {
	if update.Message != nil {
		return update.Message.Text
	} else if update.EditedMessage != nil {
		return update.EditedMessage.Text
	}
	return ""
}

func main() {
	updateData()
	go ticker(time.Tick(time.Minute * 10))

	dsp := echotron.NewDispatcher(readToken(), newBot)

	for _, id := range cc.GetSessions() {
		dsp.AddSession(id)
	}

	dsp.Poll()
}
