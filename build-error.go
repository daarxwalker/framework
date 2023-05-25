package framework

import (
	"log"

	"github.com/fatih/color"
)

func check(err error) {
	if err == nil {
		return
	}
	redColor := color.New(color.FgRed)
	log.Fatalln(redColor.Sprintf("%s", err))
}
