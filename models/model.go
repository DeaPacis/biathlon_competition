package models

import "time"

type Config struct {
	Lap         int    `json:"laps"`
	LapLen      int    `json:"lapLen"`
	PenaltyLen  int    `json:"penaltyLen"`
	FiringLines int    `json:"firingLines"`
	Start       string `json:"start"`
	StartDelta  string `json:"startDelta"`
}

type Events struct {
	Time         time.Time
	EventID      int
	CompetitorID int
	ExtraParams  string
}

type Competitor struct {
	CompetitorID            int
	AssignedStart           time.Time
	StartTime               time.Time
	FinishTime              time.Time
	StartCurrentLap         time.Time
	StartCurrentPenaltyLaps time.Time
	LapNumbers              int
	Status                  string
}

type LapInfo struct {
	Time         time.Duration
	AverageSpeed float64
}

type FinalReport struct {
	CompetitorID int
	TotalTime    time.Duration
	Laps         []LapInfo
	PenaltyLaps  LapInfo
	HitsNumber   int
	ShotsNumber  int
}

type DetectStart struct {
	CompetitorID  int
	StartDeadline time.Time
}
