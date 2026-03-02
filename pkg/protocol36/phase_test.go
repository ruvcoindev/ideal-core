package protocol36

import (
	"testing"
	"time"
)

func TestCalculatePhase(t *testing.T) {
	now := time.Now()
	
	tests := []struct {
		name         string
		lastContact  time.Time
		expectedPhase PhaseName
		expectedDay  int
	}{
		{
			name:         "Day 0 (preparation)",
			lastContact:  now,
			expectedPhase: "Preparation",
			expectedDay:  0,
		},
		{
			name:         "Day 3 (Detox)",
			lastContact:  now.AddDate(0, 0, -3),
			expectedPhase: PhaseDetox,
			expectedDay:  3,
		},
		{
			name:         "Day 15 (Rewire)",
			lastContact:  now.AddDate(0, 0, -15),
			expectedPhase: PhaseRewire,
			expectedDay:  15,
		},
		{
			name:         "Day 30 (Integration)",
			lastContact:  now.AddDate(0, 0, -30),
			expectedPhase: PhaseIntegration,
			expectedDay:  30,
		},
	}
	
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			config, days := CalculatePhase(tc.lastContact, now)
			if config.Name != tc.expectedPhase {
				t.Errorf("Phase: got %q, want %q", config.Name, tc.expectedPhase)
			}
			if days != tc.expectedDay {
				t.Errorf("Days: got %d, want %d", days, tc.expectedDay)
			}
		})
	}
}

func TestGetNextPhase(t *testing.T) {
	if GetNextPhase(PhaseDetox) != PhaseRewire {
		t.Error("GetNextPhase failed: Detox → Rewire")
	}
	if GetNextPhase(PhaseRewire) != PhaseIntegration {
		t.Error("GetNextPhase failed: Rewire → Integration")
	}
}
