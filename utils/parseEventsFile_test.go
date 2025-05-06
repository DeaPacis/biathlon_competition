package utils

import (
	"biathlon_competition/models"
	"os"
	"testing"
)

func TestParseEventsFile(t *testing.T) {
	content := `[12:00:00.000] 1 101
				[12:00:05.000] 2 101 12:01:00
				[12:00:10.000] 3 101
				[12:01:01.000] 4 101
				[12:02:00.000] 5 101 1
				[12:02:01.000] 6 101 1
				[12:02:02.000] 6 101 2
				[12:02:03.000] 7 101
				[12:02:10.000] 8 101
				[12:02:20.000] 9 101
				[12:03:00.000] 10 101
				[12:03:30.000] 11 101 injury`

	tmpFile, err := os.CreateTemp("", "test_events_*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.WriteString(content)
	if err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	_, err = tmpFile.Seek(0, 0)
	if err != nil {
		t.Fatalf("failed to seek temp file: %v", err)
	}

	config := models.Config{
		StartDelta: "00:00:30",
		Lap:        1,
	}

	finalReports, competitors := ParseEventsFile(tmpFile, config)

	if len(competitors) != 1 {
		t.Errorf("expected 1 competitor, got %d", len(competitors))
	}

	comp, ok := competitors[101]
	if !ok {
		t.Fatalf("competitor 101 not found")
	}
	if comp.CompetitorID != 101 {
		t.Errorf("expected CompetitorID 101, got %d", comp.CompetitorID)
	}

	report, ok := finalReports[101]
	if !ok {
		t.Fatalf("final report for competitor 101 not found")
	}
	if report.ShotsNumber != 5 {
		t.Errorf("expected 5 shots, got %d", report.ShotsNumber)
	}
	if report.HitsNumber != 2 {
		t.Errorf("expected 2 hits, got %d", report.HitsNumber)
	}
	if report.PenaltyLaps.Time == 0 {
		t.Errorf("expected non-zero penalty lap time")
	}
	if report.TotalTime == 0 {
		t.Errorf("expected non-zero total time")
	}
}
