package main

import (
	"github.com/NicoNex/covidtron-19000/c19"
	"time"
)

func updateData() {
	var timestamp int64

	c19.Update()
	timestamp = time.Now().Unix()

	for {
		t := time.Now().Unix()

		if t-timestamp > 86400 {
			c19.Update()
			timestamp = t
		}

		time.Sleep(time.Hour)
	}
}
