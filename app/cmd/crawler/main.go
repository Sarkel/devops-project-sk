package main

import (
	"log"

	"devops/app/internal/app"
)

func main() {
	if err := app.RunCrawler(); err != nil {
		log.Fatal(err)
	}
}
