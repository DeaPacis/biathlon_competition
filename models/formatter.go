package models

import (
	"fmt"
	"strings"
	"time"
)

func (c Competitor) String() string {
	return fmt.Sprintf(
		"%d %s %s %s %s %s %d %s",
		c.CompetitorID,
		c.AssignedStart.Format("15:04:05.000"),
		c.StartTime.Format("15:04:05.000"),
		c.FinishTime.Format("15:04:05.000"),
		c.StartCurrentLap.Format("15:04:05.000"),
		c.StartCurrentPenaltyLaps.Format("15:04:05.000"),
		c.LapNumbers,
		c.Status,
	)
}

func (l LapInfo) String() string {
	if l.Time == 0 && l.AverageSpeed == 0 {
		return "{,}"
	}
	return fmt.Sprintf(
		"{%s, %.3f}",
		FormatDuration(l.Time),
		l.AverageSpeed,
	)
}

func (f FinalReport) String() string {
	return fmt.Sprintf(
		"{%d %s %s %s %d %d}",
		f.CompetitorID,
		FormatDuration(f.TotalTime),
		FormatLapList(f.Laps),
		f.PenaltyLaps,
		f.HitsNumber,
		f.ShotsNumber,
	)
}

func FormatDuration(d time.Duration) string {
	sign := ""
	if d < 0 {
		sign = "-"
		d = -d
	}

	hours := int(d.Hours())
	d -= time.Duration(hours) * time.Hour

	minutes := int(d.Minutes())
	d -= time.Duration(minutes) * time.Minute

	seconds := int(d.Seconds())
	d -= time.Duration(seconds) * time.Second

	milliseconds := int(d.Milliseconds())

	return fmt.Sprintf(
		"%s%02d:%02d:%02d.%03d",
		sign,
		hours,
		minutes,
		seconds,
		milliseconds,
	)
}

func FormatLapList(laps []LapInfo) string {
	var result strings.Builder
	result.WriteString("[")
	for i, lap := range laps {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(lap.String())
	}
	result.WriteString("]")
	return result.String()
}
