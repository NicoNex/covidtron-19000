/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2021 Michele Dimaggio.
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
	"fmt"

	"github.com/NicoNex/covidtron-19000/apiutil"
)

func GetAndamentoMsg() InfoMsg {
	andamento := getAndamento()
	note := getNote()

	return InfoMsg{
		Generale: getAndamentoGenerale(andamento),
		Tamponi:  getAndamentoTamponi(andamento),
		Note:     getAndamentoNote(andamento, note),
	}
}

func GetRegioneMsg(regName string) InfoMsg {
	regione := getRegione(regName)

	if regione != (Regione{}) {
		return InfoMsg{
			Generale: getRegioneGenerale(regione),
			Tamponi:  getRegioneTamponi(regione),
			Note:     getRegioneNote(regione),
		}
	}

	return InfoMsg{
		Generale: "Errore: Regione non trovata.",
	}
}

func getAndamentoGenerale(andamento Andamento) string {
	return fmt.Sprintf(`*Andamento Nazionale COVID-19*
_Dati aggiornati alle %s_

Attualmente positivi: *%s* (*%s* da ieri)
Guariti: *%s*
Deceduti: *%s*
Totale positivi: *%s* (*%s* da ieri)

Ricoverati con sintomi: *%s*
In terapia intensiva: *%s*
In isolamento domiciliare: *%s*
Totale ospedalizzati: *%s*`,
		apiutil.FormatTimestamp(andamento.Data, false),
		apiutil.Ifmt(andamento.TotalePositivi),
		plus(andamento.VariazioneTotalePositivi),
		apiutil.Ifmt(andamento.DimessiGuariti),
		apiutil.Ifmt(andamento.Deceduti),
		apiutil.Ifmt(andamento.TotaleCasi),
		plus(andamento.NuoviPositivi),
		apiutil.Ifmt(andamento.RicoveratiConSintomi),
		apiutil.Ifmt(andamento.TerapiaIntensiva),
		apiutil.Ifmt(andamento.IsolamentoDomiciliare),
		apiutil.Ifmt(andamento.TotaleOspedalizzati),
	)
}

func getAndamentoTamponi(andamento Andamento) string {
	return fmt.Sprintf(`*Andamento Nazionale COVID-19*
_Dati aggiornati alle %s_

Tamponi totali: *%s*
Soggetti sottoposti al tampone: *%s*
Positivi al tampone molecolare: *%s*
Tamponi molecolari totali: *%s*
Positivi al tampone antigenico: *%s*
Tamponi antigenici totali: *%s*`,
		apiutil.FormatTimestamp(andamento.Data, false),
		apiutil.Ifmt(andamento.Tamponi),
		apiutil.Ifmt(andamento.CasiTestati),
		apiutil.Ifmt(andamento.TotalePositiviTestMol),
		apiutil.Ifmt(andamento.TamponiTestMol),
		apiutil.Ifmt(andamento.TotalePositiviTestAnt),
		apiutil.Ifmt(andamento.TamponiTestAnt),
	)
}

func getAndamentoNote(andamento Andamento, note Nota) string {
	msg := fmt.Sprintf(`*Andamento Nazionale COVID-19*
_Dati aggiornati alle %s_`,
		apiutil.FormatTimestamp(andamento.Data, false),
	)

	if note.Data == andamento.Data {
		msg += formatNote(note.Note, Note)
	} else {
		msg += "\n\nNessuna nota disponibile."
	}

	return msg
}

func getRegioneGenerale(regione Regione) string {
	return fmt.Sprintf(`*Andamento COVID-19 - Regione %s*
_Dati aggiornati alle %s_

Attualmente positivi: *%s* (*%s* da ieri)
Guariti: *%s*
Deceduti: *%s*
Totale positivi: *%s* (*%s* da ieri)

Ricoverati con sintomi: *%s*
In terapia intensiva: *%s*
In isolamento domiciliare: *%s*
Totale ospedalizzati: *%s*`,
		regione.DenominazioneRegione,
		apiutil.FormatTimestamp(regione.Data, false),
		apiutil.Ifmt(regione.TotalePositivi),
		plus(regione.VariazioneTotalePositivi),
		apiutil.Ifmt(regione.DimessiGuariti),
		apiutil.Ifmt(regione.Deceduti),
		apiutil.Ifmt(regione.TotaleCasi),
		plus(regione.NuoviPositivi),
		apiutil.Ifmt(regione.RicoveratiConSintomi),
		apiutil.Ifmt(regione.TerapiaIntensiva),
		apiutil.Ifmt(regione.IsolamentoDomiciliare),
		apiutil.Ifmt(regione.TotaleOspedalizzati),
	)
}

func getRegioneTamponi(regione Regione) string {
	return fmt.Sprintf(`*Andamento COVID-19 - Regione %s*
_Dati aggiornati alle %s_

Tamponi totali: *%s*
Soggetti sottoposti al tampone: *%s*
Positivi al tampone molecolare: *%s*
Tamponi molecolari totali: *%s*
Positivi al tampone antigenico: *%s*
Tamponi antigenici totali: *%s*`,
		regione.DenominazioneRegione,
		apiutil.FormatTimestamp(regione.Data, false),
		apiutil.Ifmt(regione.Tamponi),
		apiutil.Ifmt(regione.CasiTestati),
		apiutil.Ifmt(regione.TotalePositiviTestMol),
		apiutil.Ifmt(regione.TamponiTestMol),
		apiutil.Ifmt(regione.TotalePositiviTestAnt),
		apiutil.Ifmt(regione.TamponiTestAnt),
	)
}

func getRegioneNote(regione Regione) string {
	note := false

	msg := fmt.Sprintf(`*Andamento COVID-19 - Regione %s*
_Dati aggiornati alle %s_`,
		regione.DenominazioneRegione,
		apiutil.FormatTimestamp(regione.Data, false),
	)

	if regione.Note != "" {
		msg += formatNote(regione.Note, Note)
		note = true
	}

	if regione.NoteCasi != "" {
		msg += formatNote(regione.NoteCasi, NoteCasi)
		note = true
	}

	if regione.NoteTest != "" {
		msg += formatNote(regione.NoteTest, NoteTest)
		note = true
	}

	if !note {
		msg += "\n\nNessuna nota disponibile."
	}

	return msg
}

func GetProvinciaMsg(provincia string) string {
	var data = getProvincia(provincia)

	if data != (Provincia{}) {
		msg := fmt.Sprintf(`*Andamento COVID-19 - Provincia di %s (%s)*
_Dati aggiornati alle %s_

Totale positivi: *%s*`,
			data.DenominazioneProvincia,
			data.DenominazioneRegione,
			apiutil.FormatTimestamp(data.Data, false),
			apiutil.Ifmt(data.TotaleCasi),
		)

		if data.Note != "" {
			msg += formatNote(data.Note, Note)
		}

		return msg
	} else {
		return "Errore: Provincia non trovata."
	}
}
