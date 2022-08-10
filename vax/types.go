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

var Area = map[string]string{
	"Abruzzo":               "ABR",
	"Basilicata":            "BAS",
	"Calabria":              "CAL",
	"Campania":              "CAM",
	"Emilia-Romagna":        "EMR",
	"Friuli Venezia Giulia": "FVG",
	"Lazio":                 "LAZ",
	"Liguria":               "LIG",
	"Lombardia":             "LOM",
	"Marche":                "MAR",
	"Molise":                "MOL",
	"P.A. Bolzano":          "PAB",
	"P.A. Trento":           "PAT",
	"Piemonte":              "PIE",
	"Puglia":                "PUG",
	"Sardegna":              "SAR",
	"Sicilia":               "SIC",
	"Toscana":               "TOS",
	"Umbria":                "UMB",
	"Valle d'Aosta":         "VDA",
	"Veneto":                "VEN",
}

type LastUpdate struct {
	Timestamp string `json:"ultimo_aggiornamento"`
}

type Vaccini struct {
	Area              string  `json:"area"`
	DosiConsegnate    int     `json:"dosi_consegnate"`
	DosiSomministrate int     `json:"dosi_somministrate"`
	Percentuale       float32 `json:"percentuale_somministrazione"`
}

type Somministrazioni struct {
	DataSomministrazione string `json:"data"`
	Fornitore            string `json:"forn"`
	PrimaDose            int    `json:"d1"`
	SecondaDose          int    `json:"d2"`
}
