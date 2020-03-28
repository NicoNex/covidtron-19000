package main

import (
	"time"

	"github.com/go-co-op/gocron"
	"github.com/NicoNex/covidtron-19000/c19"
)

// func hasChanged(b []byte) bool {

// }

func updateData() {
	c19.Update()

	location, _ := time.LoadLocation("Europe/Rome")
	scheduler := gocron.NewScheduler(location)
	scheduler.Every(1).Day().At("18:30").Do(c19.Update)
	
	<- scheduler.Start()
}
