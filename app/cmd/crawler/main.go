package main

import (
	"devops-project-sk/internal/app"
	"log"
	"time"
)

func main() {
	log.Println("Crawler service started.")

	runTask()

	ticker := time.NewTicker(30 * time.Minute)
	for range ticker.C {
		runTask()
	}
}

func runTask() {
	if err := app.RunCrawler(); err != nil {
		log.Fatal(err)
	}
}
