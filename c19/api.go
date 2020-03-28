/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020 Michele Dimaggio.
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

package c19

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/NicoNex/echotron"
	"github.com/thedevsaddam/gojsonq/v2"
)

var jsonpath string

func Update() {
	var json_url = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-%s-latest.json"
	files := [3]string{"andamento-nazionale", "province", "regioni"}

	dir := fmt.Sprintf(jsonpath)
	_, err := os.Stat(dir)
	if err != nil {
		os.Mkdir(dir, 0755)
	}

	for _, value := range files {
		var url = fmt.Sprintf(json_url, value)

		var content []byte = echotron.SendGetRequest(url)

		fpath := fmt.Sprintf("%s/%s.json", jsonpath, value)
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

	fpath := fmt.Sprintf("%s/andamento-nazionale.json", jsonpath)
	search := gojsonq.New().
			File(fpath).
			First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func getRegione(regione string) *Regione {
	var data Regione

	fpath := fmt.Sprintf("%s/regioni.json", jsonpath)
	search := gojsonq.New().
			File(fpath).
			WhereContains("denominazione_regione", regione).
			First()

	if search == nil {
		return nil
	}

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return &data
}

func getProvincia(provincia string) *Provincia {
	var data Provincia

	fpath := fmt.Sprintf("%s/province.json", jsonpath)

	var search interface{}

	if len(provincia) == 2 {
		search = gojsonq.New().
				File(fpath).
				WhereContains("sigla_provincia", provincia).
				First()
	} else if search == nil {
		search = gojsonq.New().
				File(fpath).
				WhereContains("denominazione_provincia", provincia).
				First()
	}

	if search == nil {
		return nil
	}

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return &data
}

func formatTimestamp(timestamp string) string {
	tp, err := time.Parse(time.RFC3339, timestamp+"Z")

	if err != nil {
		log.Println(err)
	}

	return tp.Format("15:04 del 02/01/2006")
}

func GetAndamentoMsg() string {
	data := getAndamento()

	msg := fmt.Sprintf(`*Andamento Nazionale COVID-19*
_Dati aggiornati alle %s_

Attualmente positivi: %d (%%2B%d da ieri)
Guariti: %d
Deceduti: %d
Totale positivi: %d

Tamponi totali: %d
Ricoverati con sintomi: %d
In terapia intensiva: %d
In isolamento domiciliare: %d
Totale ospedalizzati: %d`,
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

	if data != nil {
		msg := fmt.Sprintf(`*Andamento COVID-19 - Regione %s*
_Dati aggiornati alle %s_

Attualmente positivi: %d (%%2B%d da ieri)
Guariti: %d
Deceduti: %d
Totale positivi: %d

Tamponi totali: %d
Ricoverati con sintomi: %d
In terapia intensiva: %d
In isolamento domiciliare: %d
Totale ospedalizzati: %d`,
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
	} else {
		return "Errore: Regione non trovata."
	}
}

func GetProvinciaMsg(provincia string) string {
	data := getProvincia(provincia)

	if data != nil {
		msg := fmt.Sprintf(`*Andamento COVID-19 - Provincia di %s (%s)*
_Dati aggiornati alle %s_

Totale positivi: %d`,
			data.DenominazioneProvincia,
			data.DenominazioneRegione,
			formatTimestamp(data.Data),
			data.TotaleCasi,
		)

		if data.NoteIt != "" {
			msg = fmt.Sprintf("%s\n\nNote: %s", msg, data.NoteIt)
		}

		return msg
	} else {
		return "Errore: Provincia non trovata."
	}
}

func init() {
	jsonpath = fmt.Sprintf("%s/.cache/covidtron-19000", os.Getenv("HOME"))
}
