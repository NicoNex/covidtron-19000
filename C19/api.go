/*
* covidtron-19000
* Copyright (C) 2020  Michele Dimaggio
*/

package C19

import (
	"io"
	"os"
	"fmt"
	"log"
	"bytes"
	"encoding/json"

	"github.com/NicoNex/echotron"
	"github.com/thedevsaddam/gojsonq/v2"
)

const JSON_PATH = "~/.config/covidtron-19000"

func Update() {
	var json_url = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-%s-latest.json"
	files := [3]string{"andamento-nazionale", "province", "regioni"}

	for _, value := range files {
		var url = fmt.Sprintf(json_url, value)

		var content []byte = echotron.SendGetRequest(url)

		fpath := fmt.Sprintf("%s/%s.json", JSON_PATH, value)
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

	fpath := fmt.Sprintf("%s/andamento-nazionale.json", JSON_PATH)
	search := gojsonq.New().File(fpath).First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func getRegione(regione string) Regione {
	var data Regione

	fpath := fmt.Sprintf("%s/regioni.json", JSON_PATH)
	search := gojsonq.New().File(fpath).Where("denominazione_regione", "=", regione).First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}

func getProvincia(provincia string) Provincia {
	var data Provincia

	fpath := fmt.Sprintf("%s/province.json", JSON_PATH)
	search := gojsonq.New().File(fpath).Where("denominazione_provincia", "=", provincia).First()
	bytes, _ := json.Marshal(search)
	json.Unmarshal(bytes, &data)
	return data
}
