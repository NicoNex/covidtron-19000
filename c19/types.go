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
