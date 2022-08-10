package main

import (
	"time"

	"github.com/NicoNex/covidtron-19000/c19"
	"github.com/NicoNex/covidtron-19000/cache"
	"github.com/NicoNex/covidtron-19000/vax"
	"github.com/NicoNex/echotron/v3"
)

func generateKeyboard(values []string) (kbd [][]echotron.KeyboardButton) {
	for i, v := range values {
		if i%2 == 0 {
			kbd = append(kbd, []echotron.KeyboardButton{})
		}

		kbd[len(kbd)-1] = append(kbd[len(kbd)-1], echotron.KeyboardButton{Text: v})
	}

	return append(kbd, cancelBtn)
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

func getMainKbd(chatID int64) [][]echotron.KeyboardButton {
	if isMaster(chatID) {
		return append(mainKbd, masterKbd)
	}
	return mainKbd
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

	if saved.Vax != latest.Vax {
		vax.Update()
	}

	cc.SaveCommits(latest)
}
