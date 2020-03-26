/*
* covidtron-19000
* Copyright (C) 2020  Michele Dimaggio
*/

package C19

import (
	"io"
	"os"
	"fmt"
	"log"
	"time"
	"bytes"
	"encoding/json"

	"github.com/NicoNex/echotron"
	"github.com/thedevsaddam/gojsonq/v2"
)

const JSON_PATH = "~/.config/covidtron-19000"

func Update() {
	var json_url = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-%s-latest.json"
	files := [3]string{"andamento-nazionale", "province", "regioni"}

	for _, value := range files {
		var url = fmt.Sprintf(json_url, value)

		var content []byte = echotron.SendGetRequest(url)

		fpath := fmt.Sprintf("%s/%s.json", JSON_PATH, value)
		data, err := os.Create(fpath)

		if err != nil {
			log.Println(err)
		}
		defer data.Close()

		_, err = io.Copy(data, bytes.NewReader(content))

		if err != nil {
			log.Println(err)
		}
	}
}

func getAndamento() Andamento {
	var data Andamento

	fpath := fmt.Sprintf("%s/andamento-nazionale.json", JSON_PATH)
	search := gojsonq.New().File(fpath).First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func getRegione(regione string) Regione {
	var data Regione

	fpath := fmt.Sprintf("%s/regioni.json", JSON_PATH)
	search := gojsonq.New().File(fpath).Where("denominazione_regione", "=", regione).First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func getProvincia(provincia string) Provincia {
	var data Provincia

	fpath := fmt.Sprintf("%s/province.json", JSON_PATH)
	search := gojsonq.New().File(fpath).Where("denominazione_provincia", "=", provincia).First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func formatTimestamp(timestamp string) string {
	timestamp = timestamp + "Z"

	tp, err := time.Parse(time.RFC3339, timestamp)

	if err != nil {
		log.Println(err)
	}

	return tp.Format("02/01/2006 15:04")
}

func GetAndamentoMsg() string {
	data := getAndamento()

	msg := fmt.Sprintf(`
		*Andamento Nazionale COVID-19*
		_Ultimo aggiornamento: %s_

		Attualmente positivi: %d (+%d da ieri)
		Guariti: %d
		Deceduti: %d
		Totale positivi: %d

		Tamponi totali: %d
		Ricoverati con sintomi: %d
		In terapia intensiva: %d
		In isolamento domiciliare: %d
		Totale ospedalizzati: %d
		`,
		formatTimestamp(data.Data),
		data.TotaleAttualmentePositivi,
		data.NuoviAttualmentePositivi,
		data.DimessiGuariti,
		data.Deceduti,
		data.TotaleCasi,
		data.Tamponi,
		data.RicoveratiConSintomi,
		data.TerapiaIntensiva,
		data.IsolamentoDomiciliare,
		data.TotaleOspedalizzati,
	)

	if data.NoteIt != "" {
		msg = fmt.Sprintf("%s\n\nNote: %s", msg, data.NoteIt)
	}

	return msg
}

func GetRegioneMsg(regione string) string {
	data := getRegione(regione)

	msg := fmt.Sprintf(`
		*Andamento COVID-19 - Regione %s*
		_Ultimo aggiornamento: %s_

		Attualmente positivi: %d (+%d da ieri)
		Guariti: %d
		Deceduti: %d
		Totale positivi: %d

		Tamponi totali: %d
		Ricoverati con sintomi: %d
		In terapia intensiva: %d
		In isolamento domiciliare: %d
		Totale ospedalizzati: %d
		`,
		data.DenominazioneRegione,
		formatTimestamp(data.Data),
		data.TotaleAttualmentePositivi,
		data.NuoviAttualmentePositivi,
		data.DimessiGuariti,
		data.Deceduti,
		data.TotaleCasi,
		data.Tamponi,
		data.RicoveratiConSintomi,
		data.TerapiaIntensiva,
		data.IsolamentoDomiciliare,
		data.TotaleOspedalizzati,
	)

	if data.NoteIt != "" {
		msg = fmt.Sprintf("%s\n\nNote: %s", msg, data.NoteIt)
	}

	return msg
}

func GetProvinciaMsg(provincia string) string {
	data := getProvincia(provincia)

	msg := fmt.Sprintf(`
		*Andamento COVID-19 - Provincia di %s (%s)*
		_Ultimo aggiornamento: %s_

		Totale positivi: %d
		`,
		data.DenominazioneProvincia,
		data.DenominazioneRegione,
		formatTimestamp(data.Data),
		data.TotaleCasi,
	)

	if data.NoteIt != "" {
		msg = fmt.Sprintf("%s\n\nNote: %s", msg, data.NoteIt)
	}

	return msg
}
