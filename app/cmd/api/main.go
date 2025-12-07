package main

import (
	"devops-project-sk/internal/app"
	"log"
)

func main() {
	if err := app.StartApi(); err != nil {
		log.Fatal(err)
	}
}
