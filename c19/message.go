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
	"log"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
)

var (
	andamento Andamento
	note      Nota
	regione   *Regione
	provincia Provincia
)

func GetAndamentoMsg() InfoMsg {
	andamento = getAndamento()
	note = getNote()

	return InfoMsg{
		Generale: getAndamentoGenerale(),
		Tamponi:  getAndamentoTamponi(),
		Note:     getAndamentoNote(),
	}
}

func GetRegioneMsg(regName string) InfoMsg {
	regione = getRegione(regName)

	if regione != nil {
		return InfoMsg{
			Generale: getRegioneGenerale(),
			Tamponi:  getRegioneTamponi(),
			Note:     getRegioneNote(),
		}
	}

	return InfoMsg{
		Generale: "Errore: Regione non trovata.",
	}
}

func getAndamentoGenerale() string {
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
		formatTimestamp(andamento.Data),
		ifmt(andamento.TotalePositivi),
		plus(andamento.VariazioneTotalePositivi),
		ifmt(andamento.DimessiGuariti),
		ifmt(andamento.Deceduti),
		ifmt(andamento.TotaleCasi),
		plus(andamento.NuoviPositivi),
		ifmt(andamento.RicoveratiConSintomi),
		ifmt(andamento.TerapiaIntensiva),
		ifmt(andamento.IsolamentoDomiciliare),
		ifmt(andamento.TotaleOspedalizzati),
	)
}

func getAndamentoTamponi() string {
	return fmt.Sprintf(`*Andamento Nazionale COVID-19*
_Dati aggiornati alle %s_

Tamponi totali: *%s*
Soggetti sottoposti al tampone: *%s*
Positivi al tampone molecolare: *%s*
Tamponi molecolari totali: *%s*
Positivi al tampone antigenico: *%s*
Tamponi antigenici totali: *%s*`,
		formatTimestamp(andamento.Data),
		ifmt(andamento.Tamponi),
		ifmt(andamento.CasiTestati),
		ifmt(andamento.TotalePositiviTestMol),
		ifmt(andamento.TamponiTestMol),
		ifmt(andamento.TotalePositiviTestAnt),
		ifmt(andamento.TamponiTestAnt),
	)
}

func getAndamentoNote() string {
	msg := fmt.Sprintf(`*Andamento Nazionale COVID-19*
_Dati aggiornati alle %s_`,
		formatTimestamp(andamento.Data),
	)

	if note.Data == andamento.Data {
		msg += formatNote(note.Note, Note)
	} else {
		msg += "\n\nNessuna nota disponibile."
	}

	return msg
}

func getRegioneGenerale() string {
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
		formatTimestamp(regione.Data),
		ifmt(regione.TotalePositivi),
		plus(regione.VariazioneTotalePositivi),
		ifmt(regione.DimessiGuariti),
		ifmt(regione.Deceduti),
		ifmt(regione.TotaleCasi),
		plus(regione.NuoviPositivi),
		ifmt(regione.RicoveratiConSintomi),
		ifmt(regione.TerapiaIntensiva),
		ifmt(regione.IsolamentoDomiciliare),
		ifmt(regione.TotaleOspedalizzati),
	)
}

func getRegioneTamponi() string {
	return fmt.Sprintf(`*Andamento COVID-19 - Regione %s*
_Dati aggiornati alle %s_

Tamponi totali: *%s*
Soggetti sottoposti al tampone: *%s*
Positivi al tampone molecolare: *%s*
Tamponi molecolari totali: *%s*
Positivi al tampone antigenico: *%s*
Tamponi antigenici totali: *%s*`,
		regione.DenominazioneRegione,
		formatTimestamp(regione.Data),
		ifmt(regione.Tamponi),
		ifmt(regione.CasiTestati),
		ifmt(regione.TotalePositiviTestMol),
		ifmt(regione.TamponiTestMol),
		ifmt(regione.TotalePositiviTestAnt),
		ifmt(regione.TamponiTestAnt),
	)
}

func getRegioneNote() string {
	note := false

	msg := fmt.Sprintf(`*Andamento COVID-19 - Regione %s*
_Dati aggiornati alle %s_`,
		regione.DenominazioneRegione,
		formatTimestamp(regione.Data),
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

	if data != nil {
		msg := fmt.Sprintf(`*Andamento COVID-19 - Provincia di %s (%s)*
_Dati aggiornati alle %s_

Totale positivi: *%s*`,
			data.DenominazioneProvincia,
			data.DenominazioneRegione,
			formatTimestamp(data.Data),
			ifmt(data.TotaleCasi),
		)

		if data.Note != "" {
			msg += formatNote(data.Note, Note)
		}

		return msg
	} else {
		return "Errore: Provincia non trovata."
	}
}

func formatTimestamp(timestamp string) string {
	tp, err := time.Parse(time.RFC3339, timestamp+"Z")

	if err != nil {
		log.Println(err)
	}

	return tp.Format("15:04 del 02/01/2006")
}

func formatNote(nota string, ntype NoteType) string {
	var msg strings.Builder

	msg.WriteString("\n\n*Note")

	switch ntype {
	case Note:
		msg.WriteString(" generali")
	case NoteCasi:
		msg.WriteString(" relative ai test effettuati")
	case NoteTest:
		msg.WriteString(" relative ai casi testati")
	}

	msg.WriteString(":*")

	note := strings.Split(nota, ". ")

	for i, n := range note {
		n = strings.TrimSuffix(n, "  ")

		if !strings.HasSuffix(n, ".") {
			n += "."
		}

		if strings.Contains(n, "  -") {
			spl := strings.Split(n, "  -")

			for _, s := range spl {
				if strings.HasPrefix(s, " ") {
					s = "-" + s
				}

				msg.WriteString(fmt.Sprintf("\n%s", s))
			}
		} else if strings.TrimSpace(n) != "." {
			if i == 0 || (i > 0 && len(note[i-1]) != 6) {
				msg.WriteString(fmt.Sprintf("\n- %s", n))
			} else {
				msg.WriteString(fmt.Sprintf(" %s", n))
			}
		}
	}

	return msg.String()
}

func plus(value int) string {
	if value > 0 {
		return "+" + ifmt(value)
	}
	return ifmt(value)
}

func ifmt(i int) string {
	return humanize.FormatInteger("#.###,", i)
}
