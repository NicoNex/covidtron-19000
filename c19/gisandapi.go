/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020 Alessandro Ianne.
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

	"github.com/thedevsaddam/gojsonq/v2"
)

func getNazione(nazione string) *GisandData {
	var data GisandData

	fpath := fmt.Sprintf("%s/gisanddata.json", jsonpath)
	search := gojsonq.New().
		File(fpath).
		WhereContains("country_region", nazione).
		First()

	if search == nil {
		return nil
	} 

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return &data
}

func GetNazioneMsg(nazione string) string {
	data := getNazione(nazione)

	if data == nil {
		log.Println(data)
		return "Nessun dato disponibile."
	}

	lastDay := len(data.DailyData) - 1 

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