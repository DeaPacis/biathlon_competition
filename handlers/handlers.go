package handlers

import (
	"biathlon_competition/lib"
	"biathlon_competition/models"
	"fmt"
	"log"
	"os"
	"time"
)

func HandleRegistration(event models.Events, competitors map[int]models.Competitor, finalReports map[int]models.FinalReport,
	detectStarts *[]models.DetectStart, config models.Config, outputFile *os.File) {
	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) registered\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}

	if _, exists := competitors[event.CompetitorID]; !exists {
		competitors[event.CompetitorID] = models.Competitor{
			CompetitorID: event.CompetitorID,
			Status:       "NotStarted",
		}

		finalReports[event.CompetitorID] = models.FinalReport{
			CompetitorID: event.CompetitorID,
			Laps:         make([]models.LapInfo, config.Lap),
		}

		*detectStarts = append(*detectStarts, models.DetectStart{
			CompetitorID: event.CompetitorID,
		})
	}
}

func HandleStartDraw(event models.Events, competitors map[int]models.Competitor, detectStarts *[]models.DetectStart,
	config models.Config, outputFile *os.File) {
	_, err := fmt.Fprintf(outputFile, "[%s] The start time for the competitor(%d) was set by a draw to %s\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID, event.ExtraParams)
	if err != nil {
		log.Println(err)
		return
	}

	startTime, err := time.Parse("15:04:05", event.ExtraParams)
	if err != nil {
		log.Println(err)
		return
	}
	startDelta := lib.ParseStringToDuration(config.StartDelta)

	for i := range *detectStarts {
		if (*detectStarts)[i].CompetitorID == event.CompetitorID {
			(*detectStarts)[i].StartDeadline = startTime.Add(startDelta)
		}
	}

	competitor := competitors[event.CompetitorID]
	competitor.AssignedStart, err = time.Parse("15:04:05.000", event.ExtraParams)
	if err != nil {
		log.Println(err)
		return
	}
	competitor.StartCurrentLap = competitor.AssignedStart
	competitors[event.CompetitorID] = competitor
}

func HandleStartLine(event models.Events, outputFile *os.File) {
	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) is on the start line\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}
}

func HandleStart(event models.Events, outputFile *os.File, detectStarts *[]models.DetectStart, competitors map[int]models.Competitor) {

	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) has started\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}

	for i := len(*detectStarts) - 1; i >= 0; i-- {
		if (*detectStarts)[i].CompetitorID == event.CompetitorID {
			*detectStarts = append((*detectStarts)[:i], (*detectStarts)[i+1:]...)
		}
	}

	if event.Time.Before(competitors[event.CompetitorID].AssignedStart) {
		_, err = fmt.Fprintf(outputFile, "[%s] The competitor(%d) is disqualified\n",
			event.Time.Format("15:04:05.000"), event.CompetitorID)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	competitor := competitors[event.CompetitorID]
	competitor.StartTime = event.Time
	competitor.Status = "NotFinished"
	competitors[event.CompetitorID] = competitor
}

func HandleEnterFiringRange(event models.Events, outputFile *os.File, finalReports map[int]models.FinalReport) {
	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) is on the firing range(%s)\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID, event.ExtraParams)
	if err != nil {
		log.Println(err)
		return
	}

	report := finalReports[event.CompetitorID]
	report.ShotsNumber += 5

	finalReports[event.CompetitorID] = report
}

func HandleHitTarget(event models.Events, outputFile *os.File, finalReports map[int]models.FinalReport) {
	_, err := fmt.Fprintf(outputFile, "[%s] The target(%s) has been hit by competitor(%d)\n",
		event.Time.Format("15:04:05.000"), event.ExtraParams, event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}

	report := finalReports[event.CompetitorID]
	report.HitsNumber++
	finalReports[event.CompetitorID] = report
}

func HandleLeaveFiringRange(event models.Events, outputFile *os.File) {
	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) left the firing range\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}
}

func HandleEnterPenaltyLap(event models.Events, outputFile *os.File, competitors map[int]models.Competitor) {
	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) entered the penalty laps\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}

	competitor := competitors[event.CompetitorID]
	competitor.StartCurrentPenaltyLaps = event.Time
	competitors[event.CompetitorID] = competitor
}

func HandleLeavePenaltyLap(event models.Events, outputFile *os.File, competitors map[int]models.Competitor,
	finalReports map[int]models.FinalReport, config models.Config) {

	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) left the penalty laps\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}

	report := finalReports[event.CompetitorID]
	report.PenaltyLaps.Time += event.Time.Sub(competitors[event.CompetitorID].StartCurrentPenaltyLaps)
	report.PenaltyLaps.AverageSpeed = lib.CountPenaltyLapSpeed(report, config)
	finalReports[event.CompetitorID] = report
}

func HandleFinishMainLap(event models.Events, outputFile *os.File, competitors map[int]models.Competitor,
	finalReports map[int]models.FinalReport, config models.Config) {

	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) ended the main lap\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID)
	if err != nil {
		log.Println(err)
		return
	}

	report := finalReports[event.CompetitorID]
	currentLapTime := event.Time.Sub(competitors[event.CompetitorID].StartCurrentLap)
	report.Laps[competitors[event.CompetitorID].LapNumbers].Time = currentLapTime
	report.Laps[competitors[event.CompetitorID].LapNumbers].AverageSpeed = lib.CountLapSpeed(currentLapTime, config)
	report.TotalTime += currentLapTime
	finalReports[event.CompetitorID] = report

	competitor := competitors[event.CompetitorID]
	competitor.StartCurrentLap = event.Time
	competitor.LapNumbers++

	if competitor.LapNumbers == config.Lap {
		competitor.Status = "Finished"
		_, err = fmt.Fprintf(outputFile, "[%s] The competitor(%d) has finished\n",
			event.Time.Format("15:04:05.000"), event.CompetitorID)
		if err != nil {
			log.Println(err)
			return
		}
	}

	competitors[event.CompetitorID] = competitor
}

func HandleUnableToContinue(event models.Events, outputFile *os.File) {
	_, err := fmt.Fprintf(outputFile, "[%s] The competitor(%d) can`t continue: %s\n",
		event.Time.Format("15:04:05.000"), event.CompetitorID, event.ExtraParams)
	if err != nil {
		log.Println(err)
		return
	}
}
