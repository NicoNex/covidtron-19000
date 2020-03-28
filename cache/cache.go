/*
 * Covidtron-19000 - a bot for monitoring data about COVID-19.
 * Copyright (C) 2020 Nicolò Santamaria.
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
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Cache struct {
	botName  string
	Sessions []int64 `json:"sessions"`
}

var cachepath string

func NewCache(bname string) *Cache {
	var cache = &Cache{botName: bname}

	data, err := ioutil.ReadFile(cachepath)
	if err != nil {
		log.Println(err)
		goto exit
	}

	err = json.Unmarshal(data, cache)
	if err != nil {
		log.Println(err)
	}

exit:
	return cache
}

func (c Cache) isin(s int64) bool {
	for _, v := range c.Sessions {
		if s == v {
			return true
		}
	}

	return false
}

func (c *Cache) SaveSession(s int64) {
	if !c.isin(s) {
		c.Sessions = append(c.Sessions, s)

		b, err := json.Marshal(c)
		if err != nil {
			log.Println(err)
			return
		}

		err = ioutil.WriteFile(cachepath, b, 0644)
		if err != nil {
			log.Println(err)
		}
	}
}

func (c *Cache) DelSession(s int64) {
	for k, v := range c.Sessions {
		if v == s {
			c.Sessions = append(c.Sessions[:k], c.Sessions[k+1:]...)
			break
		}
	}

	b, err := json.Marshal(c)
	if err != nil {
		log.Println(err)
		return
	}

	err = ioutil.WriteFile(cachepath, b, 0644)
	if err != nil {
		log.Println(err)
	}
}

func (c Cache) GetSessions() []int64 {
	return c.Sessions
}

func (c Cache) CountSessions() int {
	return len(c.Sessions)
}

func init() {
	ccdir := fmt.Sprintf("%s/.cache/covidtron-19000/", os.Getenv("HOME"))
	if _, err := os.Stat(ccdir); os.IsNotExist(err) {
		os.Mkdir(ccdir, 0755)
	}

	cachepath = fmt.Sprintf("%s/.cache/covidtron-19000/cache.json", os.Getenv("HOME"))
}
