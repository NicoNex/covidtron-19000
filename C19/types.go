/*
* covidtron-19000
* Copyright (C) 2020  Michele Dimaggio
*/

package C19

type Andamento struct {
	Data                      string `json:"data"`
	Deceduti                  int    `json:"deceduti"`
	DimessiGuariti            int    `json:"dimessi_guariti"`
	IsolamentoDomiciliare     int    `json:"isolamento_domiciliare"`
	NuoviAttualmentePositivi  int    `json:"nuovi_attualmente_positivi"`
	RicoveratiConSintomi      int    `json:"ricoverati_con_sintomi"`
	Tamponi                   int    `json:"tamponi"`
	TerapiaIntensiva          int    `json:"terapia_intensiva"`
	TotaleAttualmentePositivi int    `json:"totale_attualmente_positivi"`
	TotaleCasi                int    `json:"totale_casi"`
	TotaleOspedalizzati       int    `json:"totale_ospedalizzati"`
}

	
type Regione struct {
	Data                      string  `json:"data"`
	DenominazioneRegione      string  `json:"denominazione_regione"`
	RicoveratiConSintomi      int     `json:"ricoverati_con_sintomi"`
	TerapiaIntensiva          int     `json:"terapia_intensiva"`
	TotaleOspedalizzati       int     `json:"totale_ospedalizzati"`
	IsolamentoDomiciliare     int     `json:"isolamento_domiciliare"`
	TotaleAttualmentePositivi int     `json:"totale_attualmente_positivi"`
	NuoviAttualmentePositivi  int     `json:"nuovi_attualmente_positivi"`
	DimessiGuariti            int     `json:"dimessi_guariti"`
	Deceduti                  int     `json:"deceduti"`
	TotaleCasi                int     `json:"totale_casi"`
	Tamponi                   int     `json:"tamponi"`
}

type Provincia struct {
	Data                   string  `json:"data"`
	DenominazioneRegione   string  `json:"denominazione_regione"`
	DenominazioneProvincia string  `json:"denominazione_provincia"`
	TotaleCasi             int     `json:"totale_casi"`
}
