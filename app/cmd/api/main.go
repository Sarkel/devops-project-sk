package main

import (
	"devops/app/internal/app"
	"log"
)

func main() {
	if err := app.StartApi(); err != nil {
		log.Fatal(err)
	}
}
