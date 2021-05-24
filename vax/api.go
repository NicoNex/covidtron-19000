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
	"os"
	"log"

	"github.com/NicoNex/covidtron-19000/apiutil"
	jsoniter "github.com/json-iterator/go"
)

var jsonpath string
var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Update() {
	var json_url = "https://raw.githubusercontent.com/italia/covid19-opendata-vaccini/master/dati/%s"
	var files = []string{
		"last-update-dataset.json",
		"somministrazioni-vaccini-latest.json",
		"vaccini-summary-latest.json",
	}

	for _, value := range files {
		url := fmt.Sprintf(json_url, value)
		apiutil.Update(url, jsonpath, value)
	}
}

func init() {
	jsonpath = fmt.Sprintf("%s/.cache/covidtron-19000", os.Getenv("HOME"))
}

func getTimestamp() string {
	var lastUpdate LastUpdate

	data, err := os.ReadFile(fmt.Sprintf("%s/last-update-dataset.json", jsonpath))
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(data, &lastUpdate)
	return lastUpdate.Timestamp
}

func getVaxData() []Vaccini {
	var vaxData struct {
    	Data []Vaccini `json:"data"`
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/vaccini-summary-latest.json", jsonpath))
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(data, &vaxData)
	return vaxData.Data
}

func getRegione(regName string) Vaccini {
	for _, field := range getVaxData() {
		if field.Area == Area[regName] {
			return field
		}
	}

	return Vaccini{}
}

func getTotaleSomministrazioni() (sum int) {
	vaxData := getVaxData()

	for _, field := range vaxData {
		sum += field.DosiSomministrate
	}

	return
}

func getTotalePercentuale(sum int) string {
	percent := float32(sum) / 59257566 * 100
	return fmt.Sprintf("%.2f%%", percent)
}

func getTotaleVaccinati() (sum int) {
	var vaxData struct {
		Data []Somministrazioni `json:"data"`
	}

	data, err := os.ReadFile(fmt.Sprintf("%s/somministrazioni-vaccini-latest.json", jsonpath))
	if err != nil {
		log.Println(err)
	}

	json.Unmarshal(data, &vaxData)

	for _, field := range vaxData.Data {
		if field.Fornitore == "Janssen" {
			sum += field.PrimaDose
		} else {
			sum += field.SecondaDose
		}
	}

	return
}
