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

package vax

import (
	"fmt"
	"log"
	"time"

	"github.com/NicoNex/covidtron-19000/apiutil"
)

func GetAndamentoMsg() string {
	ts := getTimestamp()
	totV := getTotaleVaccinati()
	totS := getTotaleSomministrazioni()

	return fmt.Sprintf(`*Andamento Nazionale Vaccinazioni*
_Dati aggiornati alle %s_

Totale somministrazioni: *%s*
Totale persone vaccinate: *%s* (*%s*)
_(persone che hanno completato il ciclo vaccinale)_`,
		apiutil.FormatTimestamp(ts, true),
		apiutil.Ifmt(totS),
		apiutil.Ifmt(totV),
		getTotalePercentuale(totV),
	)
}

func GetRegioneMsg(regName string) string {
	ts := getTimestamp()
	data := getRegione(regName)

	return fmt.Sprintf(`*Andamento Vaccinazioni - Regione %s*
_Dati aggiornati alle %s_

Dosi consegnate: *%s*
Dosi somministrate: *%s* (*%.1f%%*)`,
		regName,
		apiutil.FormatTimestamp(ts, true),
		apiutil.Ifmt(data.DosiConsegnate),
		apiutil.Ifmt(data.DosiSomministrate),
		data.Percentuale,
	)
}
