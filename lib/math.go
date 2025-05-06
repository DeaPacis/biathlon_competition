package lib

import (
	"biathlon_competition/models"
	"math"
	"time"
)

func CountLapSpeed(currentLapTime time.Duration, config models.Config) float64 {
	lapTime := currentLapTime.Seconds()
	averageSpeed := float64(config.LapLen) / lapTime
	averageSpeed = math.Trunc(averageSpeed*1000) / 1000

	return averageSpeed
}

func CountPenaltyLapSpeed(currentFinalReport models.FinalReport, config models.Config) float64 {
	penaltyLapsTime := currentFinalReport.PenaltyLaps.Time.Seconds()
	penaltyLapsLen := config.PenaltyLen * (currentFinalReport.ShotsNumber - currentFinalReport.HitsNumber)
	averageSpeed := float64(penaltyLapsLen) / penaltyLapsTime
	averageSpeed = math.Trunc(averageSpeed*1000) / 1000

	return averageSpeed
}
