package main

import (
	"github.com/go-co-op/gocron"
	"github.com/NicoNex/covidtron-19000/c19"
)

// func hasChanged(b []byte) bool {

// }

func updateData() {
	c19.Update()
	gocron.Every(1).Day().At("18:30").Do(c19.Update)
	
	<- gocron.Start()
}
