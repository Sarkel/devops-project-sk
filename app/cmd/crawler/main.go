package main

import (
	"devops/app/internal/app"
	"log"
)

func main() {
	if err := app.RunCrawler(); err != nil {
		log.Fatal(err)
	}
}
