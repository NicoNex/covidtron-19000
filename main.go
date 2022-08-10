/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020-2022 Nicolò Santamaria, Michele Dimaggio.
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
	"github.com/NicoNex/covidtron-19000/vax"
	"github.com/NicoNex/echotron/v3"
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
	cc *cache.Cache

	andamento    c19.InfoMsg
	andamentoVax string

	regione    c19.InfoMsg
	regioneVax string

	mainKbd = [][]echotron.KeyboardButton{
		{
			{Text: "🇮🇹 Andamento nazionale"},
		},
		{
			{Text: "🏙 Dati regione"},
			{Text: "🏢 Dati provincia"},
		},
	}

	andamentoKbd = [][]echotron.InlineKeyboardButton{
		{
			{Text: "📊 Generale", CallbackData: "andamento_generale"},
			{Text: "🧪 Tamponi", CallbackData: "andamento_tamponi"},
		},
		{
			{Text: "💉 Vaccini", CallbackData: "andamento_vaccini"},
			{Text: "📋 Note", CallbackData: "andamento_note"},
		},
	}

	regioneKbd = [][]echotron.InlineKeyboardButton{
		{
			{Text: "📊 Generale", CallbackData: "regione_generale"},
			{Text: "🧪 Tamponi", CallbackData: "regione_tamponi"},
		},
		{
			{Text: "💉 Vaccini", CallbackData: "regione_vaccini"},
			{Text: "📋 Note", CallbackData: "regione_note"},
		},
	}

	cancelBtn = []echotron.KeyboardButton{
		{Text: "❌ Annulla"},
	}

	masterKbd = []echotron.KeyboardButton{
		{Text: "📊 Utenti"},
		{Text: "📥 Aggiorna dati"},
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
		regName := extractText(update)
		regione = c19.GetRegioneMsg(regName)
		regioneVax = vax.GetRegioneMsg(regName)

		if !strings.Contains(regione.Generale, "Errore") {
			b.SendMessage(
				"Caricamento...",
				b.chatID,
				&echotron.MessageOptions{
					ReplyMarkup: echotron.ReplyKeyboardMarkup{
						Keyboard:       getMainKbd(b.chatID),
						ResizeKeyboard: true,
					},
				},
			)

			resp, err := b.SendMessage(
				regione.Generale,
				b.chatID,
				&echotron.MessageOptions{
					ParseMode: echotron.Markdown,
					ReplyMarkup: echotron.InlineKeyboardMarkup{
						InlineKeyboard: regioneKbd,
					},
				},
			)

			if err != nil {
				log.Println(err)
			} else {
				b.lastMsgID = resp.Result.ID
			}
		} else {
			b.SendMessage(
				regione.Generale,
				b.chatID,
				&echotron.MessageOptions{
					ReplyMarkup: echotron.ReplyKeyboardMarkup{
						Keyboard:       getMainKbd(b.chatID),
						ResizeKeyboard: true,
					},
				},
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
		b.SendMessage(
			c19.GetProvinciaMsg(extractText(update)),
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.ReplyKeyboardMarkup{
					Keyboard:       getMainKbd(b.chatID),
					ResizeKeyboard: true,
				},
			},
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
			b.SendMessage(
				"Scegli una provincia.",
				b.chatID,
				&echotron.MessageOptions{
					ReplyMarkup: echotron.ReplyKeyboardMarkup{
						Keyboard:       generateKeyboard(kbd),
						ResizeKeyboard: true,
					},
				},
			)
			return b.handleProvincia
		}

		b.SendMessage(
			"Errore: Regione non trovata.",
			b.chatID,
			&echotron.MessageOptions{
				ReplyMarkup: echotron.ReplyKeyboardMarkup{
					Keyboard:       getMainKbd(b.chatID),
					ResizeKeyboard: true,
				},
			},
		)
		return b.handleMessage
	}
}

func (b bot) sendUpgradeNotice() {
	for _, id := range cc.GetSessions() {
		b.SendMessage(
			"Covidtron-19000 è stato aggiornato! Scopri subito le novità!",
			id,
			&echotron.MessageOptions{
				ReplyMarkup: echotron.ReplyKeyboardMarkup{
					Keyboard:       getMainKbd(b.chatID),
					ResizeKeyboard: true,
				},
			},
		)
	}
}

