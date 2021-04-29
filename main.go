/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020-2021 Nicolò Santamaria, Michele Dimaggio.
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
	"log"
	"os"
	"time"

	"github.com/NicoNex/covidtron-19000/c19"
	"github.com/NicoNex/covidtron-19000/cache"
	"github.com/NicoNex/echotron/v2"
)

const BOT_NAME = "covidtron-19000"

// Recursive definition of the state-function type.
type stateFn func(*echotron.Update) stateFn

type bot struct {
	chatID int64
	state  stateFn
	echotron.API
}

var (
	cc *cache.Cache

	mainKbd = []echotron.KbdRow{
		[]echotron.Button{
			{Text: "🇮🇹 Andamento nazionale"},
		},
		[]echotron.Button{
			{Text: "🏙 Dati regione"},
			{Text: "🏢 Dati provincia"},
		},
	}

	cancelBtn = echotron.KbdRow{
		echotron.Button{Text: "❌ Annulla"},
	}
)

func newBot(chatID int64) echotron.Bot {
	go cc.SaveSession(chatID)

	b := &bot{
		chatID: chatID,
		API:    echotron.NewAPI(readToken()),
	}
	b.state = b.handleMessage
	return b
}

func (b bot) handleRegione(update *echotron.Update) stateFn {
	switch text := extractText(update); text {
	case "❌ Annulla":
		b.sendCancel()
		return b.handleMessage
	default:
		b.SendMessageWithKeyboard(
			c19.GetRegioneMsg(extractText(update)),
			b.chatID,
			b.KeyboardMarkup(true, false, false, mainKbd...),
			echotron.ParseMarkdown,
		)
		return b.handleMessage
	}
}

func (b bot) handleProvincia(update *echotron.Update) stateFn {
	switch text := extractText(update); text {
	case "❌ Annulla":
		b.sendCancel()
		return b.handleMessage
	default:
		b.SendMessageWithKeyboard(
			c19.GetProvinciaMsg(extractText(update)),
			b.chatID,
			b.KeyboardMarkup(true, false, false, mainKbd...),
			echotron.ParseMarkdown,
		)
		return b.handleMessage
	}
}

func (b bot) chooseProvincia(update *echotron.Update) stateFn {
	switch text := extractText(update); text {
	case "❌ Annulla":
		b.sendCancel()
		return b.handleMessage
	default:
		b.SendMessageWithKeyboard(
			"Scegli una provincia.",
			b.chatID,
			b.KeyboardMarkup(true, false, false, generateKeyboard(c19.GetProvince(text))...),
		)
		return b.handleProvincia
	}
}

func (b bot) handleMessage(update *echotron.Update) stateFn {
	switch text := extractText(update); text {
	case "/start":
		b.sendIntroduction()

	case "🇮🇹 Andamento nazionale":
		b.SendMessageWithKeyboard(
			c19.GetAndamentoMsg(),
			b.chatID,
			b.KeyboardMarkup(true, false, false, mainKbd...),
			echotron.ParseMarkdown,
		)

	case "🏙 Dati regione":
		b.SendMessageWithKeyboard(
			"Scegli una regione.",
			b.chatID,
			b.KeyboardMarkup(true, false, false, generateKeyboard(c19.GetRegioni())...),
		)
		return b.handleRegione

	case "🏢 Dati provincia":
		b.SendMessageWithKeyboard(
			"Scegli una regione.",
			b.chatID,
			b.KeyboardMarkup(true, false, false, generateKeyboard(c19.GetRegioni())...),
		)
		return b.chooseProvincia

	case "/users":
		b.SendMessage(fmt.Sprintf("Utenti: %d", cc.CountSessions()), b.chatID)
	}

	return b.handleMessage
}

func (b *bot) Update(update *echotron.Update) {
	if extractText(update) == "/cancel" {
		go b.SendMessage("Operazione annullata.", b.chatID)
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
		b.chatID,
		b.InlineKbdMarkup(
			b.InlineKbdRow(
				b.InlineKbdBtnURL("☕️ Offrici un caffè", "https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=HPUYKM3VJ2QMN&source=url"),
				b.InlineKbdBtnURL("👾 GitHub Repository", "https://github.com/NicoNex/covidtron-19000"),
			),
		),
		echotron.ParseMarkdown,
	)

	b.SendMessageWithKeyboard("Seleziona un'opzione.", b.chatID, b.KeyboardMarkup(true, false, false, mainKbd...))
}

func (b bot) sendCancel() {
	b.SendMessageWithKeyboard(
		"Operazione annullata.",
		b.chatID,
		b.KeyboardMarkup(true, false, false, mainKbd...),
	)
}

func generateKeyboard(values []string) []echotron.KbdRow {
	var kbd []echotron.KbdRow

	for i, v := range values {
		if i % 2 == 0 {
			kbd = append(kbd, []echotron.Button{})
		}

		kbd[len(kbd)-1] = append(kbd[len(kbd)-1], echotron.Button{Text: v})
	}

	return append(kbd, cancelBtn)
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
	tok, err := os.ReadFile(path)
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
