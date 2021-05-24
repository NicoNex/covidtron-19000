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
	"strings"
	"time"

	"github.com/NicoNex/covidtron-19000/c19"
	"github.com/NicoNex/covidtron-19000/cache"
	"github.com/NicoNex/echotron/v2"
)

const BOT_NAME = "covidtron-19000"

// Recursive definition of the state-function type.
type stateFn func(*echotron.Update) stateFn

type bot struct {
	chatID    int64
	lastMsgID int
	state     stateFn
	echotron.API
}

var (
	andamento c19.InfoMsg
	cc        *cache.Cache
	regione   c19.InfoMsg

	mainKbd = []echotron.KbdRow{
		[]echotron.Button{
			{Text: "🇮🇹 Andamento nazionale"},
		},
		[]echotron.Button{
			{Text: "🏙 Dati regione"},
			{Text: "🏢 Dati provincia"},
		},
	}

	andamentoKbd = []echotron.InlineKbdRow{
		[]echotron.InlineButton{
			{Text: "📊 Generale", CallbackData: "andamento_generale"},
			{Text: "🧪 Tamponi", CallbackData: "andamento_tamponi"},
		},
		[]echotron.InlineButton{
			{Text: "📋 Note", CallbackData: "andamento_note"},
		},
	}

	regioneKbd = []echotron.InlineKbdRow{
		[]echotron.InlineButton{
			{Text: "📊 Generale", CallbackData: "regione_generale"},
			{Text: "🧪 Tamponi", CallbackData: "regione_tamponi"},
		},
		[]echotron.InlineButton{
			{Text: "📋 Note", CallbackData: "regione_note"},
		},
	}

	cancelBtn = echotron.KbdRow{
		echotron.Button{Text: "❌ Annulla"},
	}

	masterKbd = echotron.KbdRow{
		echotron.Button{Text: "📊 Utenti"},
		echotron.Button{Text: "📥 Aggiorna dati"},
	}

	masters = []int64{41876271, 14870908}
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

	default:
		regione = c19.GetRegioneMsg(extractText(update))

		if !strings.Contains(regione.Generale, "Errore") {
			b.SendMessageWithKeyboard(
				"Caricamento...",
				b.chatID,
				b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...),
			)

			resp, err := b.SendMessageWithKeyboard(
				regione.Generale,
				b.chatID,
				b.InlineKbdMarkup(regioneKbd...),
				echotron.ParseMarkdown,
			)

			if err != nil {
				log.Println(err)
			} else {
				b.lastMsgID = resp.Result.ID
			}
		} else {
			b.SendMessageWithKeyboard(
				regione.Generale,
				b.chatID,
				b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...),
			)
		}
	}

	return b.handleMessage
}

func (b bot) handleProvincia(update *echotron.Update) stateFn {
	switch text := extractText(update); text {
	case "❌ Annulla":
		b.sendCancel()

	default:
		b.SendMessageWithKeyboard(
			c19.GetProvinciaMsg(extractText(update)),
			b.chatID,
			b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...),
			echotron.ParseMarkdown,
		)
	}

	return b.handleMessage
}

func (b bot) chooseProvincia(update *echotron.Update) stateFn {
	switch text := extractText(update); text {
	case "❌ Annulla":
		b.sendCancel()
		return b.handleMessage

	default:
		kbd := c19.GetProvince(text)

		if kbd != nil {
			b.SendMessageWithKeyboard(
				"Scegli una provincia.",
				b.chatID,
				b.KeyboardMarkup(true, false, false, generateKeyboard(kbd)...),
			)
			return b.handleProvincia
		}

		b.SendMessageWithKeyboard(
			"Errore: Regione non trovata.",
			b.chatID,
			b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...),
		)
		return b.handleMessage
	}
}

func (b bot) sendUpgradeNotice() {
	for _, id := range cc.GetSessions() {
		b.SendMessageWithKeyboard(
			"Covidtron-19000 è stato aggiornato! Scopri subito le novità!",
			id,
			b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...),
		)
	}
}

