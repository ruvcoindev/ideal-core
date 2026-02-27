package time

import (
	"testing"
	"time"
)

func TestDetectCyclePhase(t *testing.T) {
	now := time.Now()
	lastContact := now.AddDate(0, 0, -5)
	breakupDate := time.Time{}
	
	phase := DetectCyclePhase(lastContact, breakupDate)
	
	if phase != PhasePeak {
		t.Errorf("Expected PhasePeak, got %s", phase)
	}
}

func TestGetRecommendation(t *testing.T) {
	rec := GetRecommendation(PhaseThaw)
	if rec == "" {
		t.Error("Expected non-empty recommendation")
	}
}
