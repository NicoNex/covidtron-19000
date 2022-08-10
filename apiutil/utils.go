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

package apiutil

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/dustin/go-humanize"
)

func Update(url, path, filename string) {
	dir := fmt.Sprintf(path)
	_, err := os.Stat(dir)

	if err != nil {
		os.Mkdir(dir, 0755)
	}

	resp, err := http.Get(url)
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

	fpath := fmt.Sprintf("%s/%s", path, filename)
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

func Ifmt(i int) string {
	return humanize.FormatInteger("#.###,", i)
}

func FormatTimestamp(timestamp string, tzFix bool) (fmtTime string) {
	if !tzFix {
		timestamp += "Z"
	}

	tp, err := time.Parse(time.RFC3339, timestamp)

	if err != nil {
		log.Println(err)
	}

	if tzFix {
		tz, _ := time.LoadLocation("Europe/Rome")
		fmtTime = tp.In(tz).Format("15:04 del 02/01/2006")
	} else {
		fmtTime = tp.Format("15:04 del 02/01/2006")
	}

	return fmtTime
}
