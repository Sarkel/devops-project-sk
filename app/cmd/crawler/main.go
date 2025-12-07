package main

import (
	"devops-project-sk/internal/app"
	"log"
)

func main() {
	if err := app.RunCrawler(); err != nil {
		log.Fatal(err)
	}
}
