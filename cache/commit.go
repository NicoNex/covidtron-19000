package cache

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)


func (c Cache) GetSha() string {
	resp, err := http.Get("https://api.github.com/repos/pcm-dpc/COVID-19/commits/master")

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
