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
	CasiTestati              int    `json:"casi_testati"`
	Data                     string `json:"data"`
	Deceduti                 int    `json:"deceduti"`
	DimessiGuariti           int    `json:"dimessi_guariti"`
	IsolamentoDomiciliare    int    `json:"isolamento_domiciliare"`
	NuoviPositivi            int    `json:"nuovi_positivi"`
	RicoveratiConSintomi     int    `json:"ricoverati_con_sintomi"`
	Tamponi                  int    `json:"tamponi"`
	TamponiTestMol           int    `json:"tamponi_test_molecolare"`
	TamponiTestAnt           int    `json:"tamponi_test_antigenico_rapido"`
	TerapiaIntensiva         int    `json:"terapia_intensiva"`
	TotaleCasi               int    `json:"totale_casi"`
	TotaleOspedalizzati      int    `json:"totale_ospedalizzati"`
	TotalePositivi           int    `json:"totale_positivi"`
	TotalePositiviTestAnt    int    `json:"totale_positivi_test_antigenico_rapido"`
	TotalePositiviTestMol    int    `json:"totale_positivi_test_molecolare"`
	VariazioneTotalePositivi int    `json:"variazione_totale_positivi"`
}

type Regione struct {
	CasiTestati              int    `json:"casi_testati"`
	Data                     string `json:"data"`
	Deceduti                 int    `json:"deceduti"`
	DenominazioneRegione     string `json:"denominazione_regione"`
	DimessiGuariti           int    `json:"dimessi_guariti"`
	IsolamentoDomiciliare    int    `json:"isolamento_domiciliare"`
	Note                     string `json:"note,omitempty"`
	NoteCasi                 string `json:"note_casi,omitempty"`
	NoteTest                 string `json:"note_test,omitempty"`
	NuoviPositivi            int    `json:"nuovi_positivi"`
	RicoveratiConSintomi     int    `json:"ricoverati_con_sintomi"`
	TerapiaIntensiva         int    `json:"terapia_intensiva"`
	TotaleCasi               int    `json:"totale_casi"`
	TotaleOspedalizzati      int    `json:"totale_ospedalizzati"`
	TotalePositivi           int    `json:"totale_positivi"`
	TotalePositiviTestAnt    int    `json:"totale_positivi_test_antigenico_rapido"`
	TotalePositiviTestMol    int    `json:"totale_positivi_test_molecolare"`
	Tamponi                  int    `json:"tamponi"`
	TamponiTestMol           int    `json:"tamponi_test_molecolare"`
	TamponiTestAnt           int    `json:"tamponi_test_antigenico_rapido"`
	VariazioneTotalePositivi int    `json:"variazione_totale_positivi"`
}

type Provincia struct {
	Data                   string `json:"data"`
	DenominazioneRegione   string `json:"denominazione_regione"`
	DenominazioneProvincia string `json:"denominazione_provincia"`
	Note                   string `json:"note"`
	SiglaProvincia         string `json:"sigla_provincia"`
	TotaleCasi             int    `json:"totale_casi"`
}

type Nota struct {
	Data string `json:"data"`
	Note string `json:"note"`
}
