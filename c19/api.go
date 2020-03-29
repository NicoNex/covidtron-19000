/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020 Michele Dimaggio, Alessandro Ianne.
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
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/thedevsaddam/gojsonq/v2"
)

func getNazione(nazione string) *GisandData {
	var data GisandData

	log.Println(nazione)
	fpath := fmt.Sprintf("%s/gisanddata.json", datapath)
	log.Println(fpath)
	search := gojsonq.New().
		File(fpath).
		Where("country_region", "=", nazione).
		First()

	if search == nil {
		return nil
	} 

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return &data
}

func getAndamento() Andamento {
	var data Andamento

	fpath := fmt.Sprintf("%s/andamento-nazionale.json", datapath)
	search := gojsonq.New().
		File(fpath).
		First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func getRegione(regione string) *Regione {
	var data Regione

	fpath := fmt.Sprintf("%s/regioni.json", datapath)
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

	fpath := fmt.Sprintf("%s/province.json", datapath)

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

func getNota(codice string) Nota {
	var data Nota

	fpath := fmt.Sprintf("%s/note-it.json", jsonpath)

	search := gojsonq.New().
		File(fpath).
		WhereEqual("codice", codice).
		First()

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func formatTimestamp(timestamp string) string {
	tp, err := time.Parse(time.RFC3339, timestamp+"Z")

	if err != nil {
		log.Println(err)
	}

	return tp.Format("15:04 del 02/01/2006")
}

func GetNazioneMsg(nazione string) string {
	data := getNazione(nazione)

	if data == nil {
		log.Println(data)
		return "Nessun dato disponibile."
	}

	lastDay := len(data.DailyData)

	msg := fmt.Sprintf(`*Andamento COVID-19 - %s*
_Dati aggiornati il: %s_`,
		data.Country_Region,
		data.DailyData[lastDay].Date)
	
	if data.Province_State != "" {
		msg = fmt.Sprintf("%s\nStato: %s", msg, data.Province_State)
	}

	msg = fmt.Sprintf("%s\nGuariti: %d\nDeceduti: %d\nTotale positivi: %d",
		msg,
		data.DailyData[lastDay].Recovered,
		data.DailyData[lastDay].Deaths,
		data.DailyData[lastDay].Confirmed)

	return msg
}

func formatNote(codici string) string {
	var noteData []Nota

	notes := strings.Split(codici, ";")

	for _, note := range notes {
		noteData = append(noteData, getNota(note))
	}

	msg := "\n\n*Note:*"

	for _, note := range noteData {
		msg += fmt.Sprintf("\n_%s - %s_\n%s", note.Regione, note.TipologiaAvviso, note.Note)
	}

	return msg
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
		msg += formatNote(data.NoteIt)
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
			msg += formatNote(data.NoteIt)
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
			msg += formatNote(data.NoteIt)
		}

		return msg
	} else {
		return "Errore: Provincia non trovata."
	}
}

func init() {
	datapath = fmt.Sprintf("%s/.cache/covidtron-19000", os.Getenv("HOME"))
}
