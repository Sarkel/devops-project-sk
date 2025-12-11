package main

import (
	"devops/app/internal/app"
	"log"
)

func main() {
	if err := app.StartReader(); err != nil {
		log.Fatal(err)
	}
}
