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

package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Commits struct {
	C19 string `json:"c19"`
	Vax string `json:"vax"`
}

func (c Cache) UpdateCommits() Commits {
	return Commits{
		C19: getSha("https://api.github.com/repos/pcm-dpc/COVID-19/commits/master"),
		Vax: getSha("https://api.github.com/repos/italia/covid19-opendata-vaccini/commits/master"),
	}
}

func getSha(url string) string {
	resp, err := http.Get(url)

	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	var gh struct {
		Sha string `json:"sha"`
	}
	err = json.Unmarshal(body, &gh)

	if err != nil {
		log.Println(err)
	}

	return gh.Sha
}