func (b bot) handleMessage(update *echotron.Update) stateFn {
	switch text := extractText(update); {
	case text == "/start":
		b.sendIntroduction()

	case text == "🇮🇹 Andamento nazionale":
		andamento = c19.GetAndamentoMsg()
		andamentoVax = vax.GetAndamentoMsg()

		resp, err := b.SendMessage(
			andamento.Generale,
			b.chatID,
			&echotron.MessageOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: andamentoKbd,
				},
			},
		)

		if err != nil {
			log.Println(err)
		} else {
			b.lastMsgID = resp.Result.ID
		}

	case text == "🏙 Dati regione":
		b.SendMessage(
			"Scegli una regione.",
			b.chatID,
			&echotron.MessageOptions{
				ReplyMarkup: echotron.ReplyKeyboardMarkup{
					Keyboard:       generateKeyboard(c19.GetRegioni()),
					ResizeKeyboard: true,
				},
			},
		)
		return b.handleRegione

	case text == "🏢 Dati provincia":
		b.SendMessage(
			"Scegli una regione.",
			b.chatID,
			&echotron.MessageOptions{
				ReplyMarkup: echotron.ReplyKeyboardMarkup{
					Keyboard:       generateKeyboard(c19.GetRegioni()),
					ResizeKeyboard: true,
				},
			},
		)
		return b.chooseProvincia

	case text == "📊 Utenti" && isMaster(b.chatID):
		b.SendMessage(fmt.Sprintf("Utenti: %d", cc.CountSessions()), b.chatID, nil)

	case text == "📥 Aggiorna dati" && isMaster(b.chatID):
		b.SendMessage("Aggiornamento in corso...", b.chatID, nil)
		updateData()
		b.SendMessage(
			"Aggiornamento completato.",
			b.chatID,
			&echotron.MessageOptions{
				ReplyMarkup: echotron.ReplyKeyboardMarkup{
					Keyboard:       getMainKbd(b.chatID),
					ResizeKeyboard: true,
				},
			},
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
	var (
		resp echotron.APIResponseMessage
		err  error

		msgID = echotron.NewMessageID(b.chatID, b.lastMsgID)
	)

	switch update.CallbackQuery.Data {
	case "andamento_generale":
		resp, err = b.EditMessageText(
			andamento.Generale,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: andamentoKbd,
				},
			},
		)

	case "andamento_tamponi":
		resp, err = b.EditMessageText(
			andamento.Tamponi,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: andamentoKbd,
				},
			},
		)

	case "andamento_vaccini":
		resp, err = b.EditMessageText(
			andamentoVax,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: andamentoKbd,
				},
			},
		)

	case "andamento_note":
		resp, err = b.EditMessageText(
			andamento.Note,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: andamentoKbd,
				},
			},
		)

	case "regione_generale":
		resp, err = b.EditMessageText(
			regione.Generale,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: regioneKbd,
				},
			},
		)

	case "regione_tamponi":
		resp, err = b.EditMessageText(
			regione.Tamponi,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: regioneKbd,
				},
			},
		)

	case "regione_vaccini":
		resp, err = b.EditMessageText(
			regioneVax,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: regioneKbd,
				},
			},
		)

	case "regione_note":
		resp, err = b.EditMessageText(
			regione.Note,
			msgID,
			&echotron.MessageTextOptions{
				ParseMode: echotron.Markdown,
				ReplyMarkup: echotron.InlineKeyboardMarkup{
					InlineKeyboard: regioneKbd,
				},
			},
		)
	}

	if err != nil {
		log.Println(err)
	} else if resp.Result != nil {
		b.lastMsgID = resp.Result.ID
	}

	b.AnswerCallbackQuery(update.CallbackQuery.ID, nil)
}

func (b *bot) Update(update *echotron.Update) {
	if extractText(update) == "/cancel" {
		go b.SendMessage("Operazione annullata.", b.chatID, nil)
		b.state = b.handleMessage
		return
	}

	b.state = b.state(update)
}

func (b bot) sendIntroduction() {
	b.SendMessage(`*Benvenuto su Covidtron-19000!*

Covidtron-19000 ti aiuta a monitorare in tempo reale i dati sulla diffusione del COVID-19 in Italia condivisi dalla Protezione Civile.

Bot creato da @NicoNex e @Dj\_Mike238.
Basato su [echotron](https://github.com/NicoNex/echotron).

Icona creata da [Nhor Phai](https://www.flaticon.com/authors/nhor-phai) su [Flaticon](https://www.flaticon.com).`,
		b.chatID,
		&echotron.MessageOptions{
			ParseMode:             echotron.Markdown,
			DisableWebPagePreview: true,
			ReplyMarkup: echotron.InlineKeyboardMarkup{
				InlineKeyboard: [][]echotron.InlineKeyboardButton{
					{
						{Text: "☕️ Offrici un caffè", URL: "https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=HPUYKM3VJ2QMN&source=url"},
						{Text: "👾 GitHub Repository", URL: "https://github.com/NicoNex/covidtron-19000"},
					},
				},
			},
		},
	)

	b.SendMessage(
		"Seleziona un'opzione.",
		b.chatID,
		&echotron.MessageOptions{
			ReplyMarkup: echotron.ReplyKeyboardMarkup{
				Keyboard:       getMainKbd(b.chatID),
				ResizeKeyboard: true,
			},
		},
	)
}

func (b bot) sendCancel() {
	b.SendMessage(
		"Operazione annullata.",
		b.chatID,
		&echotron.MessageOptions{
			ReplyMarkup: echotron.ReplyKeyboardMarkup{
				Keyboard:       getMainKbd(b.chatID),
				ResizeKeyboard: true,
			},
		},
	)
}

func readToken() string {
	path := fmt.Sprintf("%s/.config/covidtron-19000/token", os.Getenv("HOME"))
	tok, err := os.ReadFile(path)
	if err != nil {
		log.Println("error: could not find token file")
	}
	return string(tok)
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
