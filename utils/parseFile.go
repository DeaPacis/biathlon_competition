package utils

import (
	"biathlon_competition/handlers"
	"biathlon_competition/models"
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func ParseEventsFile(filename *os.File, config models.Config) (map[int]models.FinalReport, map[int]models.Competitor) {
	var event models.Events
	var detectStarts []models.DetectStart

	competitors := make(map[int]models.Competitor)
	finalReports := make(map[int]models.FinalReport)

	outputFile, err := os.Create("Output log.txt")
	if err != nil {
		log.Println(err)
		return nil, nil
	}
	defer outputFile.Close()

	scanner := bufio.NewScanner(filename)

	for scanner.Scan() {
		eventParts := strings.Split(scanner.Text(), " ")
		parseEvent(&event, eventParts)

		controlStarts(event, detectStarts, outputFile)

		switch event.EventID {
		case 1:
			handlers.HandleRegistration(event, competitors, finalReports, &detectStarts, config, outputFile)
		case 2:
			handlers.HandleStartDraw(event, competitors, &detectStarts, config, outputFile)
		case 3:
			handlers.HandleStartLine(event, outputFile)
		case 4:
			handlers.HandleStart(event, outputFile, &detectStarts, competitors)
		case 5:
			handlers.HandleEnterFiringRange(event, outputFile, finalReports)
		case 6:
			handlers.HandleHitTarget(event, outputFile, finalReports)
		case 7:
			handlers.HandleLeaveFiringRange(event, outputFile)
		case 8:
			handlers.HandleEnterPenaltyLap(event, outputFile, competitors)
		case 9:
			handlers.HandleLeavePenaltyLap(event, outputFile, competitors, finalReports, config)
		case 10:
			handlers.HandleFinishMainLap(event, outputFile, competitors, finalReports, config)
		case 11:
			handlers.HandleUnableToContinue(event, outputFile)
		}
	}
	return finalReports, competitors
}

func controlStarts(event models.Events, detectStarts []models.DetectStart, outputFile *os.File) {
	for i := len(detectStarts) - 1; i >= 0; i-- {
		if event.Time.After(detectStarts[i].StartDeadline) {
			_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) is disqualified\n",
				event.Time.Format("15:04:05.000"), detectStarts[i].CompetitorID)
			if err != nil {
				log.Println(err)
				return
			}
			detectStarts = append(detectStarts[:i], detectStarts[i+1:]...)
		}
	}
}

func parseEvent(event *models.Events, eventParts []string) {
	var err error
	eventParts[0] = strings.TrimSpace(eventParts[0])
	eventParts[0] = strings.Trim(eventParts[0], "[]")

	if len(eventParts[0]) == 8 {
		eventParts[0] += ".000"
	}

	event.Time, err = time.Parse("15:04:05.000", eventParts[0])
	if err != nil {
		log.Println(err)
		return
	}
	event.EventID, err = strconv.Atoi(eventParts[1])
	if err != nil {
		log.Println(err)
		return
	}
	event.CompetitorID, err = strconv.Atoi(eventParts[2])
	if err != nil {
		log.Println(err)
		return
	}
	event.ExtraParams = ""
	if len(eventParts) == 4 {
		event.ExtraParams += eventParts[3]
	} else if len(eventParts) > 4 {
		for i := 3; i < len(eventParts); i++ {
			event.ExtraParams += eventParts[i] + " "
		}
	}
}
