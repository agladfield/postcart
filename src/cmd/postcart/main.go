package main

import (
	"log"
	"os"

	"github.com/agladfield/postcart/pkg/postcart"
)

func main() {
	err := postcart.Run()
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)
}

// © Arthur Gladfield
