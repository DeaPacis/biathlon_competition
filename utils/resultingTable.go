package utils

import (
	"biathlon_competition/models"
	"fmt"
	"log"
	"os"
	"sort"
)

func ResultTable(finalReport map[int]models.FinalReport, competitors map[int]models.Competitor) {
	resultFile, err := os.Create("Resulting Table.txt")
	if err != nil {
		log.Println(err)
		return
	}
	defer resultFile.Close()

	sortedReports := sortReports(finalReport, competitors)

	for i := 0; i < len(sortedReports); i++ {
		if competitors[sortedReports[i].CompetitorID].Status == "Finished" {
			_, err = fmt.Fprintf(resultFile, "[%s] %d %s %s %d/%d\n", models.FormatDuration(sortedReports[i].TotalTime),
				sortedReports[i].CompetitorID, models.FormatLapList(sortedReports[i].Laps), sortedReports[i].PenaltyLaps, sortedReports[i].HitsNumber,
				sortedReports[i].ShotsNumber)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			_, err = fmt.Fprintf(resultFile, "[%s] %d %s %s %d/%d\n", competitors[sortedReports[i].CompetitorID].Status,
				sortedReports[i].CompetitorID, models.FormatLapList(sortedReports[i].Laps), sortedReports[i].PenaltyLaps, sortedReports[i].HitsNumber,
				sortedReports[i].ShotsNumber)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}

func sortReports(reports map[int]models.FinalReport, competitors map[int]models.Competitor) []models.FinalReport {
	finishedReports := make([]models.FinalReport, 0)
	unfinishedReports := make([]models.FinalReport, 0)
	for _, report := range reports {
		if competitors[report.CompetitorID].Status == "Finished" {
			finishedReports = append(finishedReports, report)
		} else {
			unfinishedReports = append(unfinishedReports, report)
		}
	}

	sort.Slice(finishedReports, func(i, j int) bool {
		return finishedReports[i].TotalTime < finishedReports[j].TotalTime
	})

	for _, unfinishedReport := range unfinishedReports {
		finishedReports = append(finishedReports, unfinishedReport)
	}

	return finishedReports
}