func (b bot) handleMessage(update *echotron.Update) stateFn {
	switch text := extractText(update); {
	case text == "/start":
		b.sendIntroduction()

	case text == "🇮🇹 Andamento nazionale":
		andamento = c19.GetAndamentoMsg()

		resp, err := b.SendMessageWithKeyboard(
			andamento.Generale,
			b.chatID,
			b.InlineKbdMarkup(andamentoKbd...),
			echotron.ParseMarkdown,
		)

		if err != nil {
			log.Println(err)
		} else {
			b.lastMsgID = resp.Result.ID
		}

	case text == "🏙 Dati regione":
		b.SendMessageWithKeyboard(
			"Scegli una regione.",
			b.chatID,
			b.KeyboardMarkup(true, false, false, generateKeyboard(c19.GetRegioni())...),
		)
		return b.handleRegione

	case text == "🏢 Dati provincia":
		b.SendMessageWithKeyboard(
			"Scegli una regione.",
			b.chatID,
			b.KeyboardMarkup(true, false, false, generateKeyboard(c19.GetRegioni())...),
		)
		return b.chooseProvincia

	case text == "📊 Utenti" && isMaster(b.chatID):
		b.SendMessage(fmt.Sprintf("Utenti: %d", cc.CountSessions()), b.chatID)

	case text == "📥 Aggiorna dati" && isMaster(b.chatID):
		b.SendMessage("Aggiornamento in corso...", b.chatID)
		updateData()
		b.SendMessageWithKeyboard(
			"Aggiornamento completato.",
			b.chatID,
			b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...),
		)

	case text == "/notice" && isMaster(b.chatID):
		b.sendUpgradeNotice()
	}

	if update.CallbackQuery != nil {
		b.handleCallback(update)
	}

	return b.handleMessage
}

func (b bot) handleCallback(update *echotron.Update) {
	var resp echotron.APIResponseMessage
	var err error

	switch update.CallbackQuery.Data {
	case "andamento_generale":
		resp, err = b.EditMessageTextWithKeyboard(
			b.chatID,
			b.lastMsgID,
			andamento.Generale,
			b.InlineKbdMarkup(andamentoKbd...),
			echotron.ParseMarkdown,
		)

	case "andamento_tamponi":
		resp, err = b.EditMessageTextWithKeyboard(
			b.chatID,
			b.lastMsgID,
			andamento.Tamponi,
			b.InlineKbdMarkup(andamentoKbd...),
			echotron.ParseMarkdown,
		)

	case "andamento_note":
		resp, err = b.EditMessageTextWithKeyboard(
			b.chatID,
			b.lastMsgID,
			andamento.Note,
			b.InlineKbdMarkup(andamentoKbd...),
			echotron.ParseMarkdown,
		)
	case "regione_generale":
		resp, err = b.EditMessageTextWithKeyboard(
			b.chatID,
			b.lastMsgID,
			regione.Generale,
			b.InlineKbdMarkup(regioneKbd...),
			echotron.ParseMarkdown,
		)

	case "regione_tamponi":
		resp, err = b.EditMessageTextWithKeyboard(
			b.chatID,
			b.lastMsgID,
			regione.Tamponi,
			b.InlineKbdMarkup(regioneKbd...),
			echotron.ParseMarkdown,
		)

	case "regione_note":
		resp, err = b.EditMessageTextWithKeyboard(
			b.chatID,
			b.lastMsgID,
			regione.Note,
			b.InlineKbdMarkup(regioneKbd...),
			echotron.ParseMarkdown,
		)
	}

	if err != nil {
		log.Println(err)
	} else if resp.Result != nil {
		b.lastMsgID = resp.Result.ID
	}

	b.AnswerCallbackQuery(update.CallbackQuery.ID, "", false)
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
	b.SendMessageWithKeyboard(`*Benvenuto su Covidtron-19000!*

Covidtron-19000 ti aiuta a monitorare in tempo reale i dati sulla diffusione del COVID-19 in Italia condivisi dalla Protezione Civile.

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
		echotron.DisableWebPagePreview,
	)

	b.SendMessageWithKeyboard("Seleziona un'opzione.", b.chatID, b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...))
}

func (b bot) sendCancel() {
	b.SendMessageWithKeyboard(
		"Operazione annullata.",
		b.chatID,
		b.KeyboardMarkup(true, false, false, getMainKbd(b.chatID)...),
	)
}

func generateKeyboard(values []string) []echotron.KbdRow {
	var kbd []echotron.KbdRow

	for i, v := range values {
		if i%2 == 0 {
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

	latest := cc.UpdateCommits()
	saved := cc.GetCommits()

	if saved.C19 != latest.C19 {
		c19.Update()
	}

	cc.SaveCommits(latest)
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

func isMaster(chatID int64) bool {
	for _, i := range masters {
		if i == chatID {
			return true
		}
	}
	return false
}

func getMainKbd(chatID int64) []echotron.KbdRow {
	if isMaster(chatID) {
		return append(mainKbd, masterKbd)
	}
	return mainKbd
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
