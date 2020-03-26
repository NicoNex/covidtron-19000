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
	NoteIt                    string `json:"note_it"`
	// NoteEn                    string `json:"note_en"`
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
	// Stato                     string  `json:"stato"`
	// CodiceRegione             int     `json:"codice_regione"`
	DenominazioneRegione      string  `json:"denominazione_regione"`
	// Lat                       float64 `json:"lat"`
	// Long                      float64 `json:"long"`
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
	NoteIt                    string  `json:"note_it"`
	// NoteEn                    string  `json:"note_en"`
}

type Provincia struct {
	Data                   string  `json:"data"`
	// Stato                  string  `json:"stato"`
	// CodiceRegione          int     `json:"codice_regione"`
	DenominazioneRegione   string  `json:"denominazione_regione"`
	// CodiceProvincia        int     `json:"codice_provincia"`
	DenominazioneProvincia string  `json:"denominazione_provincia"`
	// SiglaProvincia         string  `json:"sigla_provincia"`
	// Lat                    float64 `json:"lat"`
	// Long                   float64 `json:"long"`
	TotaleCasi             int     `json:"totale_casi"`
	NoteIt                 string  `json:"note_it"`
	// NoteEn                 string  `json:"note_en"`
}
