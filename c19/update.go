/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020 Michele Dimaggio, Alessandro Ianne.
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
	"os"
	"log"
	"io/ioutil"
	"encoding/csv"
	"encoding/json"
	"strconv"
	//"time"

	"github.com/NicoNex/echotron"
)

var datapath string

var json_url = "https://raw.githubusercontent.com/pcm-dpc/COVID-19/master/dati-json/dpc-covid19-ita-%s.json"

var csv_prefix = "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_time_series/"

func Update() {
	pcmdpcUpdate()
	gisandUpdate()
	gisandParser()
}

func pcmdpcUpdate() {
	
	files := [4]string{
		"andamento-nazionale-latest",
		"province-latest",
		"regioni-latest",
		"note-it"}

	dir := fmt.Sprintf(datapath)
	_, err := os.Stat(dir)
	if err != nil {
		os.Mkdir(dir, 0755)
	}

	for _, value := range files {
		
		url := fmt.Sprintf(json_url, value)
		var content []byte = echotron.SendGetRequest(url)

		fpath := fmt.Sprintf("%s/%s.json", datapath, value)

		if err := ioutil.WriteFile(fpath, content, 0755); err != nil {
			log.Println(err)
		}
	}
}

func gisandUpdate() {
	var csv_file = [3]string{
			"time_series_covid19_confirmed_global.csv",
			"time_series_covid19_deaths_global.csv",
			"time_series_covid19_recovered_global.csv"}

	dir := fmt.Sprintf("%s%s", datapath, "/gisanddata/")
	_, err := os.Stat(dir)
	if err != nil {
		os.Mkdir(dir, 0755)
	}

	for _, file := range csv_file {
		
		url := fmt.Sprintf("%s%s", csv_prefix, file)
		var content []byte = echotron.SendGetRequest(url)

		fpath := fmt.Sprintf("%s%s", dir, file)

		if err := ioutil.WriteFile(fpath, content, 0755); err != nil {
			log.Println(err)
		}
	}
}

func gisandParser() {
	csv_files := [3]string{
			"time_series_covid19_confirmed_global.csv",
			"time_series_covid19_deaths_global.csv",
			"time_series_covid19_recovered_global.csv"}

	dir := fmt.Sprintf("%s%s", datapath, "/gisanddata/")
	nations := make(map[int]GisandData)

	_, err := os.Stat(dir)
	if err != nil {
		os.Mkdir(dir, 0755)
	}

	for _, file := range csv_files {
		fpath := fmt.Sprintf("%s%s", dir, file)

		data, err := os.Open(fpath)
		if err != nil {
			log.Println(err)
		}

		reader := csv.NewReader(data)

		csv, _ := reader.ReadAll()
		day_number := len(csv[0]) - 4

		for i, csv_nation := range csv[1:]{
			
			nation := GisandData{}
			if _, value := nations[i]; !value {
				nation.Province_State = csv_nation[0]
				nation.Country_Region = csv_nation[1]
				nation.Lat, _ = strconv.ParseFloat(csv_nation[2], 64)
				nation.Long, _  = strconv.ParseFloat(csv_nation[3], 64)
				nation.DailyData = make([]DailyData, day_number)
			} else {
				nation = nations[i]
			}
			
			for j := 4; j < len(csv_nation); j++ {
				
				if nation.DailyData[j-4].Date == "" {
					nation.DailyData[j-4].Date = csv[0][j]		
				}
				
				people, _ := strconv.Atoi(csv_nation[j])
				
				switch file {
				case csv_files[0]:
					if nation.DailyData[j-4].Confirmed < people { nation.DailyData[j-4].Confirmed = people }

				case csv_files[1]:
					if nation.DailyData[j-4].Deaths < people { nation.DailyData[j-4].Deaths = people }

				case csv_files[2]:
					if nation.DailyData[j-4].Recovered < people { nation.DailyData[j-4].Recovered = people }

				}
			}
			nations[i] = nation
		}
	}

	content, _ := json.Marshal(nations)
	fpath := fmt.Sprintf("%s/gisanddata.json", datapath)

	if err := ioutil.WriteFile(fpath, content, 0755); err != nil {
		log.Println(err)
	}
	
}
