package main

import (
	"log"

	"devops/app/internal/app"
)

func main() {
	if err := app.StartApi(); err != nil {
		log.Fatal(err)
	}
}
