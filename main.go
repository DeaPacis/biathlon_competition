package main

import (
	"biathlon_competition/models"
	"biathlon_competition/utils"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {

	configFile := flag.String("c", "", "Path to configuration file (required)")
	eventsFile := flag.String("f", "", "Path to events file (required)")

	flag.Parse()

	if *configFile == "" || *eventsFile == "" {
		fmt.Println("Error: Required flags not provided")
		fmt.Println("Usage: go run main.go -c <config> -f <events>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	fmt.Printf("Config file: %s\n", *configFile)
	fmt.Printf("Events file: %s\n", *eventsFile)

	var config models.Config

	configData, err := os.ReadFile(*configFile)
	if err != nil {
		log.Println(err)
		return
	}

	err = json.Unmarshal(configData, &config)
	if err != nil {
		log.Println(err)
		return
	}

	events, err := os.Open(*eventsFile)
	if err != nil {
		log.Println(err)
		return
	}
	defer events.Close()

	finalReport, competitors := utils.ParseEventsFile(events, config)

	utils.ResultTable(finalReport, competitors)
}
