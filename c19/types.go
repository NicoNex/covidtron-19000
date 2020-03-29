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
	"fmt"
	"strconv"
	"strings"
	"encoding/json"
	"encoding/csv"
	"errors"
	"bytes"
)

type Nazione struct {
	FIPS			string		`mapstructure:"FIPS"`
	Admin2			string	`mapstructure:"Admin2"`
	Province_State	string	`mapstructure:"Province_State"`
	Country_Region	string	`mapstructure:"Country_Region"`
	Last_Update		string	`mapstructure:"Last_Update"`
	Lat 			string	`mapstructure:"Lat"`
	Long_ 			string	`mapstructure:"Long_"`
	Confirmed		string 	`mapstructure:"Confirmed"`
	Deaths			string		`mapstructure:"Deaths"`
	Recovered		string		`mapstructure:"Recovered"`
	Active			string		`mapstructure:"Active"`
	Combined_Key	string	`mapstructure:"Combined_Key"`
	
	
}

type Andamento struct {
	Data                      string `json:"data"`
	Deceduti                  int    `json:"deceduti"`
	DimessiGuariti            int    `json:"dimessi_guariti"`
	IsolamentoDomiciliare     int    `json:"isolamento_domiciliare"`
	NoteIt                    string `json:"note_it"`
	NuoviAttualmentePositivi  int    `json:"nuovi_attualmente_positivi"`
	RicoveratiConSintomi      int    `json:"ricoverati_con_sintomi"`
	Tamponi                   int    `json:"tamponi"`
	TerapiaIntensiva          int    `json:"terapia_intensiva"`
	TotaleAttualmentePositivi int    `json:"totale_attualmente_positivi"`
	TotaleCasi                int    `json:"totale_casi"`
	TotaleOspedalizzati       int    `json:"totale_ospedalizzati"`
}

type Regione struct {
	Data                      string `json:"data"`
	DenominazioneRegione      string `json:"denominazione_regione"`
	RicoveratiConSintomi      int    `json:"ricoverati_con_sintomi"`
	TerapiaIntensiva          int    `json:"terapia_intensiva"`
	TotaleOspedalizzati       int    `json:"totale_ospedalizzati"`
	IsolamentoDomiciliare     int    `json:"isolamento_domiciliare"`
	TotaleAttualmentePositivi int    `json:"totale_attualmente_positivi"`
	NuoviAttualmentePositivi  int    `json:"nuovi_attualmente_positivi"`
	DimessiGuariti            int    `json:"dimessi_guariti"`
	Deceduti                  int    `json:"deceduti"`
	TotaleCasi                int    `json:"totale_casi"`
	Tamponi                   int    `json:"tamponi"`
	NoteIt                    string `json:"note_it"`
}

type Provincia struct {
	Data                   string `json:"data"`
	DenominazioneRegione   string `json:"denominazione_regione"`
	DenominazioneProvincia string `json:"denominazione_provincia"`
	SiglaProvincia         string `json:"sigla_provincia"`
	TotaleCasi             int    `json:"totale_casi"`
	NoteIt                 string `json:"note_it"`
}

func (i *Nazione) Decode(data []byte, v interface{}) error {
	reader := csv.NewReader(bytes.NewReader(data))
	rr, err := reader.ReadAll()
	if err != nil {
		return errors.New("gojsonq: " + err.Error())
	}
	if len(rr) < 1 {
		return errors.New("gojsonq: csv data can't be empty! At least contain the header row")
	}
	var arr = make([]map[string]interface{}, 0)
	header := rr[0] // assume the very first row as header
	for i := 1; i <= len(rr)-1; i++ {
		if rr[i] == nil { // if a row is empty, skip it
			continue
		}
		mp := map[string]interface{}{}
		for j := 0; j < len(header); j++ {
			// convert data to different types
			// if header contains field like, ID|NUMBER,Name|String,IsStudent|BOOLEAN
			t := strings.Split(header[j], "|")
			var typ string
			if len(t) > 1 {
				typ = t[1]
			}
			hdr := strings.TrimSpace(t[0])
			switch typ {
			default:
				mp[hdr] = rr[i][j]

			case "STRING":
				mp[hdr] = rr[i][j]

			case "NUMBER":
				if fv, err := strconv.ParseFloat(rr[i][j], 64); err == nil {
					mp[hdr] = fv
				} else {
					mp[hdr] = 0.0
				}

			case "BOOLEAN":
				if strings.ToLower(rr[i][j]) == "true" ||
					rr[i][j] == "1" {
					mp[hdr] = true
				} else {
					mp[hdr] = false
				}

			}
		}
		arr = append(arr, mp)
	}
	bb, err := json.Marshal(arr)
	if err != nil {
		return fmt.Errorf("gojsonq: %v", err)
	}
	return json.Unmarshal(bb, &v)
}