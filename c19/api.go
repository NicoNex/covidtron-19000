/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020-2021 Michele Dimaggio.
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
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/thedevsaddam/gojsonq/v2"
)

type NoteType uint8

const (
	Note NoteType = iota
	NoteCasi
	NoteTest
)

var jsonpath string

func Update() {
	var json_url = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-%s.json"
	var files = [4]string{"andamento-nazionale-latest", "province-latest", "regioni-latest", "note"}

	dir := fmt.Sprintf(jsonpath)
	_, err := os.Stat(dir)
	if err != nil {
		os.Mkdir(dir, 0755)
	}

	for _, value := range files {
		resp, err := http.Get(fmt.Sprintf(json_url, value))
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		content, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}

		fpath := fmt.Sprintf("%s/%s.json", jsonpath, value)
		data, err := os.Create(fpath)

		if err != nil {
			log.Println(err)
		}
		defer data.Close()

		_, err = io.Copy(data, bytes.NewReader(content))

		if err != nil {
			log.Println(err)
		}
	}
}

func getAndamento() Andamento {
	var data Andamento

	fpath := fmt.Sprintf("%s/andamento-nazionale-latest.json", jsonpath)
	search := gojsonq.New().
		File(fpath).
		First()

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func getRegione(regione string) *Regione {
	var data Regione

	fpath := fmt.Sprintf("%s/regioni-latest.json", jsonpath)
	search := gojsonq.New().
		File(fpath).
		WhereContains("denominazione_regione", regione).
		First()

	if search == nil {
		return nil
	}

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return &data
}

func GetRegioni() []string {
	var data []string

	fpath := fmt.Sprintf("%s/regioni-latest.json", jsonpath)
	search := gojsonq.New().
		File(fpath).
		Pluck("denominazione_regione")

	for _, v := range search.([]interface{}) {
		data = append(data, v.(string))
	}

	sort.Strings(data)

	return data
}

func getProvincia(provincia string) *Provincia {
	var data Provincia

	fpath := fmt.Sprintf("%s/province-latest.json", jsonpath)

	if strings.Contains(provincia, "(") {
		provincia = provincia[:len(provincia)-5]
	}

	search := gojsonq.New().
		File(fpath).
		WhereContains("denominazione_provincia", provincia).
		First()

	if search == nil {
		return nil
	}

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return &data
}

func GetProvince(regione string) []string {
	var data []string

	fpath := fmt.Sprintf("%s/province-latest.json", jsonpath)

	searchProv := gojsonq.New().
		File(fpath).
		WhereContains("denominazione_regione", regione).
		Pluck("denominazione_provincia")

	searchSigle := gojsonq.New().
		File(fpath).
		WhereContains("denominazione_regione", regione).
		Pluck("sigla_provincia")

	for i, v := range searchProv.([]interface{}) {
		if v != "Fuori Regione / Provincia Autonoma" && v != "In fase di definizione/aggiornamento" {
			data = append(data, fmt.Sprintf("%s (%s)", v, searchSigle.([]interface{})[i].(string)))
		}
	}

	sort.Strings(data)

	return data
}

func getNote() Nota {
	var data Nota

	fpath := fmt.Sprintf("%s/note.json", jsonpath)

	search := gojsonq.New().
		File(fpath).
		Last()

	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func init() {
	jsonpath = fmt.Sprintf("%s/.cache/covidtron-19000", os.Getenv("HOME"))
}
